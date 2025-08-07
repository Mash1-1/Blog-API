package Repositories

import (
	"blog_api/Domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	database *mongo.Database
}


func NewUserRepository(db *mongo.Database) *UserRepository{
	return &UserRepository{
		database: db,
	}
}

func InitializeUserDB() (*mongo.Database, error){
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return &mongo.Database{}, err 
	}
	database := client.Database("user_db")
	database.CreateCollection(context.TODO(), "users")
	database.CreateCollection(context.TODO(), "pass_reset")
	
	// Clear previous uses from database
	database.Collection("users").DeleteMany(context.TODO(), bson.D{{}})
	database.Collection("pass_reset").DeleteMany(context.TODO(), bson.D{{}})
	return database, nil
}

func (usRepo *UserRepository) UpdatePassword(email, password string) error {
	_, err := usRepo.database.Collection("users").UpdateOne(context.TODO(), bson.M{"email" : email}, bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: password}}}})
	return err
}

func (usRepo *UserRepository) ForgotPassword(data Domain.ResetTokenS) error {	
	_, err := usRepo.database.Collection("pass_reset").InsertOne(context.TODO(), data)
	return err 
}

func (usRepo *UserRepository) CheckExistence(email string) error {
	var existingUser Domain.User
	return usRepo.database.Collection("users").FindOne(context.TODO(), bson.M{"email" : email}).Decode(&existingUser)
}

func (usRepo *UserRepository) GetUser(user *Domain.User) (*Domain.User, error) {
	var existingUser Domain.User 
	err := usRepo.database.Collection("users").FindOne(context.TODO(), bson.M{"email" : user.Email}).Decode(&existingUser)
	return &existingUser, err 
}

func (usRepo *UserRepository) GetTokenData(email string) (Domain.ResetTokenS, error) {
	var data Domain.ResetTokenS
	err := usRepo.database.Collection("pass_reset").FindOne(context.TODO(), bson.M{"email" : email}).Decode(&data)
	return data, err
}

func (usRepo *UserRepository) DeleteTokenData(email string) (error) {
	_, err := usRepo.database.Collection("pass_reset").DeleteMany(context.TODO(), bson.M{"email" : email})
	return err
}

func (usRepo *UserRepository) UpdateUser(user *Domain.User) error {
	updateFields := bson.M{"verified" : true}
	updateBSON := bson.D{{Key: "$set", Value: updateFields}}
	_, err := usRepo.database.Collection("users").UpdateOne(context.TODO(), bson.M{"email" : user.Email}, updateBSON)
	return err
}

func (usRepo *UserRepository) DeleteUser(email string) error {
	_, err := usRepo.database.Collection("users").DeleteMany(context.TODO(), bson.M{"email" : email})
	return err 
}

func (usRepo *UserRepository) Register(user *Domain.User) (error) {
	_, err := usRepo.database.Collection("users").InsertOne(context.TODO(), user)
	return err
}

func (usRepo *UserRepository) GetUserByEmail(email string) (*Domain.User, error) {
	var user Domain.User
	err := usRepo.database.Collection("users").FindOne(context.TODO(), bson.M{"email" : email}).Decode(&user)
	return &user, err
}