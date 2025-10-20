package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "grpc-go-demo/grpc-go-demo/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserDirectoryServiceClient(conn)
	stream, err := client.StreamDirectory(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("could not start stream: %v", err)
	}

	for {
		update, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error receiving update: %v", err)
		}

		fmt.Println("📢 Directory update:")
		for id, name := range update.Users {
			fmt.Printf("  ID: %d, Name: %s\n", id, name)
		}
	}
}
