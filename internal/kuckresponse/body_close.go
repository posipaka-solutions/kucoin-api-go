package kuckresponse

import "net/http"

func CloseBody(response *http.Response) {
	if err := response.Body.Close(); err != nil {
		panic(err.Error())
	}
}
