package surf

type Worker interface {
	Insert() error
	Load() error
	Update() error
	Delete() error
}

type Configuration struct {
	TableName string
	Fields    []Field
}

type Field struct {
	Pointer          interface{}
	Name             string
	Insertable       bool
	Updatable        bool
	UniqueIdentifier bool
	IsSet            func(interface{}) bool
}
