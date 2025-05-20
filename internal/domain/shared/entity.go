package shared

import (
	"fmt"
	"reflect"
)

type BaseEntity struct {
}

// PopulateModel populates the fields of a model from another struct or a pointer to a struct with matching field names.
func (b *BaseEntity) PopulateModel(model interface{}, data interface{}) error {
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
