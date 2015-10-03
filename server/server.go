package main

import(
    "net/http";
    "log";
    "github.com/gorilla/rpc";
    "github.com/gorilla/rpc/json";
    "fmt";
)

func main() {
    s := rpc.NewServer()
    s.RegisterCodec(json.NewCodec(), "application/json")
    //this will register the trading service use to trade for stocks and getting portfolio info
    s.RegisterService(new(TradingService), "TradingService")
    http.Handle("/rpc", s)
    //running it at port 8082
    fmt.Println("Starting server at port 8082");
    log.Fatal(http.ListenAndServe(":8082", nil))
}

