package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type DownloadJob struct {
	Source string // cli | rabbitmq | REST API | etc...
	Uri    string
}

// Returns the absolute path of the downloaded video
func DownloadAndStoreVideo(url string) string {

	fyleType := checkFileType(url)
	if fyleType != "video" {
		log.Fatal("Expected URL to download a video file. Got URL to download: ", fyleType)
	}

	os.Mkdir("downloaded-videos", 0700)
	file, err := os.CreateTemp("downloaded-videos", "video-*")
	if err != nil {
		log.Fatal("Failed when creating temp file to store video from download")
	}
	defer file.Close()

	request, err := http.Get(url)
	if err != nil {
		log.Fatal("Error during get request to download video", err)
	}
	defer request.Body.Close()

	_, err = io.Copy(file, request.Body)
	if err != nil {
		log.Fatal("Error storing body stream into file", err)
	}

	absoluteVideoPath, err := filepath.Abs(file.Name())
	if err != nil {
		log.Fatal("Error getting the absolute video path after downloading", err)
	}

	return absoluteVideoPath
}

// Downloads the first 512 bytes (or less) and check fileType
func checkFileType(url string) string {

	// Can't use http.Get because we need to modify Headers
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error downloading video ", err)
	}

	request.Header.Set("Range", "bytes=0-511")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Coudn't make request for the first 512 bytes", err)
	}
	defer response.Body.Close()

	// The code below garantees that if the server ignores the Range in the Header and
	// sends the full file, we cut the connection after downloading 512 bytes.
	// This happens becuase of how http.Response.Body works
	// Hover ".Body" and check the first few lines in the docs, its awesome.
	buffer := make([]byte, 512)
	bytesRead, err := io.ReadFull(response.Body, buffer)
	if err != nil {
		log.Fatal("Couldn't read bytes from response.Body into buffer", err)
	}

	// Need to pass :bytesRead in case less than 512 bytes were read into buffer
	MIMEType := http.DetectContentType(buffer[:bytesRead])

	fileType := strings.Split(MIMEType, "/")
	return fileType[0]
}
