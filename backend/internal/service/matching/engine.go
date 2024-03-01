// The Matching Engine pairs buy and sell orders to create mutually agreeable transactions,
// which are placed in a queue to be handled by the Order Execution component.
// Unmatched orders remain in the Matching Engines memory until matched, expired, or canceled.

package matching

import (
	"day-trading-app/backend/internal/service/models"
	"github.com/google/uuid"
	"strings"
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
var stockTxCommitQueue []models.StockMatch

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
func isExpired(tx models.StockMatch) bool {
	return time.Now().Unix()+(15*60) >= tx.Order.TimeStamp
}

func getOrderbook(tx models.StockTransaction) *orderbook {
	if bookMap.book[tx.StockID] == nil {
		bookMap.book[tx.StockID] = new(orderbook)
		bookMap.book[tx.StockID].buys = skiplist.NewCustomMap(buyIsLowerPriorityThan)
		bookMap.book[tx.StockID].sells = skiplist.NewCustomMap(sellIsLowerPriorityThan)
	}
	return bookMap.book[tx.StockID]
}

// Create child transaction based on parentTx, ordering some specified quantity of otherTx.
func createChildTx(parentTx models.StockMatch, otherTx models.StockMatch, quantityTraded int) models.StockMatch {
	var childTx = parentTx
	childTx.Order.ParentStockTxID = &parentTx.Order.StockTxID
	childTx.Order.StockTxID = strings.ToLower(childTx.Order.StockID) + "StockTxId" + uuid.New().String() // todo fix?
	childTx.PriceTx = otherTx.Order.StockPrice                                                           // TODO: SPECS DIDN'T SPECIFY WHAT TO DO IN CASE OF PRICE DIFFERENCE
	childTx.QuantityTx = quantityTraded
	childTx.Order.TimeStamp = time.Now().Unix() // todo: Should we be changing the timestamp, or keep the parent's?
	childTx.Order.OrderStatus = "COMPLETED"
	return childTx
}

// Match inserts a transaction into the matching engine, to be matched with complementary transaction(s) in its order book.
func Match(order models.StockTransaction) {
	var book = getOrderbook(order)

	var tx = models.StockMatch{Order: order, QuantityTx: 0, PriceTx: 0}

	if order.IsBuy {
		book.matchBuy(tx)
	} else {
		book.matchSell(tx)
	}

	// ExecuteOrders(stockTxCommitQueue) todo: Not sure how to best pass queue and execution flow
}

// matchBuy() and matchSell() are basically mirrors of each other, with "buy" and "sell" swapped.

func (book orderbook) matchBuy(buyTx models.StockMatch) {
	if book.sells.Len() == 0 {
		if buyTx.Order.OrderType == "LIMIT" {
			book.buys.Set(buyTx.Order, buyTx)
		} else { // "MARKET"
			stockTxCommitQueue = append(stockTxCommitQueue, buyTx)
		}

	} else {
		var sellsHasNext = true
		var sellIter = book.sells.Iterator()
		var buyQuantityRemaining = buyTx.Order.Quantity

		for buyQuantityRemaining > 0 && sellsHasNext {
			lowestSellTx := sellIter.Value().(models.StockMatch)
			sellQuantityRemaining := lowestSellTx.Order.Quantity - lowestSellTx.QuantityTx

			if isExpired(lowestSellTx) {
				book.sells.Delete(lowestSellTx)
				stockTxCommitQueue = append(stockTxCommitQueue, lowestSellTx)
			} else {
				if (buyTx.Order.OrderType == "LIMIT") && (buyTx.Order.StockPrice < lowestSellTx.Order.StockPrice) {
					break
				}

				if buyQuantityRemaining >= sellQuantityRemaining {
					if buyTx.Order.Quantity == sellQuantityRemaining { // perfect match, no children
						buyTx.PriceTx = lowestSellTx.Order.StockPrice // TODO: SPECS DIDN'T SPECIFY WHAT TO DO IN CASE OF PRICE DIFFERENCE
					} else {
						var buyChildTx = createChildTx(buyTx, lowestSellTx, sellQuantityRemaining)
						stockTxCommitQueue = append(stockTxCommitQueue, buyChildTx)
					}

					book.sells.Delete(lowestSellTx)
					lowestSellTx.Order.OrderStatus = "COMPLETED"
					lowestSellTx.PriceTx = lowestSellTx.Order.StockPrice // TODO: SPECS DIDN'T SPECIFY WHAT TO DO IN CASE OF PRICE DIFFERENCE
					lowestSellTx.QuantityTx = lowestSellTx.Order.Quantity - sellQuantityRemaining
					stockTxCommitQueue = append(stockTxCommitQueue, lowestSellTx)
				} else { // buyQuantityRemaining < sellQuantityRemaining
					var buyChildTx = createChildTx(buyTx, lowestSellTx, buyQuantityRemaining)
					stockTxCommitQueue = append(stockTxCommitQueue, buyChildTx)

					lowestSellTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
					var sellChildTx = createChildTx(lowestSellTx, buyTx, buyQuantityRemaining)
					stockTxCommitQueue = append(stockTxCommitQueue, sellChildTx)
				}
			}
			buyQuantityRemaining -= sellQuantityRemaining
			sellsHasNext = sellIter.Next()
		}

		buyTx.QuantityTx = buyTx.Order.Quantity - buyQuantityRemaining

		if buyQuantityRemaining > 0 {
			buyTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
		} else { // = 0
			buyTx.Order.OrderStatus = "COMPLETED"
		}

		if buyTx.Order.OrderType == "LIMIT" {
			book.buys.Set(buyTx.Order, buyTx)
		}

		stockTxCommitQueue = append(stockTxCommitQueue, buyTx)
	}
}

func (book orderbook) matchSell(sellTx models.StockMatch) {
	if book.buys.Len() == 0 {
		if sellTx.Order.OrderType == "LIMIT" {
			book.sells.Set(sellTx.Order, sellTx)
		} else { // "MARKET"
			stockTxCommitQueue = append(stockTxCommitQueue, sellTx)
		}

	} else {
		var buysHasNext = true
		var buyIter = book.buys.Iterator()
		var sellQuantityRemaining = sellTx.Order.Quantity

		for sellQuantityRemaining > 0 && buysHasNext {
			highestBuyTx := buyIter.Value().(models.StockMatch)
			buyQuantityRemaining := highestBuyTx.Order.Quantity - highestBuyTx.QuantityTx

			if isExpired(highestBuyTx) {
				book.buys.Delete(highestBuyTx)
				stockTxCommitQueue = append(stockTxCommitQueue, highestBuyTx)
			} else {
				if (sellTx.Order.OrderType == "LIMIT") && (sellTx.Order.StockPrice < highestBuyTx.Order.StockPrice) {
					break
				}

				if sellQuantityRemaining >= buyQuantityRemaining {
					if sellTx.Order.Quantity == buyQuantityRemaining { // perfect match, no children
						sellTx.PriceTx = highestBuyTx.Order.StockPrice // TODO: SPECS DIDN'T SPECIFY WHAT TO DO IN CASE OF PRICE DIFFERENCE
					} else {
						var sellChildTx = createChildTx(sellTx, highestBuyTx, buyQuantityRemaining)
						stockTxCommitQueue = append(stockTxCommitQueue, sellChildTx)
					}

					book.buys.Delete(highestBuyTx)
					highestBuyTx.Order.OrderStatus = "COMPLETED"
					highestBuyTx.PriceTx = highestBuyTx.Order.StockPrice // TODO: SPECS DIDN'T SPECIFY WHAT TO DO IN CASE OF PRICE DIFFERENCE
					highestBuyTx.QuantityTx = highestBuyTx.Order.Quantity - buyQuantityRemaining
					stockTxCommitQueue = append(stockTxCommitQueue, highestBuyTx)
				} else { // sellQuantityRemaining < buyQuantityRemaining
					var sellChildTx = createChildTx(sellTx, highestBuyTx, sellQuantityRemaining)
					stockTxCommitQueue = append(stockTxCommitQueue, sellChildTx)

					highestBuyTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
					var buyChildTx = createChildTx(highestBuyTx, sellTx, sellQuantityRemaining)
					stockTxCommitQueue = append(stockTxCommitQueue, buyChildTx)
				}
			}
			sellQuantityRemaining -= buyQuantityRemaining
			buysHasNext = buyIter.Next()
		}

		sellTx.QuantityTx = sellTx.Order.Quantity - sellQuantityRemaining

		if sellQuantityRemaining > 0 {
			sellTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
		} else { // = 0
			sellTx.Order.OrderStatus = "COMPLETED"
		}

		if sellTx.Order.OrderType == "LIMIT" {
			book.sells.Set(sellTx.Order, sellTx)
		}

		stockTxCommitQueue = append(stockTxCommitQueue, sellTx)
	}
}

// CancelOrder halts further activity for a limit transaction with the given stockTxID.
// If found, the matching transaction is enqueued. Basically a deliberate premature expiration.
func CancelOrder(tx models.StockMatch) (wasCancelled bool) {
	var book = bookMap.book[tx.Order.StockID]
	if book != nil { // todo rewrite normally
		if tx.Order.IsBuy {
			wasCancelled = book.cancelBuyOrder(tx)
		} else {
			wasCancelled = book.cancelSellOrder(tx)
		}
	} else {
		stockTxCommitQueue = append(stockTxCommitQueue, tx)
	}

	// ExecuteOrders(stockTxCommitQueue) todo: Not sure how to best pass queue and execution flow

	return wasCancelled
}

// cancelBuyOrder() and cancelSellOrder() are mirrors of each other, with "buy" and "sell" swapped.

func (book orderbook) cancelBuyOrder(tx models.StockMatch) (wasFound bool) {
	victimTx, wasFound := book.buys.Get(tx)
	if wasFound {
		book.buys.Delete(tx)
		victimTx := victimTx.(models.StockMatch)
		stockTxCommitQueue = append(stockTxCommitQueue, victimTx)
	}
	return wasFound
}

func (book orderbook) cancelSellOrder(tx models.StockMatch) (wasFound bool) {
	victimTx, wasFound := book.sells.Get(tx)
	if wasFound {
		book.sells.Delete(tx)
		victimTx := victimTx.(models.StockMatch)
		stockTxCommitQueue = append(stockTxCommitQueue, victimTx)
	}
	return wasFound
}
