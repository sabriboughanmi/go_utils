package storage

import (
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"sync"
	"testing"
)

var ctx context.Context

func init() {
	ctx = context.Background()
}

var storageClientOnce sync.Once
var StorageClient *storage.Client

var firebaseAppOnce sync.Once
var FirebaseApp *firebase.App

//GetFirebaseClient returns a Singleton *firebase.App
func GetFirebaseClient() (*firebase.App, error) {
	var err error
	firebaseAppOnce.Do(func() {
		opt := option.WithCredentialsFile("./../../private_data/serviceAccountKey.json")
		FirebaseApp, err = firebase.NewApp(ctx, nil, opt)
	})
	return FirebaseApp, err
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

func TestMoveFile(t *testing.T) {
	var srcBucket = "tested4you-dev.appspot.com"
	var dstBucket = "tested4you-dev.appspot.com"
	var srcStoragePath = "pages/test/img12.jpg"
	var dstStoragePath = "pages/nex test2/img12.jpg"

	storageClient, err := GetStorageClient()
	if err != nil {
		t.Errorf("Error - %v", err)
	}
	storageBatch := Batch(storageClient)

	storageBatch.Move(srcBucket, dstBucket, srcStoragePath, dstStoragePath, func(err error) {
		fmt.Printf("Eh Eh eni n9oul fama erreer %v", err)
	})

	if err = storageBatch.Commit(ctx); err != nil {
		t.Errorf("MoveFile Error - %v", err)
	}

}
