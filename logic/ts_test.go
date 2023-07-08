package logic

import (
	"context"
	"github.com/JimmyLsc/m3u8_downloader/model"
	"github.com/JimmyLsc/m3u8_downloader/util"
	"io"
	"testing"
)

func TestDownloadsTS(t *testing.T) {
	util.InitHttpClient()
	resp, _ := util.HttpGet(context.TODO(), "https://asf-doc.mushroomtrack.com/hls/mCD8BgXztjAHy5tKhW8RWQ/1688744405/34000/34192/5ca35cd2f07bfaa2.ts", make(map[string]string))
	key, _ := io.ReadAll(resp.Body)
	DownloadsTS(context.TODO(), "https://asf-doc.mushroomtrack.com/hls/mCD8BgXztjAHy5tKhW8RWQ/1688744405/34000/34192/341921789.ts", "../.cache", "test.ts", &model.EncryptionInfo{Key: string(key), IV: "0x9e70a6ac4412808a8591e3b17d68940d"}, make(map[string]string))
}
