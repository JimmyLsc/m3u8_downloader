package util

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

var httpClientOnce sync.Once
var httpClient *http.Client

func InitHttpClient() {
	httpClientOnce.Do(func() {
		httpClient = &http.Client{}
	})
}

func HttpGet(ctx context.Context, url string, headerMap map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("HttpGet error, err: %v", err)
		return nil, err
	}
	for key, value := range headerMap {
		req.Header.Add(key, value)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Errorf("HttpGet error, err: %v", err)
		return nil, err
	}
	return resp, nil
}
