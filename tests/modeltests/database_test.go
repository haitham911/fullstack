package modeltests

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/task/api/controllers"
	"github.com/task/api/models"
)

var (
	user     = "postgres"
	password = "secret"
	db       = "test"
	port     = "5433"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
)
var server = controllers.Server{}
var userInstance = models.User{}
var productInstance = models.Product{}

func TestMain(m *testing.M) {

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.3",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + db,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err.Error())
	}
	time.Sleep(30 * time.Second)
	dsn = fmt.Sprintf(dsn, user, password, port, db)
	if err = pool.Retry(func() error {
		Database(dsn)
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)

}
