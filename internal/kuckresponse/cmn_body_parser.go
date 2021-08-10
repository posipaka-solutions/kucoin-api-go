package kuckresponse

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	codeKey = "code" // error code
	msgKey  = "msg"  // error message
)

func getResponseBody(response *http.Response) (interface{}, error) {
	if response.StatusCode/100 != 2 && response.Body == nil {
		return nil, &exchangeapi.ExchangeError{
			Type:    exchangeapi.HttpErr,
			Code:    response.StatusCode,
			Message: response.Status,
		}
	}

	respondBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	err = json.Unmarshal(respondBody, &body)
	if err != nil {
		return nil, err
	}

	kucoinErrorCodeI, isOk := body[codeKey].(string)
	if !isOk {
		return nil, errors.New("[kuckresponse] -> error code key not found")
	}

	kucoinErrorCode, err := strconv.Atoi(fmt.Sprintf("%v", kucoinErrorCodeI))
	if err != nil {
		return nil, errors.New("[kuckresponse] -> error when parsing kucoinErrorCodeI to float64")
	}

	if kucoinErrorCode > 200000 && kucoinErrorCode < 900002 {
		return nil, parseKucoinError(body, kucoinErrorCode)
	}

	return body, nil
}

func parseKucoinError(body map[string]interface{}, kucoinErrorCode int) error {

	message, isOkay := body[msgKey].(string)
	if !isOkay {
		return errors.New("[kuckresponse] -> failed to parse binance error message")
	}

	return &exchangeapi.ExchangeError{
		Type:    exchangeapi.KucoinErr,
		Code:    kucoinErrorCode,
		Message: message,
	}
}
