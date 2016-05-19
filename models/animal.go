package models

type Animal struct {
	DbWorker DbWorker `json:"-"`
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Age      int      `json:"age"`
}

func (a *Animal) GetConfiguration() Configuration {
	return Configuration{
		TableName: "animals",
		Fields: []Field{
			Field{
				Pointer:          &a.Id,
				Name:             "id",
				UniqueIdentifier: true,
				IsSet: func(pointer interface{}) bool {
					pointerInt := *pointer.(*int)
					return pointerInt != 0
				},
			},
			Field{
				Pointer:    &a.Name,
				Name:       "name",
				Insertable: true,
				Updatable: true,
			},
			Field{
				Pointer:    &a.Age,
				Name:       "age",
				Insertable: true,
				Updatable: true,
			},
		},
	}
}
