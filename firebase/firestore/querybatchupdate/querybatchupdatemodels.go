package querybatchupdate

import (
	"cloud.google.com/go/firestore"
	"context"
)

//DocumentUpdateFunction gives the user a document from the ContentBatchUpdate in order to handle its custom logic.
type DocumentUpdateFunction func(documentSnapshot *firestore.DocumentSnapshot) error

type EQueryOperator string

const (
	QueryOperator_LessThan             EQueryOperator = "<"
	QueryOperator_LessThanOrEqualTo    EQueryOperator = "<="
	QueryOperator_EqualTo              EQueryOperator = "=="
	QueryOperator_GreaterThan          EQueryOperator = ">"
	QueryOperator_GreaterThanOrEqualTo EQueryOperator = ">="
	QueryOperator_NotEqualTo           EQueryOperator = "!="
	QueryOperator_ArrayContains        EQueryOperator = "array-contains"
	QueryOperator_ArrayContainsAny     EQueryOperator = "array-contains-any"
	QueryOperator_In                   EQueryOperator = "in"
	QueryOperator_NotIn                EQueryOperator = "not-in"
)

// QueryWhere defines a single Where instruction Parameter.
type QueryWhere struct {
	Path  string         `json:"p"`
	Op    EQueryOperator `json:"o"`
	Value interface{}    `json:"v"`
}

// QuerySort defines a single Sort instruction Parameter.
type QuerySort struct {
	DocumentSortKey string              `json:"dsk"`
	Direction       firestore.Direction `json:"d"`
}

// QueryPaginationParams Defines the Pagination parameters
type QueryPaginationParams struct {
	BatchCount int `json:"bc"`
}

// QuerySearchParams defines search parameters for an UpdateContentInBatch
type QuerySearchParams struct {
	CollectionID   string       `json:"cid"`
	QueryWhereKeys []QueryWhere `json:"qwk"`
	QuerySorts     []QuerySort  `json:"qs"`
}

// ContentBatchUpdate contains Parameters for a content Batch Update Operation
type ContentBatchUpdate struct {
	firestoreClient        *firestore.Client
	ctx                    context.Context
	querySearchParams      QuerySearchParams
	queryPaginationParams  QueryPaginationParams
	firestoreUpdates       []firestore.Update
	DocumentUpdateFunction DocumentUpdateFunction
	// Force document ReEncoding.it's useful for firestore document complex conversions, but comes with a little performance impact.
	ForceReEncoding bool
}

// ContentBatchUpdateSerialized contains serialized Parameters for a content Batch Update Operation
type ContentBatchUpdateSerialized struct {
	QuerySearchParams     QuerySearchParams     `json:"qsp"`
	QueryPaginationParams QueryPaginationParams `json:"qpp"`
	FirestoreUpdates      []firestore.Update    `json:"fu"`
	ForceReEncoding       bool                  `json:"fr"`
}
