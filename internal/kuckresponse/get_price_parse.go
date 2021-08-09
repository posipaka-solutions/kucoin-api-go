package kuckresponse

import (
	"errors"
	"net/http"
	"strconv"
)

func GetCurrentPriceParser(response *http.Response) (float64, error) {
	bodyI, err := getResponseBody(response)
	if err != nil {
		return 0, err
	}

	priceI, isOk := bodyI.(map[string]interface{})
	if !isOk {
		return 0, errors.New("[kuckresponse] -> error when casting bodyI to priceI GetCurrentPriceParser")
	}

	dataI := priceI["data"]

	data, isOk := dataI.(map[string]interface{})
	if isOk != true {
		return 0, errors.New("[kuckresponse] -> Error when casting dataI in GetCurrentPriceParser")
	}
	testI := data["price"]

	price, err := strconv.ParseFloat(testI.(string), 64)
	if err != nil {
		return 0, errors.New("[kuckresponse] -> error when parsing priceStr to float64 GetCurrentPriceParser")
	}

	return price, nil
}
