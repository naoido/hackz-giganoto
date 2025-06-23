package main

import (
	"fmt"
	"os"

	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	cli "object-t.com/hackz-giganoto/microservices/bff/gen/grpc/cli/bff"
)

func doGRPC(_, host string, _ int, _ bool) (goa.Endpoint, any, error) {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to gRPC server at %s: %v\n", host, err)
	}
	return cli.ParseEndpoint(
		conn,
	)
}

func grpcUsageCommands() string {
	return cli.UsageCommands()
}

func grpcUsageExamples() string {
	return cli.UsageExamples()
}
