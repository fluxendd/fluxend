package resources

type GenericResponse interface{}

func GenericResource(data interface{}) GenericResponse {
	return data
}

func GenericResourceCollection(data interface{}) GenericResponse {
	return data
}
