package modelsfixer

type ETag string

const (
	FireStoreTag = "firestore"

	Tags_Int         ETag = "int"       //converts any number/string to int64
	Tags_Float       ETag = "float"     //converts any number/string to float64
	Tags_String      ETag = "string"    //converts anything to string
	Tags_Omitempty   ETag = "omitempty" //ignores the field is nil when nullable or empty string
	Tags_SkipParsing ETag = "set"       //Takes the model as is it is without parsing it

)

var (
	supportedTags = Tags{string(Tags_Int), string(Tags_Float), string(Tags_String), string(Tags_Omitempty), string(Tags_SkipParsing)}
)
