package api

import (
	"log"
	"net/http"

	"github.com/rinosukmandityo/user-profile/helper"
)

func SetupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	if _, e := w.Write(body); e != nil {
		log.Println(e)
	}
}

func ResponseWithResult(w http.ResponseWriter, contentType string, result *helper.ResultInfo, statusCode int) {
	respBody, e := GetSerializer(contentType).EncodeResult(result)
	if e != nil {
		statusCode = http.StatusBadRequest
	}
	SetupResponse(w, contentType, respBody, statusCode)
}
