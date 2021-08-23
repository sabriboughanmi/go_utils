package ffmpeg

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/option"
	"sync"
	"testing"
)

var storageClientOnce sync.Once
var StorageClient *storage.Client
var ctx context.Context

func init() {
	ctx = context.Background()
}

//GetStorageClient returns a Singleton *storage.Client
func GetStorageClient() (*storage.Client, error) {
	var err error
	storageClientOnce.Do(func() {
		opt := option.WithCredentialsFile("./../../private_data/serviceAccountKey.json")
		// Pre-declare an err variable to avoid shadowing client.
		StorageClient, err = storage.NewClient(ctx, opt)
	})
	return StorageClient, err
}
func TestModerateVideo(t *testing.T) {
	storageClient, err := GetStorageClient()
	if err != nil {
		t.Errorf("Error - %v", err)
	}
	var temporaryStorageObject temporaryStorageObjectRef
	temporaryStorageObject.Bucket = "gs://tested4you-dev.appspot.com/"
	temporaryStorageObject.Client = storageClient
	/*
		vid, err := ffmpeg.LoadVideo("https://www.youtube.com/watch?v=w01V5FI03MQ")
		if err != nil {
				fmt.Printf("error")
		}
	*/

}
