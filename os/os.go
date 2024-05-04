package os

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//CreateTempDirectory creates a Directory in Temp path.
//This is Required cause Cloud Functions do not Create Directory Folder for us
func CreateTempDirectory(directory string) error {
	return os.Mkdir(directory, os.ModePerm)
}

//RemovePathsIfExists.
func RemovePathsIfExists(paths ...string) {
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				continue
			}
		}
		os.Remove(path)
	}
	return
}

//RemovePathIfExists .
func RemovePathIfExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	os.Remove(path)
	return true
}

// PathExists reports whether the named file or directory exists.
func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//CreateTempFile creates a file in the temp/ directory.
//If fileName contains Folders it will use the last element of path.
//note! it's required to call defer os.Remove() on the returned path if not empty to ensure the file is cleaned up
func CreateTempFile(fileName string, content []byte) (string, error) {

	fn := filepath.Base(fileName)

	tmpFile, err := ioutil.TempFile("", "*"+fn)
	if err != nil {
		return "", fmt.Errorf("err Creating TempFile %v", err)
	}

	if content != nil {
		if _, err := tmpFile.Write(content); err != nil {
			return tmpFile.Name(), err
		}
	} else {
		if _, err := tmpFile.Write([]byte{}); err != nil {
			return tmpFile.Name(), err
		}
	}

	if err := tmpFile.Close(); err != nil {
		return tmpFile.Name(), err
	}

	return tmpFile.Name(), nil
}
