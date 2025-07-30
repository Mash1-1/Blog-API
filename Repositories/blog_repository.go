package Repositories

import (
	"blog_api/Domain"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepository struct {
	Database *mongo.Collection
}

type BlogRepositoryI interface {
	UpdateBlog(updatedBlog *Domain.Blog) error
}

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

func (BlgRepo *BlogRepository) UpdateBlog(updatedBlog *Domain.Blog) error {
	// Use blog id to search and update task 
	filter := bson.D{{Key: "id", Value: updatedBlog.ID}}
	updatedBSON := bson.M{}

	// Find updatable fields
	if updatedBlog.Title != "" {
		updatedBSON["title"] = updatedBlog.Title 
	}
	if updatedBlog.Content != "" {
		updatedBSON["content"] = updatedBlog.Content
	}
	if updatedBlog.Tags != "" {
		updatedBSON["tags"] = updatedBlog.Tags
	}
	update := bson.M{"$set" : updatedBSON}
	// Do update operation in database
	updatedRes, err := BlgRepo.Database.UpdateOne(context.TODO(), filter, update)
	// Handle exceptions
	if err != nil {
		return err
	}
	if updatedRes.MatchedCount == 0 {
		return errors.New("blog not found")
	}
	return nil
}