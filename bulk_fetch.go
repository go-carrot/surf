package surf

import (
	"strings"
)

// OrderByType is an enumeration of the SQL standard order by
type OrderByType int

const (
	ORDER_BY_ASC OrderByType = iota
	ORDER_BY_DESC
)

// OrderBy is the definition of a single order by clause
type OrderBy struct {
	Field string
	Type  OrderByType
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

// BulkFetchConfig is the configuration of a Model.BulkFetch()
type BulkFetchConfig struct {
	Limit    int
	Offset   int
	OrderBys []OrderBy
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
