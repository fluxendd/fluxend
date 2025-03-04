package utils

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"runtime"
	"strings"
)

// PopulateModel populates the fields of a model from another struct or a pointer to a struct with matching field names.
func PopulateModel(model interface{}, data interface{}) error {
	// Ensure model is a pointer to a struct
	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() != reflect.Ptr || modelValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("model must be a pointer to a struct")
	}

	// Dereference the model pointer to access the struct
	modelValue = modelValue.Elem()

	// Handle data, ensuring it is a struct or a pointer to a struct
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		if dataValue.IsNil() {
			return fmt.Errorf("data is a nil pointer")
		}
		dataValue = dataValue.Elem() // Dereference the pointer
	}

	if dataValue.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a struct or pointer to a struct, got %s", dataValue.Kind())
	}

	for i := 0; i < dataValue.NumField(); i++ {
		dataField := dataValue.Type().Field(i)
		dataFieldValue := dataValue.Field(i)

		// Look for a matching field in the model
		if modelField := modelValue.FieldByName(dataField.Name); modelField.IsValid() && modelField.CanSet() {
			// Ensure the types match before setting the value
			if modelField.Type() == dataFieldValue.Type() {
				modelField.Set(dataFieldValue)
			}
		}
	}

	return nil
}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash)
}

func ComparePassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	return err == nil
}

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

func PointerToString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func BytesToKiloBytes(bytes int) int {
	return bytes / 1024
}

func GetMethodName() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	fullName := fn.Name()

	parts := strings.Split(fullName, "/")

	return parts[len(parts)-1]
}

func FormatError(err error, errType, method string) error {
	// bucketRepo.ListForProject: select.err => <error>
	return fmt.Errorf("%s: %s.err => %v", method, errType, err)
}
