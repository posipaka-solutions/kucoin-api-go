package kucoin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/posipaka-trade/kucoin-api-go/internal/kuckresponse"
	"github.com/posipaka-trade/kucoin-api-go/internal/kucrequest"
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
		return 0, err
	}

	kucrequest.SetHeader(request, kucrequest.HeaderParams{BaseUrl: BaseUrl, Endpoint: newOrderEndpoint, Method: http.MethodPost, BodyJson: requestBody, Key: manager.apiKey})

	response, err := manager.client.Do(request)
	if err != nil {
		return 0, err
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
			bodyJson["clientOid"] = uuid.New().String()
			bodyJson["side"] = orderSideAlias[params.Side]
			bodyJson["symbol"] = fmt.Sprint(params.Assets.Base, "-", params.Assets.Quote)
			bodyJson["type"] = orderTypeAlias[params.Type]
			bodyJson["funds"] = fmt.Sprintf("%f", params.Quantity)
		} else {
			bodyJson["clientOid"] = uuid.New().String()
			bodyJson["side"] = orderSideAlias[params.Side]
			bodyJson["symbol"] = fmt.Sprint(params.Assets.Base, "-", params.Assets.Quote)
			bodyJson["type"] = orderTypeAlias[params.Type]
			bodyJson["size"] = fmt.Sprintf("%f", params.Quantity)
		}
	} else {
		bodyJson["clientOid"] = uuid.New().String()
		bodyJson["side"] = orderSideAlias[params.Side]
		bodyJson["symbol"] = fmt.Sprint(params.Assets.Base, "-", params.Assets.Quote)
		bodyJson["type"] = orderTypeAlias[params.Type]
		bodyJson["size"] = fmt.Sprintf("%f", params.Quantity)
		bodyJson["price"] = fmt.Sprintf("%f", params.Price)

	}
	body, err := json.Marshal(bodyJson)
	if err != nil {
		return nil, errors.New("[kucoin] -> Error when marshaling body to bodyJson in createOrderRequestBody")
	}
	return body, nil
}
