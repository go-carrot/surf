package surf

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

// PqModel is a github.com/lib/pq implementation of a Model
type PqModel struct {
	Database *sql.DB       `json:"-"`
	Config   Configuration `json:"-"`
}

// GetConfiguration returns the configuration for the model
func (w *PqModel) GetConfiguration() *Configuration {
	return &w.Config
}

// Insert inserts the model into the database
func (w *PqModel) Insert() error {
	// Get Insertable Fields
	var insertableFields []Field
	for _, field := range w.Config.Fields {
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
	queryBuffer.WriteString(") RETURNING ")
	for i, field := range w.Config.Fields {
		queryBuffer.WriteString(field.Name)
		if (i + 1) < len(w.Config.Fields) {
			queryBuffer.WriteString(", ")
		}
	}
	queryBuffer.WriteString(";")

	// Get Value Fields
	var valueFields []interface{}
	for _, value := range insertableFields {
		valueFields = append(valueFields, value.Pointer)
	}

	// Log Query
	query := queryBuffer.String()
	PrintSqlQuery(query, valueFields...)

	// Execute Query
	row := w.Database.QueryRow(query, valueFields...)
	err := consumeRow(w, row)
	if err != nil {
		return err
	}

	// Expand foreign references
	return expandForeign(w)
}

// Load loads the model from the database from its unique identifier
// and then loads those values into the struct
func (w *PqModel) Load() error {
	// Get Unique Identifier
	uniqueIdentifierField, err := getUniqueIdentifier(w)
	if err != nil {
		return err
	}

	// Generate Query
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("SELECT ")
	for i, field := range w.Config.Fields {
		queryBuffer.WriteString(field.Name)
		if (i + 1) < len(w.Config.Fields) {
			queryBuffer.WriteString(", ")
		}
	}
	queryBuffer.WriteString(" FROM ")
	queryBuffer.WriteString(w.Config.TableName)
	queryBuffer.WriteString(" WHERE ")
	queryBuffer.WriteString(uniqueIdentifierField.Name)
	queryBuffer.WriteString("=$1;")

	// Log Query
	query := queryBuffer.String()
	PrintSqlQuery(query, uniqueIdentifierField.Pointer)

	// Execute Query
	row := w.Database.QueryRow(query, uniqueIdentifierField.Pointer)
	err = consumeRow(w, row)
	if err != nil {
		return err
	}

	// Expand foreign references
	return expandForeign(w)
}

// Update updates the model with the current values in the struct
func (w *PqModel) Update() error {
	// Get Unique Identifier
	uniqueIdentifierField, err := getUniqueIdentifier(w)
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
	queryBuffer.WriteString(" RETURNING ")
	for i, field := range w.Config.Fields {
		queryBuffer.WriteString(field.Name)
		if (i + 1) < len(w.Config.Fields) {
			queryBuffer.WriteString(", ")
		}
	}
	queryBuffer.WriteString(";")

	// Get Value Fields
	var valueFields []interface{}
	for _, value := range updatableFields {
		valueFields = append(valueFields, value.Pointer)
	}
	valueFields = append(valueFields, uniqueIdentifierField.Pointer)

	// Log Query
	query := queryBuffer.String()
	PrintSqlQuery(query, valueFields...)

	// Execute Query
	row := w.Database.QueryRow(query, valueFields...)
	err = consumeRow(w, row)
	if err != nil {
		return err
	}

	// Expand foreign references
	return expandForeign(w)
}

// Delete deletes the model
func (w *PqModel) Delete() error {
	// Get Unique Identifier
	uniqueIdentifierField, err := getUniqueIdentifier(w)
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

	// Log Query
	query := queryBuffer.String()
	PrintSqlQuery(query, uniqueIdentifierField.Pointer)

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

// BulkFetch gets an array of models
func (w *PqModel) BulkFetch(fetchConfig BulkFetchConfig, buildModel BuildModel) ([]Model, error) {
	// Set up values
	values := make([]interface{}, 0)

	// Generate query
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("SELECT ")
	for i, field := range w.Config.Fields {
		queryBuffer.WriteString(field.Name)
		if (i + 1) < len(w.Config.Fields) {
			queryBuffer.WriteString(", ")
		}
	}
	queryBuffer.WriteString(" FROM ")
	queryBuffer.WriteString(buildModel().GetConfiguration().TableName)
	if len(fetchConfig.Predicates) > 0 {
		// WHERE
		queryBuffer.WriteString(" ")
		predicatesStr, predicateValues := predicatesToString(1, fetchConfig.Predicates)

		values = append(values, predicateValues...)
		queryBuffer.WriteString(predicatesStr)
	}
	if len(fetchConfig.OrderBys) > 0 {
		queryBuffer.WriteString(" ORDER BY ")
	}
	for i, orderBy := range fetchConfig.OrderBys {
		// Validate that the orderBy.Field is a field
		valid := false
		for _, field := range w.Config.Fields {
			if field.Name == orderBy.Field {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("Could not order table '%v' by the invalid column '%v'",
				w.Config.TableName, orderBy.Field)
		}
		// Write to query
		queryBuffer.WriteString(orderBy.toString())
		if (i + 1) < len(fetchConfig.OrderBys) {
			queryBuffer.WriteString(", ")
		}
	}
	queryBuffer.WriteString(" LIMIT ")
	queryBuffer.WriteString(strconv.Itoa(fetchConfig.Limit))
	queryBuffer.WriteString(" OFFSET ")
	queryBuffer.WriteString(strconv.Itoa(fetchConfig.Offset))
	queryBuffer.WriteString(";")

	// Log Query
	query := queryBuffer.String()
	PrintSqlQuery(query, values...)

	// Execute Query
	rows, err := w.Database.Query(query, values...)
	if err != nil {
		return nil, err
	}

	// Stuff into []Model
	models := make([]Model, 0)
	for rows.Next() {
		model := buildModel()

		// Consume Rows
		fields := model.GetConfiguration().Fields
		var s []interface{}
		for _, value := range fields {
			s = append(s, value.Pointer)
		}
		err := rows.Scan(s...)
		if err != nil {
			return nil, err
		}

		models = append(models, model.(Model))
	}

	// Expand foreign references
	err = expandForeigns(buildModel, models)
	if err != nil {
		return nil, err
	}

	// OK
	return models, nil
}
