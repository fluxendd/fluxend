package table

type Repository interface {
	List(projectUUID string) ([]Table, error)
}
