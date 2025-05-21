package setting

type Repository interface {
	List() ([]Setting, error)
	Get(name string) (Setting, error)
	CreateMany(settings []Setting) (bool, error)
	Update(settings []Setting) (bool, error)
}
