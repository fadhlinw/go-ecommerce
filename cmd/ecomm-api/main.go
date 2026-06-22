package main

import (
	"log"
	"net/http"

	"github.com/fadhlinw/go-ecommerce/ecomm-api/handler"
	pb "github.com/fadhlinw/go-ecommerce/ecomm-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ecommClient := pb.NewEcommServiceClient(conn)
	apiHandler := handler.NewHandler(ecommClient)

	// Setup Routes
	r := handler.SetupRoutes(apiHandler)

	log.Println("REST API Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
