package models

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCreateItem(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16"),
		postgres.WithInitScripts(filepath.Join("..", "migrations", "1_CREATE_TABLES.sql")),
		postgres.WithDatabase("pos-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to create DB connection string")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect to DB")
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("failed to ping DB")
	}

	// testItem := Item{
	// 	Name: "Test Item",
	// 	CategoryID: 1,
	// 	Description: "Item used for integration testing.",
	// 	Price: 4.98,
	// 	Archived: false,
	// 	CreatedDate: time.Now().String(),
	// }

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(handler)
	err = CreateItem(logger, db, 1, "Test Item", "Item for integration testing.", 4.98)
	if err != nil {
		t.Fatalf("failed to create item in DB")
	}

	item, err := GetItemByID(logger, db, 1)
	if err != nil {
		t.Fatalf("Failed to retrieve item")
	}

	if item.Name != "Test Item" {
		t.Fatalf("Item name did not match expected")
	}

	// TODO
	// Create item
	// Check item equal
	// Update item
	// Validate update
	// Delete item
	// Check gone
}
