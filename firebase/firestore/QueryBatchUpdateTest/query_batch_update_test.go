package querybatchupdatetest

import (
	"cloud.google.com/go/firestore"
	"context"
	"sync"
	"testing"
)

type DocumentUpdateFunction func(documentSnapshot *firestore.DocumentSnapshot) error

func Test_Update_Document(t *testing.T) {
	const (
		Collection_Query_Update = "QueryUpdate"
	)

	var firestoreClient *firestore.Client
	var ctx context.Context

	var documentsAvailable = true
	var query = firestoreClient.Collection(Collection_Query_Update)

	downloadedDocuments, err := query.Documents(ctx).GetAll()
	var firestoreWriteBatch *firestore.WriteBatch
	var modifiedDocsCount = 0
	var lastDoc *firestore.DocumentSnapshot
	var processedDocuments = 0

	wg := sync.WaitGroup{}
	errChannel := make(chan error)
	wg.Add(1)
	for documentsAvailable {
		go func(contentBatchUpdate, waitGroup *sync.WaitGroup, errChan chan error) {
			defer waitGroup.Done()

			for _, newDoc := range downloadedDocuments {
				lastDoc = newDoc
				//Do batch updates if required
				modifiedDocsCount++
				processedDocuments++

				// Make sure the WriteBatch never exceeds the 500 document updates in a batch.
				if modifiedDocsCount == 500 {
					firestoreWriteBatch = firestoreClient.Batch()
					if _, err := firestoreWriteBatch.Commit(ctx); err != nil {
						return
					}
					processedDocuments++
					modifiedDocsCount = 0
				}
				//Set the Updates
				firestoreWriteBatch.Update(newDoc.Ref, []firestore.Update{{
					Path:  "keyToUpdate",
					Value: "newValue",
				}})
			}

			// Handle our Query DocumentUpdateFunction per document
			if DocumentUpdateFunction != nil {
				if err = DocumentUpdateFunction(lastDoc); err != nil {
					return
				}

			}
		}(contentBatchUpdate, &wg, errChannel)
	}

}
