package querybatchupdate

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
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
func (contentBatchUpdate *ContentBatchUpdate) Pagination(limit int, batchCount int) *ContentBatchUpdate {
	contentBatchUpdate.queryPaginationParams = QueryPaginationParams{
		Limit:      limit,
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
	var cursorValue interface{} = nil

	var processedDocuments = 0

	var query = contentBatchUpdate.firestoreClient.Collection(contentBatchUpdate.querySearchParams.CollectionID)

	//Declare the Query with where conditions
	for _, whereCondition := range contentBatchUpdate.querySearchParams.QueryWhereKeys {
		query.Where(whereCondition.Path, string(whereCondition.Op), whereCondition.Value)
	}

	// Declare the Query OrderBy (if required)
	if contentBatchUpdate.querySearchParams.QuerySorts != nil {
		for _, orderBy := range contentBatchUpdate.querySearchParams.QuerySorts {
			query.OrderBy(orderBy.DocumentSortKey, orderBy.Direction)
		}
	}
	//Declare the query limits
	query.Limit(contentBatchUpdate.queryPaginationParams.Limit)

	for documentsAvailable {

		//Make sure that limit doesn't exeed the batchCount
		if (contentBatchUpdate.queryPaginationParams.Limit - processedDocuments) < contentBatchUpdate.queryPaginationParams.BatchCount {
			contentBatchUpdate.queryPaginationParams.BatchCount = contentBatchUpdate.queryPaginationParams.Limit - processedDocuments
		}

		//Set a new Cursor if a fetch has already been processed in previous iterations
		if cursorValue != nil {
			query.StartAfter(cursorValue)
		}

		//Fetch all documents in parallel
		downloadedDocuments, err := query.Documents(contentBatchUpdate.ctx).GetAll()
		if err != nil {
			return err
		}

		//No documents to update
		if downloadedDocuments == nil || len(downloadedDocuments) == 0 {
			return nil
		}

		//No more documents available to fetch
		if len(downloadedDocuments) < contentBatchUpdate.queryPaginationParams.Limit {
			documentsAvailable = false
		}

		var firestoreWriteBatch = contentBatchUpdate.firestoreClient.Batch()
		var modifiedDocsCount = 0
		var lastDoc *firestore.DocumentSnapshot
		for _, newDoc := range downloadedDocuments {

			lastDoc = newDoc
			modifiedDocsCount++
			processedDocuments++
			//TODO: Make sure the WriteBatch never exceeds the 500 document updates in a batch.
			if modifiedDocsCount < 500 {
				firestoreWriteBatch.Commit(contentBatchUpdate.ctx)
				firestoreWriteBatch = contentBatchUpdate.firestoreClient.Batch()
				processedDocuments++
				modifiedDocsCount = 0
			}
			//Set the Updates
			firestoreWriteBatch.Update(newDoc.Ref, contentBatchUpdate.firestoreUpdates)
		}

		if modifiedDocsCount < contentBatchUpdate.queryPaginationParams.BatchCount {
			documentsAvailable = false
		}
		var querySort QuerySort
		if modifiedDocsCount > 0 {
			cValue, err := utils.GetValueFromSubMap(lastDoc.Data(), querySort.DocumentSortKey)
			if err != nil {
				return err
			}
			cursorValue = cValue

			if _, err := firestoreWriteBatch.Commit(contentBatchUpdate.ctx); err != nil {
				return fmt.Errorf("error Commiting Batch when updating posts %v", err)
			}
		}

		if contentBatchUpdate.queryPaginationParams.Limit != 0 && processedDocuments >= contentBatchUpdate.queryPaginationParams.Limit {
			return nil
		}
	}
	return nil
}
