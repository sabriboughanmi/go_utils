package fetchbatch

//NewIncrementCommand : Safely Creates an Increment firestoreUpdateCommand
func NewIncrementCommand(path string, value interface{}) firestoreUpdateCommand {
	return firestoreUpdateCommand{
		commandType: FirestoreCommand_Increment,
		path:        path,
		value:       value,
	}
}

//NewSetCommand : Safely Creates a Set/Update firestoreUpdateCommand
func NewSetCommand(path string, value interface{}) firestoreUpdateCommand {
	return firestoreUpdateCommand{
		commandType: FirestoreCommand_Set,
		path:        path,
		value:       value,
	}
}

//NewArrayInsertElementCommand : Safely Creates an Array Insert Element firestoreUpdateCommand
func NewArrayInsertElementCommand(path string, element interface{}) firestoreUpdateCommand {
	return firestoreUpdateCommand{
		commandType: FirestoreCommand_ArrayInsert,
		path:        path,
		value:       element,
	}
}

//NewArrayRemoveElementCommand : Safely Creates an Array Remove Element firestoreUpdateCommand
func NewArrayRemoveElementCommand(path string, element interface{}) firestoreUpdateCommand {
	return firestoreUpdateCommand{
		commandType: FirestoreCommand_ArrayRemove,
		path:        path,
		value:       element,
	}
}
