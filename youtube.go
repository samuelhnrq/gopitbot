package main

import (
	"strings"

	"github.com/otium/ytdl"
	"github.com/pkg/errors"
)

// GetVideoDownloadURL returns download url, title and error
func GetVideoDownloadURL(url string) (string, string, error) {
	info, err := ytdl.GetVideoInfo(url)
	if err != nil || !strings.Contains(url, "youtube.com/watch?v=") {
		return "", "", errors.New("Invalid URL")
	}
	if len(info.Formats) == 0 {
		return "", "", errors.New("Invalid URL")
	}
	title := info.Title
	downloadURL, err := info.GetDownloadURL(info.Formats.Best(ytdl.FormatAudioEncodingKey)[0])
	if err != nil {
		return "", "", err
	}
	return downloadURL.String(), title, nil
}
