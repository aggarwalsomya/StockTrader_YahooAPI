package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//Structure for the JSON-RPC request sent by a client.
type ClientRequest struct {
	Method string `json:"method"`
	Params [1]interface{} `json:"params"`
	Id uint64 `json:"id"`
}

// JSON-RPC response returned to a client.
type ClientResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
	Id     uint64           `json:"id"`
}

// EncodeClientRequest encodes parameters for a JSON-RPC client request.
func EncodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &ClientRequest{
		Method: method,
		Params: [1]interface{}{args},
		Id:     uint64(rand.Int63()),
	}
	return json.Marshal(c)
}

// DecodeClientResponse decodes the response body
func DecodeClientResponse(r io.Reader, reply interface{}) error {
	var c ClientResponse
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return err
	}
	if c.Error != nil {
		return fmt.Errorf("%v", c.Error)
	}
	if c.Result == nil {
	return errors.New("result is null")
	}
	return json.Unmarshal(*c.Result, reply)
}

func Execute(method string, req, res interface{}) error {
	buf, _ := EncodeClientRequest(method, req)
	body := bytes.NewBuffer(buf)

	r, _ := http.NewRequest("POST", "http://localhost:8082/rpc", body)
	r.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(r)

	if resp == nil {
		fmt.Println("response is nil")
	} else {
		return DecodeClientResponse(resp.Body, res)
	}

	return nil
}

//prints the instrcutions for the user to enter the command line params
func printHelp() {

	fmt.Println("Missing or invalid command line parameters. Use the following format to invoke client")

	fmt.Println("\n\nFor Trading stocks:")
	fmt.Println("client \"trade\" \"2000.0 (Budget)\" \"GOOG:50%,AMZN:10%,AAPL:40% (Stocks information)\"")

	fmt.Println("\n\nFor Getting portfolio:")
	fmt.Println("client \"getportfolio\" \"12322 (tradeId)\"");

	fmt.Println("\n\nExample:")
	fmt.Println("client \"trade\" \"2000.0\" \"GOOG:50%,AMZN:10%,AAPL:40%\"")
	fmt.Println("client \"getportfolio\" \"1835\"")
}

func main() {
	var reply TradingResponse
	var portfolioReply PortfolioResponse
	
	// read these from command line
	var methodName string
	var budget float32 = 0
	var tradingStocks string
	var tradeId int = -1

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	methodName = os.Args[1]
	//matching strings ignoring case
	if strings.EqualFold(methodName, "trade") {
		if len(os.Args) < 4 {
			printHelp()
			return
		}

		// strconv is unnable to give 32 float for some reason. Hence explicit typecasting!!
		b, _ := strconv.ParseFloat(os.Args[2], 32)
		budget = float32(b)
		tradingStocks = os.Args[3]

		if budget == 0 || tradingStocks == "" {
			printHelp()
			return
		}
		if err := Execute("TradingService.TradeStocks", &TradingRequest{budget, tradingStocks}, &reply); err != nil {
			fmt.Println("Error!:", err)
		} else {
			fmt.Println("tradeId: ",reply.TradeId);
			fmt.Println("stocks: ",reply.Stocks);
			fmt.Println("unvestedAmount: " ,reply.UnvestedAmount);
		}

	} else if strings.EqualFold(methodName, "getportfolio") {

		if len(os.Args) < 3 {
			printHelp()
			return
		}

		tradeId,_ = strconv.Atoi(os.Args[2])
		if tradeId < 0 {
			printHelp()
			return
		}

		if err := Execute("TradingService.GetPortfolioDetails", &PortfolioRequest{tradeId}, &portfolioReply); err != nil {
			fmt.Println("Error!: ", err)
		} else {
			fmt.Println("stocks: ", portfolioReply.Stocks);
			fmt.Println("unvestedAmount: ", portfolioReply.UnvestedAmount);
			fmt.Println("currentMarketValue: ", portfolioReply.CurrentMarketValue);
		}
	} else {
		printHelp()
		return
	}
}
