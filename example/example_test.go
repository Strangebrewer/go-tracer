package example_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Strangebrewer/go-service-template/example"
)

var testStore *example.Store

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	schema, err := os.ReadFile("../db/schema.sql")
	if err != nil {
		log.Fatalf("failed to read schema: %v", err)
	}
	if _, err := pool.Exec(ctx, string(schema)); err != nil {
		log.Fatalf("failed to apply schema: %v", err)
	}

	testStore = example.NewStore(pool)

	os.Exit(m.Run())
}

func TestExampleStore_Create(t *testing.T) {
	t.Skip("implement store methods before enabling")

	ctx := context.Background()

	req := example.CreateExampleRequest{Name: "test example"}
	result, err := testStore.Create(ctx, "test-user-id", req)

	require.NoError(t, err)
	assert.Equal(t, "test example", result.Name)
	assert.NotEmpty(t, result.ID)
}

func TestExampleStore_GetAll(t *testing.T) {
	t.Skip("implement store methods before enabling")

	ctx := context.Background()

	results, err := testStore.GetAll(ctx, "test-user-id")

	require.NoError(t, err)
	assert.NotNil(t, results)
}
