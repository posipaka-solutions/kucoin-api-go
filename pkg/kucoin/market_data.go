package kucoin

import (
	"errors"
	"fmt"
	"github.com/posipaka-trade/kucoin-api-go/internal/kuckresponse"
	"github.com/posipaka-trade/kucoin-api-go/internal/kucrequest"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi/order"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi/symbol"
	"net/http"
)

func (manager *ExchangeManager) GetCurrentPrice(symbol symbol.Assets) (float64, error) {
	params := fmt.Sprintf("symbol=%s-%s", symbol.Base, symbol.Quote)

	response, err := manager.client.Get(fmt.Sprint(BaseUrl, getPriceEndpoint, "?", params))
	if err != nil {
		return 0, errors.New("[kucoin] -> Error in GetRequest when getting current price")
	}
	defer kuckresponse.CloseBody(response)
	return kuckresponse.GetCurrentPriceParser(response)
}
func (manager *ExchangeManager) OrderInfo(orderId string, parameters order.Parameters) (float64, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprint(BaseUrl, orderInfoEndpoint, orderId), nil)
	if err != nil {
		return 0, errors.New("[kucoin] -> Error when making new request for OrderInfo")
	}
	kucrequest.SetHeader(request, kucrequest.HeaderParams{BaseUrl: BaseUrl, Endpoint: orderInfoEndpoint + orderId, Method: http.MethodGet, BodyJson: nil, Key: manager.apiKey})
	response, err := manager.client.Do(request)
	if err != nil {
		return 0, errors.New("[kucoin] -> Error when doing OrderInfo request ")
	}
	defer kuckresponse.CloseBody(response)
	return kuckresponse.OrderInfoParser(response, parameters)
}
