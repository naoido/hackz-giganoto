// Code generated by goa v3.21.1, DO NOT EDIT.
//
// bff gRPC client
//
// Command:
// $ goa gen object-t.com/hackz-giganoto/microservices/bff/design

package client

import (
	"context"

	goagrpc "goa.design/goa/v3/grpc"
	goapb "goa.design/goa/v3/grpc/pb"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
	bff "object-t.com/hackz-giganoto/microservices/bff/gen/bff"
	bffpb "object-t.com/hackz-giganoto/microservices/bff/gen/grpc/bff/pb"
)

// Client lists the service endpoint gRPC clients.
type Client struct {
	grpccli bffpb.BffClient
	opts    []grpc.CallOption
}

// StreamChatClientStream implements the bff.StreamChatClientStream interface.
type StreamChatClientStream struct {
	stream bffpb.Bff_StreamChatClient
}

// NewClient instantiates gRPC client for all the bff service servers.
func NewClient(cc *grpc.ClientConn, opts ...grpc.CallOption) *Client {
	return &Client{
		grpccli: bffpb.NewBffClient(cc),
		opts:    opts,
	}
}

// CreateRoom calls the "CreateRoom" function in bffpb.BffClient interface.
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

// History calls the "History" function in bffpb.BffClient interface.
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

// RoomList calls the "RoomList" function in bffpb.BffClient interface.
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

// JoinRoom calls the "JoinRoom" function in bffpb.BffClient interface.
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

// InviteRoom calls the "InviteRoom" function in bffpb.BffClient interface.
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

// StreamChat calls the "StreamChat" function in bffpb.BffClient interface.
func (c *Client) StreamChat() goa.Endpoint {
	return func(ctx context.Context, v any) (any, error) {
		inv := goagrpc.NewInvoker(
			BuildStreamChatFunc(c.grpccli, c.opts...),
			EncodeStreamChatRequest,
			DecodeStreamChatResponse)
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

// GetProfile calls the "GetProfile" function in bffpb.BffClient interface.
func (c *Client) GetProfile() goa.Endpoint {
	return func(ctx context.Context, v any) (any, error) {
		inv := goagrpc.NewInvoker(
			BuildGetProfileFunc(c.grpccli, c.opts...),
			EncodeGetProfileRequest,
			DecodeGetProfileResponse)
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

// UpdateProfile calls the "UpdateProfile" function in bffpb.BffClient
// interface.
func (c *Client) UpdateProfile() goa.Endpoint {
	return func(ctx context.Context, v any) (any, error) {
		inv := goagrpc.NewInvoker(
			BuildUpdateProfileFunc(c.grpccli, c.opts...),
			EncodeUpdateProfileRequest,
			DecodeUpdateProfileResponse)
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

// Recv reads instances of "bffpb.StreamChatResponse" from the "stream_chat"
// endpoint gRPC stream.
func (s *StreamChatClientStream) Recv() (*bff.EnrichedMessage, error) {
	var res *bff.EnrichedMessage
	v, err := s.stream.Recv()
	if err != nil {
		return res, err
	}
	return NewStreamChatResponseEnrichedMessage(v), nil
}

// RecvWithContext reads instances of "bffpb.StreamChatResponse" from the
// "stream_chat" endpoint gRPC stream with context.
func (s *StreamChatClientStream) RecvWithContext(ctx context.Context) (*bff.EnrichedMessage, error) {
	return s.Recv()
}

// Send streams instances of "bffpb.StreamChatStreamingRequest" to the
// "stream_chat" endpoint gRPC stream.
func (s *StreamChatClientStream) Send(res string) error {
	v := NewProtoStreamChatStreamingRequest(res)
	return s.stream.Send(v)
}

// SendWithContext streams instances of "bffpb.StreamChatStreamingRequest" to
// the "stream_chat" endpoint gRPC stream with context.
func (s *StreamChatClientStream) SendWithContext(ctx context.Context, res string) error {
	return s.Send(res)
}

func (s *StreamChatClientStream) Close() error {
	// Close the send direction of the stream
	return s.stream.CloseSend()
}
