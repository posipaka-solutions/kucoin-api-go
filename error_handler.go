package kucoinfuncs

import (
	"strconv"
)
import (
	"encoding/json"
	cmn "github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi"
	"log"
	"net/http"
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
		return &cmn.ExchangeError{
			Type:    cmn.HttpErr,
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
				return &cmn.ExchangeError{
					Type:    cmn.KucoinErr,
					Code:    bodyCode,
					Message: bodyAnswer.Msg,
				}
			}
		} else {
			return nil
		}

	} else if res != nil {
		if res.StatusCode != 200 {
			return &cmn.ExchangeError{
				Type:    cmn.HttpErr,
				Code:    res.StatusCode,
				Message: res.Status,
			}
		}
	}
	return &cmn.ExchangeError{Type: cmn.KucoinErr, Message: "Body is nil"}
}
