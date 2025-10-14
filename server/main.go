package main

import (
	"log"
	"sync"
	"time"

	pb "grpc-go-demo/grpc-go-demo/proto"
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

func (s *server) StreamUserDirectory(req *pb.Empty, stream pb.UserDirectoryService_StreamDirectoryServer) error {
	ch := make(chan *pb.UserDirectory, 1)

	s.mu.Lock()
	s.subs = append(s.subs, ch)
	s.mu.Unlock()

	// Send current directory initially
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

// Server-side functions to modify users
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
	s := newServer()
	updateChan := make(chan *pb.UserDirectory, 1)

	s.mu.Lock()
	s.subs = append(s.subs, updateChan)
	s.mu.Unlock()

	go func() {
		for update := range updateChan {
			log.Printf("client received update: %v", update.Users)
		}
	}()

	s.AddUser(1, "Alice")
	time.Sleep(1 * time.Second)

	s.AddUser(2, "Bob")
	time.Sleep(1 * time.Second)

	s.DeleteUser(1)
	time.Sleep(1 * time.Second)

	s.AddUser(3, "Charlie")
	time.Sleep(1 * time.Second)

}
