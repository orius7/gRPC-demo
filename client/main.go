package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"os"
    "os/signal"
    "syscall"

	pb "grpc-go-demo/grpc-go-demo/proto" // <-- update if your module path is different
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// LocalCache stores the latest directory received from the server.
type LocalCache struct {
	mu    sync.Mutex
	users map[int32]string
}

func (c *LocalCache) Replace(users map[int32]string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Make a copy so server map mutations won't affect us (defensive).
	copied := make(map[int32]string, len(users))
	for k, v := range users {
		copied[k] = v
	}
	c.users = copied
}

func (c *LocalCache) Print() {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("Local cache:")
	for id, name := range c.users {
		fmt.Printf("  ID: %d, Name: %s\n", id, name)
	}
}

func main() {
    const addr = "localhost:50051"

    cache := &LocalCache{users: make(map[int32]string)}
    backoffBase := time.Second

    // --- HANDLE CTRL-C CLEAN SHUTDOWN ---
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-sig
        fmt.Println("\nShutting down client... clearing cache...")
        cache.Replace(map[int32]string{}) // <-- clear
        cache.Print()
        os.Exit(0)
    }()
    // -------------------------------------

    for {
        conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
        if err != nil {
            log.Printf("failed to dial %s: %v — retrying...", addr, err)
            time.Sleep(backoffBase)
            continue
        }

        client := pb.NewUserDirectoryServiceClient(conn)
        ctx, cancel := context.WithCancel(context.Background())
        stream, err := client.StreamUserDirectory(ctx, &pb.Empty{})
        if err != nil {
            log.Printf("could not start stream: %v (will reconnect)", err)
            cancel()
            conn.Close()
            time.Sleep(backoffBase)
            continue
        }

        log.Println("Connected to server — listening for updates...")

        readErr := func() error {
            for {
                dir, err := stream.Recv()
                if err != nil {
                    return err
                }
                cache.Replace(dir.Users)
                log.Println("Received update from server")
                cache.Print()
            }
        }()

        cancel()
        conn.Close()
        log.Printf("stream ended: %v", readErr)

        time.Sleep(backoffBase)
    }
}

