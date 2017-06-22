package main

import (
	"github.com/otium/ytdl"
)

// GetVideoDownloadURL returns download url, title and error
func GetVideoDownloadURL(url string) (string, string, error) {
	info, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return "", "", err
	}

	title := info.Title
	downloadURL, err := info.GetDownloadURL(info.Formats.Best(ytdl.FormatAudioEncodingKey)[0])

	return downloadURL.String(), title, nil
}
