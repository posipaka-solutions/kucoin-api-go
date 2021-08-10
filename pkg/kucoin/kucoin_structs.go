package kucoin

import (
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi/order"
	"net/http"
)

const BaseUrl = "https://openapi-sandbox.kucoin.com" //"https://api.kucoin.com"

type ExchangeManager struct {
	apiKey exchangeapi.ApiKey

	client *http.Client
}

func New(key exchangeapi.ApiKey) *ExchangeManager {
	return &ExchangeManager{
		apiKey: key,
		client: &http.Client{},
	}
}

var orderSideAlias = map[order.Side]string{
	order.Buy:  "buy",
	order.Sell: "sell",
}

var orderTypeAlias = map[order.Type]string{
	order.Limit:  "limit",
	order.Market: "market",
}

const (
	newOrderEndpoint = "/api/v1/orders"
	//openOrdersEndpoint     = "/api/v3/openOrders"
	getPriceEndpoint = "/api/v1/market/orderbook/level1"
	//getCandlestickEndpoint = "/api/v3/klines"
	orderInfoEndpoint = "/api/v1/orders/"
)

//const goodTilCanceled = "GTC"
