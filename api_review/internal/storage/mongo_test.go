package storage

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/reversersed/LitGO-backend/tree/main/api_review/pkg/mongo"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var cfg *mongo.DatabaseConfig

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		log.Println("\t--- Integration tests are not running in short mode")
		return
	}

	ctx := context.Background()
	var err error
	var mongoContainer testcontainers.Container
	for i := 0; i < 5; i++ {
		req := testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForListeningPort("27017/tcp"),
			Env: map[string]string{
				"MONGO_INITDB_ROOT_USERNAME": "root",
				"MONGO_INITDB_ROOT_PASSWORD": "root",
			},
		}
		mongoContainer, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err == nil {
			break
		}
		log.Printf("failed to create container: %v, retry %d/5", err, i+1)
		<-time.After(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not start mongo: %s", err)
	}
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			log.Fatalf("Could not stop mongo: %s", err)
		}
	}()
	host, err := mongoContainer.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}
	port, err := mongoContainer.MappedPort(ctx, "27017/tcp")
	if err != nil {
		log.Fatal(err)
	}
	cfg = &mongo.DatabaseConfig{
		Host:     host,
		Port:     port.Int(),
		User:     "root",
		Password: "root",
		Base:     "testbase",
		AuthDb:   "admin",
	}
	exit := m.Run()
	os.Exit(exit)
}
