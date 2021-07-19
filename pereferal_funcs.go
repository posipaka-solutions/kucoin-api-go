package kucoinfuncs

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (kuCoinApi *KuCoinApi) MakeSignature(baseUrl, method string, bodyJson []byte) (string, string) {
	timestamp := strconv.Itoa(int(time.Now().Unix() * 1000))
	var signature string

	if bodyJson != nil {
		strToSign := timestamp + method + baseUrl + string(bodyJson)
		h := hmac.New(sha256.New, []byte(kuCoinApi.ApiSecret))
		h.Write([]byte(strToSign))
		signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
	} else {
		strToSign := timestamp + method + baseUrl
		h := hmac.New(sha256.New, []byte(kuCoinApi.ApiSecret))
		h.Write([]byte(strToSign))
		signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
	}

	p := hmac.New(sha256.New, []byte(kuCoinApi.ApiSecret))
	p.Write([]byte(kuCoinApi.ApiPassPh))
	passPh := base64.StdEncoding.EncodeToString(p.Sum(nil))

	return signature, passPh
}

func (kuCoinApi *KuCoinApi) DoRequest(method, endpoint, params string, bodyJson map[string]string) ([]byte, error) {

	var req *http.Request
	var err error
	bodyJsonStr, err := json.Marshal(bodyJson)
	if err != nil {
		return nil, err
	}

	if bodyJson != nil {
		req, err = http.NewRequest(method, burl+endpoint, bytes.NewBuffer(bodyJsonStr))
		if err != nil {
			return nil, err
		}
		kuCoinApi.HeaderAdd(req, method, endpoint, bodyJsonStr)
	} else {
		req, err = http.NewRequest(method, burl+endpoint+params, nil)
		if err != nil {
			return nil, err
		}
		kuCoinApi.HeaderAdd(req, method, endpoint+params, nil)
	}

	res, resErr := kuCoinApi.Client.Do(req)

	body, bodyErr := ioutil.ReadAll(res.Body)
	tradeBotError := TradeBotErrorCheck(body, res, resErr, bodyErr)
	if tradeBotError != nil {
		res.Body.Close()
		return nil, tradeBotError
	}

	return body, nil
}

func (kuCoinApi *KuCoinApi) HeaderAdd(req *http.Request, method, endpoint string, bodyJsonStr []byte) {
	signature, passPh := kuCoinApi.MakeSignature(endpoint, method, bodyJsonStr)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("KC-API-KEY", kuCoinApi.ApiKey)
	req.Header.Add("KC-API-PASSPHRASE", passPh)
	req.Header.Add("KC-API-TIMESTAMP", strconv.Itoa(int(time.Now().Unix()*1000)))
	req.Header.Add("KC-API-SIGN", signature)
	req.Header.Add("KC-API-KEY-VERSION", "2")
}

func (kuCoinApi *KuCoinApi) ReceiveData(orderIdI interface{}, whatToFind string) (float64, error) {
	orderIdByte, err := json.Marshal(orderIdI)
	if err != nil {
		return 0, err
	}
	var data map[string]interface{}

	err = json.Unmarshal(orderIdByte, &data)
	if err != nil {
		return 0, err
	}

	whatToReturnI := data[whatToFind]
	whatToReturn := fmt.Sprintf("%v", whatToReturnI)

	whatToReturnF, err := strconv.ParseFloat(whatToReturn, 64)
	if err != nil {
		return 0, err
	}
	return whatToReturnF, nil
}
