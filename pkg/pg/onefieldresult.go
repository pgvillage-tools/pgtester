package pg

import (
	"fmt"
	"strings"
	"time"
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
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("unknown datatype %v", value)
	}
}

type Result map[string]string
type Results []Result

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

func (ofr Result) String() (s string) {
	var results []string
	for key, value := range ofr {
		key = strings.Replace(key, "'", "\\'", -1)
		value = strings.Replace(value, "'", "\\'", -1)
		results = append(results, fmt.Sprintf("'%s': '%s'", key, value))
	}
	return fmt.Sprintf("{ %s }", strings.Join(results, "', '"))
}

func (ofr Result) Columns() (result []string) {
	for key := range ofr {
		result = append(result, key)
	}
	return result
}
func (ofr Result) Compare(other Result) (err error) {
	if len(ofr) != len(other) {
		return fmt.Errorf("number of columns different between row [ '%v' ] and compared row [ '%v' ]",
			strings.Join(ofr.Columns(), "', '"), strings.Join(other.Columns(), "', '"))
	}
	for key, value := range ofr {
		otherValue, exists := other[key]
		if !exists {
			return fmt.Errorf("column row ('%s') not in compared row", key)
		}
		if value != otherValue {
			return fmt.Errorf("column '%s' differs between row ('%s'), and comparedrow ('%s')", key, value, otherValue)
		}
	}
	return nil
}

func (results Results) String() (s string) {
	var arr []string
	if len(results) == 0 {
		return "[ ]"
	}
	for _, result := range results {
		arr = append(arr, result.String())
	}
	return fmt.Sprintf("[ %s ]", strings.Join(arr, "}, {"))
}

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
