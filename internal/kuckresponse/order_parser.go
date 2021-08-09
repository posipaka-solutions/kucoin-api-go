package kuckresponse

import (
	"errors"
	"fmt"
	"net/http"
)

func ParseSetOrder(response *http.Response) (string, error) {
	body, err := getResponseBody(response)
	if err != nil {
		return "", err
	}
	bodyI, isOk := body.(map[string]interface{})
	if isOk != true {
		return "", errors.New("[kuckresponse] -> Error when casting body in ParseSetOrder")
	}

	dataI := bodyI["data"]
	orderI, isOk := dataI.(map[string]interface{})
	if isOk != true {
		return "", errors.New("[kuckresponse] -> Error when casting dataI in ParseSetOrder")
	}
	orderIdI := orderI["orderId"]
	orderId := fmt.Sprintf("%v", orderIdI)

	return orderId, nil
}
