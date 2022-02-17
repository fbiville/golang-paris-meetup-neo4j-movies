package main

import (
	"encoding/json"
	"fmt"
	"github.com/fbiville/golang-paris-meetup-neo4j-movies/pkg/actors"
	"github.com/fbiville/golang-paris-meetup-neo4j-movies/pkg/errors"
	"github.com/fbiville/golang-paris-meetup-neo4j-movies/pkg/movies"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io/ioutil"
)

type config struct {
	Uri      string `json:"uri"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *config) createDriver() (neo4j.Driver, error) {
	// TODO
	return nil, fmt.Errorf("TODO")
}

func main() {
	config, err := parseConfig("config.json")
	errors.MaybePanic(err)
	driver, err := config.createDriver()
	errors.MaybePanic(err)
	defer errors.MaybePanicOnClose(driver)

	movieRepo := movies.NewMovieRepository(driver)
	directedMovies, err := movieRepo.FindMoviesDirectedBy("Ron Howard")
	errors.MaybePanic(err)
	fmt.Println("Ron Howard directed:")
	for _, movie := range directedMovies {
		fmt.Printf("\t%s in %d\n", movie.Title, movie.ReleaseYear)
	}

	actorRepo := actors.NewActorRepository(driver)
	actress, err := actorRepo.FindOneByName("Charlize Theron")
	errors.MaybePanic(err)
	fmt.Printf("Charlize Theron was born in %d\n", actress.BirthYear)
	baconPath, err := actorRepo.FindShortestPathToKevinBacon("Robert Zemeckis")
	errors.MaybePanic(err)
	fmt.Println("The Bacon path is made of")
	for _, baconSlice := range baconPath {
		fmt.Printf("\t[%s] %s\n", baconSlice.Type, baconSlice.Name)
	}
}

func parseConfig(path string) (*config, error) {
	rawConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := &config{}
	if err = json.Unmarshal(rawConfig, &config); err != nil {
		return nil, err
	}
	return config, nil
}
