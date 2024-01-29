package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/matheus-ds/day-trading-app/backend/internal/service"
	"github.com/matheus-ds/day-trading-app/backend/internal/service/transport"
	"github.com/matheus-ds/day-trading-app/backend/pkg/logger"
)

type Server struct {
	Router *gin.Engine
	Srv    *transport.HTTPEndpoint
}

func (s *Server) Initialize() {
	s.Srv = transport.NewHTTPTransport(service.New())

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
