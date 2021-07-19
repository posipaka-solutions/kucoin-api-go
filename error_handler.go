package kucoinfuncs

import (
	"github.com/posipaka-trade/KuCoin/src/cmn"
	"strconv"
)
import (
	"encoding/json"
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
		return &cmn.FakeBuyError{
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
				return &cmn.FakeBuyError{
					Type:    cmn.BinanceErr,
					Code:    bodyCode,
					Message: bodyAnswer.Msg,
				}
			}
		} else {
			return nil
		}

	} else if res != nil {
		if res.StatusCode != 200 {
			return &cmn.FakeBuyError{
				Type:    cmn.HttpErr,
				Code:    res.StatusCode,
				Message: res.Status,
			}
		}
	}
	return &cmn.FakeBuyError{Type: cmn.BinanceErr, Message: "Body is nil"}
}
