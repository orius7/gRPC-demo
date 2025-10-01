package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "grpc-go-demo/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewContractServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.ContractRequest{Name: "Aidan"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Println("Response:", resp.GetMessage())
}
