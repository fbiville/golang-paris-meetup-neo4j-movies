package movies

import (
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
	session := repo.driver.NewSession(neo4j.SessionConfig{
		BoltLogger: neo4j.ConsoleBoltLogger(),
	})
	defer session.Close()
	movies, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		results, err := tx.Run("MATCH (:Person {name: $name})-[:DIRECTED]->(m:Movie) RETURN m ORDER BY m.released ASC",
			map[string]interface{}{
				"name": director,
			})
		if err != nil {
			return nil, err
		}
		records, err := results.Collect()
		if err != nil {
			return nil, err
		}
		movies := make([]Movie, len(records))
		for i, record := range records {
			movie, _ := record.Get("m")
			movieNode := movie.(neo4j.Node)
			movies[i] = Movie{
				Title:       movieNode.Props["title"].(string),
				ReleaseYear: movieNode.Props["released"].(int64),
				Tagline:     movieNode.Props["tagline"].(string),
			}
		}
		return movies, nil
	})
	if err != nil {
		return nil, err
	}
	return movies.([]Movie), nil
}
