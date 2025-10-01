package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-go-demo/proto"
	"google.golang.org/grpc"
)

type contractServer struct {
	pb.UnimplementedContractServiceServer
}

func (s *contractServer) SayHello(ctx context.Context, req *pb.ContractRequest) (*pb.ContractResponse, error) {
	return &pb.ContractResponse{Message: "Hello, " + req.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterContractServiceServer(grpcServer, &contractServer{})

	fmt.Println("Server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
