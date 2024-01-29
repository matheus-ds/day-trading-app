package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) InializeRoutes() {
	privateConfig := cors.DefaultConfig()
	privateConfig.AllowOrigins = []string{"http://localhost:3000"}
	privateConfig.AllowCredentials = true

	s.Router.Use(cors.New(privateConfig)) // for the web endpoints, we have stricter cors policy
	s.Router.Use(gin.Recovery())

	authorized := s.Router.Group("/")
	authorized.POST("/register", s.Srv.Register)
	authorized.POST("/login", s.Srv.AuthenticateUser)
}
