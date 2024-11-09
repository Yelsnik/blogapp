package api

import (
	//"fmt"

	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	db "github.com/Yelsnik/blogapp/db/service"
	"github.com/Yelsnik/blogapp/token"

	//"github.com/Yelsnik/blogapp/token"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type PostResponse struct {
	Post  db.Post    `json:"post"`
	Image []db.Image `json:"images"`
}

func NewResponse(post db.Post, image db.Image) PostResponse {
	var responseArray []db.Image

	result := append(responseArray, image)

	return PostResponse{
		Post:  post,
		Image: result,
	}
}

func (server *Server) CreatePost(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := db.GetUserByEmail(db.MongoClient, ctx, authPayload.Email)

	formFile, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	texts := make(map[string]string)

	for k, v := range formFile.Value {
		if len(v) > 0 {
			texts[k] = v[0]
		}
	}

	postArg := &db.Post{
		Title:     texts["Title"],
		Body:      texts["Body"],
		User:      user.ID,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	postID, err := db.CreatePost(db.MongoClient, ctx, postArg)
	fmt.Println(postID)

	post, err := db.GetPostByID(db.MongoClient, ctx, postID)
	fmt.Println(post)

	var imageID primitive.ObjectID

	for _, fileHeaders := range formFile.File {

		for _, fileheader := range fileHeaders {

			file, err := fileheader.Open()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
				return
			}
			defer file.Close()

			go func(f multipart.File, fileName string) {
				_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				bucket, err := gridfs.NewBucket(db.MongoClient.Database("blog-app"))
				if err != nil {
					log.Println("Failed to create GridFS bucket:", err)
					return
				}

				uploadStream, err := bucket.OpenUploadStream(fileName)
				if err != nil {
					log.Println("Failed to open upload stream:", err)
					return
				}
				defer uploadStream.Close()

				if _, err = uploadStream.Write([]byte(fileName)); err != nil {
					log.Println("Failed to write file to GridFS:", err)
				}
			}(file, fileheader.Filename)

			imageArg := &db.Image{
				FileName:   fileheader.Filename,
				Post:       postID,
				UploadedAt: time.Now(),
			}

			imageID, err = db.InsertImage(ctx, imageArg)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
	}

	image, err := db.GetImageByID(ctx, imageID)

	response := NewResponse(post, image)

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "response": response})

}
