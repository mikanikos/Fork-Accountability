package common

import (
	"reflect"
	"strconv"
)

// Value represents a value of a message, it can contain other information if desired
type Value struct {
	Data int64 `yaml:"data"`
}

// NewValue creates a new value
func NewValue(val int64) *Value {
	return &Value{
		Data: val,
	}
}

// Equal is the equality method for values
func (value *Value) Equal(other *Value) bool {
	return reflect.DeepEqual(value, other)
}

// String representation of a value
func (value *Value) String() string {
	return strconv.FormatInt(value.Data, 10)
}
