package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/sabriboughanmi/go_utils/utils"
	"sync"
)

type storageAction int

const (
	storageRenameType storageAction = 0
	storageDeleteType storageAction = 1
	storageAddType    storageAction = 2
	storageMoveType   storageAction = 3
)

type operation struct {
	ActionType storageAction
	MetaData   interface{}
}

type addOperationMetaData struct {
	bucket       string
	fileName     string
	localPath    string
	fileMetaData map[string]string
}

type deleteOperationMetaData struct {
	bucket string
	name   string
}

type renameOperationMetaData struct {
	srcBucket string

	srcName string
	dstName string
}
type moveOperationMetaData struct {
	srcBucket string
	dstBucket string
	srcName   string
	dstName   string
}

type storageBatch struct {
	client     *storage.Client
	operations *[]operation
}

// Batch returns a storageBatch.
func Batch(storageClient *storage.Client) *storageBatch {
	return &storageBatch{client: storageClient}
}

//Add appends a storage file creation from local to the batch.
func (wb *storageBatch) Add(bucket string, fileName string, localPath string, fileMetaData map[string]string) {
	*wb.operations = append(*wb.operations, operation{
		ActionType: storageAddType,
		MetaData: addOperationMetaData{
			bucket:       bucket,
			fileName:     fileName,
			localPath:    localPath,
			fileMetaData: fileMetaData,
		},
	})
}

//Rename appends a Rename file operation to the batch.
func (wb *storageBatch) Rename(srcBucket string, srcName string, dstName string) {
	*wb.operations = append(*wb.operations, operation{
		ActionType: storageRenameType,
		MetaData: renameOperationMetaData{
			srcBucket: srcBucket,

			srcName: srcName,
			dstName: dstName,
		},
	})
}

//Move appends a move file operation to the batch.
func (wb *storageBatch) Move(srcBucket string, dstBucket string, srcName string, dstName string) {
	*wb.operations = append(*wb.operations, operation{
		ActionType: storageMoveType,
		MetaData: moveOperationMetaData{
			srcBucket: srcBucket,
			dstBucket: dstBucket,
			srcName:   srcName,
			dstName:   dstName,
		},
	})
}

//Delete is used to append an operation in which we can delete a specific file
func (wb *storageBatch) Delete(srcBucket string, name string) {
	*wb.operations = append(*wb.operations, operation{
		ActionType: storageDeleteType,
		MetaData: deleteOperationMetaData{
			bucket: srcBucket,
			name:   name,
		},
	})
}

//Commit schedules batched operations in separate goroutines
func (wb *storageBatch) Commit(ctx context.Context) error {
	//Prevent calling goroutines if no operations are cached.
	if wb.operations == nil || len(*wb.operations) == 0 {
		return nil
	}

	errorChannel := make(chan error, len(*wb.operations))
	var wg sync.WaitGroup
	for _, operation := range *wb.operations {
		switch operation.ActionType {
		case storageDeleteType:
			metadata, _ := operation.MetaData.(deleteOperationMetaData)
			wg.Add(1)
			go func(waitGroup *sync.WaitGroup, errorChan chan error) {
				defer wg.Done()
				if err := RemoveFile(metadata.bucket, metadata.name, wb.client, ctx); err != nil {
					errorChan <- err
					return
				}
			}(&wg, errorChannel)

			break
		case storageRenameType:
			metadata, _ := operation.MetaData.(renameOperationMetaData)
			wg.Add(1)
			go func(waitGroup *sync.WaitGroup, errorChan chan error) {
				defer wg.Done()
				if err := RenameFile(metadata.srcBucket, metadata.dstName, metadata.srcName, wb.client, ctx); err != nil {
					errorChan <- err
					return
				}
			}(&wg, errorChannel)

			break
		case storageMoveType:
			metadata, _ := operation.MetaData.(moveOperationMetaData)
			wg.Add(1)
			go func(waitGroup *sync.WaitGroup, errorChan chan error) {
				defer wg.Done()
				if err := MoveFile(metadata.srcBucket, metadata.dstBucket, metadata.srcName, metadata.dstName, wb.client, ctx); err != nil {
					errorChan <- err
					return
				}
			}(&wg, errorChannel)

			break
		case storageAddType:
			metadata, _ := operation.MetaData.(addOperationMetaData)
			go func(waitGroup *sync.WaitGroup, errorChan chan error) {
				defer wg.Done()
				if err := CreateStorageFileFromLocal(metadata.bucket, metadata.fileName, metadata.localPath, metadata.fileMetaData, wb.client, ctx); err != nil {
					errorChan <- err
					return
				}

			}(&wg, errorChannel)

			break
		}
	}
	wg.Wait()
	var receivedErrors []string
	func() {
		//select
		for {
			select {
			case x, closed := <-errorChannel:
				if closed {
					return
				}
				receivedErrors = append(receivedErrors, x.Error())
				break
			default:
				return

			}
		}
	}()
	if len(receivedErrors) > 0 {
		return fmt.Errorf("Got %d Errors while commit - Errors : %s  \n", len(receivedErrors), string(utils.UnsafeAnythingToJSON(receivedErrors)))
	}
	return nil

}
