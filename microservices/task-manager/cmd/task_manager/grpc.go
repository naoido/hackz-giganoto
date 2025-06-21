package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"

	"goa.design/clue/debug"
	"goa.design/clue/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	commentservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/comment_service"
	comment_servicepb "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/comment_service/pb"
	commentservicesvr "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/comment_service/server"
	label_servicepb "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/label_service/pb"
	labelservicesvr "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/label_service/server"
	task_servicepb "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/task_service/pb"
	taskservicesvr "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/task_service/server"
	labelservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/label_service"
	taskservice "object-t.com/hackz-giganoto/microservices/task-manager/gen/task_service"
)

// handleGRPCServer starts configures and starts a gRPC server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleGRPCServer(ctx context.Context, u *url.URL, taskServiceEndpoints *taskservice.Endpoints, commentServiceEndpoints *commentservice.Endpoints, labelServiceEndpoints *labelservice.Endpoints, wg *sync.WaitGroup, errc chan error, dbg bool) {

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to gRPC requests and
	// responses.
	var (
		taskServiceServer    *taskservicesvr.Server
		commentServiceServer *commentservicesvr.Server
		labelServiceServer   *labelservicesvr.Server
	)
	{
		taskServiceServer = taskservicesvr.New(taskServiceEndpoints, nil)
		commentServiceServer = commentservicesvr.New(commentServiceEndpoints, nil)
		labelServiceServer = labelservicesvr.New(labelServiceEndpoints, nil)
	}

	// Create interceptor which sets up the logger in each request context.
	chain := grpc.ChainUnaryInterceptor(log.UnaryServerInterceptor(ctx))
	if dbg {
		// Log request and response content if debug logs are enabled.
		chain = grpc.ChainUnaryInterceptor(log.UnaryServerInterceptor(ctx), debug.UnaryServerInterceptor())
	}

	// Initialize gRPC server
	srv := grpc.NewServer(chain)

	// Register the servers.
	task_servicepb.RegisterTaskServiceServer(srv, taskServiceServer)
	comment_servicepb.RegisterCommentServiceServer(srv, commentServiceServer)
	label_servicepb.RegisterLabelServiceServer(srv, labelServiceServer)

	for svc, info := range srv.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Printf(ctx, "serving gRPC method %s", svc+"/"+m.Name)
		}
	}

	// Register the server reflection service on the server.
	// See https://grpc.github.io/grpc/core/md_doc_server-reflection.html.
	reflection.Register(srv)

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start gRPC server in a separate goroutine.
		go func() {
			lis, err := net.Listen("tcp", u.Host)
			if err != nil {
				errc <- err
			}
			if lis == nil {
				errc <- fmt.Errorf("failed to listen on %q", u.Host)
			}
			log.Printf(ctx, "gRPC server listening on %q", u.Host)
			errc <- srv.Serve(lis)
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down gRPC server at %q", u.Host)
		srv.Stop()
	}()
}
