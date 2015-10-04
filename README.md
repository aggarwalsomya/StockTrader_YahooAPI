# CMPE273-Fall15-Assignment1
**Virtual Trade stocking system**

This explains how to use the virtual stock trading system developed as part of CMP273 Fall15 assignment 1

#Usage

**Start the server**

go get github.com/aggarwalsomya/273_Assignment1_YahooStock/server
cd server
go run *

**Start the client**

go get github.com/aggarwalsomya/273_Assignment1_YahooStock/client
cd client

*To trade stocks*
go run * "trade" "1000" "GOOG:50%,AMZN:10%,AAPL:40%"

*To get portfolio*
go run * "getportfolio" <tradeId>


**Using curl to make requests**

One can use the curl command too to make requests to the server

*To trade stocks*
curl  -H "Content-Type: application/json"  -d '{"method":"TradingService.TradeStocks","params":[{"Budget":2000,"StockSymbolAndPercentage":"GOOG:50%,AMZN:20%,AAPL:30%"}],"id":0}' http://localhost:8082/rpc


*To get portfolio*
curl  -H "Content-Type: application/json"  -d '{"method":"TradingService.GetPortfolioDetails","params":[{"TradeId":8081}],"id":0}' http://localhost:8082/rpc


