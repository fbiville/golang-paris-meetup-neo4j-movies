package container

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const username = "neo4j"

const password = "s3cr3t"

func StartNeo4jContainer(ctx context.Context) (testcontainers.Container, neo4j.Driver, error) {
	request := testcontainers.ContainerRequest{
		Image:        "neo4j:4.4",
		ExposedPorts: []string{"7687/tcp"},
		Env: map[string]string{"NEO4J_AUTH": fmt.Sprintf("%s/%s",
			username, password)},
		WaitingFor: wait.ForLog("Bolt enabled"),
	}
	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: request,
			Started:          true,
		})
	if err != nil {
		return nil, nil, err
	}
	driver, err := newNeo4jDriver(ctx, container)
	if err != nil {
		return container, nil, err
	}
	return container, driver, driver.VerifyConnectivity()
}

func newNeo4jDriver(ctx context.Context, container testcontainers.Container) (
	neo4j.Driver, error) {
	port, err := container.MappedPort(ctx, "7687")
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("neo4j://localhost:%d", port.Int())
	return neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
}
