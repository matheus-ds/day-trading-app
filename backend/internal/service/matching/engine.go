// The matching engine pairs buy and sell orders to create mutually agreeable transactions.

package matching

import (
	"day-trading-app/backend/internal/service/models"

	"github.com/ryszard/goskiplist/skiplist"
)

// TODO: Optimize for space. Currently stores whole transactions, as both keys and values. Smaller keys is easy.

type orderbook struct {
	buys  *skiplist.SkipList
	sells *skiplist.SkipList
}
type orderbooks struct {
	book map[string]*orderbook
}

var bookMap *orderbooks = new(orderbooks)
var stockTxCommitQueue []models.StockTransaction

// Define orderings for orderbooks. Sort first by price, then if equal price, by time.

// Prioritize buys by highest price 1st, then earliest time 2nd.
func BuyIsLowerPriorityThan(l, r interface{}) bool {
	ll := l.(models.StockTransaction)
	rr := r.(models.StockTransaction)
	if ll.StockPrice > rr.StockPrice {
		return true
	} else if ll.StockPrice == rr.StockPrice {
		return ll.TimeStamp < rr.TimeStamp
	} else {
		return false
	}
}

// Prioritize sells by lowest price 1st, then earliest time 2nd.
func SellIsLowerPriorityThan(l, r interface{}) bool {
	ll := l.(models.StockTransaction)
	rr := r.(models.StockTransaction)
	if ll.StockPrice < rr.StockPrice {
		return true
	} else if ll.StockPrice == rr.StockPrice {
		return ll.TimeStamp < rr.TimeStamp
	} else {
		return false
	}
}

func (books orderbooks) Match(tx models.StockTransaction) {
	var book = bookMap.book[tx.StockID]

	if books.book[tx.StockID] == nil {
		books.book[tx.StockID] = new(orderbook)
		books.book[tx.StockID].buys = skiplist.NewCustomMap(BuyIsLowerPriorityThan)
		books.book[tx.StockID].sells = skiplist.NewCustomMap(SellIsLowerPriorityThan)
	}

	if tx.IsBuy {
		book.MatchBuy(tx)
	} else {
		book.MatchSell(tx)
	}
}

func (book orderbook) MatchBuy(buyTx models.StockTransaction) {
	if book.sells.Len() == 0 {
		book.buys.Set(buyTx, buyTx)
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

func (book orderbook) MatchSell(tx models.StockTransaction) {
	//todo: mirror MatchBuy
}

func (book orderbook) CancelBuyOrder(tx models.StockTransaction) {
	if tx.OrderType != "LIMIT" {
		// todo: report failure/error message to user
	}

	var victimTx, wasFound = book.buys.Get(tx)
	if wasFound {
		book.buys.Delete(tx)
		victimTx := victimTx.(models.StockTransaction)
		//if !(victimTx.OrderStatus == "PARTIAL_FULFILLED") {} // todo: do we log cancelled orders that did nothing?
		stockTxCommitQueue = append(stockTxCommitQueue, victimTx)
	} else {
		// todo: report failure/error message to user
	}
}

func (book orderbook) CancelSellOrder(tx models.StockTransaction) {
	if tx.OrderType != "LIMIT" {
		// todo: report failure/error message to user
	}

	var victimTx, wasFound = book.sells.Get(tx)
	if wasFound {
		book.sells.Delete(tx)
		victimTx := victimTx.(models.StockTransaction)
		//if !(victimTx.OrderStatus == "PARTIAL_FULFILLED") {} // todo: do we log cancelled orders that did nothing?
		stockTxCommitQueue = append(stockTxCommitQueue, victimTx)
	} else {
		// todo: report failure/error message to user
	}
}
