package kuckresponse

import (
	"errors"
	"fmt"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi/order"
	"net/http"
	"strconv"
)

func OrderInfoParser(response *http.Response, parameters order.Parameters) (float64, error) {
	body, err := getResponseBody(response)
	if err != nil {
		return 0, err
	}

	bodyI, isOk := body.(map[string]interface{})
	if isOk != true {
		return 0, errors.New("[kuckresponse] -> Error when casting body in OrderInfoParser")
	}

	dataI := bodyI["data"]
	orderI, isOk := dataI.(map[string]interface{})
	if isOk != true {
		return 0, errors.New("[kuckresponse] -> Error when casting dataI in OrderInfoParser")
	}
	if parameters.Type == order.Limit {
		price, err := ReceiveData(orderI, "price")
		if err != nil {
			return 0, err
		}
		return price, nil
	} else {
		dealFunds, err := ReceiveData(orderI, "dealFunds")
		if err != nil {
			return 0, err
		}
		return dealFunds, nil
	}
}
func ReceiveData(orderIdI map[string]interface{}, whatToFind string) (float64, error) {

	whatToReturnI := orderIdI[whatToFind]
	whatToReturn := fmt.Sprintf("%v", whatToReturnI)

	whatToReturnF, err := strconv.ParseFloat(whatToReturn, 64)
	if err != nil {
		return 0, errors.New("[kuckresponse] -> Error when parsing whatToReturn to float64")
	}
	return whatToReturnF, nil
}
