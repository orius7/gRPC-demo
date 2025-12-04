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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type server struct {
	pb.UnimplementedUserDirectoryServiceServer
	users  map[int32]string
	subs   []chan *pb.UserDirectory
	mu     sync.Mutex
	client *mongo.Client // MongoDB client
}

func newServer(client *mongo.Client) *server {
	return &server{
		users:  make(map[int32]string),
		subs:   make([]chan *pb.UserDirectory, 0),
		client: client, // Initialize MongoDB client
	}
}

func (s *server) StreamUserDirectory(req *pb.Empty, stream pb.UserDirectoryService_StreamUserDirectoryServer) error {
	ch := make(chan *pb.UserDirectory, 1)

	s.mu.Lock()
	s.subs = append(s.subs, ch)
	s.mu.Unlock()

	// Send initial full directory immediately
	users, err := s.retrieveAllUsers() // Fetch from MongoDB
	if err == nil {
		ch <- &pb.UserDirectory{Users: users}
	} else {
		return err
	}

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
	// Add user to MongoDB
	collection := s.client.Database("your_database").Collection("users")
	_, err := collection.InsertOne(context.Background(), pb.User{Id: id, Name: name})
	if err != nil {
		log.Println("Failed to add user:", err)
		return
	}
	s.broadcast()
}

func (s *server) DeleteUser(id int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	collection := s.client.Database("your_database").Collection("users")
	_, err := collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		log.Println("Failed to delete user:", err)
		return
	}
	s.broadcast()
}

func (s *server) broadcast() {
	// Fetch the updated user directory from MongoDB
	users, err := s.retrieveAllUsers()
	if err != nil {
		log.Println("Failed to retrieve users:", err)
		return
	}
	dir := &pb.UserDirectory{Users: users}
	for _, sub := range s.subs {
		sub <- dir
	}
}

// New method to retrieve all users from MongoDB
func (s *server) retrieveAllUsers() (map[int32]string, error) {
	collection := s.client.Database("your_database").Collection("users")
	cursor, err := collection.Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	users := make(map[int32]string)
	for cursor.Next(context.Background()) {
		var user pb.User // Assuming your proto has a User message with ID and Name
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users[user.Id] = user.Name
	}
	return users, nil
}

//add update user details function

func connectToMongoDB() *mongo.Client {

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

	return client
}

func main() {

	client := connectToMongoDB()

	lis, err := net.Listen("tcp", ":50051") //if client connects to this port, server will accept connection
	if err != nil {
		log.Fatalf("failed to listen: %v", err) //log fatal if error occurs
	}

	//creates new gRPC server instance, not started yet

	srv := newServer(client)
	grpcServer := grpc.NewServer()

	//registers the server instance to handle incoming requests for UserDirectoryService
	//RegisterUserDirectoryServiceServer is auto-generated from proto file (helper function)
	pb.RegisterUserDirectoryServiceServer(grpcServer, srv)

	fmt.Println("Server is listening on port 50051...")
	//Serve(lis) starts gRPC server, starts listening fro incoming connections on the specified listener (lis)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//what now, i need client to continuously listen to server changes (add / delete users)

	//connect to MongoDB
	//add data in mongoDB (here)
	//create functions to add, delete, update user details in MongoDB
	//whenever a change occurs in MongoDB, server should broadcast updated user directory to all connected clients
}
