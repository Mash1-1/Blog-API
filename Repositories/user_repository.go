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

func (usRepo *UserRepository) CheckExistence(email string) error {
	var existingUser Domain.User
	return usRepo.database.FindOne(context.TODO(), bson.M{"email" : email}).Decode(&existingUser)
}

func (usRepo *UserRepository) GetUser(user *Domain.User) (*Domain.User, error) {
	var existingUser Domain.User 
	err := usRepo.database.FindOne(context.TODO(), bson.M{"email" : user.Email}).Decode(&existingUser)
	return &existingUser, err 
}

func (usRepo *UserRepository) UpdateUser(user *Domain.User) error {
	updateFields := bson.M{"verified" : true}
	updateBSON := bson.D{{Key: "$set", Value: updateFields}}
	_, err := usRepo.database.UpdateOne(context.TODO(), bson.M{"email" : user.Email}, updateBSON)
	return err
}

func (usRepo *UserRepository) DeleteUser(email string) error {
	_, err := usRepo.database.DeleteMany(context.TODO(), bson.M{"email" : email})
	return err 
}

func (usRepo *UserRepository) Register(user *Domain.User) (error) {
	_, err := usRepo.database.InsertOne(context.TODO(), user)
	return err
}