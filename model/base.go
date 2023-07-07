package model

type BaseResp struct {
	Code    int64  `json:"code,required"`
	Message string `json:"message,required"`
}
