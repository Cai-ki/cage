package sugar

import (
	"fmt"
	"reflect"
	"strconv"
)

func StrToT[T any](str string) (T, error) {
	var zero T
	targetType := reflect.TypeOf(zero)

	if targetType.Kind() == reflect.String {
		return any(str).(T), nil
	}

	var result any
	var err error

	switch targetType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		i, err = strconv.ParseInt(str, 10, 64)
		result = i
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var u uint64
		u, err = strconv.ParseUint(str, 10, 64)
		result = u
	case reflect.Float32, reflect.Float64:
		var f float64
		f, err = strconv.ParseFloat(str, 64)
		result = f
	case reflect.Bool:
		var b bool
		b, err = strconv.ParseBool(str)
		result = b
	default:
		err = fmt.Errorf("unsupported type %v", targetType)
	}

	if err != nil {
		return zero, err
	}

	rv := reflect.ValueOf(result).Convert(targetType)
	return rv.Interface().(T), nil
}

func StrToTWithDefault[T any](str string, def T) T {
	var zero T
	targetType := reflect.TypeOf(zero)

	if targetType.Kind() == reflect.String {
		return any(str).(T)
	}

	var result any
	var err error

	switch targetType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		i, err = strconv.ParseInt(str, 10, 64)
		result = i
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var u uint64
		u, err = strconv.ParseUint(str, 10, 64)
		result = u
	case reflect.Float32, reflect.Float64:
		var f float64
		f, err = strconv.ParseFloat(str, 64)
		result = f
	case reflect.Bool:
		var b bool
		b, err = strconv.ParseBool(str)
		result = b
	default:
		return def
	}

	if err != nil {
		return def
	}

	rv := reflect.ValueOf(result).Convert(targetType)
	return rv.Interface().(T)
}
