package surf

import (
	"strconv"
	"strings"
)

// OrderByType is an enumeration of the SQL standard order by
type OrderByType int

const (
	ORDER_BY_ASC OrderByType = iota
	ORDER_BY_DESC
)

type PredicateType int

const (
	WHERE_IS_NOT_NULL PredicateType = iota
	WHERE_IS_NULL
	WHERE_IN
	WHERE_LIKE
	WHERE_EQUAL
	WHERE_NOT_EQUAL
	WHERE_GREATER_THAN
	WHERE_GREATER_THAN_OR_EQUAL_TO
	WHERE_LESS_THAN
	WHERE_LESS_THAN_OR_EQUAL_TO
)

// OrderBy is the definition of a single order by clause
type OrderBy struct {
	Field string
	Type  OrderByType
}

// Predicate is the definition of a single where predicate
type Predicate struct {
	Field         string
	PredicateType PredicateType
	Values        []interface{}
}

func (p *Predicate) ToString(valueIndex int) (string, []interface{}) {
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
	if p.PredicateType != WHERE_IS_NOT_NULL && p.PredicateType != WHERE_IS_NULL {
		if len(p.Values) > 1 {
			predicate += "("
		}
		for i, value := range p.Values {
			values = append(values, value)
			predicate += "$" + strconv.Itoa(valueIndex)
			valueIndex++
			if i < len(p.Values)-1 {
				predicate += ", "
			}
		}
		if len(p.Values) > 1 {
			predicate += ")"
		}
	}

	return predicate, values
}

// ToString converts an OrderBy to SQL
func (ob *OrderBy) ToString() string {
	obType := ""
	switch ob.Type {
	case ORDER_BY_ASC:
		obType = " ASC"
	case ORDER_BY_DESC:
		obType = " DESC"
	}
	return ob.Field + obType
}

func PredicatesToString(valueIndex int, predicates []Predicate) (string, []interface{}) {
	values := make([]interface{}, 0)

	predicateStr := ""
	if len(predicates) > 0 {
		predicateStr += "WHERE "
	}
	for i, predicate := range predicates {
		iPredicateStr, iValues := predicate.ToString(valueIndex)
		valueIndex += len(iValues)
		values = append(values, iValues...)
		predicateStr += iPredicateStr
		if i < (len(predicates) - 1) {
			predicateStr += " AND "
		}
	}
	return predicateStr, values
}

// BulkFetchConfig is the configuration of a Model.BulkFetch()
type BulkFetchConfig struct {
	Limit      int
	Offset     int
	OrderBys   []OrderBy
	Predicates []Predicate
}

// ConsumeSortQuery consumes a `sort` query parameter
// and stuffs them into the OrderBys field
func (c *BulkFetchConfig) ConsumeSortQuery(sortQuery string) {
	var orderBys []OrderBy
	for _, sort := range strings.Split(sortQuery, ",") {
		if string(sort[0]) == "-" {
			orderBys = append(orderBys, OrderBy{
				Field: sort[1:],
				Type:  ORDER_BY_DESC,
			})
		} else {
			orderBys = append(orderBys, OrderBy{
				Field: sort,
				Type:  ORDER_BY_ASC,
			})
		}
	}
	c.OrderBys = orderBys
}
