package drudge

import (
	"bytes"
	"database/sql"
	"errors"
	"strconv"
)

type PqWorker struct {
	Database *sql.DB       `json:"-"`
	Config   Configuration `json:"-"`
}

func (w *PqWorker) Insert() error {
	// Get insertable fields
	var insertableFields []Field
	allFields := w.Config.Fields
	for _, field := range allFields {
		if field.Insertable {
			insertableFields = append(insertableFields, field)
		}
	}

	// Building Query
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
	row := w.Database.QueryRow(
		queryBuffer.String(),
		valueFields...,
	)

	return w.consumeRow(row)
}

func (w *PqWorker) Load() error {
	// Get unique identifier fields
	var uniqueIdentifierFields []Field
	allFields := w.Config.Fields
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
		queryBuffer.WriteString(w.Config.TableName)
		queryBuffer.WriteString(" WHERE ")
		queryBuffer.WriteString(field.Name)
		queryBuffer.WriteString("=$1;")

		row := w.Database.QueryRow(
			queryBuffer.String(),
			field.Pointer,
		)
		err := w.consumeRow(row)
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

func (w *PqWorker) Update() error {
	// Get unique identifier fields
	var uniqueIdentifierFields []Field
	for _, field := range w.Config.Fields {
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

	// Sending off query
	row := w.Database.QueryRow(queryBuffer.String(), valueFields...)
	return w.consumeRow(row)
}

func (w PqWorker) Delete() error {
	// Get unique identifier fields
	var uniqueIdentifierFields []Field
	for _, field := range w.Config.Fields {
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
	queryBuffer.WriteString(w.Config.TableName)
	queryBuffer.WriteString(" WHERE ")
	queryBuffer.WriteString(uniqueIdentifierField.Name)
	queryBuffer.WriteString("=$1;")

	// Executing Query
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

func (w *PqWorker) consumeRow(row *sql.Row) error {
	fields := w.Config.Fields
	s := make([]interface{}, 3)
	for i, value := range fields {
		s[i] = value.Pointer
	}
	return row.Scan(s...)
}

func (w *PqWorker) consumeNextRow(rows *sql.Rows) error {
	fields := w.Config.Fields
	s := make([]interface{}, 3)
	for i, value := range fields {
		s[i] = value.Pointer
	}
	return rows.Scan(s...)
}
