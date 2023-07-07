package m3u8_downloader

type DownloaderService interface {
	Download() error
}
