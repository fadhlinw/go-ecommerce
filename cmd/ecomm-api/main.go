package main

import (
	"log"

	"github.com/fadhlinw/go-ecommerce/db"
	"github.com/fadhlinw/go-ecommerce/ecomm-api/storer"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading it")
	}

	// Initialize the database connection.
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Ensure the connection is safely closed when the application shuts down
	defer dbConn.Close()

	log.Println("Database connection initialized successfully!")
	
	// Inject the raw sqlx database connection into our Storer.
	// This storer is now ready to be used by your application's logic or HTTP handlers.
	appStorer := storer.NewStorer(dbConn.GetDB())
	_ = appStorer // Placeholder to prevent "unused variable" error.
	
	log.Println("Storer initialized successfully! API is ready to serve.")
}
