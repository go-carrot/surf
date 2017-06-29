package surf

// Model is the interface that defines the type
// that will be embedded on models
type Model interface {
	Insert() error
	Load() error
	Update() error
	Delete() error
	BulkFetch(BulkFetchConfig, BuildModel) ([]Model, error)
	GetConfiguration() *Configuration
}

// Configuration is the metadata to be attached to a model
type Configuration struct {
	TableName string
	Fields    []Field
}

// Field is the definition of a single value in a model
type Field struct {
	Pointer          interface{}
	Name             string
	Insertable       bool
	Updatable        bool
	UniqueIdentifier bool
	IsSet            func(interface{}) bool
}

// BuildModel is a function that is responsible for returning a
// Model that is ready to have GetConfiguration() called
type BuildModel func() Model
