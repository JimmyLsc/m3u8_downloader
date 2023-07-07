package m3u8_downloader

import (
	"context"
	"github.com/JimmyLsc/m3u8_downloader/model"
	"github.com/JimmyLsc/m3u8_downloader/util"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestDownloadServiceImpl(t *testing.T) {
	util.InitHttpClient()
	req := &model.DownloaderReq{
		SrcURL:  "https://m3u8.wolongcdnm3u8.com:65/20230628/20ce3a69/index.m3u8",
		SrcName: "test",
	}
	resp := (&DownloadServiceImpl{}).Download(context.TODO(), req)
	log.Infof("Response: %v", util.GenJsonLog(resp))
}
