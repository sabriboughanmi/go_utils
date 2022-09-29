package querybatchupdate

import (
	"cloud.google.com/go/firestore"
	"context"
)

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
	Path  string
	Op    EQueryOperator
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
