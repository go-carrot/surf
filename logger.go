package surf

import (
	"fmt"
	"gopkg.in/guregu/null.v3"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	loggingEnabled = false
	loggingWriter  io.Writer
)

func SetLogging(enabled bool, writer io.Writer) {
	loggingEnabled = enabled
	loggingWriter = writer
}

func printQuery(query string, args ...interface{}) {
	if loggingEnabled {
		for i, arg := range args {
			query = strings.Replace(query, "$"+strconv.Itoa(i+1), pointerToLogString(arg), 1)
		}
		fmt.Fprintln(loggingWriter, "[Surf Query]: "+query)
	}
}

func pointerToLogString(pointer interface{}) string {
	switch v := pointer.(type) {
	case *null.Int:
		if v.Valid {
			return strconv.FormatInt(v.Int64, 10)
		}
		break
	case *null.String:
		if v.Valid {
			return "'" + v.String + "'"
		}
		break
	case *null.Bool:
		if v.Valid {
			return strconv.FormatBool(v.Bool)
		}
		break
	case *null.Float:
		if v.Valid {
			return strconv.FormatFloat(v.Float64, 'f', -1, 64)
		}
		break
	case *null.Time:
		if v.Valid {
			return "'" + v.Time.Format(time.RFC3339) + "'"
		}
		break
	case *string:
		return "'" + *v + "'"
	case *float32:
	case *float64:
		return strconv.FormatFloat(float64(*v), 'f', -1, 64)
	case *bool:
		return strconv.FormatBool(*v)
	case *int:
	case *int8:
	case *int16:
	case *int32:
	case *int64:
		return strconv.FormatInt(int64(*v), 10)
	case *uint:
	case *uint8:
	case *uint16:
	case *uint32:
	case *uint64:
		return strconv.FormatUint(uint64(*v), 10)
	case *time.Time:
		return "'" + (*v).Format(time.RFC3339) + "'"
	default:
		return fmt.Sprintf("%v", v)
	}
	return "null"
}
