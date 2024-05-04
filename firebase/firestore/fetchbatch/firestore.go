package fetchbatch

import (
	"cloud.google.com/go/firestore"
	numbersUtils "github.com/sabriboughanmi/go_utils/utils"
)

//AddCommand adds a firestoreUpdateCommand to the Queue
func (firestoreUpdates *FirestoreUpdatesQueue) AddCommand(command firestoreUpdateCommand) {
	firestoreUpdates.CommandsQueue = append(firestoreUpdates.CommandsQueue, command)
}

//AddCommands adds FirestoreUpdateCommands to the Queue
func (firestoreUpdates *FirestoreUpdatesQueue) AddCommands(command ...firestoreUpdateCommand) {
	firestoreUpdates.CommandsQueue = append(firestoreUpdates.CommandsQueue, command...)
}

//Merge fetch all firestoreUpdateCommand(s) from passed FirestoreUpdatesQueue(s)
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

//GetFirestoreUpdates merges all firestoreUpdateCommand(s) and returns an []firestore.Update
func (firestoreUpdates *FirestoreUpdatesQueue) GetFirestoreUpdates() []firestore.Update {

	var mergedCommands = make(map[string]firestoreUpdateCommand)

	for _, command := range firestoreUpdates.CommandsQueue {

		//Handle ArrayInsert Commands
		if command.commandType == FirestoreCommand_ArrayInsert {
			if oldCommand, found := mergedCommands[command.path]; !found {
				var newCommand = firestoreUpdateCommand{
					path:        command.path,
					commandType: FirestoreCommand_ArrayInsert,
					value:       []interface{}{command.value},
				}
				mergedCommands[command.path] = newCommand
			} else {
				var interfaceArray = oldCommand.value.([]interface{})
				interfaceArray = append(interfaceArray, command.value)
			}
			continue
		}

		//Handle ArrayRemove Commands
		if command.commandType == FirestoreCommand_ArrayRemove {

			if oldCommand, found := mergedCommands[command.path]; !found {
				var newCommand = firestoreUpdateCommand{
					path:        command.path,
					commandType: FirestoreCommand_ArrayRemove,
					value:       []interface{}{command.value},
				}
				mergedCommands[command.path] = newCommand
			} else {
				var interfaceArray = oldCommand.value.([]interface{})
				interfaceArray = append(interfaceArray, command.value)
			}
			continue
		}

		//Handle Field Set/Update Commands
		if oldCommand, found := mergedCommands[command.path]; !found || command.commandType == FirestoreCommand_Set {
			mergedCommands[command.path] = command
		} else {
			if numbersUtils.IsInteger(oldCommand.value) {
				oldValue, _ := numbersUtils.ToInt64(oldCommand.value)
				newValue, _ := numbersUtils.ToInt64(command.value)
				oldCommand.value = oldValue + newValue

			} else if numbersUtils.IsFloatingPointNumber(oldCommand.value) {
				oldValue, _ := numbersUtils.ToFloat64(oldCommand.value)
				newValue, _ := numbersUtils.ToFloat64(command.value)
				oldCommand.value = oldValue + newValue
			}

			mergedCommands[command.path] = oldCommand
		}
	}

	var updates []firestore.Update
	for path, command := range mergedCommands {
		var firestoreUpdate = firestore.Update{
			Path: path,
		}

		switch command.commandType {
		case FirestoreCommand_Set:
			firestoreUpdate.Value = command.value
			break
		case FirestoreCommand_Increment:
			firestoreUpdate.Value = firestore.Increment(command.value)
			break
		case FirestoreCommand_ArrayInsert:
			var interfaceArray = command.value.([]interface{})
			firestoreUpdate.Value = firestore.ArrayUnion(interfaceArray...)
			break
		case FirestoreCommand_ArrayRemove:
			var interfaceArray = command.value.([]interface{})
			firestoreUpdate.Value = firestore.ArrayRemove(interfaceArray...)
			break

		}

		updates = append(updates, firestoreUpdate)
	}
	return updates
}
