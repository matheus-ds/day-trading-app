package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateStockReq struct {
	StockName string `json:"stock_name"`
}

func (e HTTPEndpoint) CreateStock(c *gin.Context) {
	var stock CreateStockReq

	if err := c.BindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": "invalid request",
			},
		})
		return
	}

	createdStock, err := e.srv.CreateStock(stock.StockName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    createdStock,
	})
}

type AddStockToUserReq struct {
	StockID  string `json:"stock_id"`
	Quantity int    `json:"quantity"`
}

func (e HTTPEndpoint) AddStockToUser(c *gin.Context) {
	userName := c.GetString("user_name")
	var stock AddStockToUserReq

	if err := c.BindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": "invalid request",
			},
		})
		return
	}

	err := e.srv.AddStockToUser(userName, stock.StockID, stock.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    nil,
	})
}

func (e HTTPEndpoint) GetStockPortfolio(c *gin.Context) {
	userName := c.GetString("user_name")

	portfolio, err := e.srv.GetStockPortfolio(userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    portfolio,
	})
}

func (e HTTPEndpoint) GetStockTransactions(c *gin.Context) {

	transactions, err := e.srv.GetStockTransactions()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transactions,
	})
}

func (e HTTPEndpoint) GetStockPrices(c *gin.Context) {

	prices, err := e.srv.GetStockPrices()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prices,
	})
}

type PlaceStockReq struct {
	StockID   string `json:"stock_id"`
	IsBuy     bool   `json:"is_buy"`
	OrderType string `json:"order_type"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}

func (e HTTPEndpoint) PlaceStockOrder(c *gin.Context) {
	userName := c.GetString("user_name")
	var order PlaceStockReq

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	err := e.srv.PlaceStockOrder(userName, order.StockID, order.IsBuy, order.OrderType, order.Quantity, order.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    nil,
	})
}

type cancelTransactionReq struct {
	StockTxID string `json:"stock_tx_id"`
}

func (e HTTPEndpoint) CancelStockTransaction(c *gin.Context) {
	userName := c.GetString("user_name")
	var cancelReq cancelTransactionReq

	if err := c.BindJSON(&cancelReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": "invalid request",
			},
		})
		return
	}

	err := e.srv.CancelStockTransaction(userName, cancelReq.StockTxID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    nil,
	})
}
