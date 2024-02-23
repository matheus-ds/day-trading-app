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

	authorized.POST("/createStock", s.Srv.CreateStock)
	authorized.POST("/addStockToUser", s.Srv.AddStockToUser)

	authorized.GET("/getStockPortfolio", s.Srv.GetStockPortfolio)
	authorized.GET("/getStockTransactions", s.Srv.GetStockTransactions)
	authorized.GET("/getStockPrices", s.Srv.GetStockPrices)
	authorized.POST("/placeStockOrder", s.Srv.PlaceStockOrder)
	authorized.POST("/cancelStockTransaction", s.Srv.CancelStockTransaction)

	// wallet endpoints
	authorized.POST("/addMoneyToWallet", s.Srv.AddMoneyToWallet)
	authorized.GET("/getWalletBalance", s.Srv.GetWalletBalance)
	authorized.GET("/getWalletTransactions", s.Srv.GetWalletTransactions)
}
