package movies

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type HttpEndpointsFactory interface {
	ListMoviesEndpoint() func(w http.ResponseWriter, r *http.Request)
	GetMovieEndpoint(idParam string) func(w http.ResponseWriter, r *http.Request)
	UpdateMovieEndpoint(idParam string) func(w http.ResponseWriter, r *http.Request)
	DeleteMovieEndpoint(idParam string) func(w http.ResponseWriter, r *http.Request)
	CreateMovieEndpoint() func(w http.ResponseWriter, r *http.Request)
}

type httpEndpointsFactory struct {
	movieService MovieService
}

type customError struct {
	Message string `json:"message"`
}

func NewHttpEndpoints(movieService MovieService) HttpEndpointsFactory {
	return &httpEndpointsFactory{movieService: movieService}
}

func (httpFact *httpEndpointsFactory) ListMoviesEndpoint() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd := &ListMovieCommand{}
		data, err := cmd.Exec(httpFact.movieService)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, &customError{err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, data)
	}
}
func (httpFact *httpEndpointsFactory) GetMovieEndpoint(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd := &GetMovieCommand{}
		vars := mux.Vars(r)
		idStr, ok := vars[idParam]
		if !ok {
			respondJSON(w, http.StatusBadGateway, &customError{"id is provided"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 10)
		if err != nil {
			respondJSON(w, http.StatusBadGateway, &customError{"id should be number"})
			return
		}
		cmd.Id = id
		result, err := cmd.Exec(httpFact.movieService)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, &customError{err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, result)
	}
}
func (httpFact *httpEndpointsFactory) UpdateMovieEndpoint(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars[idParam]
		if !ok {
			respondJSON(w, http.StatusBadGateway, &customError{"id is provided"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 10)
		if err != nil {
			respondJSON(w, http.StatusBadGateway, &customError{"id should be number"})
			return
		}
		cmd := &UpdateMovieCommand{}
		if r.Header.Get("Content-Type") == "application/json" {
			err := json.NewDecoder(r.Body).Decode(cmd)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, &customError{err.Error()})
				return
			}
		}
		cmd.Id = id
		result, err := cmd.Exec(httpFact.movieService)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, &customError{err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, result)
	}
}
func (httpFact *httpEndpointsFactory) DeleteMovieEndpoint(idParam string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars[idParam]
		if !ok {
			respondJSON(w, http.StatusBadGateway, &customError{"id is provided"})
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 10)
		if err != nil {
			respondJSON(w, http.StatusBadGateway, &customError{"id should be number"})
			return
		}
		cmd := &DeleteMovieCommand{id}
		result, err := cmd.Exec(httpFact.movieService)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, &customError{err.Error()})
			return
		}
		respondJSON(w, http.StatusOK, result)
	}
}
func (httpFact *httpEndpointsFactory) CreateMovieEndpoint() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd := &CreateMovieCommand{}
		if r.Header.Get("Content-Type") == "application/json" {
			err := json.NewDecoder(r.Body).Decode(cmd)
			if err != nil {
				respondJSON(w, http.StatusInternalServerError, &customError{err.Error()})
				return
			}
		}
		result, err := cmd.Exec(httpFact.movieService)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, &customError{err.Error()})
			return
		}
		respondJSON(w, http.StatusCreated, result)
	}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}
