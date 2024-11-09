package db

import (
	"context"

	//"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

func getUserCollection(client *mongo.Client) *mongo.Collection {
	blogApp := client.Database("blog-app")
	return blogApp.Collection("users")
}

func CreateUser(client *mongo.Client, ctx context.Context, args *User) (primitive.ObjectID, error) {

	c := getUserCollection(client)

	// create index on email
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},              // Create an index on the 'email' field
		Options: options.Index().SetUnique(true), // Ensure the index is unique
	}
	_, err := c.Indexes().CreateOne(ctx, indexModel)

	// insert
	user, err := c.InsertOne(ctx, args)

	// convert *mongo.InsertOneResult to primitive.ObjectID
	object := user.InsertedID
	id := object.(primitive.ObjectID)

	return id, err
}

func GetUserByEmail(client *mongo.Client, ctx context.Context, email string) (User, error) {
	var user User
	query := bson.M{"email": email}

	col := getUserCollection(client)

	err := col.FindOne(ctx, query).Decode(&user)

	return user, err
}

func GetUserByID(client *mongo.Client, ctx context.Context, id primitive.ObjectID) (User, error) {
	var user User
	query := bson.M{"_id": id}

	col := getUserCollection(client)

	err := col.FindOne(ctx, query).Decode(&user)

	return user, err
}
