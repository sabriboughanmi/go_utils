package distributedcounters

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/sabriboughanmi/go_utils/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type onShardsCompletedUpdate func(document *firestore.DocumentRef, int64Updates map[string]int64, float64Updates map[string]float64)

//RollUp Shards of a specific Document,
//Warning! If an array of DocumentSnapshots is passed with multiple parents the first parent will get updated by all Shards
func parallelRollUpShards(waitGroup *sync.WaitGroup, errorChan chan error, client *firestore.Client, ctx context.Context, onShardsCompletedUpdate onShardsCompletedUpdate, shards ...*firestore.DocumentSnapshot) {
	defer waitGroup.Done()

	if shards == nil || len(shards) == 0 {
		fmt.Printf("Error Query: %v \n", fmt.Errorf("no documents to process"))
		errorChan <- fmt.Errorf("no documents to process")
		return
	}

	batch := client.Batch()
	/*//DEBUG:
	var ids []string
	*/
	//Collect Data from Shards
	incrementalIntFields := make(map[string]int64)
	incrementalFloatFields := make(map[string]float64)
	for i := 0; i < len(shards); i++ {

		//Cache the doc for performance reasons
		doc := shards[i]
		/*//DEBUG:
		ids = append(ids, fmt.Sprintf("Doc: %s, Parent:%s ", doc.Ref.ID, doc.Ref.Parent.Parent.ID))
		*/
		var shardStructure shardStructure
		//Collect Data
		if err := doc.DataTo(&shardStructure); err != nil {
			errorChan <- err
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

	/*//DEBUG:
	fmt.Printf("Batched Shards Count(%d): %v \n", len(ids), ids)
	*/
	var valuesToUpdate []firestore.Update

	//Collect incremental Ints
	for key, value := range incrementalIntFields {
		valuesToUpdate = append(valuesToUpdate, firestore.Update{
			Path:  key,
			Value: firestore.Increment(value),
		})
	}

	//Collect incremental Floats
	for key, value := range incrementalFloatFields {
		valuesToUpdate = append(valuesToUpdate, firestore.Update{
			Path:  key,
			Value: firestore.Increment(value),
		})
	}

	//Skip Update if Values there are no Changes
	if len(valuesToUpdate) == 0 {
		return
	}

	//Update Fields in Parent document

	if _, err := shards[0].Ref.Parent.Parent.Update(ctx, valuesToUpdate); err != nil {
		fmt.Printf("Error Updating Parent: %v \n", err)
		errorChan <- err
		return
	}

	//Delete Shards
	if _, err := batch.Commit(ctx); err != nil {
		fmt.Printf("Error batch.Commit: %v \n", err)
		errorChan <- err
		return
	}

	if onShardsCompletedUpdate != nil {
		onShardsCompletedUpdate(shards[0].Ref.Parent.Parent, incrementalIntFields, incrementalFloatFields)
	}

	return
}

//RollUp Shards of a specific Document,
//Warning! If an array of DocumentSnapshots is passed with multiple parents the first parent will get updated by all Shards
func rollUpShards(client *firestore.Client, ctx context.Context, shards ...*firestore.DocumentSnapshot) error {
	if shards == nil || len(shards) == 0 {
		return fmt.Errorf("no documents to process")
	}

	batch := client.Batch()
	/*//DEBUG:
	var ids []string
	*/
	//Collect Data from Shards
	incrementalIntFields := make(map[string]int64)
	incrementalFloatFields := make(map[string]float64)
	for i := 0; i < len(shards); i++ {

		//Cache the doc for performance reasons
		doc := shards[i]
		/*//DEBUG:
		ids = append(ids, fmt.Sprintf("Doc: %s, Parent:%s ", doc.Ref.ID, doc.Ref.Parent.Parent.ID))
		*/
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

	/*//DEBUG:
	fmt.Printf("Batched Shards Count(%d): %v \n", len(ids), ids)
	*/
	var valuesToUpdate []firestore.Update

	//Collect incremental Ints
	for key, value := range incrementalIntFields {
		valuesToUpdate = append(valuesToUpdate, firestore.Update{
			Path:  key,
			Value: firestore.Increment(value),
		})
	}

	//Collect incremental Floats
	for key, value := range incrementalFloatFields {
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
	_, err := shards[0].Ref.Parent.Parent.Update(ctx, valuesToUpdate)
	if err != nil {
		return err
	}

	//Delete Shards
	_, err = batch.Commit(ctx)
	return err
}

//ParallelRollUp RollUP all documents Shards relative to the DistributedCounters.ShardName
//This function Executes multiple RollUps in parallel. (parallelDocumentsCount will be multiplied by the ShardCount and used as Query Limiter)
func (dc *DistributedCounters) ParallelRollUp(client *firestore.Client, ctx context.Context, parallelDocumentsCount int, onShardsCompletedUpdate onShardsCompletedUpdate) error {
	wg := sync.WaitGroup{}

	//Wait for the execution of RollUps to finish even if some RollUps have failed.
	//Note! firestore Document operations only occur if the RollUp is successful.
	defer wg.Wait()

	queryLimiter := dc.ShardCount * parallelDocumentsCount
	currentTick := time.Now().Unix() / dc.RollUpTime
	ticks := make([]int64, 10)
	var i int64
	for i = 0; i < 10; i++ {
		ticks[i] = currentTick - i
	}

	//Loop Managers
	var cursor *firestore.DocumentSnapshot = nil
	var shardsInQueue []*firestore.DocumentSnapshot
	var moreShardsExists = true
	errorChan := make(chan error)

	for moreShardsExists {
		var query firestore.Query
		if cursor != nil {
			query = client.CollectionGroup(dc.ShardName).
				OrderBy(string(cursorID), firestore.Asc).
				Where(string(creationTick), "in", ticks).
				StartAfter(cursor.Data()[string(cursorID)]).
				Limit(queryLimiter)
		} else {
			query = client.CollectionGroup(dc.ShardName).
				OrderBy(string(cursorID), firestore.Asc).
				Where(string(creationTick), "in", ticks).
				Limit(queryLimiter)
		}
		it := query.Documents(ctx)
		newShards, err := it.GetAll()
		if err != nil {
			errorChan <- err
			return utils.HandleGoroutineErrors(&wg, errorChan)
		}

		/*//DEBUG:
		// Get the last document.
		var queueIds []string
		for m := 0; m < len(newShards); m++ {
			queueIds = append(queueIds, fmt.Sprintf("[Doc: %s, Parent:%s], ", newShards[m].Ref.ID, newShards[m].Ref.Parent.Parent.ID))
		}
		fmt.Printf("NewShards Count : (%d),  %v \n", len(newShards), queueIds)
		//Debug
		*/

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
					/*//DEBUG:
					var queueIds []string
					for m := 0; m < len(shardsInQueue); m++ {
						queueIds = append(queueIds, fmt.Sprintf("Doc: %s, Parent:%s", shardsInQueue[m].Ref.ID, shardsInQueue[m].Ref.Parent.Parent.ID))
					}
					fmt.Printf("Kept for later : %v \n", queueIds)
					*/
					break
				}

				//Process Remaining Shards and quit
				wg.Add(1)
				go parallelRollUpShards(&wg, errorChan, client, ctx, onShardsCompletedUpdate, shardsInQueue[firstElementToProcess:i+1]...)
				return utils.HandleGoroutineErrors(&wg, errorChan)
			}

			//Skip if Parent Still Same
			if shardsInQueue[i].Ref.Parent.Parent.ID == shardsInQueue[i+1].Ref.Parent.Parent.ID {
				continue
			}

			//Shard Parent Changed
			//Process Shards
			wg.Add(1)
			fmt.Println("parallelRollUpShards Executed")

			go parallelRollUpShards(&wg, errorChan, client, ctx, onShardsCompletedUpdate, shardsInQueue[firstElementToProcess:i+1]...)
			firstElementToProcess = i + 1
		}
	}

	return utils.HandleGoroutineErrors(&wg, errorChan)
}


//RollUp all documents Shards relative to the DistributedCounters.ShardName
func (dc *DistributedCounters) RollUp(client *firestore.Client, ctx context.Context) error {

	currentTick := time.Now().Unix() / dc.RollUpTime
	ticks := make([]int64, 10)
	var i int64
	for i = 0; i < 10; i++ {
		ticks[i] = currentTick - i
	}

	//Loop Managers
	var cursor *firestore.DocumentSnapshot = nil
	var shardsInQueue []*firestore.DocumentSnapshot
	var moreShardsExists = true

	for moreShardsExists {
		var query firestore.Query
		if cursor != nil {
			query = client.CollectionGroup(dc.ShardName).
				OrderBy(string(cursorID), firestore.Asc).
				Where(string(creationTick), "in", ticks).
				StartAfter(cursor.Data()[string(cursorID)]).
				Limit(dc.ShardCount)
		} else {
			query = client.CollectionGroup(dc.ShardName).
				OrderBy(string(cursorID), firestore.Asc).
				Where(string(creationTick), "in", ticks).
				Limit(dc.ShardCount)
		}
		it := query.Documents(ctx)
		newShards, err := it.GetAll()
		if err != nil {
			return err
		}

		/*//DEBUG:
		// Get the last document.
		var queueIds []string
		for m := 0; m < len(newShards); m++ {
			queueIds = append(queueIds, fmt.Sprintf("[Doc: %s, Parent:%s], ", newShards[m].Ref.ID, newShards[m].Ref.Parent.Parent.ID))
		}
		fmt.Printf("NewShards Count : (%d),  %v \n", len(newShards), queueIds)
		//Debug
		*/

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

					/*//DEBUG:
					var queueIds []string
					for m := 0; m < len(shardsInQueue); m++ {
						queueIds = append(queueIds, fmt.Sprintf("Doc: %s, Parent:%s", shardsInQueue[m].Ref.ID, shardsInQueue[m].Ref.Parent.Parent.ID))
					}
					fmt.Printf("Kept for later : %v \n", queueIds)
					*/
					break
				}

				//Process Remaining Shards and quit
				return rollUpShards(client, ctx, shardsInQueue[firstElementToProcess:i+1]...)
			}

			//Skip if Parent Still Same
			if shardsInQueue[i].Ref.Parent.Parent.ID == shardsInQueue[i+1].Ref.Parent.Parent.ID {
				continue
			}

			//Shard Parent Changed
			//Process Shards
			err = rollUpShards(client, ctx, shardsInQueue[firstElementToProcess:i+1]...)
			if err != nil {
				log.Fatal(err)
			}
			firstElementToProcess = i + 1
		}
	}

	return nil
}


//SingleShardRollUp Collects data from a shard document and updates it's parent document.
//This function is useful to safely rollup outdated shards. (can't be accessed by the standard functions RollUp, ParallelRollUp )
func SingleShardRollUp(shardDoc *firestore.DocumentSnapshot, ctx context.Context, onShardsCompletedUpdate onShardsCompletedUpdate) error {

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
			Path:  _shardStructureKeys.Floats + "." + key,
			Value: firestore.Increment(value),
		}
		index++
	}

	for key, value := range c.shardFields.Ints {
		updatedFields[index] = firestore.Update{
			Path:  _shardStructureKeys.Ints + "." + key,
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
