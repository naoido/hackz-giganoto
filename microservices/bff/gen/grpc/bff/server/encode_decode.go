// Code generated by goa v3.21.1, DO NOT EDIT.
//
// bff gRPC server encoders and decoders
//
// Command:
// $ goa gen object-t.com/hackz-giganoto/microservices/bff/design

package server

import (
	"context"
	"strings"

	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc/metadata"
	bff "object-t.com/hackz-giganoto/microservices/bff/gen/bff"
	bffpb "object-t.com/hackz-giganoto/microservices/bff/gen/grpc/bff/pb"
)

// EncodeCreateRoomResponse encodes responses from the "bff" service
// "create_room" endpoint.
func EncodeCreateRoomResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(string)
	if !ok {
		return nil, goagrpc.ErrInvalidType("bff", "create_room", "string", v)
	}
	resp := NewProtoCreateRoomResponse(result)
	return resp, nil
}

// DecodeCreateRoomRequest decodes requests sent to "bff" service "create_room"
// endpoint.
func DecodeCreateRoomRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		token string
		err   error
	)
	{
		if vals := md.Get("authorization"); len(vals) == 0 {
			err = goa.MergeErrors(err, goa.MissingFieldError("authorization", "metadata"))
		} else {
			token = vals[0]
		}
	}
	if err != nil {
		return nil, err
	}
	var payload *bff.CreateRoomPayload
	{
		payload = NewCreateRoomPayload(token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}
	}
	return payload, nil
}

// EncodeHistoryResponse encodes responses from the "bff" service "history"
// endpoint.
func EncodeHistoryResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.([]*bff.EnrichedMessage)
	if !ok {
		return nil, goagrpc.ErrInvalidType("bff", "history", "[]*bff.EnrichedMessage", v)
	}
	resp := NewProtoHistoryResponse(result)
	return resp, nil
}

// DecodeHistoryRequest decodes requests sent to "bff" service "history"
// endpoint.
func DecodeHistoryRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		token string
		err   error
	)
	{
		if vals := md.Get("authorization"); len(vals) == 0 {
			err = goa.MergeErrors(err, goa.MissingFieldError("authorization", "metadata"))
		} else {
			token = vals[0]
		}
	}
	if err != nil {
		return nil, err
	}
	var (
		message *bffpb.HistoryRequest
		ok      bool
	)
	{
		if message, ok = v.(*bffpb.HistoryRequest); !ok {
			return nil, goagrpc.ErrInvalidType("bff", "history", "*bffpb.HistoryRequest", v)
		}
	}
	var payload *bff.HistoryPayload
	{
		payload = NewHistoryPayload(message, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}
	}
	return payload, nil
}

// EncodeRoomListResponse encodes responses from the "bff" service "room-list"
// endpoint.
func EncodeRoomListResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.([]string)
	if !ok {
		return nil, goagrpc.ErrInvalidType("bff", "room-list", "[]string", v)
	}
	resp := NewProtoRoomListResponse(result)
	return resp, nil
}

// DecodeRoomListRequest decodes requests sent to "bff" service "room-list"
// endpoint.
func DecodeRoomListRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		token string
		err   error
	)
	{
		if vals := md.Get("authorization"); len(vals) == 0 {
			err = goa.MergeErrors(err, goa.MissingFieldError("authorization", "metadata"))
		} else {
			token = vals[0]
		}
	}
	if err != nil {
		return nil, err
	}
	var payload *bff.RoomListPayload
	{
		payload = NewRoomListPayload(token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}
	}
	return payload, nil
}

// EncodeJoinRoomResponse encodes responses from the "bff" service "join-room"
// endpoint.
func EncodeJoinRoomResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(string)
	if !ok {
		return nil, goagrpc.ErrInvalidType("bff", "join-room", "string", v)
	}
	resp := NewProtoJoinRoomResponse(result)
	return resp, nil
}

// DecodeJoinRoomRequest decodes requests sent to "bff" service "join-room"
// endpoint.
func DecodeJoinRoomRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		token string
		err   error
	)
	{
		if vals := md.Get("authorization"); len(vals) == 0 {
			err = goa.MergeErrors(err, goa.MissingFieldError("authorization", "metadata"))
		} else {
			token = vals[0]
		}
	}
	if err != nil {
		return nil, err
	}
	var (
		message *bffpb.JoinRoomRequest
		ok      bool
	)
	{
		if message, ok = v.(*bffpb.JoinRoomRequest); !ok {
			return nil, goagrpc.ErrInvalidType("bff", "join-room", "*bffpb.JoinRoomRequest", v)
		}
	}
	var payload *bff.JoinRoomPayload
	{
		payload = NewJoinRoomPayload(message, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}
	}
	return payload, nil
}

// EncodeInviteRoomResponse encodes responses from the "bff" service
// "invite-room" endpoint.
func EncodeInviteRoomResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(string)
	if !ok {
		return nil, goagrpc.ErrInvalidType("bff", "invite-room", "string", v)
	}
	resp := NewProtoInviteRoomResponse(result)
	return resp, nil
}

// DecodeInviteRoomRequest decodes requests sent to "bff" service "invite-room"
// endpoint.
func DecodeInviteRoomRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		token string
		err   error
	)
	{
		if vals := md.Get("authorization"); len(vals) == 0 {
			err = goa.MergeErrors(err, goa.MissingFieldError("authorization", "metadata"))
		} else {
			token = vals[0]
		}
	}
	if err != nil {
		return nil, err
	}
	var (
		message *bffpb.InviteRoomRequest
		ok      bool
	)
	{
		if message, ok = v.(*bffpb.InviteRoomRequest); !ok {
			return nil, goagrpc.ErrInvalidType("bff", "invite-room", "*bffpb.InviteRoomRequest", v)
		}
	}
	var payload *bff.InviteRoomPayload
	{
		payload = NewInviteRoomPayload(message, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}
	}
	return payload, nil
}

// EncodeStreamChatResponse encodes responses from the "bff" service
// "stream_chat" endpoint.
func EncodeStreamChatResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(*bff.EnrichedMessage)
	if !ok {
		return nil, goagrpc.ErrInvalidType("bff", "stream_chat", "*bff.EnrichedMessage", v)
	}
	resp := NewProtoStreamChatResponse(result)
	return resp, nil
}

// DecodeStreamChatRequest decodes requests sent to "bff" service "stream_chat"
// endpoint.
func DecodeStreamChatRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		token  string
		roomID string
		err    error
	)
	{
		if vals := md.Get("authorization"); len(vals) == 0 {
			err = goa.MergeErrors(err, goa.MissingFieldError("authorization", "metadata"))
		} else {
			token = vals[0]
		}
		if vals := md.Get("room_id"); len(vals) == 0 {
			err = goa.MergeErrors(err, goa.MissingFieldError("room_id", "metadata"))
		} else {
			roomID = vals[0]
		}
	}
	if err != nil {
		return nil, err
	}
	var payload *bff.StreamChatPayload
	{
		payload = NewStreamChatPayload(token, roomID)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}
	}
	return payload, nil
}

// EncodeGetProfileResponse encodes responses from the "bff" service
// "get_profile" endpoint.
func EncodeGetProfileResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(*bff.GetProfileResult)
	if !ok {
		return nil, goagrpc.ErrInvalidType("bff", "get_profile", "*bff.GetProfileResult", v)
	}
	resp := NewProtoGetProfileResponse(result)
	return resp, nil
}

// DecodeGetProfileRequest decodes requests sent to "bff" service "get_profile"
// endpoint.
func DecodeGetProfileRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		token string
		err   error
	)
	{
		if vals := md.Get("authorization"); len(vals) == 0 {
			err = goa.MergeErrors(err, goa.MissingFieldError("authorization", "metadata"))
		} else {
			token = vals[0]
		}
	}
	if err != nil {
		return nil, err
	}
	var (
		message *bffpb.GetProfileRequest
		ok      bool
	)
	{
		if message, ok = v.(*bffpb.GetProfileRequest); !ok {
			return nil, goagrpc.ErrInvalidType("bff", "get_profile", "*bffpb.GetProfileRequest", v)
		}
	}
	var payload *bff.GetProfilePayload
	{
		payload = NewGetProfilePayload(message, token)
		if strings.Contains(payload.Token, " ") {
			// Remove authorization scheme prefix (e.g. "Bearer")
			cred := strings.SplitN(payload.Token, " ", 2)[1]
			payload.Token = cred
		}
	}
	return payload, nil
}

// EncodeUpdateProfileResponse encodes responses from the "bff" service
// "update_profile" endpoint.
func EncodeUpdateProfileResponse(ctx context.Context, v any, hdr, trlr *metadata.MD) (any, error) {
	result, ok := v.(*bff.UpdateProfileResult)
	if !ok {
		return nil, goagrpc.ErrInvalidType("bff", "update_profile", "*bff.UpdateProfileResult", v)
	}
	resp := NewProtoUpdateProfileResponse(result)
	return resp, nil
}

// DecodeUpdateProfileRequest decodes requests sent to "bff" service
// "update_profile" endpoint.
func DecodeUpdateProfileRequest(ctx context.Context, v any, md metadata.MD) (any, error) {
	var (
		token *string
		err   error
	)
	{
		if vals := md.Get("authorization"); len(vals) > 0 {
			token = &vals[0]
		}
	}
	if err != nil {
		return nil, err
	}
	var (
		message *bffpb.UpdateProfileRequest
		ok      bool
	)
	{
		if message, ok = v.(*bffpb.UpdateProfileRequest); !ok {
			return nil, goagrpc.ErrInvalidType("bff", "update_profile", "*bffpb.UpdateProfileRequest", v)
		}
	}
	var payload *bff.UpdateProfilePayload
	{
		payload = NewUpdateProfilePayload(message, token)
		if payload.Token != nil {
			if strings.Contains(*payload.Token, " ") {
				// Remove authorization scheme prefix (e.g. "Bearer")
				cred := strings.SplitN(*payload.Token, " ", 2)[1]
				payload.Token = &cred
			}
		}
	}
	return payload, nil
}
