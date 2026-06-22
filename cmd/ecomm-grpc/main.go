package main

import (
	"log"
	"net"

	"github.com/fadhlinw/go-ecommerce/db"
	pb "github.com/fadhlinw/go-ecommerce/ecomm-grpc/pb"
	"github.com/fadhlinw/go-ecommerce/ecomm-grpc/server"
	"github.com/fadhlinw/go-ecommerce/ecomm-grpc/storer"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading it")
	}

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbConn.Close()

	log.Println("Database connection initialized successfully!")

	appStorer := storer.NewStorer(dbConn.GetDB())

	grpcServer := grpc.NewServer()

	ecommServer := server.NewEcommServer(appStorer)

	pb.RegisterEcommServiceServer(grpcServer, ecommServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Println("gRPC Server is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
