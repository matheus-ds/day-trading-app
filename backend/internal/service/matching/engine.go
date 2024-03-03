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
// Also adds transaction cost to parentTx.CostTotalTx.
func createChildTx(parentTx *models.StockMatch, quantityTraded int, priceTraded int) models.StockMatch {
	var childTx = *parentTx
	childTx.Order.ParentStockTxID = &parentTx.Order.StockTxID
	childTx.Order.StockTxID = strings.ToLower(childTx.Order.StockID) + "StockTxId" + uuid.New().String() // todo fix?
	childTx.Order.StockPrice = priceTraded
	childTx.PriceTx = priceTraded
	childTx.Order.Quantity = quantityTraded
	childTx.QuantityTx = quantityTraded
	childTx.CostTotalTx = quantityTraded * priceTraded
	childTx.Order.TimeStamp = time.Now().Unix()
	childTx.Order.OrderStatus = "COMPLETED"

	parentTx.CostTotalTx += quantityTraded * priceTraded
	return childTx
}

// Match inserts a transaction into the matching engine, to be matched with complementary transaction(s) in its order book.
func Match(order models.StockTransaction) {
	var book = getOrderbook(order)

	var tx = models.StockMatch{Order: order, QuantityTx: 0, PriceTx: 0, CostTotalTx: 0, Killed: false}

	if order.IsBuy {
		book.matchBuy(tx)
	} else {
		book.matchSell(tx)
	}

	ExecuteOrders(stockTxCommitQueue)
}

// matchBuy() and matchSell() are basically mirrors of each other, with "buy" and "sell" swapped.
// todo Specification unclear as to what price that two matched price limit-orders with different prices should trade at.
// todo   So for now we're taking the oldest limit-order's price.

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
				book.sells.Delete(lowestSellTx.Order)
				stockTxCommitQueue = append(stockTxCommitQueue, lowestSellTx)
			} else {
				if (buyTx.Order.OrderType == "LIMIT") && (buyTx.Order.StockPrice < lowestSellTx.Order.StockPrice) {
					break
				}

				if buyQuantityRemaining >= sellQuantityRemaining {
					if buyTx.Order.Quantity == sellQuantityRemaining { // perfect match, no children
						buyTx.PriceTx = lowestSellTx.Order.StockPrice
						buyTx.CostTotalTx += buyTx.Order.Quantity * lowestSellTx.Order.StockPrice
					} else {
						var buyChildTx = createChildTx(&buyTx, sellQuantityRemaining, lowestSellTx.Order.StockPrice)
						stockTxCommitQueue = append(stockTxCommitQueue, buyChildTx)
					}

					book.sells.Delete(lowestSellTx.Order)
					if lowestSellTx.Order.OrderStatus == "IN_PROGRESS" { // not parent
						lowestSellTx.PriceTx = lowestSellTx.Order.StockPrice
					}
					lowestSellTx.QuantityTx = lowestSellTx.Order.Quantity - sellQuantityRemaining
					lowestSellTx.CostTotalTx += sellQuantityRemaining * lowestSellTx.Order.StockPrice
					lowestSellTx.Order.OrderStatus = "COMPLETED"
					stockTxCommitQueue = append(stockTxCommitQueue, lowestSellTx)
				} else { // buyQuantityRemaining < sellQuantityRemaining
					var buyChildTx = createChildTx(&buyTx, buyQuantityRemaining, lowestSellTx.Order.StockPrice)
					stockTxCommitQueue = append(stockTxCommitQueue, buyChildTx)

					lowestSellTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
					var sellChildTx = createChildTx(&lowestSellTx, buyQuantityRemaining, lowestSellTx.Order.StockPrice)
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
				book.buys.Delete(highestBuyTx.Order)
				stockTxCommitQueue = append(stockTxCommitQueue, highestBuyTx)
			} else {
				if (sellTx.Order.OrderType == "LIMIT") && (sellTx.Order.StockPrice < highestBuyTx.Order.StockPrice) {
					break
				}

				if sellQuantityRemaining >= buyQuantityRemaining {
					if sellTx.Order.Quantity == buyQuantityRemaining { // perfect match, no children
						sellTx.PriceTx = highestBuyTx.Order.StockPrice
						sellTx.CostTotalTx += sellTx.Order.Quantity * highestBuyTx.Order.StockPrice
					} else {
						var sellChildTx = createChildTx(&sellTx, buyQuantityRemaining, highestBuyTx.Order.StockPrice)
						stockTxCommitQueue = append(stockTxCommitQueue, sellChildTx)
					}

					book.buys.Delete(highestBuyTx.Order)
					if highestBuyTx.Order.OrderStatus == "IN_PROGRESS" { // not parent
						highestBuyTx.PriceTx = highestBuyTx.Order.StockPrice
					}
					highestBuyTx.QuantityTx = highestBuyTx.Order.Quantity - buyQuantityRemaining
					highestBuyTx.CostTotalTx += buyQuantityRemaining * highestBuyTx.Order.StockPrice
					highestBuyTx.Order.OrderStatus = "COMPLETED"
					stockTxCommitQueue = append(stockTxCommitQueue, highestBuyTx)
				} else { // sellQuantityRemaining < buyQuantityRemaining
					var sellChildTx = createChildTx(&sellTx, sellQuantityRemaining, highestBuyTx.Order.StockPrice)
					stockTxCommitQueue = append(stockTxCommitQueue, sellChildTx)

					highestBuyTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
					var buyChildTx = createChildTx(&highestBuyTx, sellQuantityRemaining, highestBuyTx.Order.StockPrice)
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
func CancelOrder(order models.StockTransaction) (wasCancelled bool) {
	var book = bookMap.book[order.StockID]
	if book != nil {
		if order.IsBuy {
			wasCancelled = book.cancelBuyOrder(order)
		} else {
			wasCancelled = book.cancelSellOrder(order)
		}
	}

	ExecuteOrders(stockTxCommitQueue)

	return wasCancelled
}

// cancelBuyOrder() and cancelSellOrder() are mirrors of each other, with "buy" and "sell" swapped.

func (book orderbook) cancelBuyOrder(order models.StockTransaction) (wasFound bool) {
	victimTx, wasFound := book.buys.Get(order)
	if wasFound {
		book.buys.Delete(order)
		victimTx := victimTx.(models.StockMatch)
		victimTx.Killed = true
		stockTxCommitQueue = append(stockTxCommitQueue, victimTx)
	}
	return wasFound
}

func (book orderbook) cancelSellOrder(order models.StockTransaction) (wasFound bool) {
	victimTx, wasFound := book.sells.Get(order)
	if wasFound {
		book.sells.Delete(order)
		victimTx := victimTx.(models.StockMatch)
		victimTx.Killed = true
		stockTxCommitQueue = append(stockTxCommitQueue, victimTx)
	}
	return wasFound
}
