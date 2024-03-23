// The Matching Engine pairs buy and sell orders to create mutually agreeable transactions,
// which are placed in a queue to be handled by the Order Execution component.
// Unmatched orders remain in the Matching Engines memory until matched, expired, or canceled.

package matching

import (
	"day-trading-app/backend/internal/service/models"
	"github.com/google/uuid"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/adrianbrad/queue"
	"github.com/ryszard/goskiplist/skiplist"
)

func init() {
	// Start ticker goroutine to periodically flush expired limit orders
	go func() {
		for range time.Tick(1 * time.Second) {
			FlushExpired()
		}
	}()
}

// A StockTransaction plus additional info needed for processing orders.
type StockMatch struct {
	Order       models.StockTransaction `json:"order" bson:"order"`                 // original order; though matching engine will change OrderStatus
	QuantityTx  int                     `json:"quantity_tx" bson:"quantity_tx"`     // quantity actually transacted
	PriceTx     int                     `json:"price_tx" bson:"price_tx"`           // price actually transacted
	CostTotalTx int                     `json:"cost_total_tx" bson:"cost_total_tx"` // total cost transacted; needed for parent tx
	IsParent    bool                    `json:"is_parent" bson:"is_parent"`         // true if transaction has created a child
	Killed      bool                    `json:"killed" bson:"killed"`               // expired or cancelled
}

// Key for sorting orderbooks.
type bookKey struct {
	StockPrice int
	TimeStamp  int64
}

type orderbook struct {
	buys  *skiplist.SkipList
	sells *skiplist.SkipList
	m     *sync.Mutex
}

var bookMap = make(map[string]orderbook)
var bookMapLock sync.Mutex
var expireQueue = queue.NewPriority([]StockMatch{}, lessTime, queue.WithCapacity(1000000)) // todo: adjust capacity

// Define ordering comparison value for expiration queue
func lessTime(elem StockMatch, elemAfter StockMatch) bool {
	return elem.Order.TimeStamp < elemAfter.Order.TimeStamp
}

// Define orderings for orderbooks. Sort first by price, then if equal price, by time.

// Prioritize buys by highest price 1st, then earliest time 2nd.
func buyIsLowerPriorityThan(l, r interface{}) bool {
	ll := l.(bookKey)
	rr := r.(bookKey)
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
	ll := l.(bookKey)
	rr := r.(bookKey)
	if ll.StockPrice < rr.StockPrice {
		return true
	} else if ll.StockPrice == rr.StockPrice {
		return ll.TimeStamp < rr.TimeStamp
	} else {
		return false
	}
}

func makeBookKey(tx models.StockTransaction) bookKey {
	return bookKey{
		StockPrice: tx.StockPrice,
		TimeStamp:  tx.TimeStamp,
	}
}

// Returns true if transaction's timestamp is over 15 minutes old.
func isExpired(tx StockMatch) bool {
	return time.Now().UnixNano() >= tx.Order.TimeStamp+(15*time.Minute).Nanoseconds()
}

func getOrderbook(tx models.StockTransaction) orderbook {
	if bookMap[tx.StockID].buys == nil {
		bookMapLock.Lock()
		defer bookMapLock.Unlock()
		if bookMap[tx.StockID].buys == nil {
			bookMap[tx.StockID] = orderbook{
				buys:  skiplist.NewCustomMap(buyIsLowerPriorityThan),
				sells: skiplist.NewCustomMap(sellIsLowerPriorityThan),
				m:     &sync.Mutex{},
			}
		}
	}
	return bookMap[tx.StockID]
}

// Create child transaction based on parentTx, ordering some specified quantity of otherTx.
// Also adds transaction cost to parentTx.CostTotalTx.
func createChildTx(parentTx *StockMatch, quantityTraded int, priceTraded int) StockMatch {
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

	// add string "Tx" inbetween stockID's name, for example, "googleStockId" becomes "googleStockTxId"
	index := strings.Index(childTx.Order.StockID, "Stock")
	childTx.Order.StockTxID = childTx.Order.StockID[:index+len("Stock")] + "Tx" + childTx.Order.StockID[index+len("Stock"):] + uuid.New().String()

	walletTxID := strings.Replace(childTx.Order.StockTxID, "StockTxId", "WalletTxId", 1)
	childTx.Order.WalletTxID = &walletTxID

	parentTx.IsParent = true
	parentTx.QuantityTx += quantityTraded
	parentTx.CostTotalTx += quantityTraded * priceTraded
	return childTx
}

// Match inserts a transaction into the matching engine, to be matched with complementary transaction(s) in its order book.
func Match(order models.StockTransaction) (err error) {
	var book = getOrderbook(order)

	var tx = StockMatch{Order: order, QuantityTx: 0, PriceTx: 0, CostTotalTx: 0, IsParent: false, Killed: false}

	if order.OrderType == "LIMIT" {
		err = expireQueue.Offer(tx)
		if err != nil {
			return err
		}
	}

	var txCommitQueue []StockMatch

	if order.IsBuy {
		book.matchBuy(tx, &txCommitQueue)
	} else {
		book.matchSell(tx, &txCommitQueue)
	}

	err = ExecuteOrders(txCommitQueue)
	return err
}

// matchBuy() and matchSell() are basically mirrors of each other, with "buy" and "sell" swapped.
// todo Specification unclear as to what price that two matched price limit-orders with different prices should trade at.
// todo   So for now we're taking the oldest limit-order's price.

func (book orderbook) matchBuy(buyTx StockMatch, txCommitQueue *[]StockMatch) {
	book.m.Lock()
	defer book.m.Unlock()

	if book.sells.Len() == 0 {
		if buyTx.Order.OrderType == "LIMIT" {
			book.buys.Set(makeBookKey(buyTx.Order), buyTx)
		} else { // "MARKET"
			*txCommitQueue = append(*txCommitQueue, buyTx)
		}

	} else {
		var sellsHasNext = true
		var sellIter = book.sells.Iterator()
		sellIter.Next()
		var buyQuantityRemaining = buyTx.Order.Quantity

		for buyQuantityRemaining > 0 && sellsHasNext {
			lowestSellTx := sellIter.Value().(StockMatch)
			sellsHasNext = sellIter.Next()
			sellQuantityRemaining := lowestSellTx.Order.Quantity - lowestSellTx.QuantityTx

			if isExpired(lowestSellTx) {
				book.sells.Delete(makeBookKey(lowestSellTx.Order))
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

					book.sells.Delete(makeBookKey(lowestSellTx.Order))
					if lowestSellTx.Order.OrderStatus == "IN_PROGRESS" { // not parent
						lowestSellTx.PriceTx = lowestSellTx.Order.StockPrice
					}
					lowestSellTx.QuantityTx = lowestSellTx.Order.Quantity - sellQuantityRemaining
					lowestSellTx.CostTotalTx += sellQuantityRemaining * lowestSellTx.Order.StockPrice
					lowestSellTx.Order.OrderStatus = "COMPLETED"

					buyQuantityRemaining -= sellQuantityRemaining
				} else { // buyQuantityRemaining < sellQuantityRemaining
					if buyTx.IsParent {
						var buyChildTx = createChildTx(&buyTx, buyQuantityRemaining, lowestSellTx.Order.StockPrice)
						*txCommitQueue = append(*txCommitQueue, buyChildTx)
					} else {
						buyTx.PriceTx = lowestSellTx.Order.StockPrice
						buyTx.CostTotalTx += buyQuantityRemaining * lowestSellTx.Order.StockPrice
					}

					lowestSellTx.Order.OrderStatus = "PARTIAL_FULFILLED"
					var sellChildTx = createChildTx(&lowestSellTx, buyQuantityRemaining, lowestSellTx.Order.StockPrice)
					*txCommitQueue = append(*txCommitQueue, sellChildTx)
					book.sells.Set(makeBookKey(lowestSellTx.Order), lowestSellTx)

					buyQuantityRemaining = 0
				}
			}
			*txCommitQueue = append(*txCommitQueue, lowestSellTx)
		}

		buyTx.QuantityTx = buyTx.Order.Quantity - buyQuantityRemaining

		if buyTx.QuantityTx == 0 {
			//buyTx.Order.OrderStatus = "IN_PROGRESS"
		} else if buyQuantityRemaining > 0 {
			buyTx.Order.OrderStatus = "PARTIAL_FULFILLED"
		} else { // = 0
			buyTx.Order.OrderStatus = "COMPLETED"
		}

		if buyTx.Order.OrderStatus != "COMPLETED" && buyTx.Order.OrderType == "LIMIT" {
			book.buys.Set(makeBookKey(buyTx.Order), buyTx)
		}

		if !(buyTx.Order.OrderStatus == "IN_PROGRESS" && buyTx.Order.OrderType == "LIMIT") {
			*txCommitQueue = append(*txCommitQueue, buyTx)
		}
	}
}

func (book orderbook) matchSell(sellTx StockMatch, txCommitQueue *[]StockMatch) {
	book.m.Lock()
	defer book.m.Unlock()

	if book.buys.Len() == 0 {
		if sellTx.Order.OrderType == "LIMIT" {
			book.sells.Set(makeBookKey(sellTx.Order), sellTx)
		} else { // "MARKET"
			*txCommitQueue = append(*txCommitQueue, sellTx)
		}

	} else {
		var buysHasNext = true
		var buyIter = book.buys.Iterator()
		buyIter.Next()
		var sellQuantityRemaining = sellTx.Order.Quantity

		for sellQuantityRemaining > 0 && buysHasNext {
			highestBuyTx := buyIter.Value().(StockMatch)
			buysHasNext = buyIter.Next()
			buyQuantityRemaining := highestBuyTx.Order.Quantity - highestBuyTx.QuantityTx

			if isExpired(highestBuyTx) {
				book.buys.Delete(makeBookKey(highestBuyTx.Order))
			} else {
				if (sellTx.Order.OrderType == "LIMIT") && (sellTx.Order.StockPrice > highestBuyTx.Order.StockPrice) {
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

					book.buys.Delete(makeBookKey(highestBuyTx.Order))
					if highestBuyTx.Order.OrderStatus == "IN_PROGRESS" { // not parent
						highestBuyTx.PriceTx = highestBuyTx.Order.StockPrice
					}
					highestBuyTx.QuantityTx = highestBuyTx.Order.Quantity - buyQuantityRemaining
					highestBuyTx.CostTotalTx += buyQuantityRemaining * highestBuyTx.Order.StockPrice
					highestBuyTx.Order.OrderStatus = "COMPLETED"

					sellQuantityRemaining -= buyQuantityRemaining
				} else { // sellQuantityRemaining < buyQuantityRemaining
					if sellTx.IsParent {
						var sellChildTx = createChildTx(&sellTx, sellQuantityRemaining, highestBuyTx.Order.StockPrice)
						*txCommitQueue = append(*txCommitQueue, sellChildTx)
					} else {
						sellTx.PriceTx = highestBuyTx.Order.StockPrice
						sellTx.CostTotalTx += sellQuantityRemaining * highestBuyTx.Order.StockPrice
					}

					highestBuyTx.Order.OrderStatus = "PARTIAL_FULFILLED"
					var buyChildTx = createChildTx(&highestBuyTx, sellQuantityRemaining, highestBuyTx.Order.StockPrice)
					*txCommitQueue = append(*txCommitQueue, buyChildTx)
					book.buys.Set(makeBookKey(highestBuyTx.Order), highestBuyTx)

					sellQuantityRemaining = 0
				}
			}
			*txCommitQueue = append(*txCommitQueue, highestBuyTx)
		}

		sellTx.QuantityTx = sellTx.Order.Quantity - sellQuantityRemaining

		if sellTx.QuantityTx == 0 {
			//buyTx.Order.OrderStatus = "IN_PROGRESS"
		} else if sellQuantityRemaining > 0 {
			sellTx.Order.OrderStatus = "PARTIAL_FULFILLED"
		} else { // = 0
			sellTx.Order.OrderStatus = "COMPLETED"
		}

		if sellTx.Order.OrderStatus != "COMPLETED" && sellTx.Order.OrderType == "LIMIT" {
			book.sells.Set(makeBookKey(sellTx.Order), sellTx)
		}

		if !(sellTx.Order.OrderStatus == "IN_PROGRESS" && sellTx.Order.OrderType == "LIMIT") {
			*txCommitQueue = append(*txCommitQueue, sellTx)
		}
	}
}

// CancelOrder halts further activity for a limit transaction with the given stockTxID.
// If found, the matching transaction is enqueued. Basically a deliberate premature expiration.
func CancelOrder(order models.StockTransaction) (wasCancelled bool, err error) {
	var book = bookMap[order.StockID]
	book.m.Lock()

	var txCommitQueue []StockMatch
	if book.buys != nil {
		if order.IsBuy {
			wasCancelled = book.cancelBuyOrder(order, &txCommitQueue)
		} else {
			wasCancelled = book.cancelSellOrder(order, &txCommitQueue)
		}
	}
	book.m.Unlock()

	err = ExecuteOrders(txCommitQueue)
	return wasCancelled, err
}

// cancelBuyOrder() and cancelSellOrder() are mirrors of each other, with "buy" and "sell" swapped.

func (book orderbook) cancelBuyOrder(order models.StockTransaction, txCommitQueue *[]StockMatch) (wasFound bool) {
	victimTx, wasFound := book.buys.Get(makeBookKey(order))
	if wasFound {
		book.buys.Delete(makeBookKey(order))
		victimTx := victimTx.(StockMatch)
		victimTx.Killed = true
		*txCommitQueue = append(*txCommitQueue, victimTx)
	}
	return wasFound
}

func (book orderbook) cancelSellOrder(order models.StockTransaction, txCommitQueue *[]StockMatch) (wasFound bool) {
	victimTx, wasFound := book.sells.Get(makeBookKey(order))
	if wasFound {
		book.sells.Delete(makeBookKey(order))
		victimTx := victimTx.(StockMatch)
		victimTx.Killed = true
		*txCommitQueue = append(*txCommitQueue, victimTx)
	}
	return wasFound
}

// FlushExpired checks for any expired transactions based on a priority queue.
func FlushExpired() {
	var allFresh = false
	for !allFresh && !expireQueue.IsEmpty() {
		var oldest, _ = expireQueue.Peek()
		if isExpired(oldest) {
			oldest, _ = expireQueue.Get()
			_, err := CancelOrder(oldest.Order)
			if err != nil {
				log.Println(err)
			}
		} else {
			allFresh = true
		}
	}
}
