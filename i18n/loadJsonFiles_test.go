package i18n

import (
	"log"
	"os"
	"testing"
)

func Test_LoadJsonFiles(t *testing.T) {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	if err = loadJsonFiles(path + "/json files"); err != nil {
		t.Error(err)
	}
}
