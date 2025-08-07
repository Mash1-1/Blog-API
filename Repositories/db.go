package Repositories

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	db *mongo.Database
}

func InitializeDb() *mongo.Database {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Can't load environment variables")
	}

	DB_URL := os.Getenv("MONGODB_URL")

	clientOptions := options.Client().ApplyURI(DB_URL)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatalf("unable to connect to database error: %s", err)
	}

	return client.Database("Blog_DB")
}
