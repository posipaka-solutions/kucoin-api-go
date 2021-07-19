package kucoinfuncs

import (
	"encoding/json"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func TradeBotErrorCheck(body []byte, res *http.Response, resErr, bodyErr error) error {
	if resErr != nil {
		log.Print(resErr)
		return resErr

	} else if bodyErr != nil {
		log.Print(bodyErr)
		return bodyErr

	} else if res.StatusCode == 429 {
		return &exchangeapi.ExchangeError{
			Type:    exchangeapi.HttpErr,
			Code:    res.StatusCode,
			Message: res.Status,
		}
	} else if body != nil {
		if strings.Contains(string(body), "msg") {
			var bodyAnswer BodyAnswer

			jsonErr := json.Unmarshal(body, &bodyAnswer)
			if jsonErr != nil {
				return jsonErr
			}
			bodyCode, _ := strconv.Atoi(bodyAnswer.Code)
			if bodyCode != 200000 {
				return &exchangeapi.ExchangeError{
					Type:    exchangeapi.KucoinErr,
					Code:    bodyCode,
					Message: bodyAnswer.Msg,
				}
			}
		} else {
			return nil
		}

	} else if res != nil {
		if res.StatusCode != 200 {
			return &exchangeapi.ExchangeError{
				Type:    exchangeapi.HttpErr,
				Code:    res.StatusCode,
				Message: res.Status,
			}
		}
	}
	return &exchangeapi.ExchangeError{Type: exchangeapi.KucoinErr, Message: "Body is nil"}
}
