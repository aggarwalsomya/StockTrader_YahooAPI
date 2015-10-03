package main

import (
	"fmt"
//	"time"
	"sync"
)

func main1() {

	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Starting Go Routines")
	go func() {
		defer wg.Done()

//		init_server()
	}()
	
	go func() {
		defer wg.Done()
//		init_client()
	} ()
	
	fmt.Println("Waiting To Finish")
    wg.Wait()
    
    fmt.Println("\nTerminating Program")

	//	fmt.Println("IN MAIN.GO");
	//
	//	var requestInfoObj requestInfo
	//	requestInfoObj.budget = 2000
	//
	//	var r1 stockInfo_request
	//	r1.sharePercent = 70
	//	r1.stockName = "GOOGL"
	//
	//	var r2 stockInfo_request
	//	r2.sharePercent = 30
	//	r2.stockName = "AMZN"
	//
	//	requestInfoObj.stockInfo_req = append(requestInfoObj.stockInfo_req, r1);
	//	requestInfoObj.stockInfo_req = append(requestInfoObj.stockInfo_req, r2);
	//
	//	responseInfoObj := stockProcessing(requestInfoObj);
	//	fmt.Println(responseInfoObj);
	//	time.Sleep(60 * time.Second)
	//
	//	responsePortfolio := getPortfolioDetails(responseInfoObj)
	//	fmt.Println(responsePortfolio);
}
