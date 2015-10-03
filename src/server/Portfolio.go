package main

import (
	"fmt"
)

type resP struct {
	stockInfo_P []stockInfo_Portfolio
	currentMarketValue float64
	unvestedAmount float64
	errorString error //handing error cases
}

type stockInfo_Portfolio struct {
	stockName string
	noOfStocks int
	stockPrice float64
	//profit = 1; loss = 0; neutral = 2
	profitOrLoss int
}

func getPortfolioDetails(responseInfoObj responseInfo) resP {
	var resPObj resP
	
	fmt.Println("\n\nGetting you your portfolio details-");
	
	//Filling up the unvested amount
	resPObj.unvestedAmount = responseInfoObj.unvestedAmount
	resPObj.errorString = nil
	
	//Get the current symbol list from response
	var symbolList [] string
	var stockNameToOldPrices = make(map[string]float64) 
	
	for i := 0; i < len(responseInfoObj.stockInfo_res); i++ {
		var sp stockInfo_Portfolio
		symbolList = append(symbolList, responseInfoObj.stockInfo_res[i].stockName);
		
		//create a map of stock name to their old prices
		stockNameToOldPrices[responseInfoObj.stockInfo_res[i].stockName] = responseInfoObj.stockInfo_res[i].stockPrice
		
		//Fill the repsonse Portfolio strcutre with stockname and no of stocks
		sp.stockName = responseInfoObj.stockInfo_res[i].stockName
		sp.noOfStocks = responseInfoObj.stockInfo_res[i].noOfStocks
		resPObj.stockInfo_P = append(resPObj.stockInfo_P,sp)
	}
	
	
	
	//Send this list to yahoo comm to get current market price
	currentStockPrices := getAsksForSymbols(symbolList)
	
	//Filling the response Map with current price and profit/loss status
	var totalWorth float64
	for i := 0; i < len(symbolList); i++ {
		resPObj.stockInfo_P[i].stockPrice = currentStockPrices[resPObj.stockInfo_P[i].stockName]
		
		//check if there is profit or loss
		currentPrice := resPObj.stockInfo_P[i].stockPrice
		oldPrice := stockNameToOldPrices[resPObj.stockInfo_P[i].stockName]
		resPObj.stockInfo_P[i].profitOrLoss = checkProfitOrLoss(oldPrice, currentPrice)
		
		//get the total worth
		totalWorth = totalWorth + (resPObj.stockInfo_P[i].stockPrice * (float64(resPObj.stockInfo_P[i].noOfStocks)))
	}
	
	resPObj.currentMarketValue = totalWorth
	return resPObj
}

//checking profit/loss/neutral case
func checkProfitOrLoss(buyingPrice float64, currentPrice float64) int {
	if buyingPrice < currentPrice {
		return 1
		} else if buyingPrice > currentPrice{
		return 0
		} else {
		return 2
		}
}