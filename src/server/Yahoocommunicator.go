package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type stock struct {
	Query struct {
		Created string `json:"created"`
		Results struct {
			Quote quotes `json:"quote"`
		}
	}
}

type singlestock struct {
	Query struct {
		Created string `json:"created"`
		Results struct {
			Quote quoteInfo `json:"quote"`
		}
	}
}

type quotes []quoteInfo
type quoteInfo struct {
	Symbol string `json:"symbol"`
	Ask    string `json:"Ask"`
}

//yql for getting the current stock price
func getAsksForSymbols(symbols []string) map[string]float64 {
	var baseurl string = "http://query.yahooapis.com/v1/public/yql?q="
	var pricequeryString string = "select * from yahoo.finance.quotes where symbol in ('%s')"
	var symbolsString string
	for _, s := range symbols {
		symbolsString += s + ","
	}

	var finalUrlString string = baseurl + url.QueryEscape(fmt.Sprintf(pricequeryString, symbolsString)) + "&format=json&env=http://datatables.org/alltables.env"
//	fmt.Println("final url is " + finalUrlString)
	resp, err := http.Get(finalUrlString)

	var stockPrices = make(map[string]float64)
	if err != nil {
		return stockPrices
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if len(symbols) > 1 {
	
		var m stock
		_ = json.Unmarshal(body, &m)
	
		for i := 0; i < len(m.Query.Results.Quote); i++ {
			stockPrices[m.Query.Results.Quote[i].Symbol], _ = strconv.ParseFloat(m.Query.Results.Quote[i].Ask, 32)
		}
	} else {
		// parse single stock. Go json can't handle it automatically. 
		var m singlestock
		_ = json.Unmarshal(body, &m)
		stockPrices[m.Query.Results.Quote.Symbol], _ = strconv.ParseFloat(m.Query.Results.Quote.Ask, 32)

	}
	return stockPrices
}
