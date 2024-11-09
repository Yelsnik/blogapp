package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// model for posts
type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Body      string             `json:"body" bson:"body"`
	User      primitive.ObjectID `json:"user" bson:"user"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,default"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,default"`
}

func (p *Post) MarshalBSONPost() ([]byte, error) {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}

	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now()
	}

	type my Post
	return bson.Marshal((*my)(p))
}

// model for images
type Image struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileName   string             `json:"file_name" bson:"file_name"`
	Post       primitive.ObjectID `json:"post" bson:"post"`
	UploadedAt time.Time          `json:"uploaded_at" bson:"uploaded_at,default"`
}

func (i *Image) MarshalBSONImage() ([]byte, error) {
	if i.UploadedAt.IsZero() {
		i.UploadedAt = time.Now()
	}

	type my Image
	return bson.Marshal((*my)(i))
}

// model for users
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" binding:"required"`
	Email     string             `json:"email" bson:"email" binding:"required,email"`
	Password  string             `json:"password" bson:"password" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,default"`
}

func (u *User) MarshalBSONUser() ([]byte, error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	type my User
	return bson.Marshal((*my)(u))
}
