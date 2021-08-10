package main

import (
	"fmt"
	"github.com/posipaka-trade/kucoin-api-go/pkg/kucoin"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi/order"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi/symbol"
)

func main() {
	//var api kucoinfuncs.KuCoinApi
	//api.ApiKey= "a"
	mgr := kucoin.New(exchangeapi.ApiKey{
		Key:        "60f18bc2bc85c200065a4cc7",             //os.Args[1],
		Secret:     "8c382e50-a90f-4670-86dc-fe9fd6aaeb42", //os.Args[2],
		Passphrase: "25801379",
	})
	//price, err := mgr.GetCurrentPrice(symbol.Assets{
	//	Base:  "ETH",
	//	Quote: "USDT"})
	//fmt.Println(price)
	//if err != nil {
	//	panic(err.Error())
	//}
	parametrs := order.Parameters{Assets: symbol.Assets{Base: "ETH", Quote: "USDT"}, Side: order.Buy, Type: order.Market, Quantity: 55}
	priceqwe, err := mgr.SetOrder(parametrs)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(priceqwe)
}
