package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/philipstanev/Miku-stream/internal/auth"
	authdb "github.com/philipstanev/Miku-stream/internal/auth/db"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("no .env file loaded: %v", err)
	}
	fmt.Println("starting program")
	var authService auth.Service
	authService.S = auth.NewInMemorySession()
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("created pool")
	defer dbpool.Close()
	authService.Q = authdb.New(dbpool)

	mux := http.NewServeMux()
	authService.AuthRoutes(mux)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	fmt.Println("Started server")

}
