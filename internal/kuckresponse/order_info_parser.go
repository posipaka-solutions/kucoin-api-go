package kuckresponse

import (
	"errors"
	"fmt"
	"github.com/posipaka-trade/kucoin-api-go/internal/pnames"
	"net/http"
	"strconv"
)

func OrderInfoParser(response *http.Response) (float64, error) {
	body, err := getResponseBody(response)
	if err != nil {
		return 0, err
	}

	bodyI, isOk := body.(map[string]interface{})
	if isOk != true {
		return 0, errors.New("[kuckresponse] -> Error when casting body in OrderInfoParser")
	}

	orderI, isOk := bodyI[pnames.Data].(map[string]interface{})
	if isOk != true {
		return 0, errors.New("[kuckresponse] -> Error when casting dataI in OrderInfoParser")
	}

	dealFunds, err := receiveData(orderI, pnames.Size)
	if err != nil {
		return 0, err
	}

	return dealFunds, nil

}
func receiveData(orderIdI map[string]interface{}, whatToFind string) (float64, error) {

	whatToReturnI := orderIdI[whatToFind]
	whatToReturn := fmt.Sprintf("%v", whatToReturnI)

	whatToReturnF, err := strconv.ParseFloat(whatToReturn, 64)
	if err != nil {
		return 0, errors.New("[kuckresponse] -> Error when parsing whatToReturn to float64 in ReceiveData")
	}

	return whatToReturnF, nil
}
