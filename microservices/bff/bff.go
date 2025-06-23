package bffapi

import (
	"context"
	"fmt"
	"os"

	"goa.design/clue/log"
	"goa.design/goa/v3/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	bff "object-t.com/hackz-giganoto/microservices/bff/gen/bff"
	chatpb "object-t.com/hackz-giganoto/microservices/chat/gen/grpc/chat/pb"
	profilepb "object-t.com/hackz-giganoto/microservices/profile/gen/grpc/profile/pb"
	security2 "object-t.com/hackz-giganoto/pkg/security"
	"object-t.com/hackz-giganoto/pkg/telemetry"
)

type ServiceConfig struct {
	ChatServiceAddr    string
	ProfileServiceAddr string
}

type bffsrvc struct {
	chatGRPCClient    chatpb.ChatClient
	profileGRPCClient profilepb.ProfileClient
}

// NewBff returns the bff service implementation.
func NewBff() bff.Service {
	config := ServiceConfig{
		ChatServiceAddr:    getEnvOrDefault("CHAT_SERVICE_ADDR", "localhost:50053"),
		ProfileServiceAddr: getEnvOrDefault("PROFILE_SERVICE_ADDR", "localhost:50052"),
	}

	// Chat service client setup
	chatConn, err := grpc.NewClient(
		config.ChatServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		telemetry.GRPCClientInterceptor(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to chat service: %v", err))
	}
	chatGRPCClient := chatpb.NewChatClient(chatConn)

	profileConn, err := grpc.NewClient(
		config.ProfileServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		telemetry.GRPCClientInterceptor(),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to profile service: %v", err))
	}
	profileGRPCClient := profilepb.NewProfileClient(profileConn)

	return &bffsrvc{
		chatGRPCClient:    chatGRPCClient,
		profileGRPCClient: profileGRPCClient,
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
		return ctx, bff.InvalidArgument("invalid token")
	}
	return security2.HasPermission(ctx, claims, scheme)
}

// CreateRoom creates a new chat room
func (s *bffsrvc) CreateRoom(ctx context.Context, p *bff.CreateRoomPayload) (res string, err error) {
	log.Printf(ctx, "bff.create_room")
	grpcCtx := s.addJWTToContext(ctx)

	room, err := s.chatGRPCClient.CreateRoom(grpcCtx, &chatpb.CreateRoomRequest{})
	if err != nil {
		log.Printf(ctx, "failed to create room: %v", err)
		return "", err
	}
	if room == nil {
		return "", fmt.Errorf("received nil response from chat service")
	}

	return room.Field, nil
}

// History gets chat room history with enriched user names
func (s *bffsrvc) History(ctx context.Context, p *bff.HistoryPayload) (res []*bff.EnrichedMessage, err error) {
	log.Printf(ctx, "bff.history")
	grpcCtx := s.addJWTToContext(ctx)
	resp, err := s.chatGRPCClient.History(grpcCtx, &chatpb.HistoryRequest{
		RoomId: p.RoomID,
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response from chat service")
	}

	res = []*bff.EnrichedMessage{}
	for _, h := range resp.Field {
		res = append(res, &bff.EnrichedMessage{
			MessageID: &h.Id,
			RoomID:    h.RoomId,
			Message:   h.Message_,
			UserID:    h.UserId,
			UpdatedAt: &h.UpdatedAt,
			CreatedAt: &h.CreatedAt,
		})
	}

	return
}

// RoomList gets all chat rooms
func (s *bffsrvc) RoomList(ctx context.Context, p *bff.RoomListPayload) (res []string, err error) {
	log.Printf(ctx, "bff.room-list")
	grpcCtx := s.addJWTToContext(ctx)
	resp, err := s.chatGRPCClient.RoomList(grpcCtx, &chatpb.RoomListRequest{})
	if err != nil {
		return nil, bff.InternalError("InternalError")
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response from chat service")
	}

	res = resp.Field

	return
}

// JoinRoom joins a chat room
func (s *bffsrvc) JoinRoom(ctx context.Context, p *bff.JoinRoomPayload) (res string, err error) {
	log.Printf(ctx, "bff.join-room")
	grpcCtx := s.addJWTToContext(ctx)

	resp, err := s.chatGRPCClient.JoinRoom(grpcCtx, &chatpb.JoinRoomRequest{
		InviteKey: p.InviteKey,
	})
	if err != nil {
		return "", bff.InternalError("InternalError")
	}

	if resp == nil {
		return "", fmt.Errorf("received nil response from chat service")
	}

	return resp.Field, nil
}

// InviteRoom invites a user to a chat room
func (s *bffsrvc) InviteRoom(ctx context.Context, p *bff.InviteRoomPayload) (res string, err error) {
	log.Printf(ctx, "bff.invite-room")
	grpcCtx := s.addJWTToContext(ctx)
	resp, err := s.chatGRPCClient.InviteRoom(grpcCtx, &chatpb.InviteRoomRequest{
		RoomId: p.RoomID,
		UserId: p.UserID,
	})

	if err != nil {
		return "", bff.InternalError("InternalError")
	}

	if resp == nil {
		return "", fmt.Errorf("received nil response from chat service")
	}

	return resp.Field, nil
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

	chatStream, err := s.chatGRPCClient.StreamRoom(grpcCtx)
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

			c := &chatpb.StreamRoomStreamingRequest{
				Field: msg,
			}

			if err := chatStream.Send(c); err != nil {
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

			enrichedMsg := &bff.EnrichedMessage{
				RoomID:  p.RoomID,
				UserID:  chatMsg.UserId,
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

// GetProfile gets current user profile
func (s *bffsrvc) GetProfile(ctx context.Context, p *bff.GetProfilePayload) (res *bff.GetProfileResult, err error) {
	log.Printf(ctx, "bff.get_profile")
	grpcCtx := s.addJWTToContext(ctx)
	resp, err := s.profileGRPCClient.GetProfile(grpcCtx, &profilepb.GetProfileRequest{
		UserId: p.UserID,
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response from profile service")
	}

	return &bff.GetProfileResult{
		UserID: resp.UserId,
		Name:   resp.Name,
	}, nil
}

// UpdateProfile updates current user profile
func (s *bffsrvc) UpdateProfile(ctx context.Context, p *bff.UpdateProfilePayload) (res *bff.UpdateProfileResult, err error) {
	log.Printf(ctx, "bff.update_profile")
	grpcCtx := s.addJWTToContext(ctx)
	resp, err := s.profileGRPCClient.UpdateProfile(grpcCtx, &profilepb.UpdateProfileRequest{
		Name: p.Name,
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response from profile service")
	}

	return &bff.UpdateProfileResult{
		UserID: resp.UserId,
		Name:   resp.Name,
	}, nil
}
