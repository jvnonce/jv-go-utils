package maps

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/iancoleman/strcase"
	jve "github.com/jvnonce/jv-go-utils/lib/errors"
)

type Mapping interface {
	ToMap(fields ...string) M
}

/************* Map *************/

// Add-on for map
type M map[string]interface{}

// Scanner for map
func (m *M) Scan(value interface{}) error {
	if value == nil {
		*m = make(M)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return jve.ErrBadType
	}

	var result map[string]interface{}
	err := json.Unmarshal(bytes, &result)
	*m = M(result)
	return err
}

// Valuer for map
func (m M) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}

// Get raw bytes from map
func (m M) Bytes() ([]byte, error) {
	v, err := m.Value()
	return v.([]byte), err
}

// Get map as string
func (m M) AsString() string {
	b, err := m.Bytes()
	if err != nil {
		return ""
	}
	return string(b)
}

// Create Buffer from map
func (m M) Buffer() (*bytes.Buffer, error) {
	b, err := m.Bytes()
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(b)
	return buffer, nil
}

// Returns keys of map as slice
func (m M) Keys() []string {
	result := make([]string, len(m))
	index := 0
	for key := range m {
		result[index] = key
		index++
	}
	return result
}

// Copy map into new map
func (m M) Copy() M {
	result := make(M, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// Copy map into new map without some fields
func (m M) CopyExclude(fields ...string) M {
	result := make(M, len(m)-len(fields))
	for k, v := range m {
		found := false
		for _, f := range fields {
			if f == k {
				found = true
				break
			}
		}
		if !found {
			result[k] = v
		}
	}
	return result
}

// Returns values of map
func (m M) Values() []any {
	result := make([]any, len(m))
	index := 0
	for _, value := range m {
		result[index] = value
		index++
	}
	return result
}

// Returns values of listed keys
func (m M) ValuesOf(keys ...string) []any {
	result := make([]any, len(keys))
	index := 0
	for _, key := range keys {
		result[index] = m[key]
		index++
	}
	return result
}

// Transform map keys into snake format
func (m M) SnakeKeys() M {
	result := make(M, len(m))
	for key, value := range m {
		result[strcase.ToSnake(key)] = value
	}
	return result
}

// Transform map keys into camel format
func (m M) CamelKeys() M {
	result := make(M, len(m))
	for key, value := range m {
		result[strcase.ToLowerCamel(key)] = value
	}
	return result
}

// Return value of path
func (m M) ValueOf(key string) interface{} {
	path := strings.Split(key, ".")
	current := m
	count := len(path)
	var result interface{} = nil
	for index, item := range path {
		if current[item] == nil {
			return nil
		}
		if index < count-1 {
			var ok bool
			current, ok = current[item].(map[string]interface{})
			if !ok {
				return nil
			}
		} else {
			result = current[item]
		}
	}
	return result
}

/************* Array *************/

// Add-on for slice
type A []interface{}

// Scanner for slice
func (a *A) Scan(value interface{}) error {
	if value == nil {
		*a = make(A, 0)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return jve.ErrBadType
	}

	var result []interface{}
	err := json.Unmarshal(bytes, &result)
	*a = A(result)
	return err
}

// // Valuer for slice
func (a A) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}
