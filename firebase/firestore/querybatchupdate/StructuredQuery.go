package querybatchupdate

import (
	fr "cloud.google.com/go/firestore"
	"encoding/base64"
	"fmt"
	"google.golang.org/api/firestore/v1"
	"reflect"
)

//GetValue Returns a firestore.AggregationResult.Value using a key.
func GetValue(result firestore.AggregationResult, valKey string) (*firestore.Value, bool) {
	val, found := result.AggregateFields[valKey]
	if !found {
		return nil, false
	}
	return &val, true
}

//asFirestoreValue transforms any supported type to a firestore.Value
func asFirestoreValue(v interface{}) (*firestore.Value, error) {
	firestoreValue := &firestore.Value{}

	valType := reflect.TypeOf(v)
	kind := valType.Kind()
	value := reflect.ValueOf(v)

	//Handle firestore specific types conversions
	switch reflect.TypeOf(value) {
	case reflect.TypeOf(firestore.LatLng{}):
		var ptr = v.(firestore.LatLng)
		firestoreValue.GeoPointValue = &ptr
		return firestoreValue, nil

	case reflect.TypeOf([]byte{}):
		firestoreValue.BytesValue = base64.StdEncoding.EncodeToString(v.([]byte))
		return firestoreValue, nil
	}

	switch kind {
	case reflect.Bool:
		firestoreValue.BooleanValue = v.(bool)
	case reflect.Int:
		firestoreValue.IntegerValue = int64(v.(int))
	case reflect.Int8:
		firestoreValue.IntegerValue = int64(v.(int8))
	case reflect.Int16:
		firestoreValue.IntegerValue = int64(v.(int16))
	case reflect.Int32:
		firestoreValue.IntegerValue = int64(v.(int32))
	case reflect.Int64:
		firestoreValue.IntegerValue = v.(int64)
	case reflect.Uint:
		firestoreValue.IntegerValue = int64(v.(uint))
	case reflect.Uint8:
		firestoreValue.IntegerValue = int64(v.(uint8))
	case reflect.Uint16:
		firestoreValue.IntegerValue = int64(v.(uint16))
	case reflect.Uint32:
		firestoreValue.IntegerValue = int64(v.(uint32))
	case reflect.Uint64:
		firestoreValue.IntegerValue = int64(v.(uint64))
	case reflect.Float32:
		firestoreValue.DoubleValue = float64(v.(float32))
	case reflect.Float64:
		firestoreValue.DoubleValue = v.(float64)
	case reflect.String:
		firestoreValue.StringValue = v.(string)
	case reflect.Struct:
		var firestoreMap = firestore.MapValue{
			Fields: make(map[string]firestore.Value),
		}
		//Iterate over fields
		for i := 0; i < valType.NumField(); i++ {
			//Process Field
			var val, err = asFirestoreValue(value.Field(i).Interface())
			if err != nil {
				return nil, err
			}

			//Set the safe data
			firestoreMap.Fields[valType.Field(i).Name] = *val

		}
		firestoreValue.MapValue = &firestoreMap
	case reflect.Array, reflect.Slice:
		//Check null
		if value.IsNil() {
			firestoreValue.NullValue = "NULL_VALUE"
			break
		}
		//Initialize the Slice
		var arrayValue = firestore.ArrayValue{
			Values: make([]*firestore.Value, value.Len()),
		}
		//fill the slice
		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)
			elementVal, err := asFirestoreValue(item.Interface())
			if err != nil {
				return nil, err
			}
			arrayValue.Values[i] = elementVal
		}
		//set the Array Value
		firestoreValue.ArrayValue = &arrayValue
	case reflect.Map:
		//Check null
		if value.IsNil() {
			firestoreValue.NullValue = "NULL_VALUE"
			break
		}
		//Initialize the firestore MapValue
		var mapValue = firestore.MapValue{
			Fields: make(map[string]firestore.Value),
		}

		//Process map fields
		iter := value.MapRange()
		for iter.Next() {
			k := iter.Key()
			//Check if Key is a string
			if k.Kind() != reflect.String {
				return firestoreValue, fmt.Errorf("firestore doesn't support Maps using keys other than strings")
			}

			//Process the Value
			val, err := asFirestoreValue(iter.Value().Interface())
			if err != nil {
				return nil, err
			}
			mapValue.Fields[k.String()] = *val

		}

		firestoreValue.MapValue = &mapValue
	case reflect.Ptr:
		//Check null
		if value.IsNil() {
			firestoreValue.NullValue = "NULL_VALUE"
			break
		}
		//Process Field
		return asFirestoreValue(value.Elem().Interface())
	}
	return firestoreValue, nil
}

func (contentBatchUpdate *ContentBatchUpdate) createStructuredQuery() (*firestore.StructuredQuery, error) {
	structuredQuery := &firestore.StructuredQuery{
		From:    []*firestore.CollectionSelector{{CollectionId: contentBatchUpdate.querySearchParams.CollectionID, AllDescendants: false}},
		OrderBy: createOrderBy(contentBatchUpdate.querySearchParams.QuerySorts),
		Offset:  contentBatchUpdate.querySearchParams.Offset,
		Select: &firestore.Projection{
			Fields: func() []*firestore.FieldReference {

				//Return the Default documentID
				if contentBatchUpdate.querySearchParams.SelectFields == nil || len(contentBatchUpdate.querySearchParams.SelectFields) == 0 {
					return []*firestore.FieldReference{
						{FieldPath: "__name__"},
					}
				}
				//Init the selected fields slice
				var fieldsToSelect = make([]*firestore.FieldReference, len(contentBatchUpdate.querySearchParams.SelectFields)+1)

				//Set the first element
				fieldsToSelect[0] = &firestore.FieldReference{FieldPath: "__name__"}

				//Set the Other Fields
				for i, fieldPath := range contentBatchUpdate.querySearchParams.SelectFields {
					fieldsToSelect[i+1] = &firestore.FieldReference{FieldPath: fieldPath}
				}

				return fieldsToSelect
			}(),
		},
		StartAt: nil,
		EndAt:   nil,
	}

	//set the Limits if exists
	if contentBatchUpdate.queryPaginationParams.Limit != nil {
		structuredQuery.Limit = int64(*contentBatchUpdate.queryPaginationParams.Limit)
	}

	//set the Field/Composite Filters
	if contentBatchUpdate.querySearchParams.QueryWheres != nil && len(contentBatchUpdate.querySearchParams.QueryWheres) > 0 {
		where, err := createFieldOrCompositeFilter(contentBatchUpdate.querySearchParams.QueryWheres)
		if err != nil {
			return nil, err
		}
		//Set the Where param
		structuredQuery.Where = where
	}

	//Set select fields if requests
	if contentBatchUpdate.querySearchParams.SelectFields != nil {
		//Initialize the []*firestore.FieldReference
		var fields = make([]*firestore.FieldReference, len(contentBatchUpdate.querySearchParams.SelectFields))

		//fill the []*firestore.FieldReference
		for i, field := range contentBatchUpdate.querySearchParams.SelectFields {
			fields[i] = &firestore.FieldReference{
				FieldPath: field,
			}
		}
		projection := firestore.Projection{
			Fields: fields,
		}
		structuredQuery.Select = &projection
	}

	// Set StartAt with Start Doc
	if contentBatchUpdate.startDoc != nil {
		var cursor = firestore.Cursor{
			Before: contentBatchUpdate.startBefore,
			Values: make([]*firestore.Value, 1),
		}
		val, err := asFirestoreValue(contentBatchUpdate.startDoc)
		if err != nil {
			return nil, err
		}

		cursor.Values[0] = val
		//Set StartDoc Cursor
		structuredQuery.StartAt = &cursor
		// Set StartAt with Start Values
	} else if contentBatchUpdate.startVals != nil {

		var cursor = firestore.Cursor{
			Before: contentBatchUpdate.startBefore,
			Values: make([]*firestore.Value, len(contentBatchUpdate.startVals)),
		}

		//Collect the Cursor values and convert them to Firestore Value
		for i, val := range contentBatchUpdate.startVals {
			firestoreVal, err := asFirestoreValue(val)
			if err != nil {
				return nil, err
			}
			cursor.Values[i] = firestoreVal
		}
		//Set Start Values
		structuredQuery.StartAt = &cursor
	}

	//TODO: Set EndAt
	return structuredQuery, nil
}

//createOrderBy creates a []*firestore.Order from Query Sorts
func createOrderBy(querySorts []QuerySort) []*firestore.Order {
	orders := make([]*firestore.Order, len(querySorts))
	for i, querySort := range querySorts {
		orders[i] = &firestore.Order{
			Field: &firestore.FieldReference{
				FieldPath: querySort.DocumentSortKey,
			},
			Direction: func() string {
				if querySort.Direction == fr.Desc {
					return "DESCENDING"
				}
				return "Ascending"
			}(),
		}
	}
	return orders
}

//Create a Composite of Field filter
func createFieldOrCompositeFilter(queryWhereKeys []QueryWhere) (*firestore.Filter, error) {
	if len(queryWhereKeys) == 1 {
		return createSingleFilter(queryWhereKeys[0])
	}
	filters := make([]*firestore.Filter, len(queryWhereKeys))
	for i, queryWhere := range queryWhereKeys {
		filter, err := createSingleFilter(queryWhere)
		if err != nil {
			return nil, err
		}
		//Set the Filter
		filters[i] = filter
	}
	return &firestore.Filter{
		CompositeFilter: &firestore.CompositeFilter{
			Op:      "AND",
			Filters: filters,
		},
	}, nil
}

//Create a firestore.Filter from a QueryWhere model.
func createSingleFilter(queryWhere QueryWhere) (*firestore.Filter, error) {
	value, err := asFirestoreValue(queryWhere.Value)
	if err != nil {
		return nil, err
	}
	// Return a new Firestore Filter with the path, operator, and value set according to the QueryWhere
	return &firestore.Filter{
		FieldFilter: &firestore.FieldFilter{
			Field: &firestore.FieldReference{
				FieldPath: queryWhere.Path,
			},
			Op:    structuredQueryOperator[queryWhere.Op],
			Value: value,
		},
	}, nil
}
