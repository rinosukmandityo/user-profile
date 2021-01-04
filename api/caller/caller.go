package caller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type HTTPTestCaller struct {
	handler  http.Handler
	request  *http.Request
	response interface{}
}

func New(handler http.Handler) *HTTPTestCaller {
	return &HTTPTestCaller{handler: handler}
}

func (caller *HTTPTestCaller) SetRequest(request *http.Request) *HTTPTestCaller {
	caller.request = request
	return caller
}

func (caller *HTTPTestCaller) SetResponse(response interface{}) *HTTPTestCaller {
	caller.response = response
	return caller
}

func (caller *HTTPTestCaller) Exec() (*httptest.ResponseRecorder, interface{}, error) {
	recorder := httptest.NewRecorder()
	caller.handler.ServeHTTP(recorder, caller.request)

	var err error
	if caller.response != nil {
		err = json.Unmarshal(recorder.Body.Bytes(), &caller.response)
	}

	return recorder, caller.response, err
}

func (caller *HTTPTestCaller) SetHeader(param string, value string) *HTTPTestCaller {
	caller.request.Header.Set(param, value)
	return caller
}
