package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/fine-track/journals-app/db"
	"github.com/fine-track/journals-app/services"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	dbClient := db.ConnectDB()
	defer func() { dbClient.Disconnect(context.TODO()) }()

	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	log.Printf("listening on: %s", listener.Addr())

	s := grpc.NewServer()
	services.RegisterRecordsService(s)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
	defer func() { s.Stop() }()
}
