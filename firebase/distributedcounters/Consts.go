package distributedcounters

import "errors"


var (
	//Errors
	NoShardFieldSpecified = errors.New("no Shard Fields Specified")
)


var (
	//Internal Keys
	creationTick ShardField = "ct_0" // This Field is used to Track which shard has exceeded to rollup time
	cursorID     ShardField = "cd_0" // This Field is used to Track which shard has exceeded to rollup time
	internalKeys = []ShardField{creationTick, cursorID}
)



var (
	//Keys
	_shardStructureKeys = shardStructureKeys{
		Ints:         "i",
		Floats:       "f",
		CreationTick: string(creationTick),
		CursorID:     string(cursorID),
	}
)
