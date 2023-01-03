package querybatchupdate

import (
	"cloud.google.com/go/firestore"
	"context"
	"net/http"
)

type EQueryOperator string

const (
	QueryOperator_LessThan             EQueryOperator = "LESS_THAN"
	QueryOperator_LessThanOrEqualTo    EQueryOperator = "LESS_THAN_OR_EQUAL"
	QueryOperator_EqualTo              EQueryOperator = "EQUAL"
	QueryOperator_GreaterThan          EQueryOperator = "GREATER_THAN"
	QueryOperator_GreaterThanOrEqualTo EQueryOperator = "GREATER_THAN_OR_EQUAL"
	QueryOperator_NotEqualTo           EQueryOperator = "NOT_EQUAL"
	QueryOperator_ArrayContains        EQueryOperator = "ARRAY_CONTAINS"
	QueryOperator_ArrayContainsAny     EQueryOperator = "ARRAY_CONTAINS_ANY"
	QueryOperator_In                   EQueryOperator = "IN"
	QueryOperator_NotIn                EQueryOperator = "NOT_IN"
)

// QueryWhere defines a single Where instruction Parameter.
type QueryWhere struct {
	Path  string         `json:"p"`
	Op    EQueryOperator `json:"o"`
	Value interface{}    `json:"v"`
}

// QuerySort defines a single OrderBy instruction Parameter.
type QuerySort struct {
	DocumentSortKey string              `json:"dsk"`
	Direction       firestore.Direction `json:"d"`
}

// QueryPaginationParams Defines the SetBatchCount parameters
type QueryPaginationParams struct {
	BatchCount int  `json:"bc"`
	Limit      *int `json:"l"`
}

// QuerySearchParams defines search parameters for an UpdateContentInBatch
type QuerySearchParams struct {
	CollectionID string       `json:"cid"`
	QueryWheres  []QueryWhere `json:"qw"`
	QuerySorts   []QuerySort  `json:"qs"`
	SelectFields []string     `json:"sf"`
	Offset       int64        `json:"o"`
}

// ContentBatchUpdate contains Parameters for a content Batch Update Operation
type ContentBatchUpdate struct {
	firestoreClient        *firestore.Client
	ctx                    context.Context
	signedHttpClient       *http.Client
	apiEndPoint            string
	querySearchParams      QuerySearchParams
	queryPaginationParams  QueryPaginationParams
	batchOperations        map[string]BatchOperation
	startVals, endVals     []interface{}
	startDoc, endDoc       *firestore.DocumentSnapshot
	startBefore, endBefore bool
	// Force document ReEncoding.it's useful for firestore document complex conversions, but comes with a little performance impact.
	ForceReEncoding bool
}

// ContentBatchUpdateSerialized contains serialized Parameters for a content Batch Update Operation
type ContentBatchUpdateSerialized struct {
	QuerySearchParams     QuerySearchParams         `json:"qsp"`
	QueryPaginationParams QueryPaginationParams     `json:"qpp"`
	BatchOperations       map[string]BatchOperation `json:"bo"`
	ForceReEncoding       bool                      `json:"fr"`
}

type EBatchOperationType int8

const (
	BatchOperationType_Update EBatchOperationType = iota
	BatchOperationType_Delete
)

//BatchOperation contains firestore operations data and command type.
type BatchOperation struct {
	OperationType   EBatchOperationType `json:"ot"`
	FirestoreUpdate []firestore.Update  `json:"fu,omitempty"`
}
