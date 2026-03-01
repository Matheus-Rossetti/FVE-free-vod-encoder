package uploader

import (
	"log"
	"os"
)

func DeleteROT(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		log.Fatal("Error deleting a dir after uploading", err)
	}
}
