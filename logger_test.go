package surf_test

import (
	"github.com/go-carrot/surf"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v3"
	"testing"
	"time"
)

type StackWriter struct {
	Stack []string
}

func (sw *StackWriter) Write(p []byte) (n int, err error) {
	sw.Stack = append(sw.Stack, string(p))
	return 0, nil
}

func (sw *StackWriter) Peek() string {
	if len(sw.Stack) > 0 {
		return sw.Stack[len(sw.Stack)-1]
	}
	return ""
}

func TestLogger(t *testing.T) {
	// Enable logging
	stackWriter := &StackWriter{}
	surf.SetLogging(true, stackWriter)

	// float32
	var idFloat32 float32 = 8.8
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idFloat32)
	assert.Equal(t, "SELECT * FROM table WHERE id = 8.8", stackWriter.Peek())

	// float64
	var idFloat64 float64 = 8.9
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idFloat64)
	assert.Equal(t, "SELECT * FROM table WHERE id = 8.9", stackWriter.Peek())

	// bool
	var idBool bool = false
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idBool)
	assert.Equal(t, "SELECT * FROM table WHERE id = false", stackWriter.Peek())

	// int
	var idInt int = 190
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idInt)
	assert.Equal(t, "SELECT * FROM table WHERE id = 190", stackWriter.Peek())

	// int8
	var idInt8 int8 = 8
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idInt8)
	assert.Equal(t, "SELECT * FROM table WHERE id = 8", stackWriter.Peek())

	// int16
	var idInt16 int16 = 111
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idInt16)
	assert.Equal(t, "SELECT * FROM table WHERE id = 111", stackWriter.Peek())

	// int32
	var idInt32 int32 = 1110
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idInt32)
	assert.Equal(t, "SELECT * FROM table WHERE id = 1110", stackWriter.Peek())

	// int64
	var idInt64 int64 = 11100
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idInt64)
	assert.Equal(t, "SELECT * FROM table WHERE id = 11100", stackWriter.Peek())

	// uint
	var idUInt uint = 200
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idUInt)
	assert.Equal(t, "SELECT * FROM table WHERE id = 200", stackWriter.Peek())

	// uint8
	var idUInt8 uint8 = 127
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idUInt8)
	assert.Equal(t, "SELECT * FROM table WHERE id = 127", stackWriter.Peek())

	// uint16
	var idUInt16 uint16 = 1278
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idUInt16)
	assert.Equal(t, "SELECT * FROM table WHERE id = 1278", stackWriter.Peek())

	// uint32
	var idUInt32 uint32 = 12788
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idUInt32)
	assert.Equal(t, "SELECT * FROM table WHERE id = 12788", stackWriter.Peek())

	// uint64
	var idUInt64 uint64 = 127888
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idUInt64)
	assert.Equal(t, "SELECT * FROM table WHERE id = 127888", stackWriter.Peek())

	// time.Time
	const layout = "Jan 2, 2006 at 3:04pm (MST)"
	idTime, _ := time.Parse(layout, "Feb 3, 2013 at 7:54pm (PST)")
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idTime)
	assert.Equal(t, "SELECT * FROM table WHERE id = '2013-02-03T19:54:00Z'", stackWriter.Peek())

	// null.Int
	idNullInt := null.IntFrom(100)
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullInt)
	assert.Equal(t, "SELECT * FROM table WHERE id = 100", stackWriter.Peek())

	idNullIntNull := null.Int{}
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullIntNull)
	assert.Equal(t, "SELECT * FROM table WHERE id = null", stackWriter.Peek())

	// null.String
	idNullString := null.StringFrom("Hello")
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullString)
	assert.Equal(t, "SELECT * FROM table WHERE id = 'Hello'", stackWriter.Peek())

	idNullStringNull := null.String{}
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullStringNull)
	assert.Equal(t, "SELECT * FROM table WHERE id = null", stackWriter.Peek())

	// null.Bool
	idNullBool := null.BoolFrom(false)
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullBool)
	assert.Equal(t, "SELECT * FROM table WHERE id = false", stackWriter.Peek())

	idNullBoolNull := null.Bool{}
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullBoolNull)
	assert.Equal(t, "SELECT * FROM table WHERE id = null", stackWriter.Peek())

	// null.Float
	idNullFloat := null.FloatFrom(1.2)
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullFloat)
	assert.Equal(t, "SELECT * FROM table WHERE id = 1.2", stackWriter.Peek())

	idNullFloatNull := null.Float{}
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullFloatNull)
	assert.Equal(t, "SELECT * FROM table WHERE id = null", stackWriter.Peek())

	// null.Time
	idNullTime := null.TimeFrom(idTime)
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullTime)
	assert.Equal(t, "SELECT * FROM table WHERE id = '2013-02-03T19:54:00Z'", stackWriter.Peek())

	idNullTimeNull := null.Time{}
	surf.PrintSqlQuery("SELECT * FROM table WHERE id = $1", &idNullTimeNull)
	assert.Equal(t, "SELECT * FROM table WHERE id = null", stackWriter.Peek())
}
