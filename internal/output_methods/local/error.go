package local

import (
	"fmt"
	"os"
)

func (l *Local) HandleError(files []string) {

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Printf("\nFailed to delete %v", file)
			continue
		}
	}

}
