package querybatchupdate

import "cloud.google.com/go/firestore"

//constructQuery creates a Query from ContentBatchUpdate data.
func (contentBatchUpdate *ContentBatchUpdate) constructQuery() firestore.Query {

	//var updateDocumentWithLambda = contentBatchUpdate.DocumentUpdateFunction
	var query = contentBatchUpdate.firestoreClient.Collection(contentBatchUpdate.querySearchParams.CollectionID).Query

	//Declare the Query with where conditions
	if contentBatchUpdate.querySearchParams.QueryWhereKeys != nil {
		for _, whereCondition := range contentBatchUpdate.querySearchParams.QueryWhereKeys {
			query = query.Where(whereCondition.Path, string(whereCondition.Op), whereCondition.Value)
		}
	}

	// Declare the Query OrderBy
	if contentBatchUpdate.querySearchParams.QuerySorts != nil {
		for _, orderBy := range contentBatchUpdate.querySearchParams.QuerySorts {
			query = query.OrderBy(orderBy.DocumentSortKey, orderBy.Direction)
		}
	}

	//Initialize the BatchCount if not initialized.
	if contentBatchUpdate.queryPaginationParams.BatchCount == 0 {
		contentBatchUpdate.queryPaginationParams.BatchCount = 500
	}

	//Declare the query limits
	return query.Limit(contentBatchUpdate.queryPaginationParams.BatchCount)
}
