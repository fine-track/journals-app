package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB          *mongo.Database
	RecordsColl *mongo.Collection
)

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatalln("no mongodb uri found")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err.Error())
	}
	DB = client.Database("finetrack")
	RecordsColl = DB.Collection("records")

	return client
}
