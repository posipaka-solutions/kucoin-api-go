package kuckresponse

import (
	"errors"
	"github.com/posipaka-trade/kucoin-api-go/internal/pnames"
	"net/http"
	"strconv"
)

func GetCurrentPriceParser(response *http.Response) (float64, error) {
	body, err := getResponseBody(response)
	if err != nil {
		return 0, err
	}

	bodyI, isOk := body.(map[string]interface{})
	if !isOk {
		return 0, errors.New("[kuckresponse] -> Error when casting body to bodyI GetCurrentPriceParser")
	}

	dataI := bodyI[pnames.Data]

	data, isOk := dataI.(map[string]interface{})
	if isOk != true {
		return 0, errors.New("[kuckresponse] -> Error when casting dataI in GetCurrentPriceParser")
	}
	priceI := data[pnames.Price]
	priceStr, isOk := priceI.(string)
	if isOk != true {
		return 0, errors.New("[kuckresponse] -> Error when casting priceI ro string in GetCurrentPriceParser")
	}
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, errors.New("[kuckresponse] -> Error when parsing priceStr to float64 GetCurrentPriceParser")
	}

	return price, nil
}
