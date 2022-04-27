package actors

import (
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
	session := repo.driver.NewSession(neo4j.SessionConfig{
		BoltLogger: neo4j.ConsoleBoltLogger(),
	})
	defer session.Close()
	actor, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		results, err := tx.Run("MATCH (p:Person {name: $name}) RETURN p LIMIT 1",
			map[string]interface{}{
				"name": name,
			})
		if err != nil {
			return nil, err
		}
		record, err := results.Single()
		if err != nil {
			return nil, err
		}
		person, _ := record.Get("p")
		personNode := person.(neo4j.Node)
		return &Actor{
			Name:      personNode.Props["name"].(string),
			BirthYear: personNode.Props["born"].(int64),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return actor.(*Actor), nil
}

func (repo *actorNeo4jRepository) FindShortestPathToKevinBacon(actor string) (
	[]BaconSlice, error) {
	session := repo.driver.NewSession(neo4j.SessionConfig{
		BoltLogger: neo4j.ConsoleBoltLogger(),
	})
	defer session.Close()
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		results, err := tx.Run("MATCH (kevin:Person {name: 'Kevin Bacon'})\nMATCH (target:Person {name: $name})\nRETURN shortestPath((kevin)-[*0..6]-(target)) AS path",
			map[string]interface{}{
				"name": actor,
			})
		record, err := results.Single()
		if err != nil {
			return nil, err
		}
		path, _ := record.Get("path")
		baconPath := path.(neo4j.Path)
		var pathElements []BaconSlice
		pathNodes := baconPath.Nodes
		for _, node := range pathNodes {
			name, found := node.Props["name"]
			if !found {
				name = node.Props["title"]
			}
			pathElements = append(pathElements, BaconSlice{
				Type: node.Labels[0],
				Name: name.(string),
			})
		}
		return pathElements, nil
	})
	if err != nil {
		return nil, err
	}
	return results.([]BaconSlice), nil
}
