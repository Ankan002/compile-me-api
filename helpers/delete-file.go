package helpers

import (
	"log"
	"os"
)

func DeleteFile(fileName string) {
	fileDeletionError := os.Remove(fileName)

	if fileDeletionError != nil {
		log.Println("Deletion process for file " + fileName + " failed...")
	}
}
