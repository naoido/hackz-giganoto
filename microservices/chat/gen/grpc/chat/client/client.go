// Code generated by goa v3.21.1, DO NOT EDIT.
//
// chat gRPC client
//
// Command:
// $ goa gen object-t.com/hackz-giganoto/microservices/chat/design

package client

import (
	"context"

	goagrpc "goa.design/goa/v3/grpc"
	goapb "goa.design/goa/v3/grpc/pb"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
	chat "object-t.com/hackz-giganoto/microservices/chat/gen/chat"
	chatpb "object-t.com/hackz-giganoto/microservices/chat/gen/grpc/chat/pb"
)

// Client lists the service endpoint gRPC clients.
type Client struct {
	grpccli chatpb.ChatClient
	opts    []grpc.CallOption
}

// StreamRoomClientStream implements the chat.StreamRoomClientStream interface.
type StreamRoomClientStream struct {
	stream chatpb.Chat_StreamRoomClient
}

// NewClient instantiates gRPC client for all the chat service servers.
func NewClient(cc *grpc.ClientConn, opts ...grpc.CallOption) *Client {
	return &Client{
		grpccli: chatpb.NewChatClient(cc),
		opts:    opts,
	}
}

// CreateRoom calls the "CreateRoom" function in chatpb.ChatClient interface.
func (c *Client) CreateRoom() goa.Endpoint {
	return func(ctx context.Context, v any) (any, error) {
		inv := goagrpc.NewInvoker(
			BuildCreateRoomFunc(c.grpccli, c.opts...),
			EncodeCreateRoomRequest,
			DecodeCreateRoomResponse)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			resp := goagrpc.DecodeError(err)
			switch message := resp.(type) {
			case *goapb.ErrorResponse:
				return nil, goagrpc.NewServiceError(message)
			default:
				return nil, goa.Fault("%s", err.Error())
			}
		}
		return res, nil
	}
}

// History calls the "History" function in chatpb.ChatClient interface.
func (c *Client) History() goa.Endpoint {
	return func(ctx context.Context, v any) (any, error) {
		inv := goagrpc.NewInvoker(
			BuildHistoryFunc(c.grpccli, c.opts...),
			EncodeHistoryRequest,
			DecodeHistoryResponse)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			resp := goagrpc.DecodeError(err)
			switch message := resp.(type) {
			case *goapb.ErrorResponse:
				return nil, goagrpc.NewServiceError(message)
			default:
				return nil, goa.Fault("%s", err.Error())
			}
		}
		return res, nil
	}
}

// RoomList calls the "RoomList" function in chatpb.ChatClient interface.
func (c *Client) RoomList() goa.Endpoint {
	return func(ctx context.Context, v any) (any, error) {
		inv := goagrpc.NewInvoker(
			BuildRoomListFunc(c.grpccli, c.opts...),
			EncodeRoomListRequest,
			DecodeRoomListResponse)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			resp := goagrpc.DecodeError(err)
			switch message := resp.(type) {
			case *goapb.ErrorResponse:
				return nil, goagrpc.NewServiceError(message)
			default:
				return nil, goa.Fault("%s", err.Error())
			}
		}
		return res, nil
	}
}

// JoinRoom calls the "JoinRoom" function in chatpb.ChatClient interface.
func (c *Client) JoinRoom() goa.Endpoint {
	return func(ctx context.Context, v any) (any, error) {
		inv := goagrpc.NewInvoker(
			BuildJoinRoomFunc(c.grpccli, c.opts...),
			EncodeJoinRoomRequest,
			DecodeJoinRoomResponse)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			resp := goagrpc.DecodeError(err)
			switch message := resp.(type) {
			case *goapb.ErrorResponse:
				return nil, goagrpc.NewServiceError(message)
			default:
				return nil, goa.Fault("%s", err.Error())
			}
		}
		return res, nil
	}
}

// InviteRoom calls the "InviteRoom" function in chatpb.ChatClient interface.
func (c *Client) InviteRoom() goa.Endpoint {
	return func(ctx context.Context, v any) (any, error) {
		inv := goagrpc.NewInvoker(
			BuildInviteRoomFunc(c.grpccli, c.opts...),
			EncodeInviteRoomRequest,
			DecodeInviteRoomResponse)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			resp := goagrpc.DecodeError(err)
			switch message := resp.(type) {
			case *goapb.ErrorResponse:
				return nil, goagrpc.NewServiceError(message)
			default:
				return nil, goa.Fault("%s", err.Error())
			}
		}
		return res, nil
	}
}

// StreamRoom calls the "StreamRoom" function in chatpb.ChatClient interface.
func (c *Client) StreamRoom() goa.Endpoint {
	return func(ctx context.Context, v any) (any, error) {
		inv := goagrpc.NewInvoker(
			BuildStreamRoomFunc(c.grpccli, c.opts...),
			EncodeStreamRoomRequest,
			DecodeStreamRoomResponse)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			resp := goagrpc.DecodeError(err)
			switch message := resp.(type) {
			case *goapb.ErrorResponse:
				return nil, goagrpc.NewServiceError(message)
			default:
				return nil, goa.Fault("%s", err.Error())
			}
		}
		return res, nil
	}
}

// Recv reads instances of "chatpb.StreamRoomResponse" from the "stream-room"
// endpoint gRPC stream.
func (s *StreamRoomClientStream) Recv() (*chat.Chat, error) {
	var res *chat.Chat
	v, err := s.stream.Recv()
	if err != nil {
		return res, err
	}
	return NewStreamRoomResponseChat2(v), nil
}

// RecvWithContext reads instances of "chatpb.StreamRoomResponse" from the
// "stream-room" endpoint gRPC stream with context.
func (s *StreamRoomClientStream) RecvWithContext(ctx context.Context) (*chat.Chat, error) {
	return s.Recv()
}

// Send streams instances of "chatpb.StreamRoomStreamingRequest" to the
// "stream-room" endpoint gRPC stream.
func (s *StreamRoomClientStream) Send(res string) error {
	v := NewProtoStreamRoomStreamingRequest(res)
	return s.stream.Send(v)
}

// SendWithContext streams instances of "chatpb.StreamRoomStreamingRequest" to
// the "stream-room" endpoint gRPC stream with context.
func (s *StreamRoomClientStream) SendWithContext(ctx context.Context, res string) error {
	return s.Send(res)
}

func (s *StreamRoomClientStream) Close() error {
	// Close the send direction of the stream
	return s.stream.CloseSend()
}
