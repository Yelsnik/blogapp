package api

import (
	//	"fmt"

	//"github.com/Yelsnik/blogapp/token"
	"net/http"
	"time"

	db "github.com/Yelsnik/blogapp/db/service"
	"github.com/Yelsnik/blogapp/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"
)

type UserResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type LoginUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func newResponse(data db.User) UserResponse {
	return UserResponse{
		ID:        data.ID,
		Name:      data.Name,
		Email:     data.Email,
		CreatedAt: data.CreatedAt,
	}
}

type LoginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}

func (server *Server) SignUp(ctx *gin.Context) {
	// unmarshal json
	var user db.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// hash password
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	// insert user
	id, err := db.CreateUser(db.MongoClient, ctx, &arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get the inserted result from db
	result, err := db.GetUserByID(db.MongoClient, ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// return a response
	response := newResponse(result)

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "data": response})
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get user by email
	user, err := db.GetUserByEmail(db.MongoClient, ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// compare password
	err = util.ComparePassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// create token
	accessToken, err := server.tokenMaker.CreateToken(user.Email, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := LoginUserResponse{
		AccessToken: accessToken,
		User:        newResponse(user),
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "data": response})
}

func (server *Server) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	result, err := db.GetUserByID(db.MongoClient, ctx, objectId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
