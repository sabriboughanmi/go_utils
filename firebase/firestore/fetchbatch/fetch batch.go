package fetchbatch

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sabriboughanmi/go_utils/utils"
	"sync"
)

//GetFirestoreFetchBatch returns a firestore fetch batch
func GetFirestoreFetchBatch(Client *firestore.Client, context context.Context) FirestoreFetchBatch {
	return FirestoreFetchBatch{
		CommandsQueue: nil,
		Client:        Client,
		Context:       context,
	}
}

//AddCommand adds a FetchCommand to the Queue
func (ffcq *FirestoreFetchBatch) AddCommand(command FetchCommand) {
	ffcq.CommandsQueue = append(ffcq.CommandsQueue, command)
}

//AddCommands adds multiple FetchCommand to the Queue
func (ffcq *FirestoreFetchBatch) AddCommands(commands ...FetchCommand) {
	ffcq.CommandsQueue = append(ffcq.CommandsQueue, commands...)
}

//Commit fetches all commands in the Queue
func (ffcq *FirestoreFetchBatch) Commit() error {
	//Check initialization
	if ffcq.Client == nil || ffcq.Context == nil || ffcq.CommandsQueue == nil {
		return fmt.Errorf("Not initialized error. ")
	}

	//Check if queue is empty
	if len(ffcq.CommandsQueue) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	errChannel := make(chan error)

	for _, command := range ffcq.CommandsQueue {
		wg.Add(1)

		//Fetch and Get User Model in Goroutine
		go func(fetchCommand FetchCommand, waitGroup *sync.WaitGroup, errChan chan error) {
			defer waitGroup.Done()

			documentSnapshot, err := ffcq.Client.Collection(fetchCommand.Collection).Doc(fetchCommand.DocumentID).Get(ffcq.Context)
			if err != nil {
				//handle the fetch error
				if fetchCommand.FetchCommandErrorHandler != nil {
					//Handle the error
					if unhandledError := fetchCommand.FetchCommandErrorHandler(fetchCommand.AsTypePtr, err); unhandledError != nil {
						errChan <- unhandledError
						return
					}
					//handled Error but no conversion required as no documents are found
					return
				} else {
					errChan <- err
					return
				}
			}

			if fetchCommand.AsTypePtr == nil {
				errChan <- fmt.Errorf("AsTypePtr is passed as null")
				return
			}

			//Force document ReEncoding.
			if fetchCommand.ForceReEncoding {
				data, err := json.Marshal(documentSnapshot.Data())
				if err != nil {
					errChan <- err
					return
				}
				if err = json.Unmarshal(data, fetchCommand.AsTypePtr); err != nil {
					errChan <- err
					return
				}
			} else {
				if err = documentSnapshot.DataTo(fetchCommand.AsTypePtr); err != nil {
					errChan <- err
					return
				}
			}

		}(command, &wg, errChannel)
	}

	//Handle Goroutine Errors
	if err := utils.HandleGoroutineErrors(&wg, errChannel); err != nil {
		return err
	}

	return nil
}
