import React from 'react'
import StockTransactions from "./StockTransactions.js";
import StockPrices from "./StockPrices.js";
import StockPortfolio from "./StockPortfolio.js";


function MyStock() {
  return (
    <div>
    <StockTransactions></StockTransactions>
    <StockPrices></StockPrices>
    <StockPortfolio></StockPortfolio>
    </div>
  )
}

export default MyStock