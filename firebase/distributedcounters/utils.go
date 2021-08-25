package distributedcounters

import "fmt"

//Return if a string is an internal field
func isInternalFields(field ShardField) bool {
	for _, key := range internalKeys {
		if key == field {
			return true
		}
	}
	return false
}

//Panic if the key is an internal ShardField
func checkInternalFieldsUsage(field ShardField) {
	for _, key := range internalKeys {
		if key == field {
			panic(fmt.Sprintf("DistributedCounters Package do not Allow the usage of internal Keys: %v", internalKeys))
		}
	}
}