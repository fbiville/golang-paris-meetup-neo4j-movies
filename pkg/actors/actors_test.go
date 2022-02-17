package actors_test

import (
	"context"
	"github.com/fbiville/golang-paris-meetup-neo4j-movies/pkg/actors"
	"github.com/fbiville/golang-paris-meetup-neo4j-movies/pkg/container"
	"github.com/fbiville/golang-paris-meetup-neo4j-movies/pkg/errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"reflect"
	"testing"
)

func TestActors(outer *testing.T) {
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
	CREATE (kevin:Person {name: "Kevin Bacon", born: 1958})
	CREATE (meg:Person {name: "Meg Ryan", born: 1961})
	CREATE (tom:Person {name: "Tom Hanks", born: 1956})
	CREATE (apollo13:Movie {title: "Apollo 13", tagline: "Houston, we have a problem.", released: 1995})
	CREATE (sleepless:Movie {title: "Sleepless in Seattle", tagline: "[...]", released: 1993})
	CREATE (meg)-[:ACTED_IN]->(sleepless)
	CREATE (tom)-[:ACTED_IN]->(sleepless)
	CREATE (tom)-[:ACTED_IN]->(apollo13)
	CREATE (kevin)-[:ACTED_IN]->(apollo13)
`, nil)
	errors.MaybePanic(err)
	errors.MaybePanic(tx.Commit())

	actorRepository := actors.NewActorRepository(driver)

	outer.Run("Finds by name", func(t *testing.T) {
		actor, err := actorRepository.FindOneByName("Tom Hanks")

		if err != nil {
			t.Errorf("Expected nil error, got: %v", err)
		}
		if actor == nil {
			t.Errorf("Expected non-nil actor")
		}
		if actor.Name != "Tom Hanks" {
			t.Errorf("Expected name Tom Hanks, got: %s", actor.Name)
		}
		if actor.BirthYear != 1956 {
			t.Errorf("Expected birth year 1956, got: %d", actor.BirthYear)
		}
	})

	outer.Run("Finds Kevin Bacon path", func(t *testing.T) {
		baconSlices, err := actorRepository.FindShortestPathToKevinBacon("Meg Ryan")

		if err != nil {
			t.Errorf("Expected nil error, got: %v", err)
		}
		if len(baconSlices) == 0 {
			t.Errorf("Expected non-empty Bacon path")
		}
		expectedPath := []actors.BaconSlice{
			{
				Type: "Person",
				Name: "Kevin Bacon",
			},
			{
				Type: "Movie",
				Name: "Apollo 13",
			},
			{
				Type: "Person",
				Name: "Tom Hanks",
			},
			{
				Type: "Movie",
				Name: "Sleepless in Seattle",
			},
			{
				Type: "Person",
				Name: "Meg Ryan",
			},
		}
		if !reflect.DeepEqual(baconSlices, expectedPath) {
			t.Errorf("Expected Bacon path %v, but got: %v",
				expectedPath, baconSlices)
		}
	})
}
