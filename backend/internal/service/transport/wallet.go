package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletAmount struct {
	Amount float32 `json:"amount"`
}

func (e HTTPEndpoint) AddMoneyToWallet(c *gin.Context) {
	userName := c.GetString("user_name")
	var wallet WalletAmount

	if err := c.BindJSON(&wallet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data": gin.H{
				"error": "invalid request",
			},
		})
		return
	}

	err := e.srv.AddMoneyToWallet(userName, wallet.Amount)
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

func (e HTTPEndpoint) GetWalletBalance(c *gin.Context) {
	userName := c.GetString("user_name")

	balance, err := e.srv.GetWalletBalance(userName)
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
		"data":    gin.H{"balance": balance},
	})
}

func (e HTTPEndpoint) GetWalletTransactions(c *gin.Context) {
	userName := c.GetString("user_name")

	transactions, err := e.srv.GetWalletTransactions(userName)
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
