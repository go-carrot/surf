package surf

import (
	"strings"
)

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
