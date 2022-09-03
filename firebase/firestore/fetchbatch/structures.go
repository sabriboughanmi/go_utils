package fetchbatch

import (
	"cloud.google.com/go/firestore"
	"context"
)

type EFirestoreCommand int

const (
	FirestoreCommand_Set    = 0
	FirestoreCommand_Update = 1
)

type FirestoreUpdateCommand struct {
	CommandType EFirestoreCommand
	Value       int
	Path        string
}

type FirestoreUpdatesQueue struct {
	CommandsQueue []FirestoreUpdateCommand
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

type firestoreFetchBatch struct {
	CommandsQueue []FetchCommand
	Client        *firestore.Client
	Context       context.Context
}
