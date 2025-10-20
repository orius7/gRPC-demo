package main

import (
	"fmt"
	"log"
	"net"
	"sync"

	//go package generated from proto file
	pb "grpc-go-demo/grpc-go-demo/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserDirectoryServiceServer
	users map[int32]string
	subs  []chan *pb.UserDirectory
	mu    sync.Mutex
}

func main() {
	//lis: listen on port 50051
	//err: error handling
	lis, err := net.Listen("tcp", ":50051") //if client connects to this port, server will accept connection
	if err != nil {
		log.Fatalf("failed to listen: %v", err) //log fatal if error occurs
	}

	//creates new gRPC server instance, not started yet
	grpcServer := grpc.NewServer()

	//registers the server instance to handle incoming requests for UserDirectoryService
	//RegisterUserDirectoryServiceServer is auto-generated from proto file (helper function)
	pb.RegisterUserDirectoryServiceServer(grpcServer, &server{})

	fmt.Println("Server is listening on port 50051...")
	//Serve(lis) starts gRPC server, starts listening fro incoming connections on the specified listener (lis)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
