package actors

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Actor struct {
	Name      string
	BirthYear int64
}

type BaconSlice struct {
	Type string
	Name string
}

type ActorRepository interface {
	// FindOneByName returns the first actor match
	FindOneByName(name string) (*Actor, error)
	// FindShortestPathToKevinBacon returns the shortest path between an
	// actor and Kevin Bacon
	// 🇫🇷 is bacon: https://knowyourmeme.com/memes/france-is-bacon
	FindShortestPathToKevinBacon(name string) ([]BaconSlice, error)
}

type actorNeo4jRepository struct {
	driver neo4j.Driver
}

func NewActorRepository(driver neo4j.Driver) ActorRepository {
	return &actorNeo4jRepository{driver: driver}
}

func (repo *actorNeo4jRepository) FindOneByName(name string) (*Actor, error) {
	// TODO
	return nil, fmt.Errorf("TODO")
}

func (repo *actorNeo4jRepository) FindShortestPathToKevinBacon(actor string) (
	[]BaconSlice, error) {
	// TODO
	return nil, fmt.Errorf("TODO")
}
