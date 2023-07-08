package model

type M3U8Info struct {
	RootURL        string          `json:"root_url,required"`
	TsURLs         []string        `json:"ts_ur_ls,required"`
	TsLengths      []float64       `json:"ts_lengths,required"`
	Name           string          `json:"name,required"`
	ShortName      string          `json:"short_name,required"`
	EncryptionInfo *EncryptionInfo `json:"encryption_info,omitempty"`
}

type EncryptionInfo struct {
	Method string `json:"method,required"`
	Key    string `json:"key,required"`
	URI    string `json:"uri,required"`
	IV     string `json:"iv,required"`
}
