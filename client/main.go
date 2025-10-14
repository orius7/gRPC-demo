package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "grpc-go-demo/grpc-go-demo/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserDirectoryServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.AddUser(ctx, &pb.AddUserRequest{Id: 4, Name: "TESTER"})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	fmt.Println("Updated User Directory:")

	for id, name := range resp.GetUserDirectory().GetUsers() {
		fmt.Printf("  ID: %d, Name: %s\n", id, name)
	}

}
