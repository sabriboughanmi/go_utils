package go_utils

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
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
		opt := option.WithCredentialsFile("private_data/serviceAccountKey.json")
		FirebaseApp, err = firebase.NewApp(ctx, nil, opt)
	})
	return FirebaseApp, err
}

//GetStorageClient returns a Singleton *storage.Client
func GetStorageClient() (*storage.Client, error) {
	var err error
	storageClientOnce.Do(func() {
		var newFirebaseApp *firebase.App
		newFirebaseApp, err = GetFirebaseClient()
		if err != nil {
			return
		}

		// Pre-declare an err variable to avoid shadowing client.
		StorageClient, err = newFirebaseApp.Storage(ctx)
	})
	return StorageClient, err
}

func TestMoveFile(t *testing.T) {

	_, err := GetStorageClient()
	if err != nil {
		t.Errorf("Error - %v", err)
	}

}
