package main

import (
	"net/http"
	"strings";
	"strconv";
	"fmt";
	"errors";
)

type TradingService struct{
	TradeIdToResponseMap map[int]responseInfo
}

func (ts *TradingService) TradeStocks(r *http.Request, req *TradingRequest, res *TradingResponse) error {

	fmt.Println("Got the request for trading");
	var requestInfoObj requestInfo
	requestInfoObj.budget = float64(req.Budget)

	//split the string using ,
	res1 := strings.Split(req.StockSymbolAndPercentage,",")

	//split stock info by :
	for i:= range res1 {
		res2 :=  strings.Split(res1[i], ":")

		var t1 stockInfo_request
		t1.stockName = res2[0]
		t1.sharePercent, _ = strconv.ParseFloat(strings.TrimRight(res2[1], "%"), 32)
		requestInfoObj.stockInfo_req = append(requestInfoObj.stockInfo_req, t1)
	}

	//s1 is the response from the stockprocessing 
	s1 := stockProcessing(requestInfoObj)
	if(s1.errorString != "") {
		res.ErrorMsg = s1.errorString
		err := errors.New(res.ErrorMsg)
		return err
	} else {
	res.ErrorMsg = ""
	}

	if(len(ts.TradeIdToResponseMap) == 0) {
		ts.TradeIdToResponseMap = make(map[int]responseInfo)
	}
	ts.TradeIdToResponseMap[s1.tradeId] = s1
	
	//convert this response to the form in which user wants it
	res.TradeId = s1.tradeId
	res.UnvestedAmount = float32(s1.unvestedAmount)
	
	for i := 0; i < len(s1.stockInfo_res); i++ {
		if i > 0 {
			res.Stocks = res.Stocks + ","
		}
		price := strconv.FormatFloat(s1.stockInfo_res[i].stockPrice, 'f', -1, 64)
		res.Stocks = res.Stocks + s1.stockInfo_res[i].stockName + ":" + strconv.Itoa(s1.stockInfo_res[i].noOfStocks) +":$" + price
	}
	
	fmt.Println(s1); 
	return nil
}

//getting portfolio details for the user
func (ts *TradingService) GetPortfolioDetails(r *http.Request, req *PortfolioRequest, res *PortfolioResponse) error {
	var responseInfoObj responseInfo
	
	responseInfoObj,ok := ts.TradeIdToResponseMap[req.TradeId]
	if !ok {
		err:= errors.New("Invalid Trade Id")
		res.ErrorMsg = "Invalid Trade Id"
		return err
	}
	
	Presp:= getPortfolioDetails(responseInfoObj)
	fmt.Println(Presp);
	
	//convert this response into a form user wants
	res.CurrentMarketValue = float32(Presp.currentMarketValue)
	res.UnvestedAmount = float32(Presp.unvestedAmount)
	
	for i := 0; i < len(Presp.stockInfo_P); i++ {
		if i > 0 {
			res.Stocks = res.Stocks + ","
		}
		price := strconv.FormatFloat(Presp.stockInfo_P[i].stockPrice, 'f', -1, 64)
		var sym string
		if(Presp.stockInfo_P[i].profitOrLoss == 1) {
			sym = "+";
			}else if Presp.stockInfo_P[i].profitOrLoss == 0{
			sym = "-"
			} else {
			sym = "";
			}
		res.Stocks = res.Stocks + Presp.stockInfo_P[i].stockName + ":" + strconv.Itoa(Presp.stockInfo_P[i].noOfStocks) +":"+ sym+"$" + price
	}

	return nil
}
