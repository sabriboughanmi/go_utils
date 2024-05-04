package fetchbatch

import (
	"cloud.google.com/go/firestore"
	"context"
)

type EFirestoreCommand int

const (
	FirestoreCommand_Set         = 0
	FirestoreCommand_Increment   = 1
	FirestoreCommand_ArrayInsert = 2
	FirestoreCommand_ArrayRemove = 3
)

type firestoreUpdateCommand struct {
	commandType EFirestoreCommand
	path        string
	value       interface{}
}

type FirestoreUpdatesQueue struct {
	CommandsQueue []firestoreUpdateCommand
}

//FetchCommandErrorHandler defines either a fetch error can be handled or not.
type FetchCommandErrorHandler func(asTypePtr interface{}, err error) error

//FetchCommand a fetch command
type FetchCommand struct {
	Collection               string
	DocumentID               string
	AsTypePtr                interface{}
	FetchCommandErrorHandler FetchCommandErrorHandler
	// Force document ReEncoding.it's useful for firestore document complex conversions, but comes with a little performance impact.
	ForceReEncoding bool
}

type FirestoreFetchBatch struct {
	CommandsQueue []FetchCommand
	Client        *firestore.Client
	Context       context.Context
}
