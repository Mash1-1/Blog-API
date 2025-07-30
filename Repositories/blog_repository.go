package Repositories

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepository struct {
	Database *mongo.Collection
}

type BlogRepositoryI interface {}

func NewBlogRepository(db *mongo.Collection) *BlogRepository{
	return &BlogRepository{
		Database: db,
	}
}

func InitializeBlogDB() (*mongo.Collection, error){
	// Initialize collection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return &mongo.Collection{}, err 
	}

	collection := client.Database("Blog_DB").Collection("blogs")
	// Clear previous usageleftover data
	collection.DeleteMany(context.TODO(), bson.D{{}})
	return collection, nil
}