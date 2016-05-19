package models

import (
	"bytes"
	"database/sql"
	"errors"
	db "github.com/BrandonRomano/drudge/database"
	"strconv"
	"fmt"
)

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

type BaseModel interface {
	GetConfiguration() Configuration
}

type DbWorker struct {
	BaseModel BaseModel
}

func (dbw *DbWorker) Insert() error {
	configuration := dbw.BaseModel.GetConfiguration()

	// Get insertable fields
	var insertableFields []Field
	allFields := configuration.Fields
	for _, field := range allFields {
		if field.Insertable {
			insertableFields = append(insertableFields, field)
		}
	}

	// Building Query
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("INSERT INTO ")
	queryBuffer.WriteString(configuration.TableName)
	queryBuffer.WriteString("(")
	for i, field := range insertableFields {
		queryBuffer.WriteString(field.Name)
		if (i + 1) < len(insertableFields) {
			queryBuffer.WriteString(", ")
		}
	}
	queryBuffer.WriteString(") VALUES(")
	for i, _ := range insertableFields {
		queryBuffer.WriteString("$")
		queryBuffer.WriteString(strconv.Itoa(i + 1))
		if (i + 1) < len(insertableFields) {
			queryBuffer.WriteString(", ")
		}
	}
	queryBuffer.WriteString(") RETURNING *;")

	// Get Value Fields
	var valueFields []interface{}
	for _, value := range insertableFields {
		valueFields = append(valueFields, value.Pointer)
	}

	// Firing off Query
	database := db.Get()
	row := database.QueryRow(
		queryBuffer.String(),
		valueFields...,
	)

	return dbw.consumeRow(row)
}

func (dbw *DbWorker) Load() error {
	// Load Configuration
	configuration := dbw.BaseModel.GetConfiguration()

	// Get Database
	database := db.Get()

	// Get unique identifier fields
	var uniqueIdentifierFields []Field
	allFields := configuration.Fields
	for _, field := range allFields {
		if field.UniqueIdentifier {
			uniqueIdentifierFields = append(uniqueIdentifierFields, field)
		}
	}

	for i, field := range uniqueIdentifierFields {
		// If we've explicitly determined that the field isn't set, continue
		if field.IsSet != nil && !field.IsSet(field.Pointer) {
			continue
		}

		var queryBuffer bytes.Buffer
		queryBuffer.WriteString("SELECT * FROM ")
		queryBuffer.WriteString(configuration.TableName)
		queryBuffer.WriteString(" WHERE ")
		queryBuffer.WriteString(field.Name)
		queryBuffer.WriteString("=$1;")

		row := database.QueryRow(
			queryBuffer.String(),
			field.Pointer,
		)
		err := dbw.consumeRow(row)
		if err != nil {
			return nil
		} else {
			if (i + 1) >= len(uniqueIdentifierFields) {
				return err
			} else {
				continue
			}
		}
	}
	// No fields set
	return errors.New("No UniqueIdentifier fields found that are set") // TODO
}

func (dbw DbWorker) Update() error {
	// Load Configuration
	configuration := dbw.BaseModel.GetConfiguration()

	// Get Database
	database := db.Get()

	// Get unique identifier fields
	var uniqueIdentifierFields []Field
	for _, field := range configuration.Fields {
		if field.UniqueIdentifier {
			uniqueIdentifierFields = append(uniqueIdentifierFields, field)
		}
	}

	// Determine which unique identifier we will be querying with
	var uniqueIdentifierField Field
	for _, field := range uniqueIdentifierFields {
		// If we haven't explicitly determined that a field isn't set,
		// assume that this will be the field we will be using
		if field.IsSet == nil || field.IsSet(field.Pointer) {
			uniqueIdentifierField = field
			break
		}
	}
	if uniqueIdentifierField.Pointer == nil {
		return errors.New("No unique identifier was set, so we can't update this model") //TODO
	}

	// Get updatable fields
	var updatableFields []Field
	for _, field := range configuration.Fields {
		if field.Updatable {
			updatableFields = append(updatableFields, field)
		}
	}

	// Generate Query
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("UPDATE ")
	queryBuffer.WriteString(configuration.TableName)
	queryBuffer.WriteString(" SET ")
	for i, field := range updatableFields {
		queryBuffer.WriteString(field.Name)
		queryBuffer.WriteString("=$")
		queryBuffer.WriteString(strconv.Itoa(i + 1))
		if (i + 1) < len(updatableFields) {
			queryBuffer.WriteString(", ")
		}
	}
	queryBuffer.WriteString(" WHERE ")
	queryBuffer.WriteString(uniqueIdentifierField.Name)
	queryBuffer.WriteString("=$")
	queryBuffer.WriteString(strconv.Itoa(len(updatableFields) + 1))
	queryBuffer.WriteString(" RETURNING *;")

	// Get Value Fields
	var valueFields []interface{}
	for _, value := range updatableFields {
		valueFields = append(valueFields, value.Pointer)
	}
	valueFields = append(valueFields, uniqueIdentifierField.Pointer)

	// Sending off query
	row := database.QueryRow(queryBuffer.String(), valueFields...)
	return dbw.consumeRow(row)
}

func (dbw DbWorker) Delete() error {
	// Load Configuration
	configuration := dbw.BaseModel.GetConfiguration()

	// Get Database
	database := db.Get()

	// Get unique identifier fields
	var uniqueIdentifierFields []Field
	for _, field := range configuration.Fields {
		if field.UniqueIdentifier {
			uniqueIdentifierFields = append(uniqueIdentifierFields, field)
		}
	}

	// Determine which unique identifier we will be querying with
	var uniqueIdentifierField Field
	for _, field := range uniqueIdentifierFields {
		// If we haven't explicitly determined that a field isn't set,
		// assume that this will be the field we will be using
		if field.IsSet == nil || field.IsSet(field.Pointer) {
			uniqueIdentifierField = field
			break
		}
	}
	if uniqueIdentifierField.Pointer == nil {
		return errors.New("No unique identifier was set, so we can't update this model") // TODO
	}

	// Generate Query
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("DELETE FROM ")
	queryBuffer.WriteString(configuration.TableName)
	queryBuffer.WriteString(" WHERE ")
	queryBuffer.WriteString(uniqueIdentifierField.Name)
	queryBuffer.WriteString("=$1;")

	fmt.Println(queryBuffer.String())

	// Executing Query
	res, err := database.Exec(queryBuffer.String(), uniqueIdentifierField.Pointer)
	if err != nil {
		return err
	}
	numRows, _ := res.RowsAffected()
	if numRows != 1 {
		return errors.New("Nothing was deleted")
	}

	return nil
}

func (dbw *DbWorker) consumeRow(row *sql.Row) error {
	fields := dbw.BaseModel.GetConfiguration().Fields
	s := make([]interface{}, 3)
	for i, value := range fields {
		s[i] = value.Pointer
	}
	return row.Scan(s...)
}

func (dbw *DbWorker) consumeNextRow(rows *sql.Rows) error {
	fields := dbw.BaseModel.GetConfiguration().Fields
	s := make([]interface{}, 3)
	for i, value := range fields {
		s[i] = value.Pointer
	}
	return rows.Scan(s...)
}
