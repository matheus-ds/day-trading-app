package matching

import (
	"day-trading-app/backend/internal/service/models"
	"fmt"

	"github.com/ryszard/goskiplist/skiplist"
)

type Element models.StockTransaction

// Implement the interface used in skiplist
func (e Element) ExtractKey() float64 {
	return e.StockPrice
}
func (e Element) String() string {
	return fmt.Sprintf("%03d", e.StockPrice)
}

type orderbook struct {
	buys  *skiplist.SkipList
	sells *skiplist.SkipList
}
type orderbooks struct {
	book map[string]*orderbook
}

var bookMap *orderbooks = new(orderbooks)
var stockTxCommitQueue []models.StockTransaction

func (books orderbooks) Match(tx models.StockTransaction) {
	var book = bookMap.book[tx.StockID]

	if books.book[tx.StockID] == nil {
		books.book[tx.StockID] = new(orderbook)
		// todo Use NewCustomMap (with our own LessThan functions as argument) to sort by both price (1st) and time (2nd).
		// todo Would need two functions with different time ordering for buys and sells.
		books.book[tx.StockID].buys = skiplist.NewIntMap()
		books.book[tx.StockID].sells = skiplist.NewIntMap()
	}

	if tx.IsBuy {
		book.matchBuy(tx)
	} else {
		book.matchSell(tx)
	}
}

func (book orderbook) matchBuy(buyTx models.StockTransaction) {
	if book.sells.Len() == 0 {
		book.buys.Set(buyTx.StockPrice, buyTx)
	} else {
		var sellsHasNext = true
		var sellIter = book.sells.Iterator()
		var buyQuantityRemaining = buyTx.Quantity

		for buyQuantityRemaining > 0 && sellsHasNext { // todo &&
			lowestSellTx := sellIter.Value().(models.StockTransaction)
			if buyTx.Quantity >= lowestSellTx.Quantity { // todo: && !(isLimit && expired)
				// todo create and commit fulfilled child transaction

				book.sells.Delete(lowestSellTx)
				lowestSellTx.OrderStatus = "COMPLETED"
				stockTxCommitQueue = append(stockTxCommitQueue, lowestSellTx)
			}

			sellsHasNext = sellIter.Next()
		}

		// todo if buy amount remaining, change parent transaction to partial_fulfilled, and whatever else
	}

}

func (book orderbook) matchSell(tx models.StockTransaction) {
	//todo: mirror matchBuy
}

// TODO: cancel orders
func (book orderbook) cancelBuyOrder(tx models.StockTransaction) {
	if tx.OrderType != "LIMIT" {
		// todo: report failure/error message to user
	}

	var _, wasFound = book.buys.Get(tx)
	if wasFound {
		// todo
	} else {
		// todo: report failure/error message to user
	}
}
