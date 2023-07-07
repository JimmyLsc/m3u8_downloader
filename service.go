package main

type DownloaderService interface {
	Download() error
}
