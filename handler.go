package main

import (
	"context"
	"github.com/JimmyLsc/m3u8_downloader/method"
	"github.com/JimmyLsc/m3u8_downloader/model"
)

type DownloadServiceImpl struct {
}

func (d *DownloadServiceImpl) Download(ctx context.Context, req *model.DownloaderReq) *model.DownloaderResp {
	handler := method.NewDownloaderHandler(ctx, req)
	handler.Run()
	return handler.Resp
}
