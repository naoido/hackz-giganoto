package chatapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"object-t.com/hackz-giganoto/microservices/chat/gen/chat"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"goa.design/clue/log"
	"goa.design/goa/v3/security"
)

const (
	roomsKey       = "rooms"
	roomPublishKey = "room_publish"
	historyKey     = "history"
)

type chatsrvc struct {
	redis     *redis.Client
	jwtSecret []byte
}

func NewChat() chat.Service {
	// Setup Redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("secret")
	}

	return &chatsrvc{
		redis:     redisClient,
		jwtSecret: jwtSecret,
	}
}

func (s *chatsrvc) extractUserID(ctx context.Context, tokenString string) (string, error) {
	claims := make(jwt.MapClaims)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid sub claim")
	}

	return sub, nil
}

func (s *chatsrvc) JWTAuth(ctx context.Context, tokenString string, scheme *security.JWTScheme) (context.Context, error) {
	userID, err := s.extractUserID(ctx, tokenString)
	if err != nil {
		log.Print(ctx, log.KV{"chat.jwt_auth", "ERROR: failed to extract user ID"}, log.KV{"error", err.Error()})
		return ctx, chat.Unauthorized("invalid token")
	}

	ctx = context.WithValue(ctx, "user_id", userID)

	log.Info(ctx, log.KV{"chat.jwt_auth", fmt.Sprintf("authenticated user %s", userID)})
	return ctx, nil
}

func (s *chatsrvc) CreateRoom(ctx context.Context, p *chat.CreateRoomPayload) (res string, err error) {
	log.Info(ctx, log.KV{"chat.create_room", "creating new room"})

	newRoomId := uuid.New().String()
	if err := s.redis.SAdd(ctx, roomsKey, newRoomId).Err(); err != nil {
		log.Print(ctx, log.KV{"chat.create_room", "ERROR: redis SAdd failed"}, log.KV{"error", err.Error()})
		return "", chat.Internal("Internal server error")
	}

	return newRoomId, nil
}

func (s *chatsrvc) History(ctx context.Context, p *chat.HistoryPayload) (res []*chat.Chat, err error) {
	log.Info(ctx, log.KV{"chat.history", "retrieving chat history"})

	histories, err := s.redis.LRange(ctx, fmt.Sprintf("%s:%s", historyKey, p.RoomID), 0, -1).Result()
	if err != nil {
		log.Print(ctx, log.KV{"chat.history", "ERROR: redis LRange failed"}, log.KV{"error", err.Error()})
		return nil, chat.Internal("Internal server error")
	}

	if len(histories) == 0 {
		return make([]*chat.Chat, 0), nil
	}

	for _, historyJSON := range histories {
		var chatMessage chat.Chat

		if err := json.Unmarshal([]byte(historyJSON), &chatMessage); err != nil {
			log.Print(ctx, log.KV{"chat.history", "ERROR: redis Unmarshal failed"}, log.KV{"error", err.Error()})
			return nil, chat.Internal("Internal server error")
		}

		res = append(res, &chatMessage)
	}

	return res, nil
}

func (s *chatsrvc) StreamRoom(ctx context.Context, p *chat.StreamRoomPayload, stream chat.StreamRoomServerStream) (err error) {
	log.Info(ctx, log.KV{"chat.stream_room", "starting stream room"})

	// Get username from context
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return chat.Unauthorized("user not authenticated")
	}

	pubsubChannel := roomPublishKey + ":" + p.RoomID
	pubsub := s.redis.Subscribe(ctx, pubsubChannel)
	defer pubsub.Close()

	msgCh := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(msgCh)
		defer close(errCh)

		for {
			str, err := stream.Recv()
			if err == io.EOF {
				log.Info(ctx, log.KV{"chat.stream_room", "client closed connection"})
				return
			}
			if err != nil {
				log.Print(ctx, log.KV{"chat.stream_room", "ERROR: recv error"}, log.KV{"error", err.Error()})
				errCh <- err
				return
			}

			msgCh <- str
		}
	}()

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			var newChat chat.Chat
			if err := json.Unmarshal([]byte(msg.Payload), &newChat); err != nil {
				log.Print(ctx, log.KV{"chat.stream_room", "ERROR: json.Unmarshal from pubsub failed"}, log.KV{"error", err.Error()})
				continue
			}

			if err := stream.Send(&newChat); err != nil {
				log.Print(ctx, log.KV{"chat.stream_room", "ERROR: stream.Send failed"}, log.KV{"error", err.Error()})
				errCh <- err
				return
			}
		}
	}()

	for done := false; !done; {
		select {
		case msg := <-msgCh:
			if msg == "" {
				continue
			}

			newChat := chat.Chat{
				ID:        uuid.NewString(),
				UserID:    userID,
				Message:   msg,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}

			chatJSON, err := json.Marshal(newChat)
			if err != nil {
				log.Print(ctx, log.KV{"chat.stream_room", "ERROR: json.Marshal failed"}, log.KV{"error", err.Error()})
				return chat.Internal("Internal server error")
			}

			if err := s.redis.LPush(ctx, fmt.Sprintf("%s:%s", historyKey, p.RoomID), chatJSON).Err(); err != nil {
				log.Print(ctx, log.KV{"chat.stream_room", "ERROR: redis LPush failed"}, log.KV{"error", err.Error()})
				return chat.Internal("Internal server error")
			}

			if err := s.redis.Publish(ctx, pubsubChannel, chatJSON).Err(); err != nil {
				log.Print(ctx, log.KV{"chat.stream_room", "ERROR: redis Publish failed"}, log.KV{"error", err.Error()})
				return chat.Internal("Internal server error")
			}

		case err := <-errCh:
			if err != nil {
				log.Print(ctx, log.KV{"chat.stream_room", "ERROR: error in goroutine"}, log.KV{"error", err.Error()})
				return chat.Internal("Internal server error")
			}
		case <-ctx.Done():
			log.Info(ctx, log.KV{"chat.stream_room", "context done"})
			done = true
		}
	}

	return stream.Close()
}

func (s *chatsrvc) RoomList(ctx context.Context, p *chat.RoomListPayload) (res []string, err error) {
	log.Printf(ctx, "chat.room-list")
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, chat.Unauthorized("user not authenticated")
	}
	roomListKey := "rooms:" + userID
	res, err = s.redis.LRange(ctx, roomListKey, 0, -1).Result()
	if err != nil {
		return nil, chat.Internal("Internal server error")
	}

	return
}

func (s *chatsrvc) JoinRoom(ctx context.Context, p *chat.JoinRoomPayload) (res string, err error) {
	log.Printf(ctx, "chat.join-room")
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", chat.Unauthorized("user not authenticated")
	}
	inviteKey := "invite:" + p.InviteKey + ":" + userID
	roomListKey := "rooms:" + userID
	inviteRoom := s.redis.Get(ctx, inviteKey).Val()
	if inviteRoom == "" {
		return "", chat.Notfound("user not invited")
	}

	if s.redis.LPush(ctx, roomListKey, inviteRoom).Err() != nil {
		return "", chat.Internal("Internal server error")
	}

	return inviteRoom, nil
}

func (s *chatsrvc) InviteRoom(ctx context.Context, p *chat.InviteRoomPayload) (res string, err error) {
	log.Printf(ctx, "chat.invite-room")
	_, ok := ctx.Value("user_id").(string)
	if !ok {
		return "", chat.Unauthorized("user not authenticated")
	}

	newInviteId := uuid.NewString()
	inviteKey := "invite:" + newInviteId + ":" + p.UserID
	if s.redis.Set(ctx, inviteKey, p.RoomID, 0).Err() != nil {
		return "", chat.Internal("Internal server error")
	}

	return newInviteId, nil
}
