package movies

type AddMovieCommand struct {
	Name   string `json:"name"`
	Rating int64  `json:"rating"`
}

type UpdateMovieCommand struct {
	Name   *string `json:"name"`
	Rating *int64  `json:"rating"`
}

type DeleteMovieCommand struct {
	Id int64 `json:"id"`
}

type GetMovieCommand struct {
	Id int64 `json:"id"`
}

type ListMovieCommand struct {
}
