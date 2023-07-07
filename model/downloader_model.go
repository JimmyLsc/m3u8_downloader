package model

type DownloaderReq struct {
	SrcURL       string  `json:"src_url,required"`
	SrcName      string  `json:"src_name,required"`
	SrcShortName *string `json:"src_short_name,omitempty"`
	CachePath    *string `json:"cache_path,omitempty"`
	FilePath     *string `json:"file_path,omitempty"`
	WorkerNum    *int64  `json:"worker_num,omitempty"`
}

type DownloaderResp struct {
	BaseResp *BaseResp `json:"base_resp,required"`
}
