package main

import ( 
//	"fmt";
	"math/rand"
)

type requestInfo struct {
	budget     float64
	stockInfo_req []stockInfo_request
}

type stockInfo_request struct {
	sharePercent float64
	stockName    string
}

type responseInfo struct {
	tradeId        int
	errorString string 		// it will be empty if no error. else, it will put the error string in it
	stockInfo_res []stockInfo_response
	unvestedAmount float64
}

type stockInfo_response struct {
	stockName  string
	noOfStocks int
	stockPrice float64
}

func stockProcessing(requestInfoObj requestInfo) responseInfo {
	var responseInfoObj responseInfo
	checkPercent := checkFor100Percent(requestInfoObj)
	if !checkPercent {
		responseInfoObj.errorString = "Sum of all share percents should be 100."
		return responseInfoObj
	}
	
	var symbolList [] string
	for i := 0; i < len(requestInfoObj.stockInfo_req); i++ {
		symbolList = append(symbolList, requestInfoObj.stockInfo_req[i].stockName);
	}

	var investedAmount float64
	stockPrices := getAsksForSymbols(symbolList)
	if len(stockPrices) == 0 {
		responseInfoObj.errorString = "Invalid symbol. Check request."
		return responseInfoObj
	}

	for i := 0; i < len(requestInfoObj.stockInfo_req); i++ {
		var resp stockInfo_response
		resp.stockName = requestInfoObj.stockInfo_req[i].stockName 
		resp.stockPrice = getAskForStock(requestInfoObj.stockInfo_req[i].stockName, stockPrices)
		if resp.stockPrice <= 0 {
			// error scenario for invalid symbol
			responseInfoObj.errorString = "Invalid symbol- " +  requestInfoObj.stockInfo_req[i].stockName
			return responseInfoObj
		}
		
		resp.noOfStocks = int(((requestInfoObj.budget * requestInfoObj.stockInfo_req[i].sharePercent)/100)/ resp.stockPrice);
		
		investedAmount = investedAmount + (resp.stockPrice)* float64(resp.noOfStocks)
		responseInfoObj.stockInfo_res = append(responseInfoObj.stockInfo_res, resp)
	}

	responseInfoObj.tradeId = generateTradeId()
	responseInfoObj.unvestedAmount = requestInfoObj.budget - investedAmount
	responseInfoObj.errorString = ""
	return responseInfoObj
}

//sum of all % should be 100
func checkFor100Percent(requestInfoObj requestInfo) bool {
	var total float64
	for i := 0; i < len(requestInfoObj.stockInfo_req); i++ {
		total = total + requestInfoObj.stockInfo_req[i].sharePercent
	}
	if total == 100 {
		return true
	} else {
		return false
	}
}


//generating random trade id
func generateTradeId() int {
	r := rand.Intn(10000)
	return r
}

//getting the values from the map formed from yahoo query
func getAskForStock(stockName string, stockPrices map[string]float64) float64 {
	ask, present := stockPrices[stockName]
	if !present {
		return -1
	}
	
	return ask
}
