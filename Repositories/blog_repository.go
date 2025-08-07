package Repositories

import (
	"blog_api/Domain"
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepository struct {
	BlogCollection *mongo.Collection
}

func NewBlogRepository(db *mongo.Database) *BlogRepository {

	return &BlogRepository{
		BlogCollection: db.Collection("blogs"),
	}
}

func (BlgRepo *BlogRepository) Create(blog *Domain.Blog) error {
	_, err := BlgRepo.BlogCollection.InsertOne(context.TODO(), blog)
	return err
}

func (BlgRepo *BlogRepository) SearchBlog(searchBlog *Domain.Blog) ([]Domain.Blog, error) {
	// Use filter to search for tasks with the given fields
	blogs := []Domain.Blog{}
	filters := []bson.M{}
	if searchBlog.Title != "" {
		filters = append(filters, bson.M{"Title": searchBlog.Title})
	}
	if searchBlog.Owner_email != "" {
		filters = append(filters, bson.M{"Ogowner": searchBlog.Owner_email})
	}
	filter := bson.M{
		"$and": filters,
	}
	cursor, err := BlgRepo.BlogCollection.Find(context.TODO(), filter)
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
	updatedBSON["Likes"] = updatedBlog.Likes
	updatedBSON["Dislikes"] = updatedBlog.Dislikes
	updatedBSON["ViewCount"] = updatedBlog.ViewCount
	updatedBSON["Comments"] = updatedBlog.Comments
	update := bson.M{"$set": updatedBSON}
	// Do update operation in database
	updatedRes, err := BlgRepo.BlogCollection.UpdateOne(context.TODO(), filter, update)
	// Handle exceptions
	if err != nil {
		return err
	}
	if updatedRes.MatchedCount == 0 {
		return errors.New("blog not found")
	}
	return nil
}

func (BlgRepo *BlogRepository) GetAllBlogs(limit int, offset int) ([]Domain.Blog, error) {
	findOptions := options.Find()

	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	result, err := BlgRepo.BlogCollection.Find(context.TODO(), bson.D{}, findOptions)

	if err != nil {
		return nil, err
	}

	var blogs []Domain.Blog

	for result.Next(context.TODO()) {
		var blog Domain.Blog
		if err := result.Decode(&blog); err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	log.Print(blogs)

	return blogs, nil
}

func (BlgRepo *BlogRepository) DeleteBlog(ID string) error {
	filter := bson.D{{Key: "ID", Value: bson.D{{Key: "$eq", Value: ID}}}}
	result, err := BlgRepo.BlogCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("blog not found")
	}
	return nil
}

func (BlgRepo *BlogRepository) FilterBlog(filterBlog *Domain.Blog) ([]Domain.Blog, error) {
	blogs := []Domain.Blog{}
	filters := []bson.D{}

	if !filterBlog.Date.IsZero() {
		filters = append(filters, bson.D{{Key: "Date", Value: bson.D{{Key: "$eq", Value: filterBlog.Date}}}})
	}
	if len(filterBlog.Tags) > 0 {
		filters = append(filters, bson.D{{Key: "Tags", Value: bson.D{{Key: "$in", Value: filterBlog.Tags}}}})
	}

	if len(filters) == 0 {
		return nil, errors.New("at least one filter (date or tags) must be provided")
	}

	filter := bson.M{"$or": filters}
	cursor, err := BlgRepo.BlogCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find blogs: %w", err)
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var blog Domain.Blog
		if err := cursor.Decode(&blog); err != nil {
			return nil, fmt.Errorf("failed to decode blog: %w", err)
		}
		blogs = append(blogs, blog)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return blogs, nil
}

func (BlgRepo *BlogRepository) GetBlog(id string) (Domain.Blog, error) {
	var blog Domain.Blog
	filter := bson.D{{Key: "ID", Value: id}}
	err := BlgRepo.BlogCollection.FindOne(context.TODO(), filter).Decode(&blog)
	if err != nil {
		return blog, errors.New("Document with id " + id + " not found")
	}
	return blog, nil
}
