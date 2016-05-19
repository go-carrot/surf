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
			},
			Field{
				Pointer:    &a.Name,
				Name:       "name",
				Insertable: true,
			},
			Field{
				Pointer:    &a.Age,
				Name:       "age",
				Insertable: true,
			},
		},
	}
}
