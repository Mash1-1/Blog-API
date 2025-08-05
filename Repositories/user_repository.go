package Repositories

import (
	"blog_api/Domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	database *mongo.Collection
}

func NewUserRepository(db *mongo.Collection) *UserRepository{
	return &UserRepository{
		database: db,
	}
}

func InitializeUserDB() (*mongo.Collection, error){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return &mongo.Collection{}, err 
	}
	collection := client.Database("user_db").Collection("users")
	// Delete leftover data from previous uses
	collection.DeleteMany(context.TODO(), bson.D{{}})
	return collection, nil
}

func (usRepo *UserRepository) CheckExistence(email string) bool {
	var existingUser Domain.User
	err := usRepo.database.FindOne(context.TODO(), bson.M{"email" : email}).Decode(&existingUser)
	return err == nil
}

func (usRepo *UserRepository) Register(user *Domain.User) (error) {
	_, err := usRepo.database.InsertOne(context.TODO(), user)
	return err
}