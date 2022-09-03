package fetchbatch

import "cloud.google.com/go/firestore"

//AddCommand adds a FirestoreUpdateCommand to the Queue
func (firestoreUpdates *FirestoreUpdatesQueue) AddCommand(command FirestoreUpdateCommand) {
	firestoreUpdates.CommandsQueue = append(firestoreUpdates.CommandsQueue, command)
}

//AddCommands adds FirestoreUpdateCommands to the Queue
func (firestoreUpdates *FirestoreUpdatesQueue) AddCommands(command ...FirestoreUpdateCommand) {
	firestoreUpdates.CommandsQueue = append(firestoreUpdates.CommandsQueue, command...)
}

//Merge fetch all FirestoreUpdateCommand(s) from passed FirestoreUpdatesQueue(s)
func (firestoreUpdates *FirestoreUpdatesQueue) Merge(queues ...FirestoreUpdatesQueue) {
	for _, queue := range queues {
		if queue.CommandsQueue != nil {
			firestoreUpdates.CommandsQueue = append(firestoreUpdates.CommandsQueue, queue.CommandsQueue...)
		}
	}
}

//ClearQueue set the CommandsQueue to nil
func (firestoreUpdates *FirestoreUpdatesQueue) ClearQueue() {
	firestoreUpdates.CommandsQueue = nil
}

//GetFirestoreUpdates merges all FirestoreUpdateCommand(s) and returns an []firestore.Update
func (firestoreUpdates *FirestoreUpdatesQueue) GetFirestoreUpdates() []firestore.Update {

	var mergedCommands = make(map[string]FirestoreUpdateCommand)

	for _, command := range firestoreUpdates.CommandsQueue {
		if oldCommand, found := mergedCommands[command.Path]; !found || command.CommandType == FirestoreCommand_Set {
			mergedCommands[command.Path] = command
		} else {
			oldCommand.Value += command.Value
			mergedCommands[command.Path] = oldCommand
		}
	}

	var updates []firestore.Update
	for path, command := range mergedCommands {
		var firestoreUpdate = firestore.Update{
			Path: path,
		}
		if command.CommandType == FirestoreCommand_Set {
			firestoreUpdate.Value = command.Value
		} else {
			firestoreUpdate.Value = firestore.Increment(command.Value)
		}
		updates = append(updates, firestoreUpdate)
	}
	return updates
}
