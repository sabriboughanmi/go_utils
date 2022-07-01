package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//StorageFileContentType represents a file Content type in Cloud Storage in order to be recognized by the browser
type FileContentType string

const (
	ImageGif  FileContentType = "image/gif"
	ImageJPEG FileContentType = "image/jpeg"
	ImagePNG  FileContentType = "image/png"
	VideoMP4  FileContentType = "video/mp4"
	VideoMOV  FileContentType = "video/mov"
	VideoAVI  FileContentType = "video/avi"
	FilePDG  FileContentType = "file/pdf"

)



//FileExists Checks if a Storage File Exists
func FileExists(bucket, storagePath string, client *storage.Client, ctx context.Context) (bool, error) {

	bucketHandle := client.Bucket(bucket)
	objectHandle := bucketHandle.Object(storagePath)
	_, err := objectHandle.Attrs(ctx)
	if err == nil {
		return true, nil
	}

	if err == storage.ErrObjectNotExist {
		return false, nil
	}

	return false, err
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
func MoveFile(srcBucket, dstBucket, srcName, dstName string, client *storage.Client, ctx context.Context) error {
	src := client.Bucket(srcBucket).Object(srcName)
	dst := client.Bucket(dstBucket).Object(dstName)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %v", dstName, srcName, err)
	}
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", srcName, err)
	}
	return nil
}

// CopyFile copies an object into another location in Storage.
func CopyFile(srcBucket, dstBucket, srcName, dstName string, client *storage.Client, ctx context.Context) error {
	src := client.Bucket(srcBucket).Object(srcName)
	dst := client.Bucket(dstBucket).Object(dstName)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %v", dstName, srcName, err)
	}
	return nil
}

// RemoveFile Removes a file from Storage
func RemoveFile(bucket, name string, client *storage.Client, ctx context.Context) error {
	src := client.Bucket(bucket).Object(name)
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%s).Delete: %v", name, err)
	}
	return nil
}

// RemoveFilesFromBucket Removes  a list of files from Storage
func RemoveFilesFromBucket(client *storage.Client, ctx context.Context, bucket string, paths ...string) error {
	bucketHandle := client.Bucket(bucket)

	wg := sync.WaitGroup{}
	errChan := make(chan error)

	for _, storagePath := range paths {
		wg.Add(1)
		go func(errorChan chan error, waitGroup *sync.WaitGroup, path string) {
			defer waitGroup.Done()
			objHandle := bucketHandle.Object(path)
			if err := objHandle.Delete(ctx); err != nil {
				if err == storage.ErrObjectNotExist {
					fmt.Printf("Error Skipped: trying to Remove an unexisting File: %s - Error : %v \n", path, err)
					return
				}
				errorChan <- fmt.Errorf("Object(%s).Delete: %v", path, err)
				return
			}
		}(errChan, &wg, storagePath)
	}

	//Wait Removing Errors
	if err := HandleGoroutineErrors(&wg, errChan); err != nil {
		return err
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

// CreateStorageFileFromLocal creates a file in Google Cloud Storage from a Local file Path.
func CreateStorageFileFromLocal(bucket, fileName, localPath string, contentType FileContentType, fileMetaData map[string]string, client *storage.Client, ctx context.Context) (*storage.ObjectHandle, error) {
	data, err := ioutil.ReadFile(localPath)
	if err != nil {
		return nil, err
	}
	var objectHandle = client.Bucket(bucket).Object(fileName)
	wc := objectHandle.NewWriter(ctx)
	defer wc.Close()
	if string(contentType) != "" {
		wc.ContentType = string(contentType)
	}
	//defer
	if fileMetaData != nil {
		wc.Metadata = fileMetaData
	} else {
		wc.Metadata = make(map[string]string)
	}

	if _, err := wc.Write(data); err != nil {
		return nil, fmt.Errorf("CreateStorageFileFromLocal: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
	}
	return objectHandle, nil
}

//FilesInFolder Lists all Files under a folder prefix.
//
//Note!
//
//prefix must finish with "/".
//
//delimiter must be "/" otherwise all files in SubFolders will be listed.
func FilesInFolder(bucket, prefix, delimiter string, client *storage.Client, ctx context.Context) (*[]*storage.ObjectAttrs, error) {

	bucketHandle := client.Bucket(bucket)
	it := bucketHandle.Objects(ctx, &storage.Query{
		Delimiter: delimiter,
		Prefix:    prefix,
	})

	var files []*storage.ObjectAttrs
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket(%q).Objects(): %v", bucket, err)
		}
		files = append(files, attrs)
	}
	return &files, nil
}

// RenameFile rename storage file.
func RenameFile(srcBucket, srcName, dstName string, client *storage.Client, ctx context.Context) error {
	src := client.Bucket(srcBucket).Object(srcName)
	dst := client.Bucket(srcBucket).Object(dstName)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return fmt.Errorf("Object(%s).CopierFrom(%s).Run: %v", dstName, srcName, err)
	}
	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%s).Delete: %v", srcName, err)
	}
	return nil
}

// GeneratePublicUrl generates a storage object signed URL with GET method.
//Note! if the expiresDateTime is not assigned a 15minute expiration will be applied.
//
//the expiration may be no more than seven days in the future.
func GeneratePublicUrl(bucket, storageObject, serviceAccountPrivateKey, serviceAccountEmail string, expirationDateTime *time.Time) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: serviceAccountEmail,
		PrivateKey:     []byte(serviceAccountPrivateKey),
		Expires: func() time.Time {
			if expirationDateTime != nil {
				return *expirationDateTime
			}
			return time.Now().Add(15 * time.Minute)
		}(),
	}
	u, err := storage.SignedURL(bucket, storageObject, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %v", err)
	}
	return u, nil
}


//DeleteFolder Deletes all files within a Folder
func DeleteFolder(bucket, folderPath string, client *storage.Client, ctx context.Context) error{
	files,err := FilesInFolder(bucket,folderPath,"",client,ctx)
	if err != nil {
		return fmt.Errorf("error while loading files in folder path %v" ,err)
	}

	if files == nil{
		return err
	}

	batch := Batch(client)
	for _, file := range *files {
		fmt.Printf("file found : %s\n",file.Name )
		batch.Delete(bucket, file.Name,func (failError error){
			fmt.Printf("error while deleting file , filename: %v , error: %v\n" ,file.Name,failError)
		})
	}

	return batch.Commit(ctx)
}

