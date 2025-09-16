package pg

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgtype"
)

func ResultValueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case time.Duration:
		return v.String()
	case time.Time:
		return v.String()
	case int:
		return fmt.Sprintf("%d", v)
	case int8:
		return fmt.Sprintf("%d", v)
	case int16:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case uint:
		return fmt.Sprintf("%d", v)
	case uint8:
		return fmt.Sprintf("%d", v)
	case uint16:
		return fmt.Sprintf("%d", v)
	case uint32:
		return fmt.Sprintf("%d", v)
	case uint64:
		return fmt.Sprintf("%d", v)
	case []byte:
		return fmt.Sprintf("%d", v)
	case pgtype.Float4Array:
		var l []string
		for _, e := range v.Elements {
			l = append(l, fmt.Sprintf("%f", e.Float))
		}
		return fmt.Sprintf("[%s]", strings.Join(l, ","))
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("unknown datatype %v", value)
	}
}

// Result represents one row in a query result
type Result map[string]string

// Results represent all rows in a query result
type Results []Result

// NewResultFromByteArrayArray is used to convert the result of a query into a Results object
func NewResultFromByteArrayArray(cols []string, values []interface{}) (ofr Result, err error) {
	ofr = make(Result)
	if len(cols) != len(values) {
		return ofr, fmt.Errorf("number of cols different then number of values")
	}
	for i, col := range cols {
		ofr[col] = ResultValueToString(values[i])
	}
	return ofr, nil
}

// String is used to get the string value of a Result
func (ofr Result) String() (s string) {
	var results []string
	for key, value := range ofr {
		results = append(results, fmt.Sprintf("%s: %s",
			FormattedString(key),
			FormattedString(value)))
	}
	return fmt.Sprintf("{ %s }", strings.Join(results, ", "))
}

// Columns return a list of all columns of this recordset
func (ofr Result) Columns() (result []string) {
	for key := range ofr {
		result = append(result, key)
	}
	return result
}

// FormattedString is used to return a string value in a formatted fasion
func FormattedString(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "\\'"))
}

// Compare can be used to compare a result with another result
func (ofr Result) Compare(other Result) (err error) {
	if len(ofr) != len(other) {
		return fmt.Errorf("number of columns different between row %v and compared row %v",
			ofr.Columns(), other.Columns())
	}
	for key, value := range ofr {
		otherValue, exists := other[key]
		if !exists {
			return fmt.Errorf("column row (%s) not in compared row", FormattedString(key))
		}
		if matched, err := regexp.MatchString(otherValue, value); err != nil {
			if value != otherValue {
				return fmt.Errorf("comparedrow is not an re, and column %s differs between row (%s), and comparedrow (%s)",
					FormattedString(key),
					FormattedString(value),
					FormattedString(otherValue))
			}
		} else if !matched {
			return fmt.Errorf("column %s value (%s) does not match with regular expression (%s)",
				FormattedString(key),
				FormattedString(value),
				FormattedString(otherValue))
		}
	}
	return nil
}

// String can be used to get the string value of multiple results
func (results Results) String() (s string) {
	var arr []string
	if len(results) == 0 {
		return "[ ]"
	}
	for _, result := range results {
		arr = append(arr, result.String())
	}
	return fmt.Sprintf("[ %s ]", strings.Join(arr, ", "))
}

// Compare can be used to compare Results with other Results
func (results Results) Compare(other Results) (err error) {
	if len(results) != len(other) {
		return fmt.Errorf("different result (%s) then expected (%s)", results.String(),
			other.String())
	}
	for i, result := range results {
		err = result.Compare(other[i])
		if err != nil {
			return fmt.Errorf("different %d'th result: %s", i, err.Error())
		}
	}
	return nil
}
