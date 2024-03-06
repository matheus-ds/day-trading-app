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
type orderbooks map[string]orderbook

var bookMap = make(orderbooks)

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
	return time.Now().Unix()+(15*60) <= tx.Order.TimeStamp
}

func getOrderbook(tx models.StockTransaction) orderbook {
	if bookMap[tx.StockID].buys == nil {
		bookMap[tx.StockID] = orderbook{
			buys:  skiplist.NewCustomMap(buyIsLowerPriorityThan),
			sells: skiplist.NewCustomMap(sellIsLowerPriorityThan),
		}
	}
	return bookMap[tx.StockID]
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

	parentTx.IsParent = true
	parentTx.CostTotalTx += quantityTraded * priceTraded
	return childTx
}

// Match inserts a transaction into the matching engine, to be matched with complementary transaction(s) in its order book.
func Match(order models.StockTransaction) {
	var book = getOrderbook(order)

	var tx = models.StockMatch{Order: order, QuantityTx: 0, PriceTx: 0, CostTotalTx: 0, IsParent: false, Killed: false}

	var txCommitQueue []models.StockMatch
	if order.IsBuy {
		book.matchBuy(tx, &txCommitQueue)
	} else {
		book.matchSell(tx, &txCommitQueue)
	}

	ExecuteOrders(txCommitQueue)
}

// matchBuy() and matchSell() are basically mirrors of each other, with "buy" and "sell" swapped.
// todo Specification unclear as to what price that two matched price limit-orders with different prices should trade at.
// todo   So for now we're taking the oldest limit-order's price.

func (book orderbook) matchBuy(buyTx models.StockMatch, txCommitQueue *[]models.StockMatch) {
	if book.sells.Len() == 0 {
		if buyTx.Order.OrderType == "LIMIT" {
			book.buys.Set(buyTx.Order, buyTx)
		} else { // "MARKET"
			*txCommitQueue = append(*txCommitQueue, buyTx)
		}

	} else {
		var sellsHasNext = true
		var sellIter = book.sells.Iterator()
		sellIter.Next()
		var buyQuantityRemaining = buyTx.Order.Quantity

		for buyQuantityRemaining > 0 && sellsHasNext {
			lowestSellTx := sellIter.Value().(models.StockMatch)
			sellsHasNext = sellIter.Next()
			sellQuantityRemaining := lowestSellTx.Order.Quantity - lowestSellTx.QuantityTx

			if isExpired(lowestSellTx) {
				book.sells.Delete(lowestSellTx.Order)
				*txCommitQueue = append(*txCommitQueue, lowestSellTx)
			} else {
				if (buyTx.Order.OrderType == "LIMIT") && (buyTx.Order.StockPrice < lowestSellTx.Order.StockPrice) {
					break
				}

				if buyQuantityRemaining >= sellQuantityRemaining {
					if buyTx.Order.Quantity == sellQuantityRemaining {
						buyTx.PriceTx = lowestSellTx.Order.StockPrice
						buyTx.CostTotalTx += lowestSellTx.Order.Quantity * lowestSellTx.Order.StockPrice
					} else {
						if buyTx.IsParent || sellsHasNext || (!buyTx.IsParent && buyTx.Order.OrderType == "LIMIT") {
							var buyChildTx = createChildTx(&buyTx, sellQuantityRemaining, lowestSellTx.Order.StockPrice)
							*txCommitQueue = append(*txCommitQueue, buyChildTx)
						} else {
							buyTx.PriceTx = lowestSellTx.Order.StockPrice
							buyTx.CostTotalTx += lowestSellTx.Order.Quantity * lowestSellTx.Order.StockPrice
						}
					}

					book.sells.Delete(lowestSellTx.Order)
					if lowestSellTx.Order.OrderStatus == "IN_PROGRESS" { // not parent
						lowestSellTx.PriceTx = lowestSellTx.Order.StockPrice
					}
					lowestSellTx.QuantityTx = lowestSellTx.Order.Quantity - sellQuantityRemaining
					lowestSellTx.CostTotalTx += sellQuantityRemaining * lowestSellTx.Order.StockPrice
					lowestSellTx.Order.OrderStatus = "COMPLETED"
					*txCommitQueue = append(*txCommitQueue, lowestSellTx)

					buyQuantityRemaining -= sellQuantityRemaining
				} else { // buyQuantityRemaining < sellQuantityRemaining
					if buyTx.Order.OrderType == "LIMIT" || buyTx.IsParent {
						var buyChildTx = createChildTx(&buyTx, buyQuantityRemaining, lowestSellTx.Order.StockPrice)
						*txCommitQueue = append(*txCommitQueue, buyChildTx)
					} else {
						buyTx.PriceTx = lowestSellTx.Order.StockPrice
						buyTx.CostTotalTx += buyQuantityRemaining * lowestSellTx.Order.StockPrice
					}

					lowestSellTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
					var sellChildTx = createChildTx(&lowestSellTx, buyQuantityRemaining, lowestSellTx.Order.StockPrice)
					*txCommitQueue = append(*txCommitQueue, sellChildTx)

					buyQuantityRemaining = 0
				}
			}
		}

		buyTx.QuantityTx = buyTx.Order.Quantity - buyQuantityRemaining

		if buyTx.QuantityTx == 0 {
			//buyTx.Order.OrderStatus = "IN_PROGRESS" // unfulfilled
		} else if buyQuantityRemaining > 0 {
			buyTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
		} else { // = 0
			buyTx.Order.OrderStatus = "COMPLETED"
		}

		if buyTx.Order.OrderType == "LIMIT" {
			book.buys.Set(buyTx.Order, buyTx)
		}

		*txCommitQueue = append(*txCommitQueue, buyTx)
	}
}

func (book orderbook) matchSell(sellTx models.StockMatch, txCommitQueue *[]models.StockMatch) {
	if book.buys.Len() == 0 {
		if sellTx.Order.OrderType == "LIMIT" {
			book.sells.Set(sellTx.Order, sellTx)
		} else { // "MARKET"
			*txCommitQueue = append(*txCommitQueue, sellTx)
		}

	} else {
		var buysHasNext = true
		var buyIter = book.buys.Iterator()
		buyIter.Next()
		var sellQuantityRemaining = sellTx.Order.Quantity

		for sellQuantityRemaining > 0 && buysHasNext {
			highestBuyTx := buyIter.Value().(models.StockMatch)
			buysHasNext = buyIter.Next()
			buyQuantityRemaining := highestBuyTx.Order.Quantity - highestBuyTx.QuantityTx

			if isExpired(highestBuyTx) {
				book.buys.Delete(highestBuyTx.Order)
				*txCommitQueue = append(*txCommitQueue, highestBuyTx)
			} else {
				if (sellTx.Order.OrderType == "LIMIT") && (sellTx.Order.StockPrice < highestBuyTx.Order.StockPrice) {
					break
				}

				if sellQuantityRemaining >= buyQuantityRemaining {
					if sellTx.Order.Quantity == buyQuantityRemaining {
						sellTx.PriceTx = highestBuyTx.Order.StockPrice
						sellTx.CostTotalTx += highestBuyTx.Order.Quantity * highestBuyTx.Order.StockPrice
					} else {
						if sellTx.IsParent || buysHasNext || (!sellTx.IsParent && sellTx.Order.OrderType == "LIMIT") {
							var sellChildTx = createChildTx(&sellTx, buyQuantityRemaining, highestBuyTx.Order.StockPrice)
							*txCommitQueue = append(*txCommitQueue, sellChildTx)
						} else {
							sellTx.PriceTx = highestBuyTx.Order.StockPrice
							sellTx.CostTotalTx += highestBuyTx.Order.Quantity * highestBuyTx.Order.StockPrice
						}
					}

					book.buys.Delete(highestBuyTx.Order)
					if highestBuyTx.Order.OrderStatus == "IN_PROGRESS" { // not parent
						highestBuyTx.PriceTx = highestBuyTx.Order.StockPrice
					}
					highestBuyTx.QuantityTx = highestBuyTx.Order.Quantity - buyQuantityRemaining
					highestBuyTx.CostTotalTx += buyQuantityRemaining * highestBuyTx.Order.StockPrice
					highestBuyTx.Order.OrderStatus = "COMPLETED"
					*txCommitQueue = append(*txCommitQueue, highestBuyTx)

					sellQuantityRemaining -= buyQuantityRemaining
				} else { // sellQuantityRemaining < buyQuantityRemaining
					if sellTx.Order.OrderType == "LIMIT" || sellTx.IsParent {
						var sellChildTx = createChildTx(&sellTx, sellQuantityRemaining, highestBuyTx.Order.StockPrice)
						*txCommitQueue = append(*txCommitQueue, sellChildTx)
					} else {
						sellTx.PriceTx = highestBuyTx.Order.StockPrice
						sellTx.CostTotalTx += sellQuantityRemaining * highestBuyTx.Order.StockPrice
					}

					highestBuyTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
					var buyChildTx = createChildTx(&highestBuyTx, sellQuantityRemaining, highestBuyTx.Order.StockPrice)
					*txCommitQueue = append(*txCommitQueue, buyChildTx)

					sellQuantityRemaining = 0
				}
			}
		}

		sellTx.QuantityTx = sellTx.Order.Quantity - sellQuantityRemaining

		if sellTx.QuantityTx == 0 {
			//buyTx.Order.OrderStatus = "IN_PROGRESS" // unfulfilled
		} else if sellQuantityRemaining > 0 {
			sellTx.Order.OrderStatus = "PARTIALLY_FULFILLED"
		} else { // = 0
			sellTx.Order.OrderStatus = "COMPLETED"
		}

		if sellTx.Order.OrderType == "LIMIT" {
			book.sells.Set(sellTx.Order, sellTx)
		}

		*txCommitQueue = append(*txCommitQueue, sellTx)
	}
}

// CancelOrder halts further activity for a limit transaction with the given stockTxID.
// If found, the matching transaction is enqueued. Basically a deliberate premature expiration.
func CancelOrder(order models.StockTransaction) (wasCancelled bool) {
	var book = bookMap[order.StockID]
	var txCommitQueue []models.StockMatch
	if book.buys != nil {
		if order.IsBuy {
			wasCancelled = book.cancelBuyOrder(order, &txCommitQueue)
		} else {
			wasCancelled = book.cancelSellOrder(order, &txCommitQueue)
		}
	}

	ExecuteOrders(txCommitQueue)

	return wasCancelled
}

// cancelBuyOrder() and cancelSellOrder() are mirrors of each other, with "buy" and "sell" swapped.

func (book orderbook) cancelBuyOrder(order models.StockTransaction, txCommitQueue *[]models.StockMatch) (wasFound bool) {
	victimTx, wasFound := book.buys.Get(order)
	if wasFound {
		book.buys.Delete(order)
		victimTx := victimTx.(models.StockMatch)
		victimTx.Killed = true
		*txCommitQueue = append(*txCommitQueue, victimTx)
	}
	return wasFound
}

func (book orderbook) cancelSellOrder(order models.StockTransaction, txCommitQueue *[]models.StockMatch) (wasFound bool) {
	victimTx, wasFound := book.sells.Get(order)
	if wasFound {
		book.sells.Delete(order)
		victimTx := victimTx.(models.StockMatch)
		victimTx.Killed = true
		*txCommitQueue = append(*txCommitQueue, victimTx)
	}
	return wasFound
}

// FlushExpired checks all active transactions, pushing any that are expired.
func FlushExpired() {
	var txCommitQueue []models.StockMatch

	for _, book := range bookMap {
		txIter := book.buys.Iterator()
		for txIter.Next() {
			tx := txIter.Value().(models.StockMatch)
			if isExpired(tx) {
				book.sells.Delete(tx.Order)
				tx.Killed = true
				txCommitQueue = append(txCommitQueue, tx)
			}
		}
		txIter = book.sells.Iterator()
		for txIter.Next() {
			tx := txIter.Value().(models.StockMatch)
			if isExpired(tx) {
				book.sells.Delete(tx.Order)
				tx.Killed = true
				txCommitQueue = append(txCommitQueue, tx)
			}
		}
	}

	ExecuteOrders(txCommitQueue)
}
