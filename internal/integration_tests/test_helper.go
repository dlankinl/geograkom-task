package integration_tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

const (
	DbName = "postgres"
	DbUser = "postgres"
	DbPass = "postgres"
)

type TestDatabase struct {
	DbInstance *pgx.Conn
	DbAddress  string
	container  testcontainers.Container
}

func SetupTestDatabase() *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test: ", err)
	}

	err = migrateDb(dbAddr)
	if err != nil {
		log.Fatal("failed to perform db migration: ", err)
	}
	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	ctx := context.Background()
	tdb.DbInstance.Close(ctx)
	_ = tdb.container.Terminate(ctx)
}

func createContainer(ctx context.Context) (testcontainers.Container, *pgx.Conn, string, error) {
	var env = map[string]string{
		"POSTGRES_PASSWORD": DbPass,
		"POSTGRES_USER":     DbUser,
		"POSTGRES_DB":       DbName,
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)
	db, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return container, db, dbAddr, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, dbAddr, nil
}

func migrateDb(dbAddr string) error {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get path")
	}
	_ = path
	pathToMigrationFiles := filepath.Join("..", "..", "migrations")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)
	m, err := migrate.New(fmt.Sprintf("file:%s", pathToMigrationFiles), databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func SeedTestData(db *pgx.Conn) error {
	filePath := filepath.Join("..", "integration_tests", "test_data.sql")
	f, err := os.Open(filePath)
	if err != nil {
		out, iErr := exec.Command("ls", "-la").Output()
		if iErr != nil {
			fmt.Println(iErr)
		}
		fmt.Println(string(out))
		return fmt.Errorf("opening file with test data: %w", err)
	}

	err = executeTestDataScript(db, filePath)
	if err != nil {
		return fmt.Errorf("%s: %w", f.Name(), err)
	}

	return nil
}

func executeTestDataScript(db *pgx.Conn, filePath string) error {
	scriptContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("sql script reading: %w", err)
	}

	_, err = db.Exec(context.Background(), string(scriptContent))
	if err != nil {
		return fmt.Errorf("sql script execution: %w", err)
	}

	return nil
}
