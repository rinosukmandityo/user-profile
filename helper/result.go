package helper

import (
	"time"
)

type ResultStatus string
type ResultSuccess bool

const (
	Status_OK   ResultStatus  = "OK"
	Status_NOK  ResultStatus  = "NOK"
	Success_OK  ResultSuccess = true
	Success_NOK ResultSuccess = false
)

type ResultInfo struct {
	Data        interface{}
	Status      ResultStatus
	Success     ResultSuccess
	Message     string
	ArrMsg      []string
	Duration    time.Duration
	DurationTxt string
	time0       time.Time
	Total       int
}

func NewResult(defaultData interface{}) *ResultInfo {
	r := new(ResultInfo)
	r.Status = Status_OK
	r.Success = Success_OK
	r.time0 = time.Now()
	r.Data = defaultData
	return r
}

func (r *ResultInfo) SetData(o interface{}) *ResultInfo {
	r.Data = o
	r.SetDuration()
	return r
}

func (r *ResultInfo) SetTotal(o int) *ResultInfo {
	r.Total = o
	r.SetDuration()
	return r
}

func (r *ResultInfo) SetError(e error) *ResultInfo {
	r.Status = Status_NOK
	r.Success = Success_NOK
	r.Message = e.Error()
	r.SetDuration()
	return r
}

func (r *ResultInfo) SetErrMsg(str string) *ResultInfo {
	r.Status = Status_NOK
	r.Success = Success_NOK
	r.Message = str
	r.SetDuration()
	return r
}

func (r *ResultInfo) SetDuration() *ResultInfo {
	r.Duration = time.Since(r.time0)
	r.DurationTxt = r.Duration.String()
	return r
}

func (r *ResultInfo) SetMessage(message string) *ResultInfo {
	r.Message = message
	return r
}

func (r *ResultInfo) SetArrMsg(arrMsg []string) *ResultInfo {
	r.ArrMsg = arrMsg
	return r
}

func (r *ResultInfo) SetErrArrMsg(arrMsg []string) *ResultInfo {
	r.Status = Status_NOK
	r.Success = Success_NOK
	r.ArrMsg = arrMsg
	r.SetDuration()
	return r
}
