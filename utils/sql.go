package utils

import (
	"database/sql"
	"reflect"
)

func ScanToStruct(row *sql.Row, dest interface{}) error {
	// Use reflection to get pointers to the struct's fields
	v := reflect.ValueOf(dest).Elem()
	numFields := v.NumField()
	fields := make([]interface{}, numFields)
	for i := 0; i < numFields; i++ {
		fields[i] = v.Field(i).Addr().Interface()
	}
	return row.Scan(fields...)
}

func PointerValue[T any](v any) *T {
	if v == nil {
		return nil
	}
	return v.(*T)
}
