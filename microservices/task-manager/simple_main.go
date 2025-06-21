package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	taskservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/task_service"
	commentservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/comment_service"
	labelservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/label_service"
	
	tasksvr "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/task_service/server"
	commentsvr "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/comment_service/server"
	labelsvr "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/label_service/server"
	
	task_servicepb "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/task_service/pb"
	comment_servicepb "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/comment_service/pb"
	label_servicepb "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/label_service/pb"
	
	tasksvc "object-t.com/hackz-giganoto/microservices/task-manager/task_service"
	commentsvc "object-t.com/hackz-giganoto/microservices/task-manager/comment_service"
	labelsvc "object-t.com/hackz-giganoto/microservices/task-manager/label_service"
)

func main() {
	fmt.Println("🚀 Starting Task Manager Server...")
	
	// ロガー関数を作成
	logger := func(ctx context.Context, args ...any) {
		fmt.Println(args...)
	}
	
	// サービス実装を作成
	taskService := tasksvc.New(logger)
	commentService := commentsvc.New(logger)
	labelService := labelsvc.New(logger)
	
	// エンドポイントを作成
	taskEndpoints := taskservice.NewEndpoints(taskService)
	commentEndpoints := commentservice.NewEndpoints(commentService)
	labelEndpoints := labelservice.NewEndpoints(labelService)
	
	// gRPCサーバーを作成
	grpcServer := grpc.NewServer()
	
	// サービスを登録
	task_servicepb.RegisterTaskServiceServer(grpcServer, tasksvr.New(taskEndpoints, nil))
	comment_servicepb.RegisterCommentServiceServer(grpcServer, commentsvr.New(commentEndpoints, nil))
	label_servicepb.RegisterLabelServiceServer(grpcServer, labelsvr.New(labelEndpoints, nil))
	
	// リフレクション有効化
	reflection.Register(grpcServer)
	
	// リスナーを作成
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	fmt.Println("✅ gRPC server listening on :8080")
	fmt.Println("🏷️  Default labels available: bug, enhancement, documentation, good first issue")
	fmt.Println("🔧 Use grpcurl to test:")
	fmt.Println("   grpcurl -plaintext localhost:8080 list")
	fmt.Println("   grpcurl -plaintext localhost:8080 label_service.LabelService/List")
	
	// グレースフルシャットダウン
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("\n🛑 Shutting down server...")
		grpcServer.GracefulStop()
	}()
	
	// サーバー開始
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}