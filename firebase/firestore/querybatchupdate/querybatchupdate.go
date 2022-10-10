package updatedocumentaftermodifytest

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/sabriboughanmi/go_utils/utils"
	"log"
)

// CreateContentBatchUpdateInstance returns a new Content Batch Update instance
func CreateContentBatchUpdateInstance(firestoreClient *firestore.Client, ctx context.Context) *ContentBatchUpdate {
	var contentBatchUpdate = ContentBatchUpdate{
		firestoreClient: firestoreClient,
		ctx:             ctx,
	}
	return &contentBatchUpdate

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

	var query = contentBatchUpdate.firestoreClient.Collection(contentBatchUpdate.querySearchParams.CollectionID)

	//Declare the Query with where conditions
	for _, whereCondition := range contentBatchUpdate.querySearchParams.QueryWhereKeys {
		query.Where(whereCondition.Path, string(whereCondition.Op), whereCondition.Value)
	}

	var cursorValue interface{} = nil

	// Declare the Query OrderBy
	if contentBatchUpdate.querySearchParams.QuerySorts != nil {
		for _, orderBy := range contentBatchUpdate.querySearchParams.QuerySorts {
			query.OrderBy(orderBy.DocumentSortKey, orderBy.Direction)
		}
	}
	//Declare the query limits
	query.Limit(contentBatchUpdate.queryPaginationParams.BatchCount)
	query.Offset(20)

	var firestoreWriteBatch *firestore.WriteBatch
	var operationInWriteBatch = 500
	var processedDocuments = 0
	var lastDoc *firestore.DocumentSnapshot

	fetchBatchRequired := contentBatchUpdate.firestoreUpdates != nil && len(contentBatchUpdate.firestoreUpdates) > 0

	for documentsAvailable {
		if contentBatchUpdate.queryPaginationParams.BatchCount != 0 {

			if (contentBatchUpdate.queryPaginationParams.BatchCount - processedDocuments) < contentBatchUpdate.queryPaginationParams.BatchCount {
				contentBatchUpdate.queryPaginationParams.BatchCount = contentBatchUpdate.queryPaginationParams.BatchCount - processedDocuments
			}

		}
		processedDocuments++
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
		lastDoc = downloadedDocumentsSnapShots[len(downloadedDocumentsSnapShots)-1]
		//Check if Documents still available to query in next iteration
		if len(downloadedDocumentsSnapShots) < contentBatchUpdate.queryPaginationParams.BatchCount {
			documentsAvailable = false
		}

		//Set a new Cursor if a fetch has already been processed in previous iterations
		if cursorValue != nil {
			query.StartAfter(cursorValue)
			log.Println("cursorvalue", cursorValue)

		}
		//set cursor if required
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
		//Update with write Batch
		processWithWriteBatch(contentBatchUpdate, &operationInWriteBatch, firestoreWriteBatch, downloadedDocumentsSnapShots)
		//Update with lambda function
		//processWithCustomBehaviour(contentBatchUpdate, lastDoc)
		break
	}
	return nil
}
func processWithWriteBatch(contentBatchUpdate *ContentBatchUpdate, writeBatchCachedOperations *int, firestoreWriteBatch *firestore.WriteBatch, downloadedDocuments []*firestore.DocumentSnapshot) error {
	firestoreWriteBatch = contentBatchUpdate.firestoreClient.Batch()
	for _, newDoc := range downloadedDocuments {
		*writeBatchCachedOperations++
		//Set the Updates
		firestoreWriteBatch.Update(newDoc.Ref, contentBatchUpdate.firestoreUpdates)
		log.Println("update", contentBatchUpdate.firestoreUpdates)
		log.Println("downloaded_documents", downloadedDocuments)
		// Make sure the WriteBatch never exceeds the 500 document updates in a batch.
		if *writeBatchCachedOperations == 500 {
			*writeBatchCachedOperations = 0
			if _, err := firestoreWriteBatch.Commit(contentBatchUpdate.ctx); err != nil {
				return err
			}
		}
		firestoreWriteBatch.Commit(contentBatchUpdate.ctx)
	}
	return nil
}

func processWithCustomBehaviour(contentBatchUpdate *ContentBatchUpdate, downloadedDocuments *firestore.DocumentSnapshot) {
	if contentBatchUpdate.DocumentUpdateFunction != nil {
		if err := contentBatchUpdate.DocumentUpdateFunction(downloadedDocuments); err != nil {
			return
		}
	}

}
