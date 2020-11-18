package movies

import (
	"errors"
	"github.com/gocql/gocql"
	"strings"
	"time"
)

var tableCqls = []string{
	`create table if not exists movies (
		id int,
		name text,
		rating int,
		PRIMARY KEY (id)
	);`,
}

type CassandraConfig struct {
	Nodes       []string
	Database    string
	Consistency string
	Table       string
}

func consistencyLevel(consistency string) gocql.Consistency {
	switch consistency {
	case "quorum":
		return gocql.Quorum
	case "any":
		return gocql.Any
	}
	return gocql.Quorum
}

type defaultMovieStore struct {
	cql     *gocql.Session
	cluster *gocql.ClusterConfig
}

func NewCassandraMovieStore(config CassandraConfig) (MovieStore, error) {
	cluster := gocql.NewCluster(config.Nodes...)
	cluster.Keyspace = config.Database
	cluster.Consistency = consistencyLevel(config.Consistency)
	cluster.Timeout = 60 * time.Second
	store := &defaultMovieStore{cluster: cluster}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	store.cql = session
	err = store.initRepository()
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (d *defaultMovieStore) initRepository() error {
	for _, cqlQuery := range tableCqls {
		err := d.cql.Query(cqlQuery).Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultMovieStore) Create(movie *Movie) (*Movie, error) {
	data, err := d.List()
	if err != nil && err.Error() == "No data" {
		movie.Id = 1
	} else if err != nil {
		return nil, err
	} else {
		movie.Id = int64(len(data) + 1)
	}
	err = d.cql.Query("insert into movies (id,name,rating) values (?,?,?)",
		movie.Id, movie.Name, movie.Rating).Exec()
	if err != nil {
		return nil, err
	}
	return movie, nil
}
func (d *defaultMovieStore) Get(id int64) (*Movie, error) {
	iter := d.cql.Query("select * from movies where id=?", id).
		Consistency(gocql.Quorum).
		Iter()

	if iter.NumRows() <= 0 {
		return nil, errors.New("No data")
	}
	movie := &Movie{}
	for iter.Scan(&movie.Id, &movie.Name, &movie.Rating) {
	}
	return movie, nil
}
func (d *defaultMovieStore) List() ([]Movie, error) {
	iter := d.cql.Query("select * from movies").
		Consistency(gocql.Quorum).
		Iter()
	if iter.NumRows() <= 0 {
		return nil, errors.New("No data")
	}
	movies := []Movie{}
	var id int64
	var name string
	var rating int64
	for iter.Scan(&id, &name, &rating) {
		movies = append(movies, Movie{
			Id:     id,
			Name:   name,
			Rating: rating,
		})
	}
	return movies, nil
}
func (d *defaultMovieStore) Update(movie *Movie) (*Movie, error) {
	q := "update movies set"
	var parts []string
	var values []interface{}
	if movie.Name != "" {
		parts = append(parts, " name = ?")
		values = append(values, movie.Name)
	}
	if movie.Rating != 0 {
		parts = append(parts, " rating = ?")
		values = append(values, movie.Rating)
	}
	if len(parts) <= 0 {
		return nil, errors.New("Nothing to update")
	}
	q = q + strings.Join(parts, ",") + " where id = ?"
	values = append(values, movie.Id)
	err := d.cql.Query(q, values...).Exec()
	if err != nil {
		return nil, err
	}
	movieUpdated, err := d.Get(movie.Id)
	if err != nil {
		return nil, err
	}
	return movieUpdated, nil
}
func (d *defaultMovieStore) Delete(id int64) error {
	err := d.cql.Query("DELETE FROM movies WHERE id = ?", id).Exec()
	if err != nil {
		return err
	}
	return nil
}
