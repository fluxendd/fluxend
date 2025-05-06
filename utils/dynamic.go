package utils

import (
	"github.com/labstack/gommon/log"
	"reflect"
	"runtime"
	"strings"
)

func GetColumnsList[T any](alias string) []string {
	var fields []string
	var t T
	typ := reflect.TypeOf(t)

	if alias != "" {
		alias = alias + "."
	}

	if typ.Kind() != reflect.Struct {
		log.Error("Error: Expected a struct type")

		return fields
	}

	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem() // Handle pointer types
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" {
			fields = append(fields, alias+dbTag)
		}
	}

	return fields
}

func GetColumnsWithAlias[T any](alias string) string {
	columns := GetColumnsList[T](alias)

	return strings.Join(columns, ", ")
}

func GetColumns[T any]() string {
	columns := GetColumnsList[T]("")

	return strings.Join(columns, ", ")
}

func GetMethodName() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	fullName := fn.Name()

	parts := strings.Split(fullName, "/")

	return parts[len(parts)-1]
}
