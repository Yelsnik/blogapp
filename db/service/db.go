package db

import (
	"context"
	//"fmt"
	"log"

	"github.com/Yelsnik/blogapp/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var db string

func ConnectToDatabase() error {
	ctx := context.TODO()
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opt := options.Client().ApplyURI(config.MongoUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	err = client.Ping(ctx, nil)
	MongoClient = client
	db = "blog-app"

	return err
}
