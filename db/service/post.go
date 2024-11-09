package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func getPostCollection(client *mongo.Client) *mongo.Collection {
	blogApp := client.Database("blog-app")
	return blogApp.Collection("posts")
}

func getImageCollection(client *mongo.Client) *mongo.Collection {
	blogApp := client.Database("blog-app")
	return blogApp.Collection("images")
}

func CreatePost(client *mongo.Client, ctx context.Context, args *Post) (primitive.ObjectID, error) {
	col := getPostCollection(client)
	result, err := col.InsertOne(ctx, args)

	id := result.InsertedID.(primitive.ObjectID)

	return id, err
}

func InsertImage(ctx context.Context, args *Image) (primitive.ObjectID, error) {
	col := getImageCollection(MongoClient)
	result, err := col.InsertOne(ctx, args)

	id := result.InsertedID.(primitive.ObjectID)

	return id, err
}

func GetImageFromGridFS(bucket *gridfs.Bucket, fileName string) (*gridfs.DownloadStream, error) {
	db := MongoClient.Database("blog-app")

	bucket, err := gridfs.NewBucket(db)

	downloadStream, err := bucket.OpenDownloadStream(fileName)

	return downloadStream, err
}

func GetImageByID(ctx context.Context, id primitive.ObjectID) (Image, error) {
	var image Image
	col := getImageCollection(MongoClient)

	query := bson.M{"_id": id}

	err := col.FindOne(ctx, query).Decode(&image)

	return image, err
}

func GetPostByID(client *mongo.Client, ctx context.Context, id primitive.ObjectID) (Post, error) {
	var post Post
	col := getPostCollection(client)

	query := bson.M{"_id": id}

	err := col.FindOne(ctx, query).Decode(&post)

	return post, err
}
