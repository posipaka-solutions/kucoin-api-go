package kucoin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/posipaka-trade/kucoin-api-go/internal/kuckresponse"
	"github.com/posipaka-trade/kucoin-api-go/internal/kucrequest"
	"github.com/posipaka-trade/kucoin-api-go/internal/pnames"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi/order"
	"net/http"
)

func (manager *ExchangeManager) SetOrder(parameters order.Parameters) (float64, error) {
	requestBody, err := manager.createOrderRequestBody(&parameters)
	if err != nil {
		return 0, err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprint(BaseUrl, newOrderEndpoint), bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, errors.New("[kucoin] -> Error when making new request in SetOrder")
	}

	kucrequest.SetHeader(request, kucrequest.HeaderParams{BaseUrl: BaseUrl, Endpoint: newOrderEndpoint, Method: http.MethodPost, BodyJson: requestBody, Key: manager.apiKey})

	response, err := manager.client.Do(request)
	if err != nil {
		return 0, errors.New("[kucoin] -> Error when doing request in SetOrder")
	}

	defer kuckresponse.CloseBody(response)

	orderId, err := kuckresponse.ParseSetOrder(response)
	if err != nil {
		return 0, err
	}

	return manager.OrderInfo(orderId, parameters)
}

func (manager *ExchangeManager) createOrderRequestBody(params *order.Parameters) ([]byte, error) {
	bodyJson := map[string]string{}

	if params.Type == order.Market {
		if params.Side == order.Buy {
			bodyJson[pnames.OrderId] = uuid.New().String()
			bodyJson[pnames.Side] = orderSideAlias[params.Side]
			bodyJson[pnames.Symbol] = fmt.Sprint(params.Assets.Base, "-", params.Assets.Quote)
			bodyJson[pnames.Type] = orderTypeAlias[params.Type]
			bodyJson[pnames.Funds] = fmt.Sprintf("%f", params.Quantity)
		} else {
			bodyJson[pnames.OrderId] = uuid.New().String()
			bodyJson[pnames.Side] = orderSideAlias[params.Side]
			bodyJson[pnames.Symbol] = fmt.Sprint(params.Assets.Base, "-", params.Assets.Quote)
			bodyJson[pnames.Type] = orderTypeAlias[params.Type]
			bodyJson[pnames.Size] = fmt.Sprintf("%f", params.Quantity)
		}
	} else {
		bodyJson[pnames.OrderId] = uuid.New().String()
		bodyJson[pnames.Side] = orderSideAlias[params.Side]
		bodyJson[pnames.Symbol] = fmt.Sprint(params.Assets.Base, "-", params.Assets.Quote)
		bodyJson[pnames.Type] = orderTypeAlias[params.Type]
		bodyJson[pnames.Size] = fmt.Sprintf("%f", params.Quantity)
		bodyJson[pnames.Price] = fmt.Sprintf("%f", params.Price)

	}

	body, err := json.Marshal(bodyJson)
	if err != nil {
		return nil, errors.New("[kucoin] -> Error when marshaling body to bodyJson in createOrderRequestBody")
	}

	return body, nil
}
