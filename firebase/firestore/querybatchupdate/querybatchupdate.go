package querybatchupdate

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// CreateContentBatchUpdateInstance returns a new Content Batch Update instance
func CreateContentBatchUpdateInstance(firestoreClient *firestore.Client, signedHttpClient *http.Client, apiEndPoint string, ctx context.Context) *ContentBatchUpdate {
	var contentBatchUpdate = ContentBatchUpdate{
		firestoreClient:  firestoreClient,
		signedHttpClient: signedHttpClient,
		apiEndPoint:      apiEndPoint,
		batchOperations:  make(map[string]BatchOperation),
		ctx:              ctx,
		queryPaginationParams: QueryPaginationParams{
			BatchCount: 500,
		},
	}
	return &contentBatchUpdate

}

/*
//Serialize a ContentBatchUpdate Instance Data.
func (contentBatchUpdate *ContentBatchUpdate) Serialize() string {
	var serializedQuery = ContentBatchUpdateSerialized{
		QuerySearchParams:     contentBatchUpdate.querySearchParams,
		QueryPaginationParams: contentBatchUpdate.queryPaginationParams,
		BatchOperations:       contentBatchUpdate.batchOperations,
		ForceReEncoding:       contentBatchUpdate.ForceReEncoding,
	}
	return string(utils.UnsafeAnythingToJSON(serializedQuery))
}

//Deserialize loads a ContentBatchUpdate serialized Data.
func (contentBatchUpdate *ContentBatchUpdate) Deserialize(serializeQuery string) error {

	var serializedQuery ContentBatchUpdateSerialized
	if err := json.Unmarshal([]byte(serializeQuery), &serializedQuery); err != nil {
		return err
	}
	//Load Data
	contentBatchUpdate.queryPaginationParams = serializedQuery.QueryPaginationParams
	contentBatchUpdate.querySearchParams = serializedQuery.QuerySearchParams
	contentBatchUpdate.batchOperations = serializedQuery.BatchOperations
	contentBatchUpdate.ForceReEncoding = serializedQuery.ForceReEncoding
	return nil
}
*/

// Where adds a Where condition to the query
func (contentBatchUpdate *ContentBatchUpdate) Where(path string, operation EQueryOperator, value interface{}) *ContentBatchUpdate {
	contentBatchUpdate.querySearchParams.QueryWheres = append(contentBatchUpdate.querySearchParams.QueryWheres, QueryWhere{
		Path:  path,
		Op:    operation,
		Value: value,
	})
	return contentBatchUpdate
}

// OrderBy adds an OrderBy condition to the query
func (contentBatchUpdate *ContentBatchUpdate) OrderBy(documentSortKey string, direction firestore.Direction) *ContentBatchUpdate {
	contentBatchUpdate.querySearchParams.QuerySorts = append(contentBatchUpdate.querySearchParams.QuerySorts, QuerySort{
		DocumentSortKey: documentSortKey,
		Direction:       direction,
	})
	return contentBatchUpdate
}

// Collection sets the Search Parameters
func (contentBatchUpdate *ContentBatchUpdate) Collection(collectionID string) *ContentBatchUpdate {
	contentBatchUpdate.querySearchParams.CollectionID = collectionID
	return contentBatchUpdate
}

// SetBatchCount sets the SetBatchCount Parameters
func (contentBatchUpdate *ContentBatchUpdate) SetBatchCount(batchCount int) *ContentBatchUpdate {

	//Limit the batch count as firestore is limited to 500 changes at once.
	if batchCount > 500 {
		batchCount = 500
	}

	contentBatchUpdate.queryPaginationParams.BatchCount = batchCount

	return contentBatchUpdate
}

// SetQueryLimit sets the Query Limit
func (contentBatchUpdate *ContentBatchUpdate) SetQueryLimit(limit uint) *ContentBatchUpdate {
	var lmt = int(limit)
	contentBatchUpdate.queryPaginationParams.Limit = &lmt
	return contentBatchUpdate
}

// Select sets the Query Documents fields that must be downloaded Parameter
func (contentBatchUpdate *ContentBatchUpdate) Select(fields ...string) *ContentBatchUpdate {

	if contentBatchUpdate.querySearchParams.SelectFields == nil {
		contentBatchUpdate.querySearchParams.SelectFields = fields
	} else {
		contentBatchUpdate.querySearchParams.SelectFields = append(contentBatchUpdate.querySearchParams.SelectFields, fields...)
	}

	return contentBatchUpdate
}

// StartAt returns a new Query that specifies that results should start at
// the document with the given field values.
//
// StartAt may be called with a single DocumentSnapshot, representing an
// existing document within the query. The document must be a direct child of
// the location being queried (not a parent document, or document in a
// different collection, or a grandchild document, for example).
//
// Otherwise, StartAt should be called with one field value for each OrderBy clause,
// in the order that they appear. For example, in
//
//	q.OrderBy("X", Asc).OrderBy("Y", Desc).StartAt(1, 2)
//
// results will begin at the first document where X = 1 and Y = 2.
//
// If an OrderBy call uses the special DocumentID field path, the corresponding value
// should be the document ID relative to the query's collection. For example, to
// start at the document "NewYork" in the "States" collection, write
//
//	client.Collection("States").OrderBy(DocumentID, firestore.Asc).StartAt("NewYork")
//
// Calling StartAt overrides a previous call to StartAt or StartAfter.
func (contentBatchUpdate *ContentBatchUpdate) StartAt(docSnapshotOrFieldValues ...interface{}) *ContentBatchUpdate {
	var err error
	contentBatchUpdate.startBefore = true
	contentBatchUpdate.startVals, contentBatchUpdate.startDoc, err = processCursorArg("StartAt", docSnapshotOrFieldValues)
	if err != nil {
		panic(fmt.Sprintf("ContentBatchUpdate.StartAt : %v", err))
	}
	return contentBatchUpdate
}

// StartAfter returns a new Query that specifies that results should start just after
// the document with the given field values. See Query.StartAt for more information.
//
// Calling StartAfter overrides a previous call to StartAt or StartAfter.
func (contentBatchUpdate *ContentBatchUpdate) StartAfter(docSnapshotOrFieldValues ...interface{}) *ContentBatchUpdate {
	var err error
	contentBatchUpdate.startBefore = false
	contentBatchUpdate.startVals, contentBatchUpdate.startDoc, err = processCursorArg("StartAfter", docSnapshotOrFieldValues)
	if err != nil {
		panic(fmt.Sprintf("ContentBatchUpdate.StartAfter : %v", err))
	}
	return contentBatchUpdate
}

// EndAt returns a new Query that specifies that results should end at the
// document with the given field values. See Query.StartAt for more information.
//
// Calling EndAt overrides a previous call to EndAt or EndBefore.
func (contentBatchUpdate *ContentBatchUpdate) EndAt(docSnapshotOrFieldValues ...interface{}) *ContentBatchUpdate {
	var err error
	contentBatchUpdate.endBefore = false
	contentBatchUpdate.endVals, contentBatchUpdate.endDoc, err = processCursorArg("EndAt", docSnapshotOrFieldValues)
	if err != nil {
		panic(fmt.Sprintf("ContentBatchUpdate.EndAt : %v", err))
	}
	return contentBatchUpdate
}

// EndBefore returns a new Query that specifies that results should end just before
// the document with the given field values. See Query.StartAt for more information.
//
// Calling EndBefore overrides a previous call to EndAt or EndBefore.
func (contentBatchUpdate *ContentBatchUpdate) EndBefore(docSnapshotOrFieldValues ...interface{}) *ContentBatchUpdate {
	var err error
	contentBatchUpdate.endBefore = true
	contentBatchUpdate.endVals, contentBatchUpdate.endDoc, err = processCursorArg("EndBefore", docSnapshotOrFieldValues)
	if err != nil {
		panic(fmt.Sprintf("ContentBatchUpdate.EndBefore : %v", err))
	}
	return contentBatchUpdate
}

// Processes a Cursor Args
func processCursorArg(name string, docSnapshotOrFieldValues []interface{}) ([]interface{}, *firestore.DocumentSnapshot, error) {
	for _, e := range docSnapshotOrFieldValues {
		if ds, ok := e.(*firestore.DocumentSnapshot); ok {
			if len(docSnapshotOrFieldValues) == 1 {
				return nil, ds, nil
			}
			return nil, nil, fmt.Errorf("firestore: a document snapshot must be the only argument to %s", name)
		}
	}
	return docSnapshotOrFieldValues, nil, nil
}

// UpdateValues requests to UPDATE all documents with the respective IDs in the specified collection.
func (contentBatchUpdate *ContentBatchUpdate) UpdateValues(collectionID string, valuesToUpdate ...firestore.Update) *ContentBatchUpdate {
	if valuesToUpdate == nil || len(valuesToUpdate) == 0 {
		return contentBatchUpdate
	}

	//check if previous updates has been declared for this collection.
	if _, found := contentBatchUpdate.batchOperations[collectionID]; !found {
		contentBatchUpdate.batchOperations[collectionID] = BatchOperation{
			OperationType:   BatchOperationType_Update,
			FirestoreUpdate: valuesToUpdate,
		}
	} else {
		batchOperation := contentBatchUpdate.batchOperations[collectionID]
		batchOperation.FirestoreUpdate = append(batchOperation.FirestoreUpdate, valuesToUpdate...)
		contentBatchUpdate.batchOperations[collectionID] = batchOperation
	}
	return contentBatchUpdate
}

// DeleteDocument requests to DELETE all documents with the respective IDs in the specified collection.
func (contentBatchUpdate *ContentBatchUpdate) DeleteDocument(collectionID string) *ContentBatchUpdate {
	contentBatchUpdate.batchOperations[collectionID] = BatchOperation{
		OperationType:   BatchOperationType_Delete,
		FirestoreUpdate: nil,
	}
	return contentBatchUpdate
}

// SubscribeToContentBatchUpdates will be executed in parallel for every ContentBatchUpdate providing the Data processed by the Batch.
//the BatchCallback will only be executed if the WriteBatch Operation has succeeded
func (contentBatchUpdate *ContentBatchUpdate) SubscribeToContentBatchUpdates(callback BatchCallback) *ContentBatchUpdate {
	contentBatchUpdate.callback = callback
	return contentBatchUpdate
}

//TODO: Handle Error in Goroutines and Sub-Goroutines

// UpdateContentInBatch used to update User Generated Contents display name when the username has changed
func (contentBatchUpdate *ContentBatchUpdate) UpdateContentInBatch() error {

	//construct a query from ContentBatchUpdate params.
	var query = contentBatchUpdate.constructQuery()

	//Set the Query Limit
	if contentBatchUpdate.queryPaginationParams.Limit != nil {
		//Set the query limit to a higher value than the default one
		query = query.Limit(*contentBatchUpdate.queryPaginationParams.Limit)
	}

	// Start the timer
	start := time.Now()

	//Get Documents count
	var res, err = query.NewAggregationQuery().WithCount(firestore.DocumentID).Get(contentBatchUpdate.ctx)
	if err != nil {
		return err
	}

	//Get the Documents count of the query.
	queryDocumentsCount := int(res["__key__"].(*firestorepb.Value).GetIntegerValue())

	//Return if there is no documents to update
	if queryDocumentsCount == 0 {
		return nil
	}

	var wg sync.WaitGroup

	//separate the documents in order to be processed in batches
	for i := 0; i < queryDocumentsCount; i += contentBatchUpdate.queryPaginationParams.BatchCount {

		//Get the End Index
		endIndex := i + contentBatchUpdate.queryPaginationParams.BatchCount
		if endIndex > queryDocumentsCount {
			endIndex = queryDocumentsCount
		}
		//Run the firestore WriteBatch Process in parallel
		wg.Add(1)

		//Process the query fragment.
		go func(waitGroup *sync.WaitGroup, queryOffset int) {
			defer waitGroup.Done()

			structuredQuery, err := contentBatchUpdate.createStructuredQuery()
			if err != nil {
				fmt.Printf("%v", err)
				return
			}

			//Set the Offset
			structuredQuery.Offset = int64(queryOffset)

			//set currency parallel Query max documents to download
			structuredQuery.Limit = int64(endIndex - queryOffset)

			//Get StructuredQuery as request body
			structuredQueryBodyBytes, err := structuredQuery.MarshalJSON()

			//construct a structuredQuery request body.
			requestBody := []byte(fmt.Sprintf(`{"structuredQuery":%s}`, string(structuredQueryBodyBytes)))

			//Request firestore Rest API for a StructuredQuery.
			var partialSnapshots []FirestorePartialSnapshot
			if err = contentBatchUpdate.postRequestWithGoogleSignedHttpClient(requestBody, &partialSnapshots); err != nil {
				fmt.Printf("%v", err)
				return
			}

			//Filter out all skip elements that firestore adds
			for i, ps := range partialSnapshots {
				if ps.Document.DocumentFullPath != "" {
					partialSnapshots = partialSnapshots[i:]
					break
				}
			}

			//Run the the firestore batch operations only if required
			if contentBatchUpdate.batchOperations != nil && len(contentBatchUpdate.batchOperations) > 0 {

				var collectionUpdateWG sync.WaitGroup

				//Iterate over the Collections that require modifications
				for collection, batchOperations := range contentBatchUpdate.batchOperations {

					collectionUpdateWG.Add(1)

					//Process BatchOperations in separate Goroutines
					go func(fCollection string, operation *BatchOperation, firestorePartialSnapshots []FirestorePartialSnapshot, collectionsWG *sync.WaitGroup) {
						defer collectionsWG.Done()

						//Create a new firestore WriteBatch
						var firestoreWriteBatch = contentBatchUpdate.firestoreClient.Batch()

						//Create a collection Reference
						var collectionReference = contentBatchUpdate.firestoreClient.Collection(fCollection)

						//Execute the lambda function
						for _, doc := range firestorePartialSnapshots {
							//skip (firestore skipped documents objects)
							if doc.Document.DocumentFullPath == "" {
								continue
							}

							switch operation.OperationType {
							case BatchOperationType_Update:
								//Add an Update operation
								firestoreWriteBatch.Update(collectionReference.Doc(filepath.Base(doc.Document.DocumentFullPath)), operation.FirestoreUpdate)
								break
							case BatchOperationType_Delete:
								//Add a Delete operation
								firestoreWriteBatch.Delete(collectionReference.Doc(filepath.Base(doc.Document.DocumentFullPath)))
								break
							default:
								panic("unsupported BatchOperationType")
							}
						}

						if _, err = firestoreWriteBatch.Commit(contentBatchUpdate.ctx); err != nil {
							//FIXME: make it add register the errors, so we can log/handle them Later
							fmt.Printf("%v", err)
							return
						}

					}(collection, &batchOperations, partialSnapshots, &collectionUpdateWG)
				}

				//Wait the Collections Updates
				collectionUpdateWG.Wait()
			}

			// Run the CallBack if required
			if contentBatchUpdate.callback != nil {
				contentBatchUpdate.callback(partialSnapshots)
			}

		}(&wg, i)
	}

	//Wait for all goroutines to complete
	wg.Wait()

	// Stop the timer and print the elapsed time
	elapsed := time.Since(start)
	fmt.Println("Execution time:", elapsed.Seconds(), "seconds")

	return nil

}
