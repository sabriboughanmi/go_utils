package distributedcounters

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type onShardsUpdateCompleted func(document *firestore.DocumentRef, int64Updates map[string]int64, float64Updates map[string]float64)
type onShardsUpdateFailed func(err error, shards []*firestore.DocumentSnapshot)

//RollUpShards Shards of a specific Document,
//Warning! If an array of DocumentSnapshots are passed with multiple parents the first parent will get updated by all Shards.
func rollUpShards(waitGroup *sync.WaitGroup, client *firestore.Client, ctx context.Context, onShardsCompletedUpdate onShardsUpdateCompleted, onShardsUpdateFailed onShardsUpdateFailed, shards ...*firestore.DocumentSnapshot) {
	defer waitGroup.Done()

	if shards == nil || len(shards) == 0 {
		fmt.Printf("Error Query: %v \n", fmt.Errorf("no documents to process"))
		return
	}

	batch := client.Batch()

	//Collect Data from Shards
	incrementalIntFields := make(map[string]int64)
	incrementalFloatFields := make(map[string]float64)
	for i := 0; i < len(shards); i++ {

		//Cache the doc for performance reasons
		doc := shards[i]

		var shardStructure shardStructure
		//Collect Data
		if err := doc.DataTo(&shardStructure); err != nil {
			if onShardsUpdateFailed != nil {
				onShardsUpdateFailed(err, shards)
			}
			return
		}

		//Add to delete batch
		batch.Delete(doc.Ref)

		//Sum Ints
		for key, value := range shardStructure.Ints {
			//Skip internal Keys
			if isInternalFields(ShardField(key)) {
				continue
			}
			incrementalIntFields[key] += value
		}
		//Sum floats
		for key, value := range shardStructure.Floats {
			//Skip internal Keys
			if isInternalFields(ShardField(key)) {
				continue
			}
			incrementalFloatFields[key] += value
		}
	}

	//PreInitialize the valuesToUpdate in GC
	var valuesToUpdate = make([]firestore.Update, len(incrementalIntFields)+len(incrementalFloatFields))
	var updateIndex = 0

	//Create a []firestore.update from incremental Ints
	for key, value := range incrementalIntFields {
		valuesToUpdate[updateIndex] = firestore.Update{
			Path:  key,
			Value: firestore.Increment(value),
		}
		updateIndex++
	}

	//Create a []firestore.update from incremental Floats
	for key, value := range incrementalFloatFields {
		valuesToUpdate[updateIndex] = firestore.Update{
			Path:  key,
			Value: firestore.Increment(value),
		}
		updateIndex++
	}

	//Skip Update if Values there are no Changes
	if len(valuesToUpdate) == 0 {
		return
	}

	var parentDocRef = shards[0].Ref.Parent.Parent

	//Add Parent.Document update to the write batch
	batch.Update(parentDocRef, valuesToUpdate)

	//Delete Shards
	if _, err := batch.Commit(ctx); err != nil {

		/*
			//Parent Document does not exist anymore
			if status.Code(err) == codes.NotFound {
				//Clear all Sharda and omit parent document update
				batch = client.Batch()

				fmt.Printf("Parent Document '%s' doesn't exist anymore : ")
				for _, shardDocRef := range shards {
					batch.Delete(shardDocRef.Ref)
				}

				//Commit Shards delete operations
				if _, err := batch.Commit(ctx); err != nil {
					onShardsUpdateFailed(err, shards)
					return
				}
				return
			}
		*/
		onShardsUpdateFailed(err, shards)
		return
	}

	//all shards are updates successfully
	if onShardsCompletedUpdate != nil {
		onShardsCompletedUpdate(shards[0].Ref.Parent.Parent, incrementalIntFields, incrementalFloatFields)
	}

	return
}

/*
//RollUp Shards of a specific Document.
//
//Warning! If an array of DocumentSnapshots is passed with multiple parents the first parent will get updated by all Shards
func rollUpShards(client *firestore.Client, ctx context.Context, shards ...*firestore.DocumentSnapshot) error {
	if shards == nil || len(shards) == 0 {
		return fmt.Errorf("no documents to process")
	}

	batch := client.Batch()

	//Collect Data from Shards
	incrementalIntFields := make(map[string]int64)
	incrementalFloatFields := make(map[string]float64)
	for i := 0; i < len(shards); i++ {

		//Cache the doc for performance reasons
		doc := shards[i]

		var shardStructure shardStructure
		//Collect Data
		if err := doc.DataTo(&shardStructure); err != nil {
			return err
		}

		//Add to delete batch
		batch.Delete(doc.Ref)

		//Sum Ints
		for key, value := range shardStructure.Ints {
			//Skip internal Keys
			if isInternalFields(ShardField(key)) {
				continue
			}
			incrementalIntFields[key] += value
		}
		//Sum floats
		for key, value := range shardStructure.Floats {
			//Skip internal Keys
			if isInternalFields(ShardField(key)) {
				continue
			}
			incrementalFloatFields[key] += value
		}
	}


	var valuesToUpdate []firestore.Update

	//Collect incremental Ints
	for key, value := range incrementalIntFields {
		valuesToUpdate = append(valuesToUpdate, firestore.Update{
			path:  key,
			Value: firestore.Increment(value),
		})
	}

	//Collect incremental Floats
	for key, value := range incrementalFloatFields {
		valuesToUpdate = append(valuesToUpdate, firestore.Update{
			path:  key,
			Value: firestore.Increment(value),
		})
	}

	//Skip Update if Values are
	if len(valuesToUpdate) == 0 {
		return nil
	}

	//Update Fields in Parent document
	_, err := shards[0].Ref.Parent.Parent.Update(ctx, valuesToUpdate)
	if err != nil {
		return err
	}

	//Delete Shards
	_, err = batch.Commit(ctx)
	return err
}

*/

//ParallelRollUp RollUP all documents Shards relative to the DistributedCounters.ShardName.
//
//This function Executes multiple RollUps in parallel. (parallelDocumentsCount will be multiplied by the ShardCount and used as Query Limiter).
//
//If filterByTicks == true, Shards creation time will be ignored. (useful to update olf shards)
func (dc *DistributedCounters) ParallelRollUp(client *firestore.Client, ctx context.Context, parallelDocumentsCount int, filterByTicks bool, onCompleted onShardsUpdateCompleted, onFailed onShardsUpdateFailed) error {
	wg := sync.WaitGroup{}

	//Wait for the execution of RollUps to finish even if some RollUps have failed.
	//Note! firestore Document operations only occur if the RollUp is successful.
	defer wg.Wait()

	queryLimiter := dc.ShardCount * parallelDocumentsCount

	var ticks []int64

	if filterByTicks {
		currentTick := time.Now().Unix() / dc.RollUpTime
		ticks = make([]int64, 10)
		var i int64
		for i = 0; i < 10; i++ {
			ticks[i] = currentTick - i
		}
	}

	//Loop Managers
	var cursor *firestore.DocumentSnapshot = nil
	var shardsInQueue []*firestore.DocumentSnapshot
	var moreShardsExists = true

	for moreShardsExists {
		var query firestore.Query
		if cursor != nil {
			query = client.CollectionGroup(dc.ShardName).OrderBy(key_shardStructureModel.CursorID, firestore.Asc)
			//Filter with Ticks
			if filterByTicks {
				query = query.Where(key_shardStructureModel.CreationTick, "in", ticks)
			}
			query = query.StartAfter(cursor.Data()[key_shardStructureModel.CursorID]).Limit(queryLimiter)
		} else {
			query = client.CollectionGroup(dc.ShardName).OrderBy(key_shardStructureModel.CursorID, firestore.Asc)
			//Filter with Ticks
			if filterByTicks {
				query = query.Where(key_shardStructureModel.CreationTick, "in", ticks)
			}
			query.Limit(queryLimiter)
		}
		it := query.Documents(ctx)
		newShards, err := it.GetAll()
		if err != nil {
			return err
		}

		//Prepare Exit/Cursor
		if len(newShards) < dc.ShardCount {
			moreShardsExists = false
		} else {
			cursor = newShards[len(newShards)-1]
		}

		//Append new Shards
		shardsInQueue = append(shardsInQueue, newShards...)

		firstElementToProcess := 0
		//Process Shards Queue
		for i := 0; i < len(shardsInQueue); i++ {

			//Last Shard in Queue
			if i+1 == len(shardsInQueue) {
				if moreShardsExists {
					//Remove Processed Shards from shardsInQueue
					shardsInQueue = shardsInQueue[firstElementToProcess : i+1]
					break
				}

				//Process Remaining Shards and quit
				wg.Add(1)
				go rollUpShards(&wg, client, ctx, onCompleted, onFailed, shardsInQueue[firstElementToProcess:i+1]...)
				return nil
			}

			//Skip if Parent Still Same
			if shardsInQueue[i].Ref.Parent.Parent.ID == shardsInQueue[i+1].Ref.Parent.Parent.ID {
				continue
			}

			//Shard Parent Changed
			//Process Shards
			wg.Add(1)
			fmt.Println("parallelRollUpShards Executed")

			go rollUpShards(&wg, client, ctx, onCompleted, onFailed, shardsInQueue[firstElementToProcess:i+1]...)
			firstElementToProcess = i + 1
		}
	}

	return nil
}

//ParallelRollUp Collects data from a shard document and updates it's parent document.
func singleShardRollUp(shardDoc *firestore.DocumentSnapshot, ctx context.Context, onShardsCompletedUpdate onShardsUpdateCompleted) error {

	var shardStructure shardStructure
	//Collect Data
	if err := shardDoc.DataTo(&shardStructure); err != nil {
		return err
	}

	var valuesToUpdate []firestore.Update

	//Collect incremental Ints
	for key, value := range shardStructure.Ints {
		valuesToUpdate = append(valuesToUpdate, firestore.Update{
			Path:  key,
			Value: firestore.Increment(value),
		})
	}

	//Collect incremental Floats
	for key, value := range shardStructure.Floats {
		valuesToUpdate = append(valuesToUpdate, firestore.Update{
			Path:  key,
			Value: firestore.Increment(value),
		})
	}

	//Skip Update if Values are
	if len(valuesToUpdate) == 0 {
		return nil
	}

	//Update Fields in Parent document
	if _, err := shardDoc.Ref.Parent.Parent.Update(ctx, valuesToUpdate); err != nil {
		return err
	}

	if onShardsCompletedUpdate != nil {
		//Execute the lambda passed with Shard Document parent  and collected data
		onShardsCompletedUpdate(shardDoc.Ref.Parent.Parent, shardStructure.Ints, shardStructure.Floats)
	}

	return nil
}

//CreateDistributedCounter returns a CreateDistributedCounter to manage Shards
func (dc *DistributedCounters) CreateDistributedCounter() DistributedCounterInstance {
	return DistributedCounterInstance{
		shardName: dc.ShardName,
		numShards: dc.ShardCount,
		shardFields: shardStructure{
			Ints:   make(map[string]int64),
			Floats: make(map[string]float64),
		},
		rollUpTime: dc.RollUpTime,
	}
}

//IncrementField Increments a ShardField for updated.
//Note! The Shard supported values are:
//   int, int8, int16, int32, int64
//   uint8, uint16, uint32
//   float32, float64
func (c *DistributedCounterInstance) IncrementField(field ShardField, value interface{}) {
	checkInternalFieldsUsage(field)
	switch value.(type) {
	case int:
		c.shardFields.Ints[string(field)] += int64(value.(int))
		break

	case int8:
		c.shardFields.Ints[string(field)] += int64(value.(int8))
		break

	case int32:
		c.shardFields.Ints[string(field)] += int64(value.(int32))
		break

	case int64:
		c.shardFields.Ints[string(field)] += value.(int64)
		break

	case uint:
		c.shardFields.Ints[string(field)] += int64(value.(uint))
		break

	case uint8:
		c.shardFields.Ints[string(field)] += int64(value.(uint8))
		break

	case uint16:
		c.shardFields.Ints[string(field)] += int64(value.(uint16))
		break

	case uint32:
		c.shardFields.Ints[string(field)] += int64(value.(uint32))
		break

	case uint64:
		c.shardFields.Ints[string(field)] += int64(value.(uint64))
		break

		//Handle Floats
	case float32:
		c.shardFields.Floats[string(field)] += float64(value.(float32))
		break
	case float64:
		c.shardFields.Floats[string(field)] += value.(float64)
		break
	default:
		panic("IncrementField supported values are: int, int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32,float64!")
	}
}

// UpdateCounters updates a randomly picked shard of a Document.
//If no ShardField specified, an NoShardFieldSpecified will be returned
func (c *DistributedCounterInstance) UpdateCounters(ctx context.Context, docRef *firestore.DocumentRef) (*firestore.WriteResult, error) {
	updateCount := len(c.shardFields.Ints) + len(c.shardFields.Floats)

	if updateCount == 0 {
		return nil, NoShardFieldSpecified
	}
	rand.Seed(time.Now().UnixNano())
	docID := strconv.Itoa(rand.Intn(c.numShards))
	shardRef := docRef.Collection(c.shardName).Doc(docID)

	//preallocate the slice for performance reasons
	updatedFields := make([]firestore.Update, len(c.shardFields.Floats)+len(c.shardFields.Ints))

	index := 0
	for key, value := range c.shardFields.Floats {
		updatedFields[index] = firestore.Update{
			Path:  key_shardStructureModel.Floats + "." + key,
			Value: firestore.Increment(value),
		}
		index++
	}

	for key, value := range c.shardFields.Ints {
		updatedFields[index] = firestore.Update{
			Path:  key_shardStructureModel.Ints + "." + key,
			Value: firestore.Increment(value),
		}
		index++
	}

	wr, err := shardRef.Update(ctx, updatedFields)

	//Create New Shard if not existing (add missing Default Fields)
	if status.Code(err) == codes.NotFound {

		//Add Next Tick for roll-up.
		c.shardFields.CreationTick = (time.Now().Unix() / c.rollUpTime) + 2

		//Add DocumentID for roll-up Updates
		c.shardFields.DocumentID = docRef.ID

		//Add Pagination Cursor
		c.shardFields.CursorID = docRef.ID + "_" + docID

		return shardRef.Set(ctx, c.shardFields)
	}

	return wr, err
}
