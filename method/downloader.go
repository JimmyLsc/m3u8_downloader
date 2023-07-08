package method

import (
	"context"
	"github.com/JimmyLsc/m3u8_downloader/logic"
	"github.com/JimmyLsc/m3u8_downloader/model"
	"github.com/JimmyLsc/m3u8_downloader/util"
	log "github.com/sirupsen/logrus"
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
	if !d.checkParam() {
		return
	}
	ctx := d.Ctx
	srcUrl := d.Req.SrcURL
	info, err := logic.GetM3U8Info(ctx, srcUrl, make(map[string]string))
	if err != nil {
		log.Errorf("DownloaderHandler error, err: %v", err)
		return
	}
	info.Name = d.Req.SrcName
	info.ShortName = *d.Req.SrcShortName
	logic.DownloadVideo(ctx, info, *d.Req.CachePath, *d.Req.FilePath, *d.Req.WorkerNum)
}

func (d *DownloaderHandler) checkParam() bool {
	if d.Req.SrcURL == "" || d.Req.SrcName == "" {
		d.Resp.BaseResp = util.GenParamErrorResp()
		return false
	}
	if d.Req.SrcShortName == nil || *(d.Req.SrcShortName) == "" {
		d.Req.SrcShortName = &(d.Req.SrcName)
	}

	if d.Req.CachePath == nil || *(d.Req.CachePath) == "" {
		cachePath := "./.cache"
		d.Req.CachePath = &(cachePath)
	}

	if d.Req.FilePath == nil || *(d.Req.FilePath) == "" {
		filePath := "./video"
		d.Req.FilePath = &(filePath)
	}

	if d.Req.WorkerNum == nil || *(d.Req.WorkerNum) == 0 {
		workerNum := int64(12)
		d.Req.WorkerNum = &(workerNum)
	}
	return true
}
