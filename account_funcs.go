package kucoinfuncs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (kuCoinApi *KuCoinApi) GetCurrentPrice(symbol exchangeapi.AssetsSymbol) (float64, error) {
	var price float64
	endPoint := "/api/v1/market/orderbook/level1?"
	params := fmt.Sprintf("symbol=%s-%s", symbol.Base, symbol.Quote)

	body, tradeBotError := kuCoinApi.DoRequest(http.MethodGet, endPoint, params, nil)
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

func (kuCoinApi *KuCoinApi) GetSymbolLimits(symbol exchangeapi.AssetsSymbol) (exchangeapi.SymbolLimits, error) {
	resp, err := kuCoinApi.client.Get(fmt.Sprint(burl, "/api/v1/symbols?market=ETH-USDT"))

	if err != nil {
		return exchangeapi.SymbolLimits{}, err
	}

	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return exchangeapi.SymbolLimits{}, &exchangeapi.ExchangeError{
			Type:    exchangeapi.HttpErr,
			Code:    resp.StatusCode,
			Message: resp.Status,
		}
	}

	var body map[string]interface{}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return exchangeapi.SymbolLimits{}, err
	}
	err = json.Unmarshal(respBytes, &body)
	if err != nil {
		return exchangeapi.SymbolLimits{}, err
	}

	//for idx, _ := range body {
	//	parsedSymbol, isOkay := body[idx]["symbol"].(string)
	//	if !isOkay {
	//		continue
	//	}
	//
	//	if parsedSymbol != fmt.Sprint(symbol.Base, "-", symbol.Quote) {
	//		continue
	//	}
	//
	//	baseMinSize, isOkay := body[idx]["baseMinSize"].(string)
	//	if !isOkay {
	//		return exchangeapi.SymbolLimits{}, errors.New("baseMinSize parse failed")
	//	}
	//
	//	baseMaxSize, isOkay := body[idx]["baseMaxSize"].(string)
	//	if !isOkay {
	//		return exchangeapi.SymbolLimits{}, errors.New("baseMaxSize parse failed")
	//	}
	//
	//	baseIncrement, isOkay := body[idx]["baseIncrement"].(string)
	//	if !isOkay {
	//		return exchangeapi.SymbolLimits{}, errors.New("baseIncrement parse failed")
	//	}
	//
	//	priceIncrement, isOkay := body[idx]["priceIncrement"].(string)
	//	if !isOkay {
	//		return exchangeapi.SymbolLimits{}, errors.New("priceIncrement parse failed")
	//	}
	//
	//	limits := exchangeapi.SymbolLimits{
	//		Symbol: symbol,
	//	}
	//
	//	limits.BaseMinSize, _ = strconv.ParseFloat(baseMinSize, 64)
	//	limits.BaseMaxSize, _ = strconv.ParseFloat(baseMaxSize, 64)
	//	limits.BaseIncrement, _ = strconv.ParseFloat(baseIncrement, 64)
	//	limits.PriceIncrement, _ = strconv.ParseFloat(priceIncrement, 64)
	//	return limits, nil
	//}

	return exchangeapi.SymbolLimits{}, errors.New("failed to get symbol limits")
}

func (kuCoinApi *KuCoinApi) GetServerTime() (uint64, error) {
	resp, err := kuCoinApi.client.Get(fmt.Sprint(burl, "/api/v1/timestamp"))
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return 0, &exchangeapi.ExchangeError{
			Type:    exchangeapi.HttpErr,
			Code:    resp.StatusCode,
			Message: resp.Status,
		}
	}

	var body map[string]interface{}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(respBytes, &body)
	if err != nil {
		return 0, err
	}

	time, isOkay := body["data"].(uint64)
	if !isOkay {
		return 0, errors.New("failed to parse server time response")
	}
	return time, nil
}

func (kuCoinApi *KuCoinApi) SetOrder(orderParams exchangeapi.OrderParameters) (float64, error) {
	endpoint := "/api/v1/orders"
	bodyJson := map[string]string{}
	if orderParams.Type == exchangeapi.Market {
		if orderParams.Side == exchangeapi.Buy {
			bodyJson["clientOid"] = uuid.New().String()
			bodyJson["side"] = orderSideAlias[orderParams.Side]
			bodyJson["symbol"] = fmt.Sprint(orderParams.Symbol.Base, "-", orderParams.Symbol.Quote)
			bodyJson["type"] = orderTypeAlias[orderParams.Type]
			bodyJson["funds"] = fmt.Sprintf("%f", orderParams.Quantity)
		} else {
			bodyJson["clientOid"] = uuid.New().String()
			bodyJson["side"] = orderSideAlias[orderParams.Side]
			bodyJson["symbol"] = fmt.Sprint(orderParams.Symbol.Base, "-", orderParams.Symbol.Quote)
			bodyJson["type"] = orderTypeAlias[orderParams.Type]
			bodyJson["size"] = fmt.Sprintf("%f", orderParams.Quantity)
		}
	} else {
		bodyJson["clientOid"] = uuid.New().String()
		bodyJson["side"] = orderSideAlias[orderParams.Side]
		bodyJson["symbol"] = fmt.Sprint(orderParams.Symbol.Base, "-", orderParams.Symbol.Quote)
		bodyJson["type"] = orderTypeAlias[orderParams.Type]
		bodyJson["size"] = fmt.Sprintf("%f", orderParams.Quantity)
		bodyJson["price"] = fmt.Sprintf("%f", orderParams.Price)

	}
	body, tradeBotError := kuCoinApi.DoRequest(http.MethodPost, endpoint, "", bodyJson)
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

	if orderParams.Type == exchangeapi.Limit {
		price, tradeBotError := kuCoinApi.ReceiveData(orderIdI, "price")
		if tradeBotError != nil {
			return 0, tradeBotError
		}
		return price, nil
	} else {
		dealFunds, tradeBotError := kuCoinApi.ReceiveData(orderIdI, "dealFunds")
		if tradeBotError != nil {
			return 0, tradeBotError
		}
		return dealFunds, nil
	}
}

func (kuCoinApi *KuCoinApi) OrderInfo(orderId string) (interface{}, error) {
	endPoint := "/api/v1/orders/"

	body, tradeBotError := kuCoinApi.DoRequest(http.MethodGet, endPoint, orderId, nil)
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

	body, tradeBotError := kuCoinApi.DoRequest(http.MethodGet, endPoint, params, nil)
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
