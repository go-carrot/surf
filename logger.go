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

// SetLogging adjusts the configuration for logging. You can enable
// and disable the logging here. By default, logging is disabled.
//
// Most calls to this function will be called like SetLogging(true, os.Stdout)
func SetLogging(enabled bool, writer io.Writer) {
	loggingEnabled = enabled
	loggingWriter = writer
}

// printQuery prints a query if the user has enabled logging
func PrintSqlQuery(query string, args ...interface{}) {
	if loggingEnabled {
		for i, arg := range args {
			query = strings.Replace(query, "$"+strconv.Itoa(i+1), pointerToLogString(arg), 1)
		}
		fmt.Fprint(loggingWriter, query)
	}
}

// pointerToLogString converts a value pointer to the string
// that should be logged for it
func pointerToLogString(pointer interface{}) string {
	switch v := pointer.(type) {
	case *string:
		return "'" + *v + "'"
	case *float32:
		return strconv.FormatFloat(float64(*v), 'f', -1, 32)
	case *float64:
		return strconv.FormatFloat(*v, 'f', -1, 64)
	case *bool:
		return strconv.FormatBool(*v)
	case *int:
		return strconv.Itoa(*v)
	case *int8:
		return strconv.FormatInt(int64(*v), 10)
	case *int16:
		return strconv.FormatInt(int64(*v), 10)
	case *int32:
		return strconv.FormatInt(int64(*v), 10)
	case *int64:
		return strconv.FormatInt(*v, 10)
	case *uint:
		return strconv.FormatUint(uint64(*v), 10)
	case *uint8:
		return strconv.FormatUint(uint64(*v), 10)
	case *uint16:
		return strconv.FormatUint(uint64(*v), 10)
	case *uint32:
		return strconv.FormatUint(uint64(*v), 10)
	case *uint64:
		return strconv.FormatUint(*v, 10)
	case *time.Time:
		return "'" + (*v).Format(time.RFC3339) + "'"
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
	default:
		return fmt.Sprintf("%v", v)
	}
	return "null"
}
