package ffmpeg

import (
	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io/ioutil"
	"sync"
	"testing"
	"time"
)

var storageClientOnce sync.Once
var StorageClient *storage.Client
var ctx context.Context

func init() {
	ctx = context.Background()
}

const (
	serviceAccountPath = "./../private_data/serviceAccountKey.json"
)

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

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	fmt.Printf("%v: %v\n", msg, time.Since(start))
}

func TestModerateVideo(t *testing.T) {
	storageClient, err := GetStorageClient()
	if err != nil {
		t.Errorf("Error - %v", err)
	}

	jsonKey, err := ioutil.ReadFile(serviceAccountPath)
	if err != nil {
		t.Errorf("ioutil.ReadFile: %v", err)
		return
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		t.Errorf("google.JWTConfigFromJSON: %v", err)
		return
	}

	var temporaryStorageObject = GetModerateVideoMetadata(storageClient, "tested4you-dev.appspot.com", string(conf.PrivateKey), conf.Email)
	vid, err := LoadVideo("C:/Users/sabri/Downloads/vd.mp4")
	if err != nil {
		t.Errorf("Error load video  - %v", err)
	}

	opt := option.WithCredentialsFile(serviceAccountPath)
	// Pre-declare an err variable to avoid shadowing client.
	AnnotationClient, err := vision.NewImageAnnotatorClient(ctx, opt)

	defer duration(track("\nModeration Took :"))

	err, ok := vid.ModerateVideo(5, ctx, 3, &temporaryStorageObject, AnnotationClient)
	if err != nil {
		t.Errorf("Error moderate video  - %v", err)
	}

	fmt.Printf("testting %v", ok)
}
