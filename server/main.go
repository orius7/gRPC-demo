package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-go-demo/grpc-go-demo/proto"
	"google.golang.org/grpc"
)

type contractServer struct {
	pb.UnimplementedContractServiceServer
}

func (s *contractServer) SayHello(ctx context.Context, req *pb.ContractRequest) (*pb.ContractResponse, error) {
	return &pb.ContractResponse{Message: "Hello, " + req.GetName()}, nil
}

func main() {
	//create a default User Directory using proto schema 
	dir := &pb.UserDirectory{
		Users: map[int32]string{
			1: "Alice",
			2: "Bob",
			3: "Charlie",
		},
	}

	//next steps:
	//allows server to update dictionary
	//allows client to query dictionary
	
	_ = dir // intentionally tells compiler to ignore unused variable, "im aware of dir but im not using it right now"
	//go doesn't allow unused variables

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
