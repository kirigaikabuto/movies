package movies

import "errors"

type MovieService interface {
	List(cmd *ListMovieCommand) ([]Movie, error)
	Get(cmd *GetMovieCommand) (*Movie, error)
	Create(cmd *CreateMovieCommand) (*Movie, error)
	Update(cmd *UpdateMovieCommand) (*Movie, error)
	Delete(cmd *DeleteMovieCommand) (*Movie, error)
}

type movieService struct {
	store MovieStore
}

func NewMovieService(store MovieStore) MovieService {
	return &movieService{
		store: store,
	}
}

func (m *movieService) List(cmd *ListMovieCommand) ([]Movie, error) {
	movies, err := m.store.List()
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (m *movieService) Get(cmd *GetMovieCommand) (*Movie, error) {
	movie, err := m.store.Get(cmd.Id)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (m *movieService) Create(cmd *CreateMovieCommand) (*Movie, error) {
	if cmd.Name == "" {
		return nil, errors.New("name field is empty")
	} else if cmd.Rating == 0 {
		return nil, errors.New("rating field is empty")
	}
	movie := &Movie{
		Name:   cmd.Name,
		Rating: cmd.Rating,
	}
	newMovie, err := m.store.Create(movie)
	if err != nil {
		return nil, err
	}
	return newMovie, nil
}

func (m *movieService) Update(cmd *UpdateMovieCommand) (*Movie, error) {
	movieUpdate := &Movie{}
	if cmd.Name != nil {
		movieUpdate.Name = *cmd.Name
	}
	if cmd.Rating != nil {
		movieUpdate.Rating = *cmd.Rating
	}
	if cmd.Id == 0 {
		movieUpdate.Id = cmd.Id
	}
	updatedMovie, err := m.store.Update(movieUpdate)
	if err != nil {
		return nil, err
	}
	return updatedMovie, nil
}

func (m *movieService) Delete(cmd *DeleteMovieCommand) (*Movie, error) {
	if cmd.Id == 0 {
		return nil, errors.New("id is empty")
	}
	movie, err := m.store.Get(cmd.Id)
	if err != nil {
		return nil, err
	}
	err = m.store.Delete(movie.Id)
	if err != nil {
		return nil, err
	}
	return movie, nil
}
