package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"net"
	"os"
	"turps/internal/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
	url  = flag.String("url", "", "The PostgreSQL database URL")
)

func main() {
	flag.Parse()

	if *url == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()
	s := grpc.NewTurpsServer(ctx, *url)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("failed to serve: %v", err)
	}
}
