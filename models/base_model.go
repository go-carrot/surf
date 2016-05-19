package models

import (
	"bytes"
	"database/sql"
	db "github.com/BrandonRomano/drudge/database"
	"strconv"
)

type Configuration struct {
	TableName string
	Fields    []Field
}

type Field struct {
	Pointer          interface{}
	Name             string
	Insertable       bool
	UniqueIdentifier bool
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
		// TODO add test to see if we should run this, don't blindly run every query
		// else slugs will always take 2 querys from the database
		// need a way to tell if these are null or not before doing this...  Needs
		// a unified nullable interface with a isNull() method
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
	// Impossible to reach, just satisfying the compiler
	return nil
}

func (dbw DbWorker) consumeRow(row *sql.Row) error {
	fields := dbw.BaseModel.GetConfiguration().Fields
	s := make([]interface{}, 3)
	for i, value := range fields {
		s[i] = value.Pointer
	}
	return row.Scan(s...)
}

func (dbw DbWorker) consumeNextRow(rows *sql.Rows) error {
	fields := dbw.BaseModel.GetConfiguration().Fields
	s := make([]interface{}, 3)
	for i, value := range fields {
		s[i] = value.Pointer
	}
	return rows.Scan(s...)
}
