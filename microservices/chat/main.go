package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	chatservice "object-t.com/hackz-giganoto/microservices/chat/gen/chat"
	chatpb "object-t.com/hackz-giganoto/microservices/chat/gen/grpc/chat/pb"
	chatgrpc "object-t.com/hackz-giganoto/microservices/chat/gen/grpc/chat/server"
	profileservice "object-t.com/hackz-giganoto/microservices/chat/gen/profile"
	profilepb "object-t.com/hackz-giganoto/microservices/chat/gen/grpc/profile/pb"
	profilegrpc "object-t.com/hackz-giganoto/microservices/chat/gen/grpc/profile/server"
)

func main() {
	var (
		chatSvc    = NewChatService()
		profileSvc = NewProfileService()
	)

	chatEndpoints := chatservice.NewEndpoints(chatSvc)
	profileEndpoints := profileservice.NewEndpoints(profileSvc)

	chatServer := chatgrpc.New(chatEndpoints, nil, nil)
	profileServer := profilegrpc.New(profileEndpoints, nil)

	grpcServer := grpc.NewServer()
	chatpb.RegisterChatServer(grpcServer, chatServer)
	profilepb.RegisterProfileServer(grpcServer, profileServer)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on :50052")

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("shutting down gRPC server...")
	grpcServer.GracefulStop()
}