# Drudge

Drudge is a high level datastore worker that provides CRUD operations for your models.

## Getting Started

### Setting up a basic Model

Before we start anything, we'll need to create a model.  For example, we will use an Animal.

```go
type Animal struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

### Embed a drudge.Worker

After this is set up, we can now [embed](https://golang.org/doc/effective_go.html#embedding) a `drudge.Worker` into our model.

```go
type Animal struct {
    drudge.Worker
    Id   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

A `drudge.Worker` is actually just an interface, so we'll need to decide what type of worker we want to use!

> There's more information [below](#workers) on `drudge.Worker`.

For this example, we are going to be using a `drudge.PqWorker`.

You can make the decision of what `drudge.Worker` to pack into your model at run time, but it will probably be easiest to create a constructor type function and create your models through that.  Here I provide a constructor type fuction, but also a `Prep` function which will provide you the flexibility to not use the constructor.

```go
func NewAnimal() *Animal {
    animal := new(Animal)
    return animal.Prep()
}

func (a *Animal) Prep() *Animal {
    a.Worker = &drudge.PqWorker{
		Database: db.Get(), // This is a *sql.DB, with github.com/lib/pq as a driver
		Config: // TODO
	}
    return a
}
```

### Setting up Config

You'll notice in the last code snippet in the previous section, there is a `// TODO` mark.

In that `Prep` method, we still need to define the Config.

Before going into detail, here is the `Prep` method with a fully filled out Config.

```go
func (a *Animal) Prep() *Animal {
	a.Worker = &drudge.PqWorker{
		Database: db.Get(),
		Config: drudge.Configuration{
			TableName: "animals",
			Fields: []drudge.Field{
				drudge.Field{
					Pointer:          &a.Id,
					Name:             "id",
					UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerInt := *pointer.(*int)
						return pointerInt != 0
					},
				},
				drudge.Field{
					Pointer:    &a.Name,
					Name:       "name",
					Insertable: true,
					Updatable:  true,
				},
				drudge.Field{
					Pointer:    &a.Age,
					Name:       "age",
					Insertable: true,
					Updatable:  true,
				},
			},
		},
	}
	return a
}
```

A `druge.Configuration` has two fields, `TableName` and `Fields`.

`TableName` is simply a string that represents your table/collection name in your datastore.

`Fields` is an array of `drudge.Field` (which is explained in detail [below](#drudgefield)).

## drudge.Field

A `drudge.Field` defines how a `drudge.Worker` will interact with a field.

a `drudge.Field` contains a few values that determine this interaction:

#### Pointer

This is a pointer to the `drudge.BaseModel`'s field.

#### Name

This is the name of the field as specified in the datastore.

#### Insertable

This value specifies if this `drudge.Field` is to be considered by the `Insert()` method of our worker.

#### Updatable

This value specifies if this `drudge.Field` is to be considered by the `Update()` method of our worker.

#### UniqueIdentifier

> **Note**: `IsSet` is also required if `UniqueIdentifier` is set to true.

This value specifies that this field can unique identify an entry in the datastore.

You do not _need_ to set this to true for all of your `UNIQUE` fields in your datastore, but you can.

Setting `UniqueIdentifier` to true gives you the following:
    - The ability to set that fields value in the `drudge.BaseModel` and call `Load()` against it.
    - Call `Update()` with this field in the where clause / filter
    - Call `Delete()` with this field in the where clause / filter.

> If you are using a `drudge.Worker` that is backed by a relational database, it is strongly recommended that column is indexed.

#### IsSet

This is a function that determines if the value in the struct is set or not.

This is only required if `UniqueIdentifier` is set.

This function will likely look something like this:

```go
// ...
IsSet: func(pointer interface{}) bool {
    pointerInt := *pointer.(*int)
    return pointerInt != 0
},
// ...
```

## Workers

Workers are simply implementations that adhere to this interface:

```go
type Worker interface {
    Insert() error
    Load() error
    Update() error
    Delete() error
}
```

> Right now in this library there is only `drudge.PqWorker` written, but I plan to at minimum write a MySQL worker in the near future.

### drudge.PqWorker

`drudge.PqWorker` is written on top of [github.com/lib/pq](https://github.com/lib/pq).  This provides high level PostgreSQL CRUD operations to your models.

## Acknowledgements

Thanks to [@roideuniverse](https://github.com/roideuniverse) for some early guidance that ultimately lead into the creation of this library.

## License

[MIT](LICENSE.md)
