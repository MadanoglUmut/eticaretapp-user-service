package repositories

import (
	"UserService/pkg/psql"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var db *gorm.DB

type TestDB struct {
	Container testcontainers.Container
}

func (t *TestDB) Setup() error {
	ctx := context.Background()

	dbConfig := map[string]string{
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_DB":       "users",
	}

	defaultPort := nat.Port("5432/tcp")
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{defaultPort.Port()},
			Env:          dbConfig,
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort(defaultPort),
			),
		},
		Started: true,
	})
	if err != nil {
		return err
	}
	t.Container = container
	port, err := container.MappedPort(ctx, defaultPort)
	if err != nil {
		return err
	}
	fmt.Println("Veritabanı başladı port numarasi:", port)
	db = psql.Connect("0.0.0.0", dbConfig["POSTGRES_USER"], dbConfig["POSTGRES_PASSWORD"], dbConfig["POSTGRES_DB"], port.Port())
	return t.loadSQLFiles()
}

func (t *TestDB) loadSQLFiles() error {
	fileCreate, err := os.ReadFile("../../psql/create_tables.sql")
	if err != nil {
		return err
	}
	if err := db.Exec(string(fileCreate)).Error; err != nil {
		return err
	}

	fileFill, err := os.ReadFile("../../psql/fill_tables.sql")
	if err != nil {
		return err
	}
	return db.Exec(string(fileFill)).Error
}

func (t *TestDB) CleanUp() {
	t.Container.Terminate(context.Background())
}

func TestMain(m *testing.M) {
	testDB := &TestDB{}
	if err := testDB.Setup(); err != nil {
		fmt.Println("Veritabanı bağlantısı başarısız", err)
		os.Exit(1)
	}
	defer testDB.CleanUp()

	os.Exit(m.Run())
}
