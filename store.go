package movies

type MovieStore interface {
	Create(movie *Movie) (*Movie, error)
	Get(id int64) (*Movie, error)
	List() ([]Movie, error)
	Update(movie *Movie) (*Movie, error)
	Delete(id int64) error
}
