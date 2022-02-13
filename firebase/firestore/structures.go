package firestore

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

//FetchCommand a fetch command
type FetchCommand struct {
	Collection string
	DocumentID string
	AsTypePtr  interface{}
}

type firestoreFetchBatch struct {
	CommandsQueue []FetchCommand
	Client        *firestore.Client
	Context       context.Context
}
