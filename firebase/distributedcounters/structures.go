package distributedcounters

type ShardField string

//DistributedCounters used to initialize Distributed Counters
type DistributedCounters struct {
	ShardCount int
	ShardName  string
	RollUpTime int64 //how many seconds before the next rollup
}

//shardStructure is the structure in which the Shard is saved
//Please make Sure to Change _shardStructureKeys Instance if Json Keys are modified
type shardStructure struct {
	Ints         map[string]int64   `json:"i,omitempty" firestore:"i,omitempty"`
	Floats       map[string]float64 `json:"f,omitempty" firestore:"f,omitempty"`
	DocumentID   string             `json:"di_0,omitempty" firestore:"di_0,omitempty"`
	CreationTick int64              `json:"ct_0,omitempty" firestore:"ct_0,omitempty"`
	CursorID     string             `json:"cd_0,omitempty" firestore:"cd_0,omitempty"`
}
//shardStructureKeys Represents json Keys of shardStructure
type shardStructureKeys struct {
	Ints         string
	Floats       string
	DocumentID   string
	CreationTick string
	CursorID     string
}


//DistributedCounterInstance is a collection of documents (shards)
//to realize counter with high frequency.
//This Struct will be created by every Incremental Section (Videos Likes, Comments Likes ..)
type DistributedCounterInstance struct {
	shardName   string
	numShards   int
	shardFields shardStructure
	rollUpTime  int64
}


