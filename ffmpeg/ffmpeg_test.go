package ffmpeg

import (
	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	"google.golang.org/api/option"
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

	vid, err := LoadVideo("C:/Users/T4ULabs/Downloads/vd.mp4")
	if err != nil {
		t.Errorf("Error load video  - %v", err)
	}

	opt := option.WithCredentialsFile(serviceAccountPath)
	// Pre-declare an err variable to avoid shadowing client.
	AnnotationClient, err := vision.NewImageAnnotatorClient(ctx, opt)

	defer duration(track("\nModeration Took :"))

	err, ok := vid.ModerateVideo(5, ctx, 3, AnnotationClient)
	if err != nil {
		t.Errorf("Error moderate video  - %v", err)
	}

	fmt.Printf("testting %v", ok)
}

func TestLoadVideoFromReEncodedFragments(t *testing.T) {

	video, err := LoadVideoFromReEncodedFragments("C:\\Users\\Sabri\\Downloads\\Video\\output.mp4", false,
		"C:\\Users\\Sabri\\Downloads\\Video\\1.mp4", "C:\\Users\\Sabri\\Downloads\\Video\\2.mp4")

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Video duration: %v", video.GetDuration())
}
