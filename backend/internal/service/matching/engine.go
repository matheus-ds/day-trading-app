// The Matching Engine pairs buy and sell orders to create mutually agreeable transactions,
// which are placed in a queue to be handled by the Order Execution component.
// Unmatched orders remain in the Matching Engines memory until matched, expired, or canceled.

package matching

import (
	"day-trading-app/backend/internal/service/models"
	"time"

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

// Returns true if transaction's timestamp is over 15 minutes old.
func isExpired(tx models.StockTransaction) bool {
	if tx.OrderType == "LIMIT" { // todo: this check is probably unnecessary
		return time.Now().Unix()+(15*60) >= tx.TimeStamp
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

// Create child transaction based on parentTx, ordering full quantity of otherTx.
func createChildTx(parentTx models.StockTransaction, otherTx models.StockTransaction, quantityTraded int) models.StockTransaction {
	var childTx = parentTx
	childTx.ParentStockTxID = &parentTx.StockTxID
	// childTx.StockTxID = TODO child StockTxID scheme? Should this be decided by Order Execution?
	childTx.StockPrice = otherTx.StockPrice
	childTx.Quantity = quantityTraded
	// childTx.TimeStamp = TODO do we change timestamp?

	childTx.OrderStatus = "COMPLETED"
	return childTx
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
		if buyTx.OrderType == "LIMIT" {
			book.buys.Set(buyTx, buyTx)
		} else { // "MARKET"
			buyTx.OrderStatus = "UNFULFILLED" // todo: should this be handled differently?
			stockTxCommitQueue = append(stockTxCommitQueue, buyTx)
		}

	} else {
		var sellsHasNext = true
		var sellIter = book.sells.Iterator()
		var buyQuantityRemaining = buyTx.Quantity
		var childTxCount = 0 // todo: do we need this?

		for buyQuantityRemaining > 0 && sellsHasNext { // todo &&
			lowestSellTx := sellIter.Value().(models.StockTransaction)
			if isExpired(lowestSellTx) {
				book.sells.Delete(lowestSellTx)
				lowestSellTx.OrderStatus = "EXPIRED" // todo
				stockTxCommitQueue = append(stockTxCommitQueue, lowestSellTx)
			} else {
				if buyQuantityRemaining >= lowestSellTx.Quantity {
					if buyTx.Quantity == lowestSellTx.Quantity { // perfect match, no children
						buyTx.OrderStatus = "COMPLETED"
					} else {
						var childTx = createChildTx(buyTx, lowestSellTx, lowestSellTx.Quantity)
						stockTxCommitQueue = append(stockTxCommitQueue, childTx)
						childTxCount++
					}

					book.sells.Delete(lowestSellTx)
					lowestSellTx.OrderStatus = "COMPLETED"
					stockTxCommitQueue = append(stockTxCommitQueue, lowestSellTx)
				} else { // buyQuantityRemaining < lowestSellTx.Quantity
					var childTx = createChildTx(buyTx, lowestSellTx, buyQuantityRemaining)
					stockTxCommitQueue = append(stockTxCommitQueue, childTx)
					childTxCount++

					// todo handle sell LIMIT if IN_PROGRESS vs PARTIALLY_FULFILLED

				}
			}
			buyQuantityRemaining -= lowestSellTx.Quantity
			sellsHasNext = sellIter.Next()
		}

		if (buyQuantityRemaining > 0) && (buyTx.OrderType == "MARKET") {
			buyTx.OrderStatus = "PARTIALLY_FULFILLED" // todo: anything else?
		}
		stockTxCommitQueue = append(stockTxCommitQueue, buyTx)

		// todo else if LIMIT order, store in engine (with info about amount remaining to fulfill,
		//      or change Quantity field and let the Order Execution contrast with the already stored IN_PROGRESS tx in database)?
	}

}

func (book orderbook) matchSell(tx models.StockTransaction) {
	//todo: mirror matchBuy
}

func CancelOrder(tx models.StockTransaction) {
	if tx.OrderType != "LIMIT" {
		// todo: report failure/error message to user. Actually, this should probably be handled earlier than the matcher.
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
