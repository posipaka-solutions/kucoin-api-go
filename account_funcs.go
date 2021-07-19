package kucoinfuncs

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"strconv"
	//"strconv"
)

func (kuCoinApi *KuCoinApi) GetPrice(currency string, fiat string) (float64, error) {
	var price float64
	endPoint := "/api/v1/market/orderbook/level1?"
	params := fmt.Sprintf("symbol=%s-%s", currency, fiat)

	body, tradeBotError := kuCoinApi.DoRequest(Get, endPoint, params, nil)
	if tradeBotError != nil {
		return 0, tradeBotError
	}

	priceMap := new(PriceMap)
	err := json.Unmarshal(body, priceMap)
	if err != nil {
		return 0, err
	}
	priceStr := priceMap.Data.Price
	price, err = strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}

	return price, nil
	//return string(body),nil
}

func (kuCoinApi *KuCoinApi) NewMarketOrder(orderParams MarketOrdersParams) (float64, error) {
	endpoint := "/api/v1/orders"
	bodyJson := map[string]string{}
	if orderParams.Side == Buy {
		bodyJson["clientOid"] = orderParams.Uuid
		bodyJson["side"] = orderParams.Side
		bodyJson["symbol"] = orderParams.Symbol
		bodyJson["type"] = orderParams.OrderType
		bodyJson["funds"] = fmt.Sprintf("%f", orderParams.FiatCount)
	} else {
		bodyJson["clientOid"] = orderParams.Uuid
		bodyJson["side"] = orderParams.Side
		bodyJson["symbol"] = orderParams.Symbol
		bodyJson["type"] = orderParams.OrderType
		bodyJson["size"] = fmt.Sprintf("%f", orderParams.AssetCount)
	}
	body, tradeBotError := kuCoinApi.DoRequest(Post, endpoint, "", bodyJson)
	if tradeBotError != nil {
		return 0, tradeBotError
	}
	var bodyAnswer map[string]interface{}

	err := json.Unmarshal(body, &bodyAnswer)
	if err != nil {
		return 0, err
	}
	dataI := bodyAnswer["data"]
	var dataMap map[string]string

	dataByte, err := json.Marshal(dataI)
	err = json.Unmarshal(dataByte, &dataMap)
	if err != nil {
		return 0, err
	}

	orderIdI, tradeBotError := kuCoinApi.OrderInfo(dataMap["orderId"])
	if tradeBotError != nil {
		return 0, tradeBotError
	}

	orderIdByte, err := json.Marshal(orderIdI)
	if err != nil {
		return 0, err
	}
	var data map[string]interface{}

	err = json.Unmarshal(orderIdByte, &data)
	if err != nil {
		return 0, err
	}

	quantity := data["dealFunds"]
	dealFunds := fmt.Sprintf("%v", quantity)

	dealFundsFloat, err := strconv.ParseFloat(dealFunds, 64)
	if err != nil {
		return 0, err
	}
	return dealFundsFloat, nil
}

func (kuCoinApi *KuCoinApi) OrderInfo(orderId string) (interface{}, error) {
	endPoint := "/api/v1/orders/"

	body, tradeBotError := kuCoinApi.DoRequest(Get, endPoint, orderId, nil)
	if tradeBotError != nil {
		return "", tradeBotError
	}
	var boughtPrice map[string]interface{}

	err := json.Unmarshal(body, &boughtPrice)
	if err != nil {
		return "", err
	}

	dataI := boughtPrice["data"]

	return dataI, nil

}

func (kuCoinApi *KuCoinApi) GetAllOrders() (string, error) {
	endPoint := "/api/v1/orders?"
	params := "status=done"

	body, tradeBotError := kuCoinApi.DoRequest(Get, endPoint, params, nil)
	if tradeBotError != nil {
		return "", tradeBotError
	}
	return string(body), nil

}

func (kuCoinApi KuCoinApi) NewLimitOrder(orderParams LimitOrdersParams) (bool, error) {
	endpoint := "/api/v1/orders"
	bodyJson := map[string]string{}
	bodyJson["clientOid"] = orderParams.Uuid
	bodyJson["side"] = orderParams.Side
	bodyJson["symbol"] = orderParams.Symbol
	bodyJson["type"] = orderParams.OrderType
	bodyJson["size"] = fmt.Sprintf("%f", orderParams.FiatCount)
	bodyJson["price"] = fmt.Sprintf("%f", orderParams.Price)

	body, tradeBotError := kuCoinApi.DoRequest(Post, endpoint, "", bodyJson)
	if tradeBotError != nil {
		return false, tradeBotError
	}
	var bodyAnswer map[string]interface{}

	err := json.Unmarshal(body, &bodyAnswer)
	if err != nil {
		return false, err
	}

	dataI := bodyAnswer["data"]
	var dataMap map[string]string

	dataByte, err := json.Marshal(dataI)

	err = json.Unmarshal(dataByte, &dataMap)
	if err != nil {
		return false, nil
	}
	orderIdI, tradeBotError := kuCoinApi.OrderInfo(dataMap["orderId"])
	if tradeBotError != nil {
		return false, tradeBotError
	}
	orderIdByte, err := json.Marshal(orderIdI)
	if err != nil {
		return false, err
	}

	var data map[string]interface{}

	err = json.Unmarshal(orderIdByte, &data)
	if err != nil {
		return false, err
	}

	isActiveI := data["isActive"]
	isActiveStr := fmt.Sprintf("%v", isActiveI)
	isActive, err := strconv.ParseBool(isActiveStr)
	if err != nil {
		return false, err
	}
	return isActive, nil

}
