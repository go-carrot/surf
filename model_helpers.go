package surf

import (
	"database/sql"
	"errors"
	"fmt"
	"gopkg.in/guregu/null.v3"
)

// getUniqueIdentifier Returns the unique identifier that this model will
// query against.
//
// This will return the first field in Configuration.Fields that:
//   - Has `UniqueIdentifier` set to true
//   - Returns true from `IsSet`
//
// This function will panic in the event that it encounters a field that is a
// `UniqueIdentifier`, and doesn't have `IsSet` implemented.
func getUniqueIdentifier(w Model) (Field, error) {
	// Get all unique identifier fields
	var uniqueIdentifierFields []Field
	for _, field := range w.GetConfiguration().Fields {
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

// expandForeign expands all foreign references for a single Model
func expandForeign(model Model) error {
	// Load all foreign references
	for _, field := range model.GetConfiguration().Fields {

		// If it's a set foreign reference
		if field.GetReference != nil && field.SetReference != nil && field.IsSet(field.Pointer) {

			// Get the reference type
			modelBuilder, identifier := field.GetReference()
			model := modelBuilder()

			// Set the identifier on the foreign reference
			// The foreign reference value may only be a `null.Int` or an `int64`
			// The identifier on the foreign model may only be of type `int64`
			for _, modelField := range model.GetConfiguration().Fields {
				if modelField.Name == identifier {
					switch tv := field.Pointer.(type) {
					case *int64:
						*(modelField.Pointer.(*int64)) = *tv
						break
					case *null.Int:
						*(modelField.Pointer.(*int64)) = tv.Int64
						break
					}
					break
				}
			}

			// Load
			err := model.Load()
			if err != nil {
				return err
			}

			// Set reference
			err = field.SetReference(model)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// expandForeigns expands all foreign references for an array of Model
func expandForeigns(modelBuilder BuildModel, models []Model) error {
	// Expand all foreign references
	for _, field := range modelBuilder().GetConfiguration().Fields {
		// If the field is a foreign key
		if field.GetReference != nil && field.SetReference != nil {
			builder, foreignField := field.GetReference()
			err := expandForeignsByField(field.Name, builder, foreignField, models)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// expandForeignsByField expands a single foreign key for an array of Model
func expandForeignsByField(fieldName string, foreignBuilder BuildModel, foreignField string, models []Model) error {
	// Get Foreign IDs
	ids := make([]interface{}, 0)
	for _, model := range models {
		for _, field := range model.GetConfiguration().Fields {
			if field.Name == fieldName {
				switch tv := field.Pointer.(type) {
				case *int64:
					ids = appendIfMissing(ids, *tv)
					break
				case *null.Int:
					if tv.Valid {
						ids = appendIfMissing(ids, tv.Int64)
					}
					break
				}
			}
		}
	}

	// Load Foreign models
	foreignModels, err := foreignBuilder().BulkFetch(
		BulkFetchConfig{
			Limit: len(ids),
			Predicates: []Predicate{{
				Field:         foreignField,
				PredicateType: WHERE_IN,
				Values:        ids,
			}},
		},
		foreignBuilder,
	)
	if err != nil {
		return err
	}

	// Stuff foreign models into models
	for _, model := range models {
		for _, field := range model.GetConfiguration().Fields {
			if field.Name == fieldName {
				var toMatch int64
				switch tv := field.Pointer.(type) {
				case *int64:
					toMatch = *tv
					break
				case *null.Int:
					toMatch = tv.Int64
					break
				}

			MatchForeignModel:
				for _, foreignModel := range foreignModels {
				FindField:
					for _, foreignModelField := range foreignModel.GetConfiguration().Fields {
						if foreignModelField.Name == foreignField {

							if *(foreignModelField.Pointer.(*int64)) == toMatch {
								field.SetReference(foreignModel)
								break MatchForeignModel
							}
							break FindField
						}
					}
				}
				break
			}
		}
	}
	return nil
}

// appendIfMissing functions like append(), but will only add the
// int64 to the slice if it doesn't exist in the slice already
func appendIfMissing(slice []interface{}, i int64) []interface{} {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// consumeRow Scans a *sql.Row into our struct
// that is using this model
func consumeRow(w Model, row *sql.Row) error {
	fields := w.GetConfiguration().Fields
	var s []interface{}
	for _, value := range fields {
		s = append(s, value.Pointer)
	}
	return row.Scan(s...)
}
