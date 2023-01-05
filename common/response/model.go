package response

import "fmt"

type Response struct {
	// 数据集
	Code   int32  `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Msg    string `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	Status string `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
}

type response struct {
	Response
	Data interface{} `json:"data"`
}

type Page struct {
	Count     int `json:"count"`
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
}

type page struct {
	Page
	List interface{} `json:"list"`
}

func (e *response) SetData(data interface{}) {
	e.Data = data
}

func (e *response) Clone() Responses {
	res := *e
	return &res
}

func (e *response) SetMsg(s string) {
	if e.Msg == "" {
		e.Msg = s
	} else {
		e.Msg += fmt.Sprintf(", %s", s)
	}
}

func (e *response) SetCode(code int32) {
	e.Code = code
}

func (e *response) SetSuccess(success bool) {
	if !success {
		e.Status = "error"
	}
}
