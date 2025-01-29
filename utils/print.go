package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"
)

func Dump(data interface{}) {
	val := reflect.ValueOf(data)

	// Check if it's a pointer and resolve it
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		dumpStruct(val)
	case reflect.Slice, reflect.Array:
		dumpSliceOrArray(val)
	case reflect.Map:
		dumpMap(val)
	default:
		dumpBasic(val)
	}
}

func Dd(data interface{}) {
	Dump(data)
	os.Exit(0)
}

func DumpJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

func DdJSON(data interface{}) {
	DumpJSON(data)

	os.Exit(1)
}

func dumpStruct(val reflect.Value) {
	typ := val.Type()
	fmt.Printf("Struct: %s\n", typ.Name())

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Skip time.Time fields
		if fieldType.Type == reflect.TypeOf(time.Time{}) {
			fmt.Printf("%s: %v\n", fieldType.Name, field)

			continue
		}

		fmt.Printf("%s: %v\n", fieldType.Name, field.Interface())

		// Recurse into nested structs
		if field.Kind() == reflect.Struct {
			Dump(field.Interface())
		}
	}
}

func dumpSliceOrArray(val reflect.Value) {
	fmt.Println("Slice/Array:")

	for i := 0; i < val.Len(); i++ {
		fmt.Printf("[%d]: %v\n", i, val.Index(i).Interface())
	}
}

func dumpMap(val reflect.Value) {
	fmt.Println("Map:")

	for _, key := range val.MapKeys() {
		fmt.Printf("%v: %v\n", key.Interface(), val.MapIndex(key).Interface())
	}
}

func dumpBasic(val reflect.Value) {
	switch val.Kind() {

	case reflect.String:
		fmt.Printf("String: %s\n", val.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Printf("Integer: %d\n", val.Int())
	case reflect.Float32, reflect.Float64:
		fmt.Printf("Float: %f\n", val.Float())
	case reflect.Bool:
		fmt.Printf("Boolean: %v\n", val.Bool())
	default:
		fmt.Printf("Type: %s, Value: %v\n", val.Kind(), val.Interface())
	}
}
