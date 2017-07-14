package surf

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
func (ob *OrderBy) toString() string {
	obType := ""
	switch ob.Type {
	case ORDER_BY_ASC:
		obType = " ASC"
	case ORDER_BY_DESC:
		obType = " DESC"
	}
	return ob.Field + obType
}
