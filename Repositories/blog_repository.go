package Repositories

import (
	"blog_api/Domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepository struct {
	Database *mongo.Collection
}

type BlogRepositoryI interface {
	Create(Domain.Blog) error
}

func NewBlogRepository(db *mongo.Collection) *BlogRepository {
	return &BlogRepository{
		Database: db,
	}
}

func InitializeBlogDB() (*mongo.Collection, error) {
	// Initialize collection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return &mongo.Collection{}, err
	}
	db := client.Database("Blog_DB")

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"title":    "Blog object Validation",
			"required": []string{"_id", "title", "content", "owner"},
			"properties": bson.M{
				"_id": bson.M{
					"bsonType": "objectId",
				},
				"title": bson.M{
					"bsonType":    "string",
					"description": "Name must be a string of lenght 20 and is a required field",
					"maxLength":   20,
				},
				"content": bson.M{
					"bsonType":    "string",
					"minLength":   20,
					"description": "content must be a string with minimum length of 20 characters",
				},
				"owner": bson.M{
					"bsonType":    "object",
					"description": "owner is User type with fields",
					"required":    []string{"username", "email", "password", "role"},
					"properties": bson.M{
						"username": bson.M{
							"bsonType":    "string",
							"maxLength":   20,
							"description": "username should be length of less than 20 characters",
						},
						"email": bson.M{
							"bsonType":  "string",
							"minLength": 8,
							"pattern":   `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
						},
						"password": bson.M{
							"bsonType":  "string",
							"minLength": 8,
							"maxLength": 20,
							"pattern": `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#\$%\^&\*])[A-Za-z\d!@#\$%\^&\*]{8,}$
`,
						},
						"role": bson.M{
							"bsonType":    "string",
							"enum":        []string{"user", "admin"},
							"description": "Role must be one of 'admin',or 'user'",
						},
					},
				},
			},
		},
	}

	opts := options.CreateCollection().SetValidator(validator)

	err = db.CreateCollection(context.TODO(), "blogs", opts)
	if err != nil {
		return nil, err
	}

	collection := db.Collection("blogs")
	// Clear previous usageleftover data
	collection.DeleteMany(context.TODO(), bson.D{{}})
	return collection, nil
}

func (BlgRepo *BlogRepository) Create(blog Domain.Blog) error {
	_, err := BlgRepo.Database.InsertOne(context.TODO(), blog)
	return err
}
