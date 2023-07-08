package logic

import (
	"context"
	"github.com/JimmyLsc/m3u8_downloader/model"
	"github.com/JimmyLsc/m3u8_downloader/util"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type urlItem struct {
	Index int
	URL   string
}

func MDownloadsTS(ctx context.Context, info *model.M3U8Info, cachePath string, workerNum int64) error {
	workChan := make(chan urlItem, len(info.TsURLs))
	done := make(chan int, workerNum)
	defer func() {
		close(workChan)
		close(done)
	}()
	var wg sync.WaitGroup
	for index, URL := range info.TsURLs {
		workChan <- urlItem{
			Index: index,
			URL:   URL,
		}
		wg.Add(1)
	}
	log.Info(len(info.TsURLs))
	var err error
	for i := 0; i < int(workerNum); i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				case item := <-workChan:
					err = DownloadsTS(ctx, item.URL, cachePath, strconv.Itoa(item.Index)+".ts", info.EncryptionInfo, nil)
					if err != nil {
						log.Errorf("DownloadTS error, err:%v", err)
					}
					wg.Done()
				}
			}
		}()
	}
	wg.Wait()
	for i := 0; i < int(workerNum); i++ {
		done <- 1
	}

	return err
}

func DownloadsTS(ctx context.Context, url, cachePath, tsName string, encrypt *model.EncryptionInfo, header map[string]string) error {
	resp, err := util.HttpGet(ctx, url, header)
	if err != nil {
		log.Errorf("DownloadTS error, err:%v", err)
		return err
	}
	if err != nil {
		log.Errorf("DownloadTS error, err:%v", err)
		return err
	}
	cacheURL := filepath.Join(cachePath, tsName)
	file, err := os.Create(cacheURL)
	if err != nil {
		log.Errorf("DownloadTS error, err:%v", err)
		return err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("DownloadTS error, err:%v", err)
		return err
	}
	originData := data
	if encrypt != nil {
		originData, err = util.AES128Decrypt(data, []byte(encrypt.Key), []byte(encrypt.IV))
		if err != nil {
			log.Errorf("DownloadTS error, err:%v", err)
			return err
		}
	}
	_, err = file.Write(originData)
	if err != nil {
		log.Errorf("DownloadTS error, err:%v", err)
		return err
	}
	return nil
}

func CombineTS(ctx context.Context, info *model.M3U8Info, cachePath string, videoURL string) error {
	log.Info("Start Combine")
	outFile, err := os.Create(videoURL)
	if err != nil {
		log.Errorf("CombineTS error, err:%v", err)
		return err
	}
	for index := 0; index < len(info.TsURLs); index++ {
		tsURL := filepath.Join(cachePath, strconv.Itoa(index)+".ts")
		tsFile, err := os.Open(tsURL)
		if err != nil {
			log.Errorf("CombineTS error, err:%v", err)
			return err
		}
		data, err := io.ReadAll(tsFile)
		if err != nil {
			log.Errorf("CombineTS error, err:%v", err)
			return err
		}
		_, err = outFile.Write(data)
		if err != nil {
			log.Errorf("CombineTS error, err:%v", err)
			return err
		}
	}
	return nil
}