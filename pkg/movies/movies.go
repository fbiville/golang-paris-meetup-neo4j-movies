package movies

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Movie struct {
	Title       string
	ReleaseYear int64
	Tagline     string
}

type MovieRepository interface {
	FindMoviesDirectedBy(director string) ([]Movie, error)
}

type movieNeo4jRepository struct {
	driver neo4j.Driver
}

func NewMovieRepository(driver neo4j.Driver) MovieRepository {
	return &movieNeo4jRepository{driver: driver}
}

func (repo *movieNeo4jRepository) FindMoviesDirectedBy(director string) ([]Movie, error) {
	// TODO
	return nil, fmt.Errorf("TODO")
}
