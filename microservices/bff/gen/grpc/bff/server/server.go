// Code generated by goa v3.21.1, DO NOT EDIT.
//
// bff gRPC server
//
// Command:
// $ goa gen object-t.com/hackz-giganoto/microservices/bff/design

package server

import (
	"context"
	"errors"

	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc/codes"
	bff "object-t.com/hackz-giganoto/microservices/bff/gen/bff"
	bffpb "object-t.com/hackz-giganoto/microservices/bff/gen/grpc/bff/pb"
)

// Server implements the bffpb.BffServer interface.
type Server struct {
	CreateRoomH    goagrpc.UnaryHandler
	HistoryH       goagrpc.UnaryHandler
	RoomListH      goagrpc.UnaryHandler
	JoinRoomH      goagrpc.UnaryHandler
	InviteRoomH    goagrpc.UnaryHandler
	StreamChatH    goagrpc.StreamHandler
	GetProfileH    goagrpc.UnaryHandler
	UpdateProfileH goagrpc.UnaryHandler
	bffpb.UnimplementedBffServer
}

// StreamChatServerStream implements the bff.StreamChatServerStream interface.
type StreamChatServerStream struct {
	stream bffpb.Bff_StreamChatServer
}

// New instantiates the server struct with the bff service endpoints.
func New(e *bff.Endpoints, uh goagrpc.UnaryHandler, sh goagrpc.StreamHandler) *Server {
	return &Server{
		CreateRoomH:    NewCreateRoomHandler(e.CreateRoom, uh),
		HistoryH:       NewHistoryHandler(e.History, uh),
		RoomListH:      NewRoomListHandler(e.RoomList, uh),
		JoinRoomH:      NewJoinRoomHandler(e.JoinRoom, uh),
		InviteRoomH:    NewInviteRoomHandler(e.InviteRoom, uh),
		StreamChatH:    NewStreamChatHandler(e.StreamChat, sh),
		GetProfileH:    NewGetProfileHandler(e.GetProfile, uh),
		UpdateProfileH: NewUpdateProfileHandler(e.UpdateProfile, uh),
	}
}

// NewCreateRoomHandler creates a gRPC handler which serves the "bff" service
// "create_room" endpoint.
func NewCreateRoomHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeCreateRoomRequest, EncodeCreateRoomResponse)
	}
	return h
}

// CreateRoom implements the "CreateRoom" method in bffpb.BffServer interface.
func (s *Server) CreateRoom(ctx context.Context, message *bffpb.CreateRoomRequest) (*bffpb.CreateRoomResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "create_room")
	ctx = context.WithValue(ctx, goa.ServiceKey, "bff")
	resp, err := s.CreateRoomH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "unauthorized":
				return nil, goagrpc.NewStatusError(codes.Unauthenticated, err, goagrpc.NewErrorResponse(err))
			case "internal_error":
				return nil, goagrpc.NewStatusError(codes.Internal, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*bffpb.CreateRoomResponse), nil
}

// NewHistoryHandler creates a gRPC handler which serves the "bff" service
// "history" endpoint.
func NewHistoryHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeHistoryRequest, EncodeHistoryResponse)
	}
	return h
}

// History implements the "History" method in bffpb.BffServer interface.
func (s *Server) History(ctx context.Context, message *bffpb.HistoryRequest) (*bffpb.HistoryResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "history")
	ctx = context.WithValue(ctx, goa.ServiceKey, "bff")
	resp, err := s.HistoryH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "unauthorized":
				return nil, goagrpc.NewStatusError(codes.Unauthenticated, err, goagrpc.NewErrorResponse(err))
			case "permission-denied":
				return nil, goagrpc.NewStatusError(codes.PermissionDenied, err, goagrpc.NewErrorResponse(err))
			case "internal_error":
				return nil, goagrpc.NewStatusError(codes.Internal, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*bffpb.HistoryResponse), nil
}

// NewRoomListHandler creates a gRPC handler which serves the "bff" service
// "room-list" endpoint.
func NewRoomListHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeRoomListRequest, EncodeRoomListResponse)
	}
	return h
}

// RoomList implements the "RoomList" method in bffpb.BffServer interface.
func (s *Server) RoomList(ctx context.Context, message *bffpb.RoomListRequest) (*bffpb.RoomListResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "room-list")
	ctx = context.WithValue(ctx, goa.ServiceKey, "bff")
	resp, err := s.RoomListH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "unauthorized":
				return nil, goagrpc.NewStatusError(codes.InvalidArgument, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*bffpb.RoomListResponse), nil
}

// NewJoinRoomHandler creates a gRPC handler which serves the "bff" service
// "join-room" endpoint.
func NewJoinRoomHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeJoinRoomRequest, EncodeJoinRoomResponse)
	}
	return h
}

// JoinRoom implements the "JoinRoom" method in bffpb.BffServer interface.
func (s *Server) JoinRoom(ctx context.Context, message *bffpb.JoinRoomRequest) (*bffpb.JoinRoomResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "join-room")
	ctx = context.WithValue(ctx, goa.ServiceKey, "bff")
	resp, err := s.JoinRoomH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "notfound":
				return nil, goagrpc.NewStatusError(codes.NotFound, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*bffpb.JoinRoomResponse), nil
}

// NewInviteRoomHandler creates a gRPC handler which serves the "bff" service
// "invite-room" endpoint.
func NewInviteRoomHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeInviteRoomRequest, EncodeInviteRoomResponse)
	}
	return h
}

// InviteRoom implements the "InviteRoom" method in bffpb.BffServer interface.
func (s *Server) InviteRoom(ctx context.Context, message *bffpb.InviteRoomRequest) (*bffpb.InviteRoomResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "invite-room")
	ctx = context.WithValue(ctx, goa.ServiceKey, "bff")
	resp, err := s.InviteRoomH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "invalid_argument":
				return nil, goagrpc.NewStatusError(codes.InvalidArgument, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*bffpb.InviteRoomResponse), nil
}

// NewStreamChatHandler creates a gRPC handler which serves the "bff" service
// "stream_chat" endpoint.
func NewStreamChatHandler(endpoint goa.Endpoint, h goagrpc.StreamHandler) goagrpc.StreamHandler {
	if h == nil {
		h = goagrpc.NewStreamHandler(endpoint, DecodeStreamChatRequest)
	}
	return h
}

// StreamChat implements the "StreamChat" method in bffpb.BffServer interface.
func (s *Server) StreamChat(stream bffpb.Bff_StreamChatServer) error {
	ctx := stream.Context()
	ctx = context.WithValue(ctx, goa.MethodKey, "stream_chat")
	ctx = context.WithValue(ctx, goa.ServiceKey, "bff")
	p, err := s.StreamChatH.Decode(ctx, nil)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "unauthorized":
				return goagrpc.NewStatusError(codes.Unauthenticated, err, goagrpc.NewErrorResponse(err))
			case "permission-denied":
				return goagrpc.NewStatusError(codes.PermissionDenied, err, goagrpc.NewErrorResponse(err))
			case "internal_error":
				return goagrpc.NewStatusError(codes.Internal, err, goagrpc.NewErrorResponse(err))
			}
		}
		return goagrpc.EncodeError(err)
	}
	ep := &bff.StreamChatEndpointInput{
		Stream:  &StreamChatServerStream{stream: stream},
		Payload: p.(*bff.StreamChatPayload),
	}
	err = s.StreamChatH.Handle(ctx, ep)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "unauthorized":
				return goagrpc.NewStatusError(codes.Unauthenticated, err, goagrpc.NewErrorResponse(err))
			case "permission-denied":
				return goagrpc.NewStatusError(codes.PermissionDenied, err, goagrpc.NewErrorResponse(err))
			case "internal_error":
				return goagrpc.NewStatusError(codes.Internal, err, goagrpc.NewErrorResponse(err))
			}
		}
		return goagrpc.EncodeError(err)
	}
	return nil
}

// NewGetProfileHandler creates a gRPC handler which serves the "bff" service
// "get_profile" endpoint.
func NewGetProfileHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeGetProfileRequest, EncodeGetProfileResponse)
	}
	return h
}

// GetProfile implements the "GetProfile" method in bffpb.BffServer interface.
func (s *Server) GetProfile(ctx context.Context, message *bffpb.GetProfileRequest) (*bffpb.GetProfileResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "get_profile")
	ctx = context.WithValue(ctx, goa.ServiceKey, "bff")
	resp, err := s.GetProfileH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "unauthorized":
				return nil, goagrpc.NewStatusError(codes.Unauthenticated, err, goagrpc.NewErrorResponse(err))
			case "not_found":
				return nil, goagrpc.NewStatusError(codes.NotFound, err, goagrpc.NewErrorResponse(err))
			case "internal_error":
				return nil, goagrpc.NewStatusError(codes.Internal, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*bffpb.GetProfileResponse), nil
}

// NewUpdateProfileHandler creates a gRPC handler which serves the "bff"
// service "update_profile" endpoint.
func NewUpdateProfileHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeUpdateProfileRequest, EncodeUpdateProfileResponse)
	}
	return h
}

// UpdateProfile implements the "UpdateProfile" method in bffpb.BffServer
// interface.
func (s *Server) UpdateProfile(ctx context.Context, message *bffpb.UpdateProfileRequest) (*bffpb.UpdateProfileResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "update_profile")
	ctx = context.WithValue(ctx, goa.ServiceKey, "bff")
	resp, err := s.UpdateProfileH.Handle(ctx, message)
	if err != nil {
		var en goa.GoaErrorNamer
		if errors.As(err, &en) {
			switch en.GoaErrorName() {
			case "unauthorized":
				return nil, goagrpc.NewStatusError(codes.Unauthenticated, err, goagrpc.NewErrorResponse(err))
			case "bad_request":
				return nil, goagrpc.NewStatusError(codes.InvalidArgument, err, goagrpc.NewErrorResponse(err))
			case "internal_error":
				return nil, goagrpc.NewStatusError(codes.Internal, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*bffpb.UpdateProfileResponse), nil
}

// Send streams instances of "bffpb.StreamChatResponse" to the "stream_chat"
// endpoint gRPC stream.
func (s *StreamChatServerStream) Send(res *bff.EnrichedMessage) error {
	v := NewProtoEnrichedMessageStreamChatResponse(res)
	return s.stream.Send(v)
}

// SendWithContext streams instances of "bffpb.StreamChatResponse" to the
// "stream_chat" endpoint gRPC stream with context.
func (s *StreamChatServerStream) SendWithContext(ctx context.Context, res *bff.EnrichedMessage) error {
	return s.Send(res)
}

// Recv reads instances of "bffpb.StreamChatStreamingRequest" from the
// "stream_chat" endpoint gRPC stream.
func (s *StreamChatServerStream) Recv() (string, error) {
	var res string
	v, err := s.stream.Recv()
	if err != nil {
		return res, err
	}
	return NewStreamChatStreamingRequestStreamChatStreamingRequest(v), nil
}

// RecvWithContext reads instances of "bffpb.StreamChatStreamingRequest" from
// the "stream_chat" endpoint gRPC stream with context.
func (s *StreamChatServerStream) RecvWithContext(ctx context.Context) (string, error) {
	return s.Recv()
}

func (s *StreamChatServerStream) Close() error {
	// nothing to do here
	return nil
}
