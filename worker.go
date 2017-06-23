package surf

// Worker is the interface that defines the type
// to be embedded on models
type Worker interface {
	Insert() error
	Load() error
	Update() error
	Delete() error
	GetConfiguration() *Configuration
}

// Configuration is the definition of a model
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
