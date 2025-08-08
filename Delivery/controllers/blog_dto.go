package controllers

import "time"

// Types to use for binding (entities with Json Tags) and also bson format for storing
type BlogDTO struct {
	ID        string  
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Owner_email string      `json:"owner"`
	Tags      []string  `json:"tags"`
	Date      time.Time `json:"date"`
	ViewCount int       `json:"viewCount"`
	Comments  []string  `json:"comments"`
}