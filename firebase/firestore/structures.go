package firestore

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
