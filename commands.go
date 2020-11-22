package movies

type CreateMovieCommand struct {
	Name   string `json:"name"`
	Rating int64  `json:"rating"`
}

func (cmd *CreateMovieCommand) Exec(service MovieService) (interface{}, error) {
	return service.Create(cmd)
}

type UpdateMovieCommand struct {
	Id     int64   `json:"id"`
	Name   *string `json:"name"`
	Rating *int64  `json:"rating"`
}

func (cmd *UpdateMovieCommand) Exec(service MovieService) (interface{}, error) {
	return service.Update(cmd)
}

type DeleteMovieCommand struct {
	Id int64 `json:"id"`
}

func (cmd *DeleteMovieCommand) Exec(service MovieService) (interface{}, error) {
	return service.Delete(cmd)
}

type GetMovieCommand struct {
	Id int64 `json:"id"`
}

func (cmd *GetMovieCommand) Exec(service MovieService) (interface{}, error) {
	return service.Get(cmd)
}

type ListMovieCommand struct {
}

func (cmd *ListMovieCommand) Exec(service MovieService) (interface{}, error) {
	return service.List(cmd)
}
