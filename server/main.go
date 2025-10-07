package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-go-demo/grpc-go-demo/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserDirectoryServiceServer
	users map[int32]string
}

func (s *server) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	s.users[req.Id] = req.Name
	fmt.Printf("Added user: %v (%d)\n", req.Name, req.Id)

	return &pb.AddUserResponse{
		Message:   fmt.Sprintf("Added user %s", req.Name),
		Directory: &pb.UserDirectory{Users: s.users},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	/*
		pb.RegisterUserDirectoryServiceServer(s, &server{
			users: map[int32]string{
				1: "Alice",
				2: "Bob",
				3: "Charlie",
			},
		})

	*/
	log.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
