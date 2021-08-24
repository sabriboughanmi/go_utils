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
		opt := option.WithCredentialsFile("./../private_data/serviceAccountKey.json")
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

	var temporaryStorageObject =  GetTemporaryStorageObjectRef(storageClient, "gs://tested4you-dev.appspot.com/")
	vid , err :=  LoadVideo("C:/Users/T4ULabs/Downloads/vd.mp4")
	if  err != nil{
		t.Errorf("Error load video  - %v", err)
	}
		err = vid.ModerateVideo(5, ctx, 3, &temporaryStorageObject) ;	if err != nil {
		t.Errorf("Error moderate video  - %v", err)
	}
}
