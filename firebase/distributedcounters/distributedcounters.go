package distributedcounters

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/sabriboughanmi/go_utils/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"strconv"
	"time"
)

const (
	documentID       ShardField = "did" //This field is used to Track/Order Shards by their Parent Document for roll-up Process
	lastRollUpUpdate ShardField = "lru" // This Field is used to Track Last roll-up Updates to skip un-updated Shards
)

var (
	internalKeys = []ShardField{documentID, lastRollUpUpdate}
)

func checkInternalFieldsUsage(field ShardField) {
	for _, key := range internalKeys {
		if key == field {
			panic(fmt.Sprintf("DistributedCounters Package do not Allow the usage of internal Keys: %v", internalKeys))
		}
	}
}

var (
	//Errors
	NoShardFieldSpecified = errors.New("no Shard Fields Specified")
)

type ShardField string

type DistributedShard map[ShardField]interface{}

// distributedCounter is a collection of documents (shards)
// to realize counter with high frequency.
//This Struct will be created by every Incremental Section (Videos Likes, Comments Likes ..)
type distributedCounter struct {
	ShardName             string
	NumShards             int
	shardFields           DistributedShard
	defaultShardStructure interface{}
}

//CreateDistributedCounter returns a CreateDistributedCounter to manage Shards
func CreateDistributedCounter(shardName string, numShards int, defaultShard interface{}) distributedCounter {
	return distributedCounter{
		ShardName:             shardName,
		NumShards:             numShards,
		shardFields:           make(map[ShardField]interface{}),
		defaultShardStructure: defaultShard,
	}
}

//AddFieldForUpdate Adds a Shard.Field for updated
func (c *distributedCounter) AddFieldForUpdate(field ShardField, value interface{}) {
	checkInternalFieldsUsage(field)
	c.shardFields[field] = value
}

//Increments a Shard.Field for updated.
//Note! The Shard
//The supported values are:
//   int, int8, int16, int32, int64
//   uint8, uint16, uint32
//   float32, float64
func (c *distributedCounter) IncrementField(field ShardField, value interface{}) {
	checkInternalFieldsUsage(field)
	c.shardFields[field] = firestore.Increment(value)
}

// CreateShards creates a given number of shards as sub-collection of the specified document.
//(This operation need to be done once per Document or it will reinitialize all shards Data )
func (c *distributedCounter) CreateShards(ctx context.Context, docRef *firestore.DocumentRef, shardData interface{}) error {
	colRef := docRef.Collection(c.ShardName)

	// Initialize each shard with count=0
	for num := 0; num < c.NumShards; num++ {
		if _, err := colRef.Doc(strconv.Itoa(num)).Set(ctx, shardData); err != nil {
			return err
		}
	}
	return nil
}

// UpdateCounters updates a randomly picked shard of a Document.
//If no ShardField specified, an NoShardFieldSpecified will be returned
func (c *distributedCounter) UpdateCounters(ctx context.Context, docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	updateCount := len(c.shardFields)
	if updateCount == 0 {
		return nil, NoShardFieldSpecified
	}
	rand.Seed(time.Now().UnixNano())
	docID := strconv.Itoa(rand.Intn(c.NumShards))
	shardRef := docRef.Collection(c.ShardName).Doc(docID)

	//preallocate the slice for performance reasons
	updatedFields := make([]firestore.Update, updateCount+2)
	index := 0
	for key, value := range c.shardFields {
		updatedFields[index] = firestore.Update{
			Path:  string(key),
			Value: value,
		}
		index++
	}

	//Add DocumentID for roll-up Updates
	updatedFields[index] = firestore.Update{
		Path:  string(documentID),
		Value: docRef.ID,
	}

	//Add DocumentID for roll-up Updates
	updatedFields[index+1] = firestore.Update{
		Path:  string(lastRollUpUpdate),
		Value: time.Now(),
	}

	wr, err := shardRef.Update(ctx, updatedFields)

	//Create New Shard if not existing (add missing Default Fields)
	if status.Code(err) == codes.NotFound {
		defaultStructure := make(map[ShardField]interface{})
		if err = utils.InterfaceAs(c.defaultShardStructure, &defaultStructure); err != nil {
			return nil, fmt.Errorf("error mapping default shard structure for creation!: %v", err)
		}

		//Update Fields in defaultStructure
		for i:=0; i<len(updatedFields); i++{
			updatedField := updatedFields[i]
			defaultStructure[ShardField(updatedField.Path)] = updatedField.Value
		}

		return shardRef.Set(ctx, updatedFields)
	}

	return wr, err
}

/*
// UpdateCounter increments a randomly picked shard.
func (c *distributedCounter) UpdateCounter(ctx context.Context, docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	rand.Seed(time.Now().UnixNano())

	docID := strconv.Itoa(rand.Intn(c.NumShards))

	shardRef := docRef.Collection(c.ShardName).Doc(docID)
	return shardRef.Update(ctx, []firestore.Update{
		{Path: "c", Value: firestore.Increment(1)},
		{Path: "lu", Value: time.Now()},
	})
}
*/
