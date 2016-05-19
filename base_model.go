package drudge

type BaseModel interface {
	GetConfiguration() Configuration
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
