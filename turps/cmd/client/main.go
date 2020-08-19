package main

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "turps/api"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTurpsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.UpsertChangeList(ctx, &pb.UpsertChangeListRequest{ChangeList: &pb.ChangeList{
		Tz:      ptypes.TimestampNow(),
		TestIds: []string{"test1", "test2"},
	}})
	if err != nil {
		log.Fatalf("could not create: %v", err)
	}
	log.Printf("Change list id is %s", r.ChangeList.ChangeListId)

}
