package main

import (
	"context"
	"github.com/JimmyLsc/m3u8_downloader/model"
	"github.com/JimmyLsc/m3u8_downloader/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	req := &model.DownloaderReq{}
	resp := (&DownloadServiceImpl{}).Download(context.TODO(), req)
	log.Infof("Response: %v", util.GenJsonLog(resp))
}
