package Repositories

import (
	"blog_api/Domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	UserCollection *mongo.Collection
	ResetPassword  *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		UserCollection: db.Collection("users"),
		ResetPassword:  db.Collection("pass_reset"),
	}
}

func (usRepo *UserRepository) UpdatePassword(email, password string) error {
	_, err := usRepo.UserCollection.UpdateOne(context.TODO(), bson.M{"email": email}, bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: password}}}})
	return err
}

func (usRepo *UserRepository) ForgotPassword(data Domain.ResetTokenS) error {
	_, err := usRepo.ResetPassword.InsertOne(context.TODO(), data)
	return err
}

func (usRepo *UserRepository) CheckExistence(email string) error {
	var existingUser Domain.User
	return usRepo.UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&existingUser)
}

func (usRepo *UserRepository) GetUser(user *Domain.User) (*Domain.User, error) {
	var existingUser Domain.User
	err := usRepo.UserCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	return &existingUser, err
}

func (usRepo *UserRepository) GetTokenData(email string) (Domain.ResetTokenS, error) {
	var data Domain.ResetTokenS
	err := usRepo.ResetPassword.FindOne(context.TODO(), bson.M{"email": email}).Decode(&data)
	return data, err
}

func (usRepo *UserRepository) DeleteTokenData(email string) error {
	_, err := usRepo.ResetPassword.DeleteMany(context.TODO(), bson.M{"email": email})
	return err
}

func (usRepo *UserRepository) UpdateUser(user *Domain.User) error {
	updateFields := bson.M{"verified": true}
	updateBSON := bson.D{{Key: "$set", Value: updateFields}}
	_, err := usRepo.UserCollection.UpdateOne(context.TODO(), bson.M{"email": user.Email}, updateBSON)
	return err
}

func (usRepo *UserRepository) DeleteUser(email string) error {
	_, err := usRepo.UserCollection.DeleteMany(context.TODO(), bson.M{"email": email})
	return err
}

func (usRepo *UserRepository) Register(user *Domain.User) error {
	_, err := usRepo.UserCollection.InsertOne(context.TODO(), user)
	return err
}

func (usRepo *UserRepository) GetUserByEmail(email string) (*Domain.User, error) {
	var user Domain.User
	err := usRepo.UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (usRepo *UserRepository) UpdateUserProfile(user *Domain.User) (*Domain.User, error) {
	updateFields := bson.M{
		"username": user.Username,
		"bio":      user.Bio,
	}

	updateBSON := bson.D{{Key: "$set", Value: updateFields}}

	// Perform the update
	_, err := usRepo.UserCollection.UpdateOne(
		context.TODO(),
		bson.M{"email": user.Email},
		updateBSON,
	)
	if err != nil {
		return nil, err
	}

	// Fetch the updated user
	var updatedUser Domain.User
	err = usRepo.UserCollection.FindOne(
		context.TODO(),
		bson.M{"email": user.Email},
	).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (usRepo *UserRepository) UpdateUserRole(email string, role string) (*Domain.User, error) {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"role": role}}

	_, err := usRepo.UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	var updatedUser Domain.User
	err = usRepo.UserCollection.FindOne(context.TODO(), filter).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
