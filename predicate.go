package surf

import (
	"strconv"
)

type PredicateType int

const (
	WHERE_IS_NOT_NULL PredicateType = iota
	WHERE_IS_NULL
	WHERE_IN
	WHERE_NOT_IN
	WHERE_LIKE
	WHERE_EQUAL
	WHERE_NOT_EQUAL
	WHERE_GREATER_THAN
	WHERE_GREATER_THAN_OR_EQUAL_TO
	WHERE_LESS_THAN
	WHERE_LESS_THAN_OR_EQUAL_TO
)

// getPredicateTypeString returns the predicate type string from it's value
func getPredicateTypeString(predicateType PredicateType) string {
	switch predicateType {
	case WHERE_IS_NOT_NULL:
		return "WHERE_IS_NOT_NULL"
	case WHERE_IS_NULL:
		return "WHERE_IS_NULL"
	case WHERE_IN:
		return "WHERE_IN"
	case WHERE_NOT_IN:
		return "WHERE_NOT_IN"
	case WHERE_LIKE:
		return "WHERE_LIKE"
	case WHERE_EQUAL:
		return "WHERE_EQUAL"
	case WHERE_NOT_EQUAL:
		return "WHERE_NOT_EQUAL"
	case WHERE_GREATER_THAN:
		return "WHERE_GREATER_THAN"
	case WHERE_GREATER_THAN_OR_EQUAL_TO:
		return "WHERE_GREATER_THAN_OR_EQUAL_TO"
	case WHERE_LESS_THAN:
		return "WHERE_LESS_THAN"
	case WHERE_LESS_THAN_OR_EQUAL_TO:
		return "WHERE_LESS_THAN_OR_EQUAL_TO"
	}
	return ""
}

// Predicate is the definition of a single where SQL predicate
type Predicate struct {
	Field         string
	PredicateType PredicateType
	Values        []interface{}
}

// toString will convert a predicate to it's query string, along with its values
// to be passed along with the query
//
// This function will panic in the event that this is called on a malformed predicate
func (p *Predicate) toString(valueIndex int) (string, []interface{}) {
	// Field
	predicate := p.Field

	// Type
	switch p.PredicateType {
	case WHERE_IS_NOT_NULL:
		predicate += " IS NOT NULL"
		break
	case WHERE_IS_NULL:
		predicate += " IS NULL"
		break
	case WHERE_IN:
		predicate += " IN "
		break
	case WHERE_NOT_IN:
		predicate += " NOT IN "
		break
	case WHERE_LIKE:
		predicate += " LIKE "
		break
	case WHERE_EQUAL:
		predicate += " = "
		break
	case WHERE_NOT_EQUAL:
		predicate += " != "
		break
	case WHERE_GREATER_THAN:
		predicate += " > "
		break
	case WHERE_GREATER_THAN_OR_EQUAL_TO:
		predicate += " >= "
		break
	case WHERE_LESS_THAN:
		predicate += " < "
		break
	case WHERE_LESS_THAN_OR_EQUAL_TO:
		predicate += " <= "
		break
	}

	// Values
	values := make([]interface{}, 0)
	switch p.PredicateType {
	case WHERE_IN,
		WHERE_NOT_IN:
		if len(p.Values) == 0 {
			panic("`" + getPredicateTypeString(p.PredicateType) + "` predicates require at least one value.")
		}
		predicate += "("
		for i, value := range p.Values {
			values = append(values, value)
			predicate += "$" + strconv.Itoa(valueIndex)
			valueIndex++
			if i < len(p.Values)-1 {
				predicate += ", "
			}
		}
		predicate += ")"
		break
	case WHERE_LIKE,
		WHERE_EQUAL,
		WHERE_NOT_EQUAL,
		WHERE_GREATER_THAN,
		WHERE_GREATER_THAN_OR_EQUAL_TO,
		WHERE_LESS_THAN,
		WHERE_LESS_THAN_OR_EQUAL_TO:
		if len(p.Values) != 1 {
			panic("`" + getPredicateTypeString(p.PredicateType) + "` predicates require exactly one value.")
		}
		values = append(values, p.Values[0])
		predicate += "$" + strconv.Itoa(valueIndex)
		break
	}

	return predicate, values
}

// predicatesToString converts an array of predicates to a query string, along with its values
// to be passed along with the query
//
// This function will panic in the event that it encounters a malformed predicate
func predicatesToString(valueIndex int, predicates []Predicate) (string, []interface{}) {
	values := make([]interface{}, 0)

	predicateStr := ""
	if len(predicates) > 0 {
		predicateStr += "WHERE "
	}
	for i, predicate := range predicates {
		iPredicateStr, iValues := predicate.toString(valueIndex)
		valueIndex += len(iValues)
		values = append(values, iValues...)
		predicateStr += iPredicateStr
		if i < (len(predicates) - 1) {
			predicateStr += " AND "
		}
	}
	return predicateStr, values
}
