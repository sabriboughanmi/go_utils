package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//unmarshalLanguagesJsons loads languages from the Json.
func unmarshalLanguagesJsons() {
	if err := json.Unmarshal([]byte(languagesJson), &languagesJsonModel); err != nil {
		panic(err)
	}
}

//Must be called to load Languages Json files and stores them to the languagesJson variable.
func loadJsonFiles(directory string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	var jsonFiles = make(map[ELanguageCode]string)
	for _, f := range files {
		if strings.Contains(f.Name(), ".json") {
			//Load file bytes
			bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", directory, f.Name()))
			if err != nil {
				return err
			}
			//Set the file bytes as string
			jsonFiles[ELanguageCode(strings.Replace(f.Name(), ".json", "", 1))] = string(bytes)
		}
	}

	languages_json_script := "package languages " +
		"\n var languagesJson string = `%s`"

	mapB, _ := json.Marshal(jsonFiles)

	goFile := []byte(fmt.Sprintf(languages_json_script, string(mapB)))

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	if err = ioutil.WriteFile(path+"/languages_json.go", goFile, 0666); err != nil {
		return err
	}

	return nil

}
