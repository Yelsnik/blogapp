package api

import (
	"fmt"

	"github.com/Yelsnik/blogapp/token"
	"github.com/Yelsnik/blogapp/util"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func (s *Server) setUpRouter() {
	router := gin.Default()

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	router.POST("/v1/sign-up", s.SignUp)
	router.GET("/v1/user/:id", s.GetUser)
	router.POST("/v1/login", s.LoginUser)

	authRoutes.POST("v1/create-post", s.CreatePost)

	s.router = router
}

func NewServer(config util.Config) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		tokenMaker: tokenMaker,
		config:     config,
	}

	server.setUpRouter()
	return server, nil
}

func (s *Server) StartServer(address string) error {
	return s.router.Run(address)
}
