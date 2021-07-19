package kucoinfuncs

import (
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi"
	"net/http"

)

//type MarketOrdersParams struct {
//	Uuid       string
//	Symbol     string
//	Side       string
//	OrderType  string
//	AssetCount float64
//	FiatCount  float64
//}
//type LimitOrdersParams struct {
//	Uuid      string
//	Symbol    string
//	Side      string
//	OrderType string
//	FiatCount float64
//	Price     float64
//}

type KuCoinApi struct {
	ApiKey    string
	ApiSecret string
	ApiPassPh string
	Client    http.Client
}
type BodyAnswer struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
type PriceMap struct {
	Code string     `json:"code"`
	Data TickerData `json:"data"`
}
type TickerData struct {
	Time        int64  `json:"time"`
	Sequence    string `json:"sequence"`
	Price       string `json:"price"`
	Size        string `json:"size"`
	BestBid     string `json:"bestBid"`
	BestBidSize string `json:"bestBidSize"`
	BestAsk     string `json:"bestAsk"`
	BestAskSize string `json:"bestAskSize"`
}

var orderSideAlias = map[exchangeapi.OrderSide]string{
	exchangeapi.Buy: "BUY",
	exchangeapi.Sell: "SELL",
}

var orderTypeAlias = map[exchangeapi.OrderType]string{
	exchangeapi.Limit: "LIMIT",
	exchangeapi.Market: "MARKET",
}

const (
	Buy  = "buy"
	Sell = "sell"
)
const (
	Gtc = "GTC"
)

// trading order type
const (
	Limit  = "limit"
	Market = "market"
)

// Request methods
const (
	Get    = "GET"
	Post   = "POST"
	Delete = "DELETE"
)
const burl = "https://openapi-sandbox.kucoin.com" //"https://api.kucoin.com" //"https://openapi-sandbox.kucoin.com"
