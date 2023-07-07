package logic

import (
	"context"
	"github.com/JimmyLsc/m3u8_downloader/model"
	"github.com/JimmyLsc/m3u8_downloader/util"
	log "github.com/sirupsen/logrus"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func GetM3U8Info(ctx context.Context, srcURL string) (*model.M3U8Info, error) {
	header := make(map[string]string)
	resp, err := util.HttpGet(ctx, srcURL, header)
	if err != nil {
		log.Errorf("GetM3U8Info error, err: %v", err)
		return nil, err
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("GetM3U8Info error, err: %v", err)
		return nil, err
	}
	info, err := parseM3U8(ctx, srcURL, string(respBytes))
	if err != nil {
		log.Errorf("GetM3U8Info error, err: %v", err)
		return nil, err
	}
	return info, nil
}

func parseM3U8(ctx context.Context, srcUrl, content string) (*model.M3U8Info, error) {
	info := &model.M3U8Info{}

	// Parse root path
	paths := strings.Split(srcUrl, "/")
	info.RootURL = strings.Join(paths[:len(paths)-1], "/")

	// Parse tsURL
	tsRe, err := regexp.Compile(`.+\.ts`)
	if err != nil {
		log.Errorf("parseM3U8 error, err: %v", err)
		return nil, err
	}
	tsNameSlice := tsRe.FindAllString(content, -1)
	tsURLs := tsNameSlice
	if !strings.Contains(tsNameSlice[0], "http") {
		for index, value := range tsNameSlice {
			tsNameSlice[index] = info.RootURL + "/" + value
		}
	}
	info.TsURLs = tsURLs

	// Parse ts length
	tsLengthRe, err := regexp.Compile(`#EXTINF:.*,`)
	if err != nil {
		log.Errorf("parseM3U8 error, err: %v", err)
		return nil, err
	}
	tsLengthSlice := tsLengthRe.FindAllString(content, -1)
	var tsLengths []float64
	for _, value := range tsLengthSlice {
		length, err := strconv.ParseFloat(value[len("#EXTINF:"):len(value)-1], 64)
		if err != nil {
			log.Errorf("parseM3U8 error, err: %v", err)
			return nil, err
		}
		tsLengths = append(tsLengths, length)
	}
	info.TsLengths = tsLengths
	return info, nil
}

func DownloadVideo(ctx context.Context, info *model.M3U8Info, cachePath, videoPath string) error {

	return nil
}
