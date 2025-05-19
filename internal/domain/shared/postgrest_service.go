package shared

type PostgrestService interface {
	StartContainer(dbName string)
	RemoveContainer(dbName string)
	HasContainer(dbName string) bool
}
