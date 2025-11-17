package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	//go package generated from proto file
	pb "grpc-go-demo/grpc-go-demo/proto"

	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
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

func (s *server) RetreiveAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.broadcast
}

//add update user details function

func main() {

	//get username and password from environment variables
	uri := os.Getenv("URI")
	if uri == "" {
		log.Fatal("Environment variable URI is not set")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	//lis: listen on port 50051
	//err: error handling
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	lis, err := net.Listen("tcp", ":50051") //if client connects to this port, server will accept connection
	if err != nil {
		log.Fatalf("failed to listen: %v", err) //log fatal if error occurs
	}

	//creates new gRPC server instance, not started yet

	srv := newServer()
	grpcServer := grpc.NewServer()

	//registers the server instance to handle incoming requests for UserDirectoryService
	//RegisterUserDirectoryServiceServer is auto-generated from proto file (helper function)
	pb.RegisterUserDirectoryServiceServer(grpcServer, &server{})

	fmt.Println("Server is listening on port 50051...")
	//Serve(lis) starts gRPC server, starts listening fro incoming connections on the specified listener (lis)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//what now, i need client to continuously listen to server changes (add / delete users)

}
