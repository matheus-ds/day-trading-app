// The Matching Engine pairs buy and sell orders to create mutually agreeable transactions,
// which are placed in a queue to be handled by the Order Execution component.
// Unmatched orders remain in the Matching Engines memory until matched, expired, or canceled.

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
func buyIsLowerPriorityThan(l, r interface{}) bool {
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
func sellIsLowerPriorityThan(l, r interface{}) bool {
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

func getOrderbook(tx models.StockTransaction) *orderbook {
	if bookMap.book[tx.StockID] == nil {
		bookMap.book[tx.StockID] = new(orderbook)
		bookMap.book[tx.StockID].buys = skiplist.NewCustomMap(buyIsLowerPriorityThan)
		bookMap.book[tx.StockID].sells = skiplist.NewCustomMap(sellIsLowerPriorityThan)
	}
	return bookMap.book[tx.StockID]
}

func Match(tx models.StockTransaction) {
	var book = getOrderbook(tx)

	if tx.IsBuy {
		book.matchBuy(tx)
	} else {
		book.matchSell(tx)
	}

	// todo ExecuteOrders(stockTxCommitQueue)
}

func (book orderbook) matchBuy(buyTx models.StockTransaction) {
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

func (book orderbook) matchSell(tx models.StockTransaction) {
	//todo: mirror matchBuy
}

func CancelOrder(tx models.StockTransaction) {
	if tx.OrderType != "LIMIT" {
		// todo: report failure/error message to user
	}

	var book = bookMap.book[tx.StockID]
	if book == nil {
		// todo: report failure/error message to user
	}

	if tx.IsBuy {
		book.cancelBuyOrder(tx)
	} else {
		book.cancelSellOrder(tx)
	}

	// todo ExecuteOrders(stockTxCommitQueue)
}

func (book orderbook) cancelBuyOrder(tx models.StockTransaction) {
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

func (book orderbook) cancelSellOrder(tx models.StockTransaction) {
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
