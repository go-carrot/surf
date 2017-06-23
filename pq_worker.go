package surf

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// PqWorker is a github.com/lib/pq implementation of a Worker
type PqWorker struct {
	Database *sql.DB       `json:"-"`
	Config   Configuration `json:"-"`
}

// GetConfiguration returns the configuration for the model
func (w *PqWorker) GetConfiguration() *Configuration {
	return &w.Config
}

// Insert inserts the model into the database
func (w *PqWorker) Insert() error {
	// Get Insertable Fields
	var insertableFields []Field
	allFields := w.Config.Fields
	for _, field := range allFields {
		if field.Insertable {
			insertableFields = append(insertableFields, field)
		}
	}

	// Generate Query
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("INSERT INTO ")
	queryBuffer.WriteString(w.Config.TableName)
	queryBuffer.WriteString("(")
	for i, field := range insertableFields {
		queryBuffer.WriteString(field.Name)
		if (i + 1) < len(insertableFields) {
			queryBuffer.WriteString(", ")
		}
	}
	queryBuffer.WriteString(") VALUES(")
	for i := range insertableFields {
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

	// Execute Query
	row := w.Database.QueryRow(queryBuffer.String(), valueFields...)
	return w.consumeRow(row)
}

// Load loads the model from the database from its unique identifier
// and then loads those values into the struct
func (w *PqWorker) Load() error {
	// Get Unique Identifier
	uniqueIdentifierField, err := w.getUniqueIdentifier()
	if err != nil {
		return err
	}

	// Generate Query
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("SELECT * FROM ")
	queryBuffer.WriteString(w.Config.TableName)
	queryBuffer.WriteString(" WHERE ")
	queryBuffer.WriteString(uniqueIdentifierField.Name)
	queryBuffer.WriteString("=$1;")

	// Execute Query
	row := w.Database.QueryRow(queryBuffer.String(), uniqueIdentifierField.Pointer)
	return w.consumeRow(row)
}

// Update updates the model with the current values in the struct
func (w *PqWorker) Update() error {
	// Get Unique Identifier
	uniqueIdentifierField, err := w.getUniqueIdentifier()
	if err != nil {
		return err
	}

	// Get updatable fields
	var updatableFields []Field
	for _, field := range w.Config.Fields {
		if field.Updatable {
			updatableFields = append(updatableFields, field)
		}
	}

	// Generate Query
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("UPDATE ")
	queryBuffer.WriteString(w.Config.TableName)
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

	// Execute Query
	row := w.Database.QueryRow(queryBuffer.String(), valueFields...)
	return w.consumeRow(row)
}

// Delete deletes the model
func (w PqWorker) Delete() error {
	// Get Unique Identifier
	uniqueIdentifierField, err := w.getUniqueIdentifier()
	if err != nil {
		return err
	}

	// Generate Query
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("DELETE FROM ")
	queryBuffer.WriteString(w.Config.TableName)
	queryBuffer.WriteString(" WHERE ")
	queryBuffer.WriteString(uniqueIdentifierField.Name)
	queryBuffer.WriteString("=$1;")

	// Execute Query
	res, err := w.Database.Exec(queryBuffer.String(), uniqueIdentifierField.Pointer)
	if err != nil {
		return err
	}
	numRows, _ := res.RowsAffected()
	if numRows != 1 {
		return errors.New("Nothing was deleted")
	}
	return nil
}

// getUniqueIdentifier Returns the unique identifier that this worker will
// query against.
//
// This will return the first field in Configuration.Fields that:
//   - Has `UniqueIdentifier` set to true
//   - Returns true from `IsSet`
//
// This function will panic in the event that it encounters a field that is a
// `UniqueIdentifier`, and doesn't have `IsSet` implemented.
func (w *PqWorker) getUniqueIdentifier() (Field, error) {
	// Get all unique identifier fields
	var uniqueIdentifierFields []Field
	for _, field := range w.Config.Fields {
		if field.UniqueIdentifier {
			uniqueIdentifierFields = append(uniqueIdentifierFields, field)
		}
	}

	// Determine which unique identifier we will be querying with
	var uniqueIdentifierField Field
	for _, field := range uniqueIdentifierFields {
		if field.IsSet == nil {
			panic(fmt.Sprintf("Field `%v` must implement IsSet, as it is a `UniqueIdentifier`", field.Name))
		} else if field.IsSet(field.Pointer) {
			uniqueIdentifierField = field
			break
		}
	}

	// Return
	if uniqueIdentifierField.Pointer == nil {
		return uniqueIdentifierField, errors.New("There is no UniqueIdentifier Field that is set")
	}
	return uniqueIdentifierField, nil
}

// consumeRow Scans a *sql.Row into our struct
// that is using this worker
func (w *PqWorker) consumeRow(row *sql.Row) error {
	fields := w.Config.Fields
	var s []interface{}
	for _, value := range fields {
		s = append(s, value.Pointer)
	}
	return row.Scan(s...)
}
