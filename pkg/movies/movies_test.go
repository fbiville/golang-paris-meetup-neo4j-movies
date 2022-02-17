package movies_test

import (
	"context"
	"github.com/fbiville/golang-paris-meetup-neo4j-movies/pkg/container"
	"github.com/fbiville/golang-paris-meetup-neo4j-movies/pkg/errors"
	"github.com/fbiville/golang-paris-meetup-neo4j-movies/pkg/movies"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"reflect"
	"testing"
)

func TestMovies(outer *testing.T) {
	ctx := context.Background()
	neo4jContainer, driver, err := container.StartNeo4jContainer(ctx)
	defer func() {
		errors.MaybePanic(neo4jContainer.Terminate(ctx))
	}()
	defer errors.MaybePanicOnClose(driver)
	errors.MaybePanic(err)
	session := driver.NewSession(neo4j.SessionConfig{})
	defer errors.MaybePanicOnClose(session)
	tx, err := session.BeginTransaction()
	errors.MaybePanic(err)
	defer errors.MaybePanicOnClose(tx)
	_, err = tx.Run(`
	CREATE (lana:Person {name: "Lana Wachowski", born: 1965})
	CREATE (lilly:Person {name: "Lilly Wachowski", born: 1967})
	CREATE (theMatrix:Movie {title: "The Matrix", tagline: "Welcome to the Real World", released: 1999})
	CREATE (theMatrixReloaded:Movie {title: "The Matrix Reloaded", tagline: "Free your mind", released: 2003})
	CREATE (lana)-[:DIRECTED]->(theMatrix)
	CREATE (lana)-[:DIRECTED]->(theMatrixReloaded)
	CREATE (lilly)-[:DIRECTED]->(theMatrix)
	CREATE (lilly)-[:DIRECTED]->(theMatrixReloaded)
`, nil)
	errors.MaybePanic(err)
	errors.MaybePanic(tx.Commit())

	repository := movies.NewMovieRepository(driver)

	outer.Run("Finds movies directed by someone", func(t *testing.T) {
		results, err := repository.FindMoviesDirectedBy("Lana Wachowski")

		if err != nil {
			t.Errorf("Expected non-nil error, got: %v", err)
		}
		if len(results) != 2 {
			t.Errorf("Expected 2 movies, got: %d", len(results))
		}
		firstMovie := movies.Movie{Title: "The Matrix", ReleaseYear: 1999,
			Tagline: "Welcome to the Real World",
		}
		if !reflect.DeepEqual(results[0], firstMovie) {
			t.Errorf("Expected %s as first movie, got: %s", firstMovie.Title, results[0].Title)
		}
		secondMovie := movies.Movie{Title: "The Matrix Reloaded",
			ReleaseYear: 2003,
			Tagline:     "Free your mind",
		}
		if !reflect.DeepEqual(results[1], secondMovie) {
			t.Errorf("Expected %s as second movie, got: %s", firstMovie.Title,
				results[0].Title)
		}
	})
}
