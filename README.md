<a href="https://engineering.carrot.is/"><p align="center"><img src="https://cloud.githubusercontent.com/assets/2105067/24525319/d3d26516-1567-11e7-9506-7611b3287d53.png" alt="Go Carrot" width="350px" align="center;" /></p></a>
# Surf

[![Build Status](https://travis-ci.org/go-carrot/surf.svg?branch=master)](https://travis-ci.org/go-carrot/surf) [![codecov](https://codecov.io/gh/go-carrot/surf/branch/master/graph/badge.svg)](https://codecov.io/gh/go-carrot/surf) [![Go Report Card](https://goreportcard.com/badge/github.com/go-carrot/surf)](https://goreportcard.com/report/github.com/go-carrot/surf) [![Gitter](https://img.shields.io/gitter/room/nwjs/nw.js.svg)](https://gitter.im/go-carrot/surf)

Surf is a high level datastore worker that provides CRUD operations for your models.

## In Use

Before I dive into explaining how to use this library, let me first show an example of how you will interface with your models after everthing is set up:

```go
// Inserts
myAnimal := models.NewAnimal()
myAnimal.Name = "Rigby"
myAnimal.Age = 3
myAnimal.Insert()

// Loads
myAnimalCopy := models.NewAnimal()
myAnimalCopy.Id = myAnimal.Id
myAnimalCopy.Load() // After this, myAnimalCopy's Name will be "Rigby" and Age will be 3 (pulled from the database)

// Updates
myAnimalCopy.Age = 4
myAnimalCopy.Update() // Updates Age in the database

// Deletes
myAnimalCopy.Delete() // Deletes from the database
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

### Embed a surf.Worker

After this is set up, we can now [embed](https://golang.org/doc/effective_go.html#embedding) a `surf.Worker` into our model.

```go
type Animal struct {
    surf.Worker
    Id   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

A `surf.Worker` is actually just an interface, so we'll need to decide what type of worker we want to use!

> There's more information [below](#workers) on `surf.Worker`.

For this example, we are going to be using a `surf.PqWorker`.

You can make the decision of what `surf.Worker` to pack into your model at run time, but it will probably be easiest to create a constructor type function and create your models through that.  Here I provide a constructor type fuction, but also a `Prep` function which will provide you the flexibility to not use the constructor.

```go
func NewAnimal() *Animal {
    animal := new(Animal)
    return animal.Prep()
}

func (a *Animal) Prep() *Animal {
    a.Worker = &surf.PqWorker{
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
	a.Worker = &surf.PqWorker{
		Database: db.Get(),
		Config: surf.Configuration{
			TableName: "animals",
			Fields: []surf.Field{
				surf.Field{
					Pointer:          &a.Id,
					Name:             "id",
					UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerInt := *pointer.(*int)
						return pointerInt != 0
					},
				},
				surf.Field{
					Pointer:    &a.Name,
					Name:       "name",
					Insertable: true,
					Updatable:  true,
				},
				surf.Field{
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

A `surf.Configuration` has two fields, `TableName` and `Fields`.

`TableName` is simply a string that represents your table/collection name in your datastore.

`Fields` is an array of `surf.Field` (which is explained in detail [below](#surffield)).

## surf.Field

A `surf.Field` defines how a `surf.Worker` will interact with a field.

a `surf.Field` contains a few values that determine this interaction:

#### Pointer

This is a pointer to the `surf.BaseModel`'s field.

#### Name

This is the name of the field as specified in the datastore.

#### Insertable

This value specifies if this `surf.Field` is to be considered by the `Insert()` method of our worker.

#### Updatable

This value specifies if this `surf.Field` is to be considered by the `Update()` method of our worker.

#### UniqueIdentifier

> **Note**: `IsSet` is also required if `UniqueIdentifier` is set to true.

This value specifies that this field can unique identify an entry in the datastore.

You do not _need_ to set this to true for all of your `UNIQUE` fields in your datastore, but you can.

Setting `UniqueIdentifier` to true gives you the following:

- The ability to set that fields value in the `surf.BaseModel` and call `Load()` against it.
- Call `Update()` with this field in the where clause / filter
- Call `Delete()` with this field in the where clause / filter.

> If you are using a `surf.Worker` that is backed by a relational database, it is strongly recommended that column is indexed.

#### IsSet

This is a function that determines if the value in the struct is set or not.

This is only required if `UniqueIdentifier` is set to `true`.

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

Workers are simply implementations that adhere to the following interface:

```go
type Worker interface {
    Insert() error
    Load() error
    Update() error
    Delete() error
}
```

> Right now in this library there is only `surf.PqWorker` written, but I plan to at minimum write a MySQL worker in the near future.

### surf.PqWorker

`surf.PqWorker` is written on top of [github.com/lib/pq](https://github.com/lib/pq).  This provides high level PostgreSQL CRUD operations to your models.

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

```sh
SERF_TEST_DATABASE_URL=""
```

After this is all set up you can run `go test` the following to run the tests.  To check the coverage run `go test -cover`

## Acknowledgements

Thanks to [@roideuniverse](https://github.com/roideuniverse) for some early guidance that ultimately lead into the creation of this library.

## License

[MIT](LICENSE.md)
