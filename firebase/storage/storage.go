package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
)

//StorageFileContentType represents a file Content type in Cloud Storage in order to be recognized by the browser
type FileContentType string

const (
	ImageGif  FileContentType = "image/gif"
	ImageJPEG FileContentType = "image/jpeg"

	VideoMP4 FileContentType = "video/mp4"
	VideoMOV FileContentType = "video/mov"
	VideoAVI FileContentType = "video/avi"
)

/*
func GetDownloadURL(bucket, storagePath string) string {

	staticUrl := "https://firebasestorage.googleapis.com/v0/b/{Bucket}/o/{FilePath}?alt=media&token={AccessToken}"

	keyValues := make(map[string]string)
	keyValues["Bucket"] = bucket
	keyValues["AccessToken"] = ""
	keyValues["FilePath"] = strings.ReplaceAll(storagePath, "/", "%2F")

	return ReplaceKeys(staticUrl, keyValues)
}*/

//FileExists Checks if a Storage File Exists
func FileExists(bucket, storagePath string, client *storage.Client, ctx context.Context) (bool, error) {

	bucketHandle := client.Bucket(bucket)
	objectHandle := bucketHandle.Object(storagePath)
	if _, err := objectHandle.Attrs(ctx); err != nil {
		return false, nil
	}

	return true, nil
}

//LoadFileInTempPath Loads a file from Storage to an OS path
//Note! it's the caller responsibility to Remove to defer os.Remove() on the returned path if not empty to ensure the file is cleaned up
func LoadFileInTempPath(bucket, storagePath string, client *storage.Client, ctx context.Context) (string, error) {

	bucketHandle := client.Bucket(bucket)
	objectHandle := bucketHandle.Object(storagePath)
	_, err := objectHandle.Attrs(ctx)
	if err != nil {
		return "", fmt.Errorf("StorageLoadFileAtPath :  Error Loading File from Bucket : %s, File %s: with Error :  %v", bucket, storagePath, err)
	}

	reader, err := objectHandle.NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("StorageLoadFileAtPath : Error Creating Reader for objectHandle %v", err)
	}
	defer reader.Close() // Clean up

	fileContent, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("unable to read data from bucket %s, file %s: %v", bucket, storagePath, err)
	}

	fn := filepath.Base(storagePath)

	tmpFile, err := ioutil.TempFile("", "*"+fn)
	if err != nil {
		return "", fmt.Errorf("err Creating TempFile %v", err)
	}

	if err = ioutil.WriteFile(tmpFile.Name(), fileContent, os.ModePerm); err != nil {
		return tmpFile.Name(), fmt.Errorf("StorageLoadFileAtPath : Error Writing to File, Bucket: %s, FilePath: %s, Error:  %v", bucket, storagePath, err)
	}

	return tmpFile.Name(), nil
}

//LoadFileAtPath Loads a file from Storage to an OS path
func LoadFileAtPath(bucket, storagePath, dstPath string, client *storage.Client, ctx context.Context) error {

	bucketHandle := client.Bucket(bucket)
	objectHandle := bucketHandle.Object(storagePath)
	_, err := objectHandle.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("StorageLoadFileAtPath :  Error Loading File from Bucket : %s, File %s: with Error :  %v", bucket, storagePath, err)
	}

	reader, err := objectHandle.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("StorageLoadFileAtPath : Error Creating Reader for objectHandle %v", err)
	}
	defer reader.Close() // Clean up

	fileContent, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("unable to read data from bucket %s, file %s: %v", bucket, storagePath, err)
	}

	if err = ioutil.WriteFile(dstPath, fileContent, os.ModePerm); err != nil {
		return fmt.Errorf("StorageLoadFileAtPath : Error Writing to File, Bucket: %s, FilePath: %s, Error:  %v", bucket, storagePath, err)
	}

	return nil
}

// MoveFile moves an object into another location.
func MoveFile(bucket, srcName, dstName string, client *storage.Client, ctx context.Context) error {
	src := client.Bucket(bucket).Object(srcName)
	dst := client.Bucket(bucket).Object(dstName)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %v", dstName, srcName, err)
	}
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", srcName, err)
	}
	return nil
}

// CreateFile creates a file in Google Cloud Storage.
func CreateFile(bucket, fileName string, content []byte, contentType FileContentType, fileMetaData map[string]string, client *storage.Client, ctx context.Context) error {
	wc := client.Bucket(bucket).Object(fileName).NewWriter(ctx)
	defer wc.Close()

	wc.ContentType = string(contentType)
	if fileMetaData != nil {
		wc.Metadata = fileMetaData
	} else {
		wc.Metadata = make(map[string]string)
	}

	if _, err := wc.Write(content); err != nil {
		return fmt.Errorf("createFile: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
	}
	return nil
}
