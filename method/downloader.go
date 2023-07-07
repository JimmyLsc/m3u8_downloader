package method

import (
	"context"
	"github.com/JimmyLsc/m3u8_downloader/model"
	"github.com/JimmyLsc/m3u8_downloader/util"
)

type DownloaderHandler struct {
	Ctx  context.Context
	Req  *model.DownloaderReq
	Resp *model.DownloaderResp
}

func NewDownloaderHandler(ctx context.Context, req *model.DownloaderReq) *DownloaderHandler {
	return &DownloaderHandler{
		Ctx: ctx,
		Req: req,
		Resp: &model.DownloaderResp{
			BaseResp: util.GenSuccessBaseResp(),
		},
	}
}

func (d *DownloaderHandler) Run() {

}
