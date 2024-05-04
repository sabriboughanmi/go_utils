package fetchbatch

//NewIncrementCommand : Safely Creates an Increment FirestoreUpdateCommand
func NewIncrementCommand(path string, value interface{}) FirestoreUpdateCommand {
	return FirestoreUpdateCommand{
		commandType: FirestoreCommand_Increment,
		path:        path,
		value:       value,
	}
}

//NewSetCommand : Safely Creates a Set/Update FirestoreUpdateCommand
func NewSetCommand(path string, value interface{}) FirestoreUpdateCommand {
	return FirestoreUpdateCommand{
		commandType: FirestoreCommand_Set,
		path:        path,
		value:       value,
	}
}

//NewArrayInsertElementCommand : Safely Creates an Array Insert Element FirestoreUpdateCommand
func NewArrayInsertElementCommand(path string, element interface{}) FirestoreUpdateCommand {
	return FirestoreUpdateCommand{
		commandType: FirestoreCommand_ArrayInsert,
		path:        path,
		value:       element,
	}
}

//NewArrayRemoveElementCommand : Safely Creates an Array Remove Element FirestoreUpdateCommand
func NewArrayRemoveElementCommand(path string, element interface{}) FirestoreUpdateCommand {
	return FirestoreUpdateCommand{
		commandType: FirestoreCommand_ArrayRemove,
		path:        path,
		value:       element,
	}
}
