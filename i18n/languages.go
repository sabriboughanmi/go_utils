package i18n

import "fmt"

var (
	languagesJsonModel map[ELanguageCode]string = nil
	packageInitialized                          = false
)

//InitPackageForCloudFunctions Initialize the Package, so it can be used in Firebase Cloud Functions
func InitPackageForCloudFunctions() {
	if !packageInitialized {
		unmarshalLanguagesJsons()
	}
}

//GetLanguage returns a languageCode file as []byte and the json file format.
func GetLanguage(languageCode ELanguageCode) ([]byte, string) {

	//Return the Correct Language file as []byte
	if bytes, ok := languagesJsonModel[languageCode]; ok {
		return []byte(bytes), fmt.Sprintf("%s.json", languageCode)
	}

	//Return the English version of file as []byte
	if bytes, ok := languagesJsonModel[LanguageCode_En_Us]; ok {
		return []byte(bytes), fmt.Sprintf("%s.json", languageCode)
	}

	//No translations exists
	return nil, ""
}

//GetSupportedLanguages returns an [] ELanguageCode .
func GetSupportedLanguages() []ELanguageCode {
	var supportedLanguages = make([]ELanguageCode, len(languagesJsonModel))
	var i = 0
	for sl, _ := range languagesJsonModel {
		supportedLanguages[i] = sl
		i++
	}
	return supportedLanguages
}
