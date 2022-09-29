package querybatchupdate

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/sabriboughanmi/go_utils/utils"
	"google.golang.org/api/iterator"
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

// SetSearchParameters sets the Search Parameters
func (contentBatchUpdate *ContentBatchUpdate) SetSearchParameters(collectionID string) {

	contentBatchUpdate.querySearchParams.CollectionID = collectionID
}

// SetPaginationParameters sets the Pagination Parameters
func (contentBatchUpdate *ContentBatchUpdate) SetPaginationParameters(limit *int, batchCount int) {
	contentBatchUpdate.queryPaginationParams = QueryPaginationParams{
		Limit:      limit,
		BatchCount: batchCount,
	}
}

// SetValuesToUpdate sets all document  Values To Update
func (contentBatchUpdate *ContentBatchUpdate) SetValuesToUpdate(valuesToUpdate ...ValueToUpdate) {

	if valuesToUpdate == nil || len(valuesToUpdate) == 0 {
		panic("valueToUpdate is empty or nil")
	}
	var valueToUpdateLength = len(valuesToUpdate)

	//declare the ValueToUpdate slice in advance for performance reasons
	contentBatchUpdate.valuesToUpdate = make([]ValueToUpdate, valueToUpdateLength)

	//Fill the values that require updates
	for i, valueToUpdate := range valuesToUpdate {
		contentBatchUpdate.valuesToUpdate[i] = valueToUpdate
	}
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
	if contentBatchUpdate.queryPaginationParams.Limit != nil {
		query.Limit(*contentBatchUpdate.queryPaginationParams.Limit)
	}

	for documentsAvailable {
		var iter *firestore.DocumentIterator

		if contentBatchUpdate.queryPaginationParams.Limit != nil {

			if (*contentBatchUpdate.queryPaginationParams.Limit - processedDocuments) < contentBatchUpdate.queryPaginationParams.BatchCount {
				contentBatchUpdate.queryPaginationParams.BatchCount = *contentBatchUpdate.queryPaginationParams.Limit - processedDocuments
			}
		}

		//Set a new Cursor if a fetch has already been processed in previous iterations
		if cursorValue != nil {
			query.StartAfter(cursorValue)
		}

		iter = query.Documents(contentBatchUpdate.ctx)

		var batch = contentBatchUpdate.firestoreClient.Batch()
		var modifiedDocsCount = 0
		var lastDoc *firestore.DocumentSnapshot
		for {
			var newDoc, err = iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				// Handle error, possibly by returning the error to the caller. Break the loop or return.
				return fmt.Errorf("Error Iterating Posts: %v\n", err)
			}
			var valueToUpdate ValueToUpdate
			lastDoc = newDoc
			modifiedDocsCount++
			processedDocuments++
			batch.Update(newDoc.Ref, []firestore.Update{{
				Path:  valueToUpdate.Key,
				Value: contentBatchUpdate,
			}})
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

			if _, err := batch.Commit(contentBatchUpdate.ctx); err != nil {
				return fmt.Errorf("error Commiting Batch when updating posts %v", err)
			}
		}

		if contentBatchUpdate.queryPaginationParams.Limit != nil && processedDocuments >= *contentBatchUpdate.queryPaginationParams.Limit {
			return nil
		}
	}
	return nil
}
