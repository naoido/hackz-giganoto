package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"goa.design/clue/debug"
	"goa.design/clue/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// 自分で実装したサービスロジック
	commentservice "object-t.com/hackz-giganoto/microservices/task-manager/comment_service"
	labelservice "object-t.com/hackz-giganoto/microservices/task-manager/label_service"
	taskservice "object-t.com/hackz-giganoto/microservices/task-manager/task_service"

	// Goaが生成したgRPCサーバー関連
	commentpb "object-t.com/hackz-giganoto/microservices/task-manager/gen/comment_service"
	commentsvr "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/comment_service/server"
	labelsvr "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/label_service/server"
	tasksvr "object-t.com/hackz-giganoto/microservices/task-manager/gen/grpc/task_service/server"
	labelpb "object-t.com/hackz-giganoto/microservices/task-manager/gen/label_service"
	taskpb "object-t.com/hackz-giganoto/microservices/task-manager/gen/task_service"
)

func main() {
	var (
		hostF     = flag.String("host", "localhost", "Server host (valid values: localhost)")
		domainF   = flag.String("domain", "", "Host domain name (overrides host domain specified in service design)")
		grpcPortF = flag.String("grpc-port", "8080", "gRPC port (overrides host gRPC port specified in service design)")
		secureF   = flag.Bool("secure", false, "Use secure scheme (https or grpcs)")
		dbgF      = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	format := log.FormatJSON
	if log.IsTerminal() {
		format = log.FormatTerminal
	}
	ctx := log.Context(context.Background(), log.WithFormat(format))
	if *dbgF {
		ctx = log.Context(ctx, log.WithDebug())
		log.Debugf(ctx, "debug logs enabled")
	}

	// 1. サービスを初期化するところを修正！
	var (
		taskServiceSvc    taskpb.Service
		commentServiceSvc commentpb.Service
		labelServiceSvc   labelpb.Service
	)
	{
		// 自分で実装したサービスのコンストラクタ（New...）を呼ぶよ
		taskServiceSvc = taskservice.New(log.Context(ctx, log.WithPrefix("taskService")))
		commentServiceSvc = commentservice.New(log.Context(ctx, log.WithPrefix("commentService")))
		labelServiceSvc = labelservice.New(log.Context(ctx, log.WithPrefix("labelService")))
	}

	// Wrap the services in endpoints that can be invoked from other services.
	var (
		taskServiceEndpoints    *taskpb.Endpoints
		commentServiceEndpoints *commentpb.Endpoints
		labelServiceEndpoints   *labelpb.Endpoints
	)
	{
		taskServiceEndpoints = taskpb.NewEndpoints(taskServiceSvc)
		taskServiceEndpoints.Use(debug.LogPayloads())
		taskServiceEndpoints.Use(log.Endpoint)

		commentServiceEndpoints = commentpb.NewEndpoints(commentServiceSvc)
		commentServiceEndpoints.Use(debug.LogPayloads())
		commentServiceEndpoints.Use(log.Endpoint)

		labelServiceEndpoints = labelpb.NewEndpoints(labelServiceSvc)
		labelServiceEndpoints.Use(debug.LogPayloads())
		labelServiceEndpoints.Use(log.Endpoint)
	}

	errc := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)

	switch *hostF {
	case "localhost":
		{
			addr := "grpc://localhost:" + *grpcPortF
			u, err := url.Parse(addr)
			if err != nil {
				log.Fatalf(ctx, err, "invalid URL %#v\n", addr)
			}
			if *secureF {
				u.Scheme = "grpcs"
			}
			if *domainF != "" {
				u.Host = *domainF
			}
			handleGRPCServer(ctx, u, taskServiceEndpoints, commentServiceEndpoints, labelServiceEndpoints, &wg, errc, *dbgF)
		}

	default:
		log.Fatal(ctx, fmt.Errorf("invalid host argument: %q (valid hosts: localhost)", *hostF))
	}

	log.Printf(ctx, "exiting (%v)", <-errc)

	cancel()

	wg.Wait()
	log.Printf(ctx, "exited")
}

// 2. 足りなかった handleGRPCServer 関数を追加！
func handleGRPCServer(ctx context.Context, u *url.URL, taskEnpts *taskpb.Endpoints, commentEnpts *commentpb.Endpoints, labelEnpts *labelpb.Endpoints, wg *sync.WaitGroup, errc chan error, debug bool) {
	// Setup gRPC server.
	var opts []grpc.ServerOption
	if debug {
		opts = append(opts, grpc.UnaryInterceptor(debug.UnaryServerInterceptor()))
	}
	grpcServer := grpc.NewServer(opts...)

	// Register the services.
	tasksvr.Register(grpcServer, tasksvr.New(taskEnpts, nil))
	commentsvr.Register(grpcServer, commentsvr.New(commentEnpts, nil))
	labelsvr.Register(grpcServer, labelsvr.New(labelEnpts, nil))

	// For debugging purposes, allows tools like grpcurl to query the server.
	reflection.Register(grpcServer)

	wg.Add(1)
	go func() {
		defer wg.Done()

		// Start gRPC server in a separate goroutine.
		go func() {
			lis, err := net.Listen("tcp", u.Host)
			if err != nil {
				errc <- err
				return
			}
			log.Printf(ctx, "gRPC server listening on %q", u.Host)
			errc <- grpcServer.Serve(lis)
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down gRPC server at %q", u.Host)
		grpcServer.GracefulStop()
	}()
}
