package server

import (
	"os"

	"day-trading-app/backend/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"day-trading-app/backend/internal/service"
	"day-trading-app/backend/internal/service/store"
	"day-trading-app/backend/internal/service/transport"
	"day-trading-app/backend/pkg/logger"
)

type Server struct {
	cfg    *config.Config
	Router *gin.Engine
	Srv    *transport.HTTPEndpoint
	DB     service.Database
}

func (s *Server) Initialize() {
	s.DB = store.GetMongoHandler()
	s.Srv = transport.NewHTTPTransport(service.New(s.DB))

	s.Router = gin.Default()

	s.InializeRoutes()
}

func (s *Server) Run(port string) {
	err := s.Router.Run(port)
	if err != nil {
		logger.Error("failed to run server", logger.ErrorType(err))
	}
}

func Run() {
	httpServer := Server{}

	err := godotenv.Load()
	if err != nil {
		logger.Debug("No '.env' file found, using global env vars")
	}
	port := ":" + os.Getenv("GIN_PORT")
	httpServer.Initialize()
	httpServer.Run(port)
}
