package logic

import (
	"context"
	"github.com/JimmyLsc/m3u8_downloader/model"
	"github.com/JimmyLsc/m3u8_downloader/util"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GetM3U8Info(ctx context.Context, srcURL string, header map[string]string) (*model.M3U8Info, error) {
	resp, err := util.GetHttpClient().HttpGet(ctx, srcURL, header)
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
	lines := strings.Split(content, "\n")
	var tsNameSlice []string
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			tsNameSlice = append(tsNameSlice, line)
		}
	}
	tsURLs := tsNameSlice
	if !strings.Contains(tsNameSlice[0], "http") {
		for index, value := range tsNameSlice {
			tsNameSlice[index] = info.RootURL + "/" + value
		}
	}
	info.TsURLs = tsURLs[:len(tsURLs)-1]

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
	info.TsLengths = tsLengths[:len(tsLengths)-1]

	// Parse Encryption
	if err = parseEncryption(ctx, info, content); err != nil {
		log.Errorf("parseM3U8 error, err: %v", err)
		return nil, err
	}

	return info, nil
}

func parseEncryption(ctx context.Context, info *model.M3U8Info, content string) error {
	encryptRe, err := regexp.Compile(`#EXT-X-KEY:.*`)
	if err != nil {
		log.Errorf("parseM3U8 error, err: %v", err)
		return err
	}
	encryptStr := encryptRe.FindString(content)
	if len(encryptStr) != 0 {
		encryptInfoStr := encryptStr[len("#EXT-X-KEY:"):]
		encryptMap := make(map[string]string)
		for _, kv := range strings.Split(encryptInfoStr, ",") {
			encryptMap[strings.Split(kv, "=")[0]] = strings.Split(kv, "=")[1]
		}
		info.EncryptionInfo = &model.EncryptionInfo{}
		if _, ok := encryptMap["METHOD"]; ok {
			info.EncryptionInfo.Method = strings.Trim(encryptMap["METHOD"], "\"")
		}
		if _, ok := encryptMap["URI"]; ok {
			info.EncryptionInfo.URI = strings.Trim(encryptMap["URI"], "\"")
		}
		if _, ok := encryptMap["IV"]; ok {
			info.EncryptionInfo.IV = strings.Trim(encryptMap["IV"], "\"")
		}
		if _, ok := encryptMap["KEY"]; ok {
			info.EncryptionInfo.Key = strings.Trim(encryptMap["KEY"], "\"")
		}
		if len(info.EncryptionInfo.Key) == 0 {
			resp, err := util.GetHttpClient().HttpGet(ctx, info.RootURL+"/"+info.EncryptionInfo.URI, make(map[string]string))
			if err != nil {
				log.Errorf("parseM3U8 error, err: %v", err)
				return err
			}
			key, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("parseM3U8 error, err: %v", err)
				return err
			}
			info.EncryptionInfo.Key = string(key)
		}
	}
	return nil
}

func DownloadVideo(ctx context.Context, info *model.M3U8Info, cachePath, videoPath string, workNum int64) error {
	cacheDir := filepath.Join(cachePath, info.ShortName)
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		log.Errorf("DownloadVideo error, err: %v", err)
		return err
	}
	finishChan := make(chan int, workNum)
	if err := MDownloadsTS(ctx, info, cacheDir, workNum, finishChan); err != nil {
		log.Errorf("DownloadVideo error, err: %v", err)
		return err
	}
	if err := os.MkdirAll(videoPath, os.ModePerm); err != nil {
		log.Errorf("DownloadVideo error, err: %v", err)
		return err
	}
	if err := CombineTS(ctx, info, cacheDir, filepath.Join(videoPath, info.Name+".mp4")); err != nil {
		log.Errorf("DownloadVideo error, err: %v", err)
		return err
	}
	for i := 0; i < int(workNum); i++ {
		<-finishChan
	}
	if err := os.RemoveAll(cacheDir); err != nil {
		log.Errorf("DownloadVideo error, err: %v", err)
		return err
	}
	return nil
}
