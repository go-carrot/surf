# Serf

[![Build Status](https://travis-ci.org/BrandonRomano/serf.svg?branch=br.travis)](https://travis-ci.org/BrandonRomano/serf)

Serf is a high level datastore worker that provides CRUD operations for your models.

## In Use

Before I dive into explaining how to use this library, let me first show an example of how you will interface with your models after everthing is set up:

```go
// Inserts
rigby := models.NewAnimal()
rigby.Name = "Rigby"
rigby.Age = 3
rigby.Insert()

// Loads
rigbyCopy := models.NewAnimal()
rigbyCopy.Id = rigby.Id
rigbyTwo.Load() // After this, rigbyCopy's Name will be "Rigby" and Age will be 3 (pulled from the database)

// Updates
rigbyCopy.Age = 4
rigbyCopy.Update() // Updates Age in the database

// Deletes
rigbyCopy.Delete() // Deletes from the database
```

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

### Embed a serf.Worker

After this is set up, we can now [embed](https://golang.org/doc/effective_go.html#embedding) a `serf.Worker` into our model.

```go
type Animal struct {
    serf.Worker
    Id   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

A `serf.Worker` is actually just an interface, so we'll need to decide what type of worker we want to use!

> There's more information [below](#workers) on `serf.Worker`.

For this example, we are going to be using a `serf.PqWorker`.

You can make the decision of what `serf.Worker` to pack into your model at run time, but it will probably be easiest to create a constructor type function and create your models through that.  Here I provide a constructor type fuction, but also a `Prep` function which will provide you the flexibility to not use the constructor.

```go
func NewAnimal() *Animal {
    animal := new(Animal)
    return animal.Prep()
}

func (a *Animal) Prep() *Animal {
    a.Worker = &serf.PqWorker{
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
	a.Worker = &serf.PqWorker{
		Database: db.Get(),
		Config: serf.Configuration{
			TableName: "animals",
			Fields: []serf.Field{
				serf.Field{
					Pointer:          &a.Id,
					Name:             "id",
					UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerInt := *pointer.(*int)
						return pointerInt != 0
					},
				},
				serf.Field{
					Pointer:    &a.Name,
					Name:       "name",
					Insertable: true,
					Updatable:  true,
				},
				serf.Field{
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

A `serf.Configuration` has two fields, `TableName` and `Fields`.

`TableName` is simply a string that represents your table/collection name in your datastore.

`Fields` is an array of `serf.Field` (which is explained in detail [below](#serffield)).

## serf.Field

A `serf.Field` defines how a `serf.Worker` will interact with a field.

a `serf.Field` contains a few values that determine this interaction:

#### Pointer

This is a pointer to the `serf.BaseModel`'s field.

#### Name

This is the name of the field as specified in the datastore.

#### Insertable

This value specifies if this `serf.Field` is to be considered by the `Insert()` method of our worker.

#### Updatable

This value specifies if this `serf.Field` is to be considered by the `Update()` method of our worker.

#### UniqueIdentifier

> **Note**: `IsSet` is also required if `UniqueIdentifier` is set to true.

This value specifies that this field can unique identify an entry in the datastore.

You do not _need_ to set this to true for all of your `UNIQUE` fields in your datastore, but you can.

Setting `UniqueIdentifier` to true gives you the following:
    - The ability to set that fields value in the `serf.BaseModel` and call `Load()` against it.
    - Call `Update()` with this field in the where clause / filter
    - Call `Delete()` with this field in the where clause / filter.

> If you are using a `serf.Worker` that is backed by a relational database, it is strongly recommended that column is indexed.

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

> Right now in this library there is only `serf.PqWorker` written, but I plan to at minimum write a MySQL worker in the near future.

### serf.PqWorker

`serf.PqWorker` is written on top of [github.com/lib/pq](https://github.com/lib/pq).  This provides high level PostgreSQL CRUD operations to your models.

## Running Tests

Before running tests, you must set up a database with a single table.

```sql
CREATE TABLE animals(
    id    serial          PRIMARY KEY,
    slug  TEXT            UNIQUE NOT NULL,
    name  TEXT            NOT NULL,
    age   int             NOT NULL
);
```

You'll then need to have an environment variable set pointing to the database URL:

```
SERF_TEST_DATABASE_URL=""
```

After this is all set up you can run `go test` the following to run the tests.  To check the coverage run `go test -cover`

## Acknowledgements

Thanks to [@roideuniverse](https://github.com/roideuniverse) for some early guidance that ultimately lead into the creation of this library.

## License

[MIT](LICENSE.md)
