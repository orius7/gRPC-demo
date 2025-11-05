package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

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

func newServer() *server {
	return &server{
		users: make(map[int32]string),
		subs:  make([]chan *pb.UserDirectory, 0),
	}
}

func (s *server) StreamUserDirectory(req *pb.Empty, stream pb.UserDirectoryService_StreamUserDirectoryServer) error {
	ch := make(chan *pb.UserDirectory, 1)

	s.mu.Lock()
	s.subs = append(s.subs, ch)
	s.mu.Unlock()

	// Send initial full directory immediately
	ch <- &pb.UserDirectory{Users: s.users}

	for {
		select {
		case update := <-ch:
			if err := stream.Send(update); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *server) AddUser(id int32, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[id] = name
	s.broadcast()
}

func (s *server) DeleteUser(id int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, id)
	s.broadcast()
}

func (s *server) broadcast() {
	dir := &pb.UserDirectory{Users: s.users}
	for _, sub := range s.subs {
		sub <- dir
	}
}



func main() {
	//lis: listen on port 50051
	//err: error handling
	lis, err := net.Listen("tcp", ":50051") //if client connects to this port, server will accept connection
	if err != nil {
		log.Fatalf("failed to listen: %v", err) //log fatal if error occurs
	}

	//creates new gRPC server instance, not started yet
	
	srv := newServer()
	grpcServer := grpc.NewServer()

	//registers the server instance to handle incoming requests for UserDirectoryService
	//RegisterUserDirectoryServiceServer is auto-generated from proto file (helper function)
	pb.RegisterUserDirectoryServiceServer(grpcServer, srv)


	//testing
	go func() {
		time.Sleep(2* time.Second)

		log.Println("Adding initial users to directory...")
		srv.AddUser(1, "Alice")
		srv.AddUser(2, "Bob")
		srv.AddUser(3, "Charlie")

		time.Sleep(5 * time.Second)
		log.Println("Updating...")
		srv.AddUser(4, "Diana")

		time.Sleep(5 * time.Second)
		log.Println("Deleting...")
		srv.DeleteUser(2) // Bob removed
	}()

	fmt.Println("Server is running on port 50051...")
	//Serve(lis) starts gRPC server, starts listening fro incoming connections on the specified listener (lis)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//what now, i need client to continuously listen to server changes (add / delete users)

}
