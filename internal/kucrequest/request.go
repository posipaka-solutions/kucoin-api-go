package kucrequest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/posipaka-trade/posipaka-trade-cmn/exchangeapi"
	"net/http"
	"strconv"
	"time"
)

type HeaderParams struct {
	BaseUrl  string
	Endpoint string
	Method   string
	BodyJson []byte
	Key      exchangeapi.ApiKey
}

func SetHeader(req *http.Request, headP HeaderParams) {

	signature, passPh := Sign(headP)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("KC-API-KEY", headP.Key.Key)
	req.Header.Add("KC-API-PASSPHRASE", passPh)
	req.Header.Add("KC-API-TIMESTAMP", strconv.Itoa(int(time.Now().Unix()*1000)))
	req.Header.Add("KC-API-SIGN", signature)
	req.Header.Add("KC-API-KEY-VERSION", "2")
}

func Sign(headP HeaderParams) (string, string) {
	timestamp := strconv.Itoa(int(time.Now().Unix() * 1000))
	var signature string

	if headP.BodyJson != nil {
		strToSign := timestamp + headP.Method + headP.Endpoint + string(headP.BodyJson)
		h := hmac.New(sha256.New, []byte(headP.Key.Secret))
		h.Write([]byte(strToSign))
		signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
	} else {
		strToSign := timestamp + headP.Method + headP.Endpoint
		h := hmac.New(sha256.New, []byte(headP.Key.Secret))
		h.Write([]byte(strToSign))
		signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
	}

	p := hmac.New(sha256.New, []byte(headP.Key.Secret))
	p.Write([]byte(headP.Key.Passphrase))
	passPh := base64.StdEncoding.EncodeToString(p.Sum(nil))

	return signature, passPh
}
