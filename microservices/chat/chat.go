package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"goa.design/goa/v3/security"

	chatservice "object-t.com/hackz-giganoto/microservices/chat/gen/chat"
)

type ChatServiceImpl struct {
	rooms   map[string]*ChatRoom
	mu      sync.RWMutex
	clients map[string]map[string]chan *chatservice.JoinChatResult
}

type ChatRoom struct {
	ID       string
	Messages []*chatservice.ChatMessage
	mu       sync.RWMutex
}

func NewChatService() *ChatServiceImpl {
	return &ChatServiceImpl{
		rooms:   make(map[string]*ChatRoom),
		clients: make(map[string]map[string]chan *chatservice.JoinChatResult),
	}
}

func (c *ChatServiceImpl) SendMessage(ctx context.Context, p *chatservice.SendMessagePayload) (*chatservice.SendMessageResult, error) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		return nil, chatservice.Unauthorized("invalid token")
	}

	if p.RoomID == "" || p.Message == "" {
		return nil, chatservice.BadRequest("room_id and message are required")
	}

	messageID := uuid.New().String()
	timestamp := time.Now().Format(time.RFC3339)

	c.mu.Lock()
	room, exists := c.rooms[p.RoomID]
	if !exists {
		room = &ChatRoom{
			ID:       p.RoomID,
			Messages: make([]*chatservice.ChatMessage, 0),
		}
		c.rooms[p.RoomID] = room
	}
	c.mu.Unlock()

	messageType := "text"
	if p.MessageType != nil {
		messageType = *p.MessageType
	}

	userName := c.getUserName(userID)
	
	message := &chatservice.ChatMessage{
		MessageID:   messageID,
		UserID:      userID,
		UserName:    userName,
		Message:     p.Message,
		MessageType: &messageType,
		Timestamp:   timestamp,
	}

	room.mu.Lock()
	room.Messages = append(room.Messages, message)
	room.mu.Unlock()

	c.broadcastMessage(p.RoomID, &chatservice.JoinChatResult{
		MessageID:   messageID,
		UserID:      userID,
		UserName:    userName,
		Message:     p.Message,
		MessageType: &messageType,
		Timestamp:   timestamp,
	})

	return &chatservice.SendMessageResult{
		MessageID: messageID,
		Timestamp: timestamp,
	}, nil
}

func (c *ChatServiceImpl) JoinChat(ctx context.Context, p *chatservice.JoinChatPayload, stream chatservice.JoinChatServerStream) error {
	_, err := c.getUserIDFromContext(ctx)
	if err != nil {
		return chatservice.Unauthorized("invalid token")
	}

	if p.RoomID == "" {
		return chatservice.BadRequest("room_id is required")
	}

	clientID := uuid.New().String()
	messageChan := make(chan *chatservice.JoinChatResult, 100)

	c.mu.Lock()
	if c.clients[p.RoomID] == nil {
		c.clients[p.RoomID] = make(map[string]chan *chatservice.JoinChatResult)
	}
	c.clients[p.RoomID][clientID] = messageChan
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.clients[p.RoomID], clientID)
		if len(c.clients[p.RoomID]) == 0 {
			delete(c.clients, p.RoomID)
		}
		c.mu.Unlock()
		close(messageChan)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-messageChan:
			if !ok {
				return nil
			}
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
	}
}

func (c *ChatServiceImpl) GetChatHistory(ctx context.Context, p *chatservice.GetChatHistoryPayload) (*chatservice.GetChatHistoryResult, error) {
	_, err := c.getUserIDFromContext(ctx)
	if err != nil {
		return nil, chatservice.Unauthorized("invalid token")
	}

	if p.RoomID == "" {
		return nil, chatservice.BadRequest("room_id is required")
	}

	c.mu.RLock()
	room, exists := c.rooms[p.RoomID]
	c.mu.RUnlock()

	if !exists {
		return nil, chatservice.NotFound("room not found")
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	totalCount := len(room.Messages)
	messages := make([]*chatservice.ChatMessage, 0)

	limit := 50
	if p.Limit != nil && *p.Limit > 0 {
		limit = *p.Limit
	}

	offset := 0
	if p.Offset != nil && *p.Offset > 0 {
		offset = *p.Offset
	}

	start := offset
	end := offset + limit
	if start >= totalCount {
		return &chatservice.GetChatHistoryResult{
			Messages:   messages,
			TotalCount: totalCount,
		}, nil
	}

	if end > totalCount {
		end = totalCount
	}

	for i := start; i < end; i++ {
		messages = append(messages, room.Messages[i])
	}

	return &chatservice.GetChatHistoryResult{
		Messages:   messages,
		TotalCount: totalCount,
	}, nil
}

func (c *ChatServiceImpl) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	if token == "" {
		return ctx, fmt.Errorf("missing token")
	}
	
	userID := c.extractUserIDFromToken(token)
	if userID == "" {
		return ctx, fmt.Errorf("invalid token")
	}
	
	return context.WithValue(ctx, "userID", userID), nil
}

func (c *ChatServiceImpl) broadcastMessage(roomID string, message *chatservice.JoinChatResult) {
	c.mu.RLock()
	clients, exists := c.clients[roomID]
	c.mu.RUnlock()

	if !exists {
		return
	}

	var wg sync.WaitGroup
	for _, clientChan := range clients {
		wg.Add(1)
		go func(ch chan *chatservice.JoinChatResult) {
			defer wg.Done()
			select {
			case ch <- message:
			case <-time.After(5 * time.Second):
			}
		}(clientChan)
	}
	wg.Wait()
}

func (c *ChatServiceImpl) getUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func (c *ChatServiceImpl) extractUserIDFromToken(token string) string {
	return "user_" + token[:8]
}

func (c *ChatServiceImpl) getUserName(userID string) string {
	return "User " + userID[5:]
}