package main

import (
	"fmt"
	"log"

	//"github.com/Yelsnik/blogapp/api"
	"github.com/Yelsnik/blogapp/api"
	db "github.com/Yelsnik/blogapp/db/service"
	"github.com/Yelsnik/blogapp/util"
	//"github.com/gin-gonic/gin"
)

func init() {
	if err := db.ConnectToDatabase(); err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	} else {
		fmt.Println("Connection established")
	}

}

func main() {
	// load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// create new server
	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// start the server
	err = server.StartServer(config.Port)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
