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
	BlogCollection  *mongo.Collection
	LikesCollection *mongo.Collection
}

type LikeTrackerDTO struct {
	BlogID    string `bson:"id"`
	UserEmail string `bson:"email"`
	Liked     int    `bson:"liked"`
}

func NewBlogRepository(db *mongo.Database) *BlogRepository {
	return &BlogRepository{
		BlogCollection:  db.Collection("blogs"),
		LikesCollection: db.Collection("likes"),
	}
}

func (BlgRepo *BlogRepository) FindLiked(user_email, blog_id string) (*Domain.LikeTracker, error) {
	var tmp LikeTrackerDTO
	filter := bson.M{"id": blog_id, "email": user_email}
	err := BlgRepo.LikesCollection.FindOne(context.TODO(), filter).Decode(&tmp)
	return ChangeToDomain(&tmp), err
}

func (BlgRepo *BlogRepository) CreateLikeTk(lt Domain.LikeTracker) error {
	_, err := BlgRepo.LikesCollection.InsertOne(context.TODO(), ChangeToDTO(lt))
	return err
}

func (BlgRepo *BlogRepository) DeleteLikeTk(lt Domain.LikeTracker) error {
	_, err := BlgRepo.LikesCollection.DeleteMany(context.TODO(), bson.M{"email": lt.UserEmail, "id": lt.BlogID})
	return err
}

func (BlgRepo *BlogRepository) NumberOfLikes(id string) (int64, error) {
	filter := bson.M{"id": id, "liked": 1}
	return BlgRepo.LikesCollection.CountDocuments(context.TODO(), filter)
}

func (BlgRepo *BlogRepository) NumberOfDislikes(id string) (int64, error) {
	filter := bson.M{"id": id, "liked": -1}
	return BlgRepo.LikesCollection.CountDocuments(context.TODO(), filter)
}

func (BlgRepo *BlogRepository) Create(blog *Domain.Blog) error {
	_, err := BlgRepo.BlogCollection.InsertOne(context.TODO(), blog)
	return err
}

func (BlgRepo *BlogRepository) SearchBlog(searchBlog *Domain.Blog) ([]Domain.Blog, error) {
	var blogs []Domain.Blog
	filters := bson.M{}

	if searchBlog.Title != "" {
		filters["Title"] = searchBlog.Title
	}
	if searchBlog.Owner_email != "" {
		filters["Owner_email"] = searchBlog.Owner_email
	}

	// If no filters, return empty slice instead of querying everything
	if len(filters) == 0 {
		return []Domain.Blog{}, nil
	}

	cursor, err := BlgRepo.BlogCollection.Find(context.TODO(), filters)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var elem Domain.Blog
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}
		blogs = append(blogs, elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (BlgRepo *BlogRepository) UpdateBlog(updatedBlog *Domain.Blog) error {
	// Use blog ID to search and update task
	filter := bson.D{{Key: "id", Value: updatedBlog.ID}}
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
	// updatedBSON["Likes"] = updatedBlog.Likes
	// updatedBSON["Dislikes"] = updatedBlog.Dislikes
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
	filter := bson.D{{Key: "id", Value: id}}
	err := BlgRepo.BlogCollection.FindOne(context.TODO(), filter).Decode(&blog)
	if err != nil {
		return blog, errors.New("Document with id " + id + " not found")
	}
	return blog, nil
}

func (BlgRepo *BlogRepository) GetPopularBlogs() ([]Domain.Blog, error) {
	pipeline := mongo.Pipeline{
	bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "likes"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "BlogID"},
		{Key: "as", Value: "reactions"},
	}}},
	bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "likes", Value: bson.D{
			{Key: "$size", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$reactions"},
					{Key: "as", Value: "reaction"},
					{Key: "cond", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$$reaction.Liked", 1}},
					}},
				}},
			}},
		}},
		{Key: "dislikes", Value: bson.D{
			{Key: "$size", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$reactions"},
					{Key: "as", Value: "reaction"},
					{Key: "cond", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$$reaction.Liked", -1}},
					}},
				}},
			}},
		}},
	}}},
	bson.D{{Key: "$addFields", Value: bson.D{
		{Key: "score", Value: bson.D{
			{Key: "$add", Value: bson.A{
				bson.D{{Key: "$subtract", Value: bson.A{"$likes", "$dislikes"}}},
				"$ViewCount",
			}},
		}},
	}}},
	bson.D{{Key: "$sort", Value: bson.D{
		{Key: "score", Value: -1},
	}}},
}

	cursor, err := BlgRepo.BlogCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate error: %w", err)
	}
	defer cursor.Close(context.TODO())

	var blogs []Domain.Blog
	if err := cursor.All(context.TODO(), &blogs); err != nil {
		return nil, fmt.Errorf("cursor decode error: %w", err)
	}
	return blogs, nil
}

func ChangeToDTO(t Domain.LikeTracker) LikeTrackerDTO {
	return LikeTrackerDTO{
		BlogID:    t.BlogID,
		UserEmail: t.UserEmail,
		Liked:     t.Liked,
	}
}

func ChangeToDomain(t *LikeTrackerDTO) *Domain.LikeTracker {
	return &Domain.LikeTracker{
		BlogID:    t.BlogID,
		UserEmail: t.UserEmail,
		Liked:     t.Liked,
	}
}
