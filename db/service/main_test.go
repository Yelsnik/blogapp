package db

import (
	"fmt"
	"os"
	"testing"

	"context"
	//"fmt"
	"log"

	"github.com/Yelsnik/blogapp/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testClient *mongo.Client

func connectTestDB() (*mongo.Client, error) {
	// load config
	ctx := context.TODO()
	config, err := util.LoadConfig("../..")
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

	return client, err
}

func getTestUserCollection(client *mongo.Client) *mongo.Collection {
	blogApp := client.Database("blog-app")
	return blogApp.Collection("users")
}
func getTestPostCollection(client *mongo.Client) *mongo.Collection {
	blogApp := client.Database("blog-app")
	return blogApp.Collection("posts")
}

func TestMain(m *testing.M) {

	client, err := connectTestDB()
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	testClient = client
	getTestUserCollection(testClient)
	getTestPostCollection(testClient)
	fmt.Println("Connection established")

	os.Exit(m.Run())
}
