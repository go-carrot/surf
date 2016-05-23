package serf_test

import (
	"database/sql"
	"github.com/BrandonRomano/serf"
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
	serf.Worker
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewPlace(dbConnection *sql.DB) *Place {
	place := new(Place)
	return place.Prep(dbConnection)
}

func (p *Place) Prep(dbConnection *sql.DB) *Place {
	p.Worker = &serf.PqWorker{
		Database: dbConnection,
		Config: serf.Configuration{
			TableName: "place",
			Fields: []serf.Field{
				serf.Field{
					Pointer:          &p.Id,
					Name:             "id",
					UniqueIdentifier: true,
					IsSet: func(pointer interface{}) bool {
						pointerInt := *pointer.(*int)
						return pointerInt != 0
					},
				},
				serf.Field{
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
	serf.Worker
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewPerson(dbConnection *sql.DB) *Person {
	person := new(Person)
	return person.Prep(dbConnection)
}

func (p *Person) Prep(dbConnection *sql.DB) *Person {
	p.Worker = &serf.PqWorker{
		Database: dbConnection,
		Config: serf.Configuration{
			TableName: "people",
			Fields: []serf.Field{
				serf.Field{
					Pointer:          &p.Id,
					Name:             "id",
					UniqueIdentifier: true,
					// Intentionally leaving out IsSet function
					// this isn't proper usage, but we still need
					// to test the code-path for this to make sure
					// we panic!
				},
				serf.Field{
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
	serf.Worker
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
	a.Worker = &serf.PqWorker{
		Database: dbConnection,
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

// ================================
// ========== Test Suite ==========
// ================================

type PqWorkerTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (suite *PqWorkerTestSuite) SetupTest() {
	databaseUrl := os.Getenv("SERF_TEST_URL")

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

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPqWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(PqWorkerTestSuite))
}
