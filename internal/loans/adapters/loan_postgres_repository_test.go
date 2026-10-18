package adapters_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/wreckitral/pinjamean/internal/loans/adapters"
)

func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set: repository tests require a dedicated test database (e.g. postgres://postgres:postgres@localhost:5432/pinjamean_test?sslmode=disable)")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)

	err = pool.Ping(ctx)
	require.NoError(t, err, "failed to ping test database")

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestPostgresLoanRepository(t *testing.T) {
	t.Parallel()

	pool := testDB(t)
	repo := adapters.NewPostgresLoanRepository(pool)
	testLoanRepositoryContract(t, repo)
}


func TestMigrations_ApplyUpDownUp(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Fatal("TEST_DATABASE_URL is not set: repository tests require a dedicated test database (e.g. postgres://postgres:postgres@localhost:5432/pinjamean_test?sslmode=disable)")
	}

	_, currentFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	migrationsDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "migrations")

	m, err := migrate.New("file://"+filepath.ToSlash(migrationsDir), dbURL)
	require.NoError(t, err)
	defer func() {
		sourceErr, dbErr := m.Close()
		require.NoError(t, sourceErr)
		require.NoError(t, dbErr)
	}()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		require.NoError(t, err)
	}

	// 2. Down: roll back all migrations
	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		require.NoError(t, err)
	}
	version, dirty, err := m.Version()
	require.ErrorIs(t, err, migrate.ErrNilVersion)
	require.False(t, dirty)
	require.Equal(t, uint(0), version)

	// 3. Up: reapply all migrations
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		require.NoError(t, err)
	}
	version, dirty, err = m.Version()
	require.NoError(t, err)
	require.False(t, dirty)
	require.Equal(t, uint(6), version)
}

