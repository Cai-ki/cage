package jsondb

import (
	"encoding/json"
	"reflect"
)

// unmarshalRecords 反序列化记录到目标切片
func unmarshalRecords(records []*Record, result interface{}) error {
	resultVal := reflect.ValueOf(result)
	if resultVal.Kind() != reflect.Ptr || resultVal.Elem().Kind() != reflect.Slice {
		return ErrInvalidResultType
	}

	sliceType := resultVal.Elem().Type()
	elemType := sliceType.Elem()
	resultSlice := reflect.MakeSlice(sliceType, len(records), len(records))

	for i, rec := range records {
		elem := reflect.New(elemType)
		if err := json.Unmarshal(rec.RawData, elem.Interface()); err != nil {
			return err
		}
		resultSlice.Index(i).Set(elem.Elem())
	}

	resultVal.Elem().Set(resultSlice)
	return nil
}

// 预定义错误
var (
	ErrMissingTimeField  = &jsonError{"time field missing in JSON"}
	ErrInvalidTimeFormat = &jsonError{"invalid time format"}
	ErrInvalidResultType = &jsonError{"result must be a pointer to slice"}
)

type jsonError struct{ msg string }

func (e *jsonError) Error() string { return e.msg }
