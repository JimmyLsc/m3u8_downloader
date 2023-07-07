package util

import "github.com/JimmyLsc/m3u8_downloader/model"

func GenSuccessBaseResp() *model.BaseResp {
	return &model.BaseResp{
		Code:    0,
		Message: "Success",
	}
}

func GenParamErrorResp() *model.BaseResp {
	return &model.BaseResp{
		Code:    1,
		Message: "Param Error",
	}
}
