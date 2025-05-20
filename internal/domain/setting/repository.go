package setting

type Repository interface {
	List() ([]Setting, error)
	CreateMany(settings []Setting) (bool, error)
	Update(settings []Setting) (bool, error)
}
