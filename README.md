## **gRPC Demo (Go)**

A simple gRPC demo project written in Go that establishes a server–client relationship using Protocol Buffers (protobuf) and gRPC.

The goal of this project is to build a system where:

The server maintains a dictionary-like database of user data (e.g., names, IDs, phone numbers).

The client always keeps a complete, up-to-date copy of that database locally.

Whenever the server updates, deletes, or modifies data, the client detects the change and automatically synchronizes its copy to match the server’s version.

## Project Overview

Server

- Stores user records (ID, name, phone number, etc.) in an in-memory or external database (TBD).

- Exposes gRPC methods to:

  - Add a new user

  - Update user details

  - Delete a user

  - Retrieve all users

- Notifies connected clients when updates occur.

Client

- Connects to the gRPC server.

- Maintains a local copy of the server’s database.

- Listens for update notifications and automatically synchronizes when changes are detected.

## Tech Stack

Language: Go

RPC Framework: gRPC

Schema Definition: Protocol Buffers (.proto files)

Database: TBD (currently in-memory dictionary)

## Getting Started
Prerequisites

- Go 1.20+

- protoc (Protocol Buffers compiler)

- gRPC and protobuf plugins for Go:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Clone the Repository

```
git clone https://github.com/orius7/gRPC-demo.git
cd gRPC-demo
```

## Compile the Proto File

```
protoc --go_out=. --go-grpc_out=. proto/contract.proto
```

## Run the Server

```
go run ./server
```
## Run the Client

```
go run ./client
```


## Project Structure

```
gRPC-demo/
├── proto/                # .proto files and generated Go code
│   └── contract.proto
├── server/               # gRPC server implementation
│   └── main.go
├── client/               # gRPC client implementation
│   └── main.go
├── go.mod
└── README.md
```

## Next Steps

- Add persistent database support (e.g., SQLite, PostgreSQL, or Firestore).

- Implement real-time client updates using gRPC streams.

- Add authentication and data validation.

- Dockerize the server for easier deployment.

## License

This project is licensed under the [MIT License](https://github.com/orius7/gRPC-demo?tab=MIT-1-ov-file).
