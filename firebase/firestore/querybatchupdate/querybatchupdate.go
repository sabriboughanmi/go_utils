package querybatchupdate

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/sabriboughanmi/go_utils/utils"
)

// CreateContentBatchUpdateInstance returns a new Content Batch Update instance
func CreateContentBatchUpdateInstance(firestoreClient *firestore.Client, ctx context.Context) *ContentBatchUpdate {
	var contentBatchUpdate = ContentBatchUpdate{
		firestoreClient: firestoreClient,
		ctx:             ctx,
	}
	return &contentBatchUpdate

}

//Serialize a ContentBatchUpdate Instance Data.
func (contentBatchUpdate *ContentBatchUpdate) Serialize() string {
	var serializedQuery = ContentBatchUpdateSerialized{
		QuerySearchParams:     contentBatchUpdate.querySearchParams,
		QueryPaginationParams: contentBatchUpdate.queryPaginationParams,
		FirestoreUpdates:      contentBatchUpdate.firestoreUpdates,
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
	contentBatchUpdate.firestoreUpdates = serializedQuery.FirestoreUpdates
	contentBatchUpdate.ForceReEncoding = serializedQuery.ForceReEncoding
	return nil
}

// Where adds a Where condition to the query
func (contentBatchUpdate *ContentBatchUpdate) Where(path string, operation EQueryOperator, value interface{}) *ContentBatchUpdate {
	contentBatchUpdate.querySearchParams.QueryWhereKeys = append(contentBatchUpdate.querySearchParams.QueryWhereKeys, QueryWhere{
		Path:  path,
		Op:    operation,
		Value: value,
	})
	return contentBatchUpdate
}

// Sort adds a Sort condition to the query
func (contentBatchUpdate *ContentBatchUpdate) Sort(documentSortKey string, direction firestore.Direction) *ContentBatchUpdate {
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

// Pagination sets the Pagination Parameters
func (contentBatchUpdate *ContentBatchUpdate) Pagination(batchCount int) *ContentBatchUpdate {

	//Limit the batch count as firestore is limited to 500 changes at once.
	if batchCount > 500 {
		batchCount = 500
	}

	contentBatchUpdate.queryPaginationParams = QueryPaginationParams{
		BatchCount: batchCount,
	}

	return contentBatchUpdate
}

// ValueToUpdate add a value to update in batch
func (contentBatchUpdate *ContentBatchUpdate) ValueToUpdate(key string, value interface{}) *ContentBatchUpdate {
	contentBatchUpdate.firestoreUpdates = append(contentBatchUpdate.firestoreUpdates, firestore.Update{
		Path:  key,
		Value: value,
	})
	return contentBatchUpdate
}

// ValuesToUpdate adds multiple values to update in batch
func (contentBatchUpdate *ContentBatchUpdate) ValuesToUpdate(valuesToUpdate ...firestore.Update) *ContentBatchUpdate {
	if valuesToUpdate == nil || len(valuesToUpdate) == 0 {
		return contentBatchUpdate
	}
	contentBatchUpdate.firestoreUpdates = append(contentBatchUpdate.firestoreUpdates, valuesToUpdate...)
	return contentBatchUpdate
}

// UpdateContentInBatch used to update User Generated Contents display name when the username has changed
func (contentBatchUpdate *ContentBatchUpdate) UpdateContentInBatch() error {

	var documentsAvailable = true

	//construct a query from ContentBatchUpdate params.
	var query = contentBatchUpdate.constructQuery()

	var cursorValue interface{} = nil

	var firestoreWriteBatch *firestore.WriteBatch
	var operationInWriteBatch = 0
	var lastDoc *firestore.DocumentSnapshot

	fetchBatchRequired := contentBatchUpdate.firestoreUpdates != nil && len(contentBatchUpdate.firestoreUpdates) > 0

	for documentsAvailable {

		//The WriteBatch is only required if the firestore static updates are required
		if fetchBatchRequired {
			firestoreWriteBatch = contentBatchUpdate.firestoreClient.Batch()
		}
		//Fetch all documents in parallel
		downloadedDocumentsSnapShots, err := query.Documents(contentBatchUpdate.ctx).GetAll()
		if err != nil {
			return err
		}

		//No documents to update
		if downloadedDocumentsSnapShots == nil || len(downloadedDocumentsSnapShots) == 0 {
			return nil
		}

		//Set last document (used for pagination)
		lastDoc = downloadedDocumentsSnapShots[len(downloadedDocumentsSnapShots)-1]

		//Check if Documents still available to query in next iteration
		if len(downloadedDocumentsSnapShots) < contentBatchUpdate.queryPaginationParams.BatchCount {
			documentsAvailable = false
		}

		//set cursor if more documents can be fetched.
		if documentsAvailable {
			if contentBatchUpdate.querySearchParams.QuerySorts != nil {
				if len(contentBatchUpdate.querySearchParams.QuerySorts) > 1 {
					if cursorValue == nil {
						cursorValue = make([]interface{}, len(contentBatchUpdate.querySearchParams.QuerySorts))
					}

					multipleValue := cursorValue.([]interface{})
					for i, orderKeys := range contentBatchUpdate.querySearchParams.QuerySorts {
						if operationInWriteBatch > 0 {
							value, err := utils.GetValueFromSubMap(lastDoc.Data(), orderKeys.DocumentSortKey)
							if err != nil {
								return err
							}

							multipleValue[i] = value

							if _, err := firestoreWriteBatch.Commit(contentBatchUpdate.ctx); err != nil {
								return err
							}
						}
					}
					cursorValue = multipleValue
				} else {
					//firestoreWriteBatch.Commit(contentBatchUpdate.ctx)
					cursorValue = lastDoc.Data()[contentBatchUpdate.querySearchParams.QuerySorts[0].DocumentSortKey]
				}

			}
		}

		//Set a new Cursor if a fetch has already been processed in previous iterations
		if cursorValue != nil {
			query = query.StartAfter(cursorValue)
		}

		//update with lambda function
		if contentBatchUpdate.DocumentUpdateFunction != nil {
			//Execute the lambda function
			for _, doc := range downloadedDocumentsSnapShots {
				if err = contentBatchUpdate.DocumentUpdateFunction(doc); err != nil {
					return err
				}
			}

		} else {

			//Execute the lambda function
			for _, doc := range downloadedDocumentsSnapShots {

				//increment the operations in the write batch
				operationInWriteBatch++

				firestoreWriteBatch.Update(doc.Ref, contentBatchUpdate.firestoreUpdates)

				//update with write batch
				if operationInWriteBatch == contentBatchUpdate.queryPaginationParams.BatchCount {
					if _, err := firestoreWriteBatch.Commit(contentBatchUpdate.ctx); err != nil {
						return err
					}

					//Reset the operations count
					operationInWriteBatch = 0

					//Reset the Write batch
					firestoreWriteBatch = contentBatchUpdate.firestoreClient.Batch()

				}
			}

			//if no documents available but some operations are still in queue
			if !documentsAvailable && operationInWriteBatch > 0 {
				if _, err := firestoreWriteBatch.Commit(contentBatchUpdate.ctx); err != nil {
					return err
				}
			}

		}

	}

	return nil
}
