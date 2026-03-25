package dispatcher

import (
	"os"
)

func (d *dispatcher) DeleteROT(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		d.log.Fatal("Error deleting a dir after uploading", err)
	}
}
