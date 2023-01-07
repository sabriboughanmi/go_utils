package querybatchupdate

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

//constructQuery creates a Query from ContentBatchUpdate data.
func (contentBatchUpdate *ContentBatchUpdate) constructQuery() firestore.Query {

	//var updateDocumentWithLambda = contentBatchUpdate.DocumentUpdateFunction
	var query = contentBatchUpdate.firestoreClient.Collection(contentBatchUpdate.querySearchParams.CollectionID).Query

	//Declare the Query with where conditions
	if contentBatchUpdate.querySearchParams.QueryWheres != nil {
		for _, whereCondition := range contentBatchUpdate.querySearchParams.QueryWheres {
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

	//Set the Limits
	if contentBatchUpdate.queryPaginationParams.Limit != nil {
		//Set the Query Limit Parameter
		query = query.Limit(*contentBatchUpdate.queryPaginationParams.Limit)
	}

	return query
}

//postRequestWithGoogleSignedHttpClient returns a Post request response as type pointer.
//eg: apiEndpoint := "https://firestore.googleapis.com/v1/projects/<project_id>/databases/(default)/documents:runQuery"
func (contentBatchUpdate *ContentBatchUpdate) postRequestWithGoogleSignedHttpClient(requestBody []byte, typeRef interface{}) error {
	// Send the HTTP request to the Cloud Firestore API.
	resp, err := contentBatchUpdate.signedHttpClient.Post(contentBatchUpdate.apiEndPoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print the HTTP response status and body.
	if resp.StatusCode != 200 {
		return fmt.Errorf("status : %s body: %s", resp.Status, string(b))
	}

	err = json.Unmarshal(b, typeRef)
	if err != nil {
		return err
	}
	return nil
}

//GetDocumentID returns the current DocumentID
func (firestorePartialSnapshot *FirestorePartialSnapshot) GetDocumentID() string {
	return filepath.Base(firestorePartialSnapshot.Document.DocumentFullPath)
}

//GetDocumentRef returns a *firestore.DocumentRef, to facilitate hierarchy access.
func (firestorePartialSnapshot *FirestorePartialSnapshot) GetDocumentRef(client *firestore.Client) *firestore.DocumentRef {
	return client.Doc(firestorePartialSnapshot.Document.DocumentFullPath)
}
