package querybatchupdate

import (
	"cloud.google.com/go/firestore"
	"context"
)

// QueryWhere defines a single Where instruction Parameter.
type QueryWhere struct {
	Path  string
	Op    string
	Value interface{}
}

// QuerySort defines a single Sort instruction Parameter.
type QuerySort struct {
	DocumentSortKey string
	Direction       firestore.Direction
}

// QueryPaginationParams Defines the Pagination parameters
type QueryPaginationParams struct {
	Limit      *int
	BatchCount int
}

// ValueToUpdate Defines a document single value to Update.
type ValueToUpdate struct {
	Key   string
	Value interface{}
}

// QuerySearchParams defines search parameters for an UpdateContentInBatch
type QuerySearchParams struct {
	CollectionID   string
	QueryWhereKeys []QueryWhere
	QuerySorts     []QuerySort
}

// ContentBatchUpdate contains Parameters for a content Batch Update Operation
type ContentBatchUpdate struct {
	firestoreClient       *firestore.Client
	ctx                   context.Context
	querySearchParams     QuerySearchParams
	queryPaginationParams QueryPaginationParams
	valuesToUpdate        []ValueToUpdate
}
