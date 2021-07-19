package kucoinfuncs

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	cmn "github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi"
)

func (kuCoinApi *KuCoinApi) GetCurrentPrice(currency string, fiat string) (float64, error) {
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
}

func (kuCoinApi *KuCoinApi) SetOrder(orderParams cmn.OrderParameters) (float64, error) {
	endpoint := "/api/v1/orders"
	bodyJson := map[string]string{}
	if orderParams.Type == "MARKET" {
		if orderParams.Side == {
			bodyJson["clientOid"] = uuid.New().String()
			bodyJson["side"] = orderParams.Side
			bodyJson["symbol"] = orderParams.Symbol
			bodyJson["type"] = orderParams.Type
			bodyJson["funds"] = fmt.Sprintf("%f", orderParams.Quantity)
		} else {
			bodyJson["clientOid"] = uuid.New().String()
			bodyJson["side"] = orderParams.Side
			bodyJson["symbol"] = orderParams.Symbol
			bodyJson["type"] = orderParams.Type
			bodyJson["size"] = fmt.Sprintf("%f", orderParams.Quantity)
		}
	}else {
		bodyJson["clientOid"] = uuid.New().String()
		bodyJson["side"] = orderParams.Side
		bodyJson["symbol"] = orderParams.Symbol
		bodyJson["type"] = orderParams.Type
		bodyJson["size"] = fmt.Sprintf("%f", orderParams.Quantity)
		bodyJson["price"] = fmt.Sprintf("%f", orderParams.Price)

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

	if orderParams.Type == "LIMIT"{
		price,tradeBotError := kuCoinApi.ReceiveData(orderIdI,"price")
		if tradeBotError != nil {
			return 0, tradeBotError
		}
		return price,nil
	} else{
		dealFunds, tradeBotError := kuCoinApi.ReceiveData(orderIdI,"dealFunds")
		if tradeBotError != nil {
			return 0, tradeBotError
		}
		return dealFunds,nil
	}
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

func (kuCoinApi *KuCoinApi) GetAllOrders() (string, error) { //не смотри на нее
	endPoint := "/api/v1/orders?"
	params := "status=done"

	body, tradeBotError := kuCoinApi.DoRequest(Get, endPoint, params, nil)
	if tradeBotError != nil {
		return "", tradeBotError
	}

	return string(body), nil

}

//func (kuCoinApi *KuCoinApi) NewLimitOrder(orderParams LimitOrdersParams) (bool, error) {
//	endpoint := "/api/v1/orders"
//	bodyJson := map[string]string{}
//	bodyJson["clientOid"] = orderParams.Uuid
//	bodyJson["side"] = orderParams.Side
//	bodyJson["symbol"] = orderParams.Symbol
//	bodyJson["type"] = orderParams.OrderType
//	bodyJson["size"] = fmt.Sprintf("%f", orderParams.FiatCount)
//	bodyJson["price"] = fmt.Sprintf("%f", orderParams.Price)
//
//	body, tradeBotError := kuCoinApi.DoRequest(Post, endpoint, "", bodyJson)
//	if tradeBotError != nil {
//		return false, tradeBotError
//	}
//	var bodyAnswer map[string]interface{}
//
//	err := json.Unmarshal(body, &bodyAnswer)
//	if err != nil {
//		return false, err
//	}
//
//	dataI := bodyAnswer["data"]
//	var dataMap map[string]string
//
//	dataByte, err := json.Marshal(dataI)
//
//	err = json.Unmarshal(dataByte, &dataMap)
//	if err != nil {
//		return false, nil
//	}
//
//	orderIdI, tradeBotError := kuCoinApi.OrderInfo(dataMap["orderId"])
//	if tradeBotError != nil {
//		return false, tradeBotError
//	}
//
//	orderIdByte, err := json.Marshal(orderIdI)
//	if err != nil {
//		return false, err
//	}
//
//	var data map[string]interface{}
//
//	err = json.Unmarshal(orderIdByte, &data)
//	if err != nil {
//		return false, err
//	}
//
//	isActiveI := data["isActive"]
//	isActiveStr := fmt.Sprintf("%v", isActiveI)
//	isActive, err := strconv.ParseBool(isActiveStr)
//	if err != nil {
//		return false, err
//	}
//	return isActive, nil
//
//}
