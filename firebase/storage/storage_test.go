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

func TestDeleteFolder(t *testing.T) {
	appStorage,err := GetStorageClient()
	if err != nil {
		t.Error(fmt.Errorf("error getting storage %v" ,err))
	}

	if err =DeleteFolder("tested4you-dev.appspot.com","pages/testdelete/",appStorage,ctx);err != nil {
		t.Error(fmt.Errorf("error getting storage %v" ,err))
	}
}


