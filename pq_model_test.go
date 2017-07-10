package surf_test

import (
	"database/sql"
	"github.com/go-carrot/surf"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

// =================================
// ========== Place Model ==========
// =================================

// There is no need to create a model for this in the database,
// as we are only using this model to intentionally create any error
// in the database.
type Place struct {
	surf.Model
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewPlace(dbConnection *sql.DB) *Place {
	place := new(Place)
	return place.Prep(dbConnection)
}

func (p *Place) Prep(dbConnection *sql.DB) *Place {
	p.Model = &surf.PqModel{
		Database: dbConnection,
		Config: surf.Configuration{
			TableName: "place",
			Fields: []surf.Field{
				{
					Pointer:          &p.Id,
					Name:             "id",
					UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerInt := *pointer.(*int)
						return pointerInt != 0
					},
				},
				{
					Pointer:    &p.Name,
					Name:       "name",
					Insertable: true,
					Updatable:  true,
				},
			},
		},
	}
	return p
}

// ==================================
// ========== Person Model ==========
// ==================================

// There is no need to create a model for this in the database
// as we are using this to test a failure before we even hit the db
type Person struct {
	surf.Model
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewPerson(dbConnection *sql.DB) *Person {
	person := new(Person)
	return person.Prep(dbConnection)
}

func (p *Person) Prep(dbConnection *sql.DB) *Person {
	p.Model = &surf.PqModel{
		Database: dbConnection,
		Config: surf.Configuration{
			TableName: "people",
			Fields: []surf.Field{
				{
					Pointer:          &p.Id,
					Name:             "id",
					UniqueIdentifier: true,
					// Intentionally leaving out IsSet function
					// this isn't proper usage, but we still need
					// to test the code-path for this to make sure
					// we panic!
				},
				{
					Pointer:    &p.Name,
					Name:       "name",
					Insertable: true,
					Updatable:  true,
				},
			},
		},
	}
	return p
}

// ==================================
// ========== Animal Model ==========
// ==================================

/*
Represents:

CREATE TABLE animals(
    id    serial          PRIMARY KEY,
    slug  TEXT            UNIQUE NOT NULL,
    name  TEXT            NOT NULL,
    age   int             NOT NULL
);
*/
type Animal struct {
	surf.Model
	Id   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func NewAnimal(dbConnection *sql.DB) *Animal {
	animal := new(Animal)
	return animal.Prep(dbConnection)
}

func (a *Animal) Prep(dbConnection *sql.DB) *Animal {
	a.Model = &surf.PqModel{
		Database: dbConnection,
		Config: surf.Configuration{
			TableName: "animals",
			Fields: []surf.Field{
				{
					Pointer:          &a.Id,
					Name:             "id",
					UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerInt := *pointer.(*int)
						return pointerInt != 0
					},
				},
				{
					Pointer:          &a.Slug,
					Name:             "slug",
					UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerStr := *pointer.(*string)
						return pointerStr != ""
					},
					Insertable: true,
					Updatable:  true,
				},
				{
					Pointer:    &a.Name,
					Name:       "name",
					Insertable: true,
					Updatable:  true,
				},
				{
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

// ===============================
// ========== Toy Model ==========
// ===============================

/**
Represents:

CREATE TABLE toys(
  id serial PRIMARY KEY,
  name text NOT NULL,
  owner bigint NOT NULL REFERENCES animals(id) ON DELETE CASCADE
);
*/
type Toy struct {
	surf.Model
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	OwnerId int     `json:"-"`
	Owner   *Animal `json:"animal"`
}

func NewToy(dbConnection *sql.DB) *Toy {
	toy := new(Toy)
	return toy.Prep(dbConnection)
}

func (t *Toy) Prep(dbConnection *sql.DB) *Toy {
	t.Model = &surf.PqModel{
		Database: dbConnection,
		Config: surf.Configuration{
			TableName: "toys",
			Fields: []surf.Field{
				{Pointer: &t.Id, Name: "id", UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerInt := *pointer.(*int)
						return pointerInt != 0
					}},
				{Pointer: &t.Name, Name: "name", Insertable: true, Updatable: true},
				{Pointer: &t.OwnerId, Name: "owner", Insertable: true, Updatable: true,
					GetReference: func() (surf.BuildModel, string) {
						return func() surf.Model {
							return NewAnimal(dbConnection)
						}, "id"
					},
					SetReference: func(model surf.Model) error {
						t.Owner = model.(*Animal)
						return nil
					},
					IsSet: func(pointer interface{}) bool {
						pointerInt := *pointer.(*int)
						return pointerInt != 0
					},
				},
			},
		},
	}
	return t
}

// ==================================================
// ========== Animal Consume Failure Model ==========
// ==================================================

// This model exists purely so we can test that ConsumeRows
// properly throws an error
//
// Intentionally using the wrong type for Name
type AnimalConsume struct {
	surf.Model
	Id   int    `json:"id"`
	Slug string `json:"slug"`
	Name int    `json:"name"`
	Age  int    `json:"age"`
}

func NewAnimalConsume(dbConnection *sql.DB) *AnimalConsume {
	animalConsume := new(AnimalConsume)
	return animalConsume.Prep(dbConnection)
}

func (ac *AnimalConsume) Prep(dbConnection *sql.DB) *AnimalConsume {
	ac.Model = &surf.PqModel{
		Database: dbConnection,
		Config: surf.Configuration{
			TableName: "animals",
			Fields: []surf.Field{
				{
					Pointer:          &ac.Id,
					Name:             "id",
					UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerInt := *pointer.(*int)
						return pointerInt != 0
					},
				},
				{
					Pointer:          &ac.Slug,
					Name:             "slug",
					UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerStr := *pointer.(*string)
						return pointerStr != ""
					},
					Insertable: true,
					Updatable:  true,
				},
				{
					Pointer:    &ac.Name,
					Name:       "name",
					Insertable: true,
					Updatable:  true,
				},
				{
					Pointer:    &ac.Age,
					Name:       "age",
					Insertable: true,
					Updatable:  true,
				},
			},
		},
	}
	return ac
}

// ================================
// ========== Test Suite ==========
// ================================

type PqWorkerTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (suite *PqWorkerTestSuite) SetupTest() {
	databaseUrl := os.Getenv("SERF_TEST_DATABASE_URL")

	// Opening + storing the connection
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		suite.Fail("Failed to open database connection")
	}

	// Pinging the database
	err = db.Ping()
	if err != nil {
		suite.Fail("Failed to communicate with database")
	}

	suite.db = db
}

func (suite *PqWorkerTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *PqWorkerTestSuite) TestInsert() {
	// Create an Animal
	rigby := NewAnimal(suite.db)
	rigby.Name = "Rigby"
	rigby.Slug = "rigby"
	rigby.Age = 3
	rigby.Insert()
	assert.NotEqual(suite.T(), 0, rigby.Id)

	// Cause a conflict to test errors are being thrown
	rigbyTwo := NewAnimal(suite.db)
	rigbyTwo.Name = "Rigby Two"
	rigbyTwo.Slug = "rigby" // This should cause a conflict
	rigbyTwo.Age = 3
	err := rigbyTwo.Insert()
	assert.NotEqual(suite.T(), nil, err) // We expect an error here

	// Clean up Rigby
	rigby.Delete()
}

func (suite *PqWorkerTestSuite) TestLoad() {
	// Create an Animal
	rigby := NewAnimal(suite.db)
	rigby.Name = "Rigby"
	rigby.Slug = "rigby"
	rigby.Age = 3
	rigby.Insert()
	assert.NotEqual(suite.T(), 0, rigby.Id)

	// Verify it loads from id
	rigbyIdLoad := NewAnimal(suite.db)
	rigbyIdLoad.Id = rigby.Id
	rigbyIdLoad.Load()
	assert.Equal(suite.T(), rigby.Slug, rigbyIdLoad.Slug)

	// Verify it loads from slug
	rigbySlugLoad := NewAnimal(suite.db)
	rigbySlugLoad.Slug = rigby.Slug
	rigbySlugLoad.Load()
	assert.Equal(suite.T(), rigby.Id, rigbyIdLoad.Id)

	// Verify the ID is used before the slug
	idBeforeSlugLoad := NewAnimal(suite.db)
	idBeforeSlugLoad.Id = rigby.Id
	idBeforeSlugLoad.Slug = "asdf"
	err := idBeforeSlugLoad.Load()
	assert.Equal(suite.T(), rigby.Slug, idBeforeSlugLoad.Slug)

	// Make sure an error is thrown when nothing is set
	dummyAnimal := NewAnimal(suite.db)
	err = dummyAnimal.Load()
	assert.NotEqual(suite.T(), nil, err)

	// Make sure an error is thrown when trying to load something that doesn't exist
	anotherDummyAnimal := NewAnimal(suite.db)
	anotherDummyAnimal.Slug = "wow-cool-cat"
	err = anotherDummyAnimal.Load()
	assert.NotEqual(suite.T(), nil, err)

	// Clean up Rigby
	rigby.Delete()
}

func (suite *PqWorkerTestSuite) TestBulkLoad() {
	// Create some Animals
	luna := NewAnimal(suite.db)
	luna.Name = "Luna"
	luna.Slug = "luna"
	luna.Age = 2
	luna.Insert()
	assert.NotEqual(suite.T(), 0, luna.Id)

	rae := NewAnimal(suite.db)
	rae.Name = "Rae"
	rae.Slug = "rae"
	rae.Age = 2
	rae.Insert()
	assert.NotEqual(suite.T(), 0, rae.Id)

	// ====== Load ASC

	// Bulk load
	animals, err := NewAnimal(suite.db).BulkFetch(surf.BulkFetchConfig{
		Limit:  10,
		Offset: 0,
		OrderBys: []surf.OrderBy{
			{Field: "name", Type: surf.ORDER_BY_ASC},
		},
	}, func() surf.Model { return NewAnimal(suite.db) })
	assert.Nil(suite.T(), err)

	// Test that they were loaded
	lunaFound, raeFound := false, false
	for _, animal := range animals {
		switch animal.(*Animal).Name {
		case "Luna":
			lunaFound = true
		case "Rae":
			// Luna should come first in the list
			assert.True(suite.T(), lunaFound)
			raeFound = true
		}
	}
	assert.True(suite.T(), lunaFound && raeFound)

	// ====== Load DESC

	// Bulk load
	config := surf.BulkFetchConfig{
		Limit:  10,
		Offset: 0,
	}
	config.ConsumeSortQuery("-name,id")
	animals, err = NewAnimal(suite.db).BulkFetch(config, func() surf.Model {
		return NewAnimal(suite.db)
	})
	assert.Nil(suite.T(), err)

	// Test that they were loaded
	lunaFound, raeFound = false, false
	for _, animal := range animals {
		switch animal.(*Animal).Name {
		case "Rae":
			raeFound = true
		case "Luna":
			// Rae should come first in the list
			assert.True(suite.T(), raeFound)
			lunaFound = true
		}
	}
	assert.True(suite.T(), lunaFound && raeFound)

	// Clean Up
	luna.Delete()
	rae.Delete()
}

func (suite *PqWorkerTestSuite) TestBulkLoadFailures() {
	// Create an Animal
	luna := NewAnimal(suite.db)
	luna.Name = "Luna"
	luna.Slug = "luna"
	luna.Age = 2
	luna.Insert()
	assert.NotEqual(suite.T(), 0, luna.Id)

	// Load via invalid Order By
	config := surf.BulkFetchConfig{
		Limit:  10,
		Offset: 0,
	}
	config.ConsumeSortQuery("-helloworld")
	_, err := NewAnimal(suite.db).BulkFetch(config, func() surf.Model {
		return NewAnimal(suite.db)
	})
	assert.NotNil(suite.T(), err)

	// Test a failed DB query
	_, err = NewPlace(suite.db).BulkFetch(surf.BulkFetchConfig{
		Limit:  10,
		Offset: 0,
		OrderBys: []surf.OrderBy{
			{Field: "id", Type: surf.ORDER_BY_ASC},
		},
	}, func() surf.Model { return NewPlace(suite.db) })
	assert.NotNil(suite.T(), err)

	// Test that ConsumeRow will fail
	_, err = NewAnimalConsume(suite.db).BulkFetch(surf.BulkFetchConfig{
		Limit:  10,
		Offset: 0,
		OrderBys: []surf.OrderBy{
			{Field: "id", Type: surf.ORDER_BY_ASC},
		},
	}, func() surf.Model { return NewAnimalConsume(suite.db) })
	assert.NotNil(suite.T(), err)

	// Clean up
	luna.Delete()
}

func (suite *PqWorkerTestSuite) TestUpdate() {
	// Create an Animal
	rigby := NewAnimal(suite.db)
	rigby.Name = "Rigby"
	rigby.Slug = "rigby"
	rigby.Age = 3
	rigby.Insert()
	assert.NotEqual(suite.T(), 0, rigby.Id)

	// Update
	rigby.Age = 4
	err := rigby.Update()
	assert.Equal(suite.T(), nil, err)

	// Verify the update happened in the DB
	rigbyVerification := NewAnimal(suite.db)
	rigbyVerification.Id = rigby.Id
	rigbyVerification.Load()
	assert.Equal(suite.T(), 4, rigbyVerification.Age)

	// Create
	norbert := NewAnimal(suite.db)
	norbert.Name = "Norbert"
	norbert.Slug = "norbert"
	norbert.Age = 1
	norbert.Insert()

	// Update and cause conflict
	norbert.Slug = "rigby"
	err = norbert.Update()
	assert.NotEqual(suite.T(), nil, err)

	// Update without a unique, there should be an error
	dummyAnimal := NewAnimal(suite.db)
	err = dummyAnimal.Update()
	assert.NotEqual(suite.T(), nil, err)

	// Clean up
	rigby.Delete()
	norbert.Delete()
}

func (suite *PqWorkerTestSuite) TestDelete() {
	// Create an Animal
	rigby := NewAnimal(suite.db)
	rigby.Name = "Rigby"
	rigby.Slug = "rigby"
	rigby.Age = 3
	rigby.Insert()
	assert.NotEqual(suite.T(), 0, rigby.Id)

	// Delete
	err := rigby.Delete()
	assert.Equal(suite.T(), nil, err)

	// Try to delete an animal that doesn't exist
	fakeAnimal := NewAnimal(suite.db)
	fakeAnimal.Slug = "some-fake-animal"
	err = fakeAnimal.Delete()
	assert.NotEqual(suite.T(), nil, err)

	// Make sure we can't delete without a unique set
	emptyAnimal := NewAnimal(suite.db)
	err = emptyAnimal.Delete()
	assert.NotEqual(suite.T(), nil, err)
}

func (suite *PqWorkerTestSuite) TestUniqueMissingIsSet() {
	defer func() {
		recover()
	}()

	// Load a person
	brandon := NewPerson(suite.db)
	brandon.Name = "Brandon"
	brandon.Load()

	// We should not be able to reach this point, as brandon.Load()
	// should cause a panic.
	suite.Fail("brandon.Load() should have caused a panic!")
}

func (suite *PqWorkerTestSuite) TestDeleteSqlError() {
	nyc := NewPlace(suite.db)
	nyc.Id = 1
	nyc.Name = "New York City"

	err := nyc.Delete()
	assert.NotEqual(suite.T(), nil, err)
}

func (suite *PqWorkerTestSuite) TestGetConfiguration() {
	config := NewAnimal(suite.db).GetConfiguration()
	var hasId, hasSlug, hasName, hasAge bool
	for _, field := range config.Fields {
		switch field.Name {
		case "id":
			hasId = true
		case "slug":
			hasSlug = true
		case "name":
			hasName = true
		case "age":
			hasAge = true
		}
	}
	assert.True(suite.T(), (hasId && hasSlug && hasName && hasAge))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPqWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(PqWorkerTestSuite))
}
