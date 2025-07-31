package Repositories

import (
	"blog_api/Domain"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepository struct {
	Database *mongo.Collection
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
	//check if the collection exists befor crearing one
	created, err := collectionExists(db, "blogs")
	if err != nil {
		log.Fatalln(err)
	}
	if created {
		// Clear previous usageleftover data
		collection := db.Collection("blogs")
		collection.DeleteMany(context.TODO(), bson.D{{}})
		return collection, nil
	}

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"title":    "Blog object Validation",
			"required": []string{"ID", "Title", "Content"},
			"properties": bson.M{
				"ID": bson.M{
					"bsonType": "string",
				},
				"Title": bson.M{
					"bsonType":    "string",
					"description": "Name must be a string of lenght 20 and is a required field",
					"maxLength":   20,
				},
				"Content": bson.M{
					"bsonType":    "string",
					"minLength":   20,
					"description": "Content must be a string with minimum length of 20 characters",
				},
				"Owner": bson.M{
					"bsonType":    "object",
					"description": "Owner is User type with fields",
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
	return collection, nil
}
func (BlgRepo *BlogRepository) Create(blog *Domain.Blog) error {
	_, err := BlgRepo.Database.InsertOne(context.TODO(), blog.ToBlogDTO())
	return err
}

func (BlgRepo *BlogRepository) SearchBlog(searchBlog *Domain.Blog) ([]Domain.Blog, error) {
	// Use filter to search for tasks with the given fields
	searchBSON := bson.M{}
	blogs := []Domain.Blog{}
	if searchBlog.Title != "" {
		searchBSON["Title"] = searchBlog.Title
	}
	var tmp = Domain.User{}
	if searchBlog.Owner != tmp {
		searchBSON["Owner"] = searchBlog.Owner
	}

	filter := bson.M{
		"$or": []bson.M{
			{"Title": searchBSON["Title"]},
			{"Owner": searchBSON["Owner"]},
		},
	}
	cursor, err := BlgRepo.Database.Find(context.TODO(), filter)
	if err != nil {
		return []Domain.Blog{}, err
	}
	for cursor.Next(context.TODO()) {
		var elem Domain.Blog
		err = cursor.Decode(&elem)
		if err != nil {
			return []Domain.Blog{}, err
		}
		blogs = append(blogs, elem)
	}
	if err = cursor.Err(); err != nil {
		return []Domain.Blog{}, err
	}
	return blogs, nil
}

func (BlgRepo *BlogRepository) UpdateBlog(updatedBlog *Domain.Blog) error {
	// Use blog ID to search and update task
	filter := bson.D{{Key: "ID", Value: updatedBlog.ID}}
	updatedBSON := bson.M{}

	// Find updatable fields
	if updatedBlog.Title != "" {
		updatedBSON["Title"] = updatedBlog.Title
	}
	if updatedBlog.Content != "" {
		updatedBSON["Content"] = updatedBlog.Content
	}
	if updatedBlog.Tags != nil {
		updatedBSON["Tags"] = updatedBlog.Tags
	}
	update := bson.M{"$set": updatedBSON}
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

func (BlgRepo *BlogRepository) DeleteBlog(ID string) error {
	filter := bson.D{{Key: "ID", Value: bson.D{{Key: "$eq", Value: ID}}}}
	result, err := BlgRepo.Database.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("blog not found")
	}
	return nil
}

// checks if the given db is created or not
func collectionExists(client *mongo.Database, collname string) (bool, error) {
	dbs, err := client.ListCollectionNames(context.TODO(), bson.D{{Key: "name", Value: collname}})
	if err != nil {
		return false, err
	}

	for _, name := range dbs {
		if name == collname {
			return true, nil
		}
	}
	return false, nil
}
