package bffapi

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	chatpb "object-t.com/hackz-giganoto/microservices/chat/gen/grpc/chat/pb"
	"object-t.com/hackz-giganoto/microservices/profile/gen/profile"
	security2 "object-t.com/hackz-giganoto/pkg/security"
	"object-t.com/hackz-giganoto/pkg/telemetry"
	"os"

	"goa.design/clue/log"
	"goa.design/goa/v3/security"
	bff "object-t.com/hackz-giganoto/microservices/bff/gen/bff"
)

type ServiceConfig struct {
	ChatServiceAddr    string
	ProfileServiceAddr string
}

type bffsrvc struct {
	chatGrpcClient chatpb.ChatClient
}

// NewBff returns the bff service implementation.
func NewBff() bff.Service {
	config := ServiceConfig{
		ChatServiceAddr: getEnvOrDefault("CHAT_SERVICE_ADDR", "localhost:50053"),
	}

	chatConn, err := grpc.NewClient(
		config.ChatServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		telemetry.GRPCClientInterceptor(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to chat service: %v", err))
	}
	chatGRPCClient := chatpb.NewChatClient(chatConn)

	return &bffsrvc{
		chatGrpcClient: chatGRPCClient,
	}
}

func getEnvOrDefault(envVar, defaultValue string) string {
	if value := os.Getenv(envVar); value != "" {
		return value
	}
	return defaultValue
}

func (s *bffsrvc) addJWTToContext(ctx context.Context) context.Context {
	if token, ok := ctx.Value("jwt_token").(string); ok {
		md := metadata.Pairs("authorization", "Bearer "+token)
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

// JWTAuth implements the authorization logic for service "bff" for the "jwt"
// security scheme.
func (s *bffsrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	claims, err := security2.ValidToken(token)
	if err != nil {
		log.Printf(ctx, "invalid token: %v", err)
		return ctx, profile.InvalidToken("invalid token")
	}
	return security2.HasPermission(ctx, claims, scheme)
}

// Create a new chat room
func (s *bffsrvc) CreateRoom(ctx context.Context, p *bff.CreateRoomPayload) (res string, err error) {
	log.Printf(ctx, "bff.create_room")
	return s.CreateRoom(ctx, &bff.CreateRoomPayload{})
}

// Get chat room history with enriched user names
func (s *bffsrvc) History(ctx context.Context, p *bff.HistoryPayload) (res []*bff.EnrichedMessage, err error) {
	log.Printf(ctx, "bff.history")
	history, err := s.chatGrpcClient.History(ctx, &chatpb.HistoryRequest{})
	if err != nil {
		return nil, err
	}

	res = []*bff.EnrichedMessage{}
	for _, h := range history.Field {
		res = append(res, &bff.EnrichedMessage{
			MessageID: &h.Id,
			Message:   h.Message_,
			UserID:    h.UserId,
			CreatedAt: &h.CreatedAt,
			UpdatedAt: &h.UpdatedAt,
		})
	}

	return
}

// Get all chat rooms history
func (s *bffsrvc) RoomList(ctx context.Context, p *bff.RoomListPayload) (res []string, err error) {
	log.Printf(ctx, "bff.room-list")
	list, err := s.chatGrpcClient.RoomList(ctx, &chatpb.RoomListRequest{})
	if err != nil {
		return nil, err
	}

	return list.Field, nil
}

// Creates a new chat room
func (s *bffsrvc) JoinRoom(ctx context.Context, p *bff.JoinRoomPayload) (res string, err error) {
	log.Printf(ctx, "bff.join-room")
	return s.JoinRoom(ctx, p)
}

// Creates a new chat room
func (s *bffsrvc) InviteRoom(ctx context.Context, p *bff.InviteRoomPayload) (res string, err error) {
	log.Printf(ctx, "bff.invite-room")
	invRes, err := s.chatGrpcClient.InviteRoom(ctx, &chatpb.InviteRoomRequest{
		RoomId: p.RoomID,
		UserId: p.UserID,
	})

	return invRes.String(), err
}

func (s *bffsrvc) StreamChat(ctx context.Context, p *bff.StreamChatPayload, stream bff.StreamChatServerStream) (err error) {

	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		return bff.Unauthorized("user not authenticated")
	}

	grpcCtx := s.addJWTToContext(ctx)

	md := metadata.Pairs(
		"authorization", p.Token,
		"room_id", p.RoomID,
	)
	grpcCtx = metadata.NewOutgoingContext(grpcCtx, md)

	chatStream, err := s.chatGrpcClient.StreamRoom(grpcCtx)
	if err != nil {
		log.Print(ctx, log.KV{"bff.stream_chat", "ERROR: failed to connect to chat service"}, log.KV{"error", err.Error()})
		return bff.InternalError("failed to connect to chat service")
	}
	defer func() {
		if err := chatStream.CloseSend(); err != nil {
			log.Print(ctx, log.KV{"bff.stream_chat", "ERROR: error closing chat stream"}, log.KV{"error", err.Error()})
		}
	}()

	errCh := make(chan error, 3)

	go func() {
		defer func() {
			if err := chatStream.CloseSend(); err != nil {
				log.Print(ctx, log.KV{"bff.stream_chat", "ERROR: error closing chat stream send"}, log.KV{"error", err.Error()})
			}
		}()

		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Print(ctx, log.KV{"bff.stream_chat", "ERROR: client recv error"}, log.KV{"error", err.Error()})
				errCh <- err
				return
			}

			chatReq := &chatpb.StreamRoomStreamingRequest{
				Field: msg,
			}

			if err := chatStream.Send(chatReq); err != nil {
				log.Print(ctx, log.KV{"bff.stream_chat", "ERROR: chat send error"}, log.KV{"error", err.Error()})
				errCh <- err
				return
			}

			log.Print(ctx, log.KV{"bff.stream_chat", fmt.Sprintf("DEBUG: forwarded message from user %s to chat service", userID)})
		}
	}()

	go func() {
		for {
			chatMsg, err := chatStream.Recv()
			if err != nil {
				log.Print(ctx, log.KV{"bff.stream_chat", "ERROR: chat recv error"}, log.KV{"error", err.Error()})
				errCh <- err
				return
			}

			userId := chatMsg.UserId

			enrichedMsg := &bff.EnrichedMessage{
				RoomID:  p.RoomID,
				UserID:  userId,
				Message: chatMsg.Message_,
			}

			if err := stream.Send(enrichedMsg); err != nil {
				log.Print(ctx, log.KV{"bff.stream_chat", "ERROR: client send error"}, log.KV{"error", err.Error()})
				errCh <- err
				return
			}

			log.Print(ctx, log.KV{"bff.stream_chat", fmt.Sprintf("DEBUG: sent enriched message to client")})
		}
	}()

	// Wait for context cancellation or error
	select {
	case err := <-errCh:
		if err != nil {
			log.Print(ctx, log.KV{"bff.stream_chat", "ERROR: stream error"}, log.KV{"error", err.Error()})
			return bff.InternalError("stream error")
		}
	case <-ctx.Done():
		log.Info(ctx, log.KV{"bff.stream_chat", "context cancelled"})
	}

	return stream.Close()
}

// Get current user profile
func (s *bffsrvc) GetProfile(ctx context.Context, p *bff.GetProfilePayload) (res *bff.GetProfileResult, err error) {
	res = &bff.GetProfileResult{}
	log.Printf(ctx, "bff.get_profile")
	return
}

// Update current user profile
func (s *bffsrvc) UpdateProfile(ctx context.Context, p *bff.UpdateProfilePayload) (res *bff.UpdateProfileResult, err error) {
	res = &bff.UpdateProfileResult{}
	log.Printf(ctx, "bff.update_profile")
	return
}
