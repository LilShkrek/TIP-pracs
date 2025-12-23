package notes

import (
	"context"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) *mongo.Database {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://root:secret@localhost:27017/?authSource=admin"
	}

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.NewClient(opts)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	db := client.Database("pz8_test")
	t.Cleanup(func() {
		db.Drop(ctx)
		client.Disconnect(ctx)
	})

	return db
}

func TestCreateAndGet(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	repo, err := NewRepo(db)
	if err != nil {
		t.Fatal(err)
	}

	created, err := repo.Create(ctx, "Test Note", "Test content")
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.ByID(ctx, created.ID.Hex())
	if err != nil {
		t.Fatal(err)
	}

	if got.Title != "Test Note" {
		t.Fatalf("want 'Test Note', got '%s'", got.Title)
	}
	if got.Content != "Test content" {
		t.Fatalf("want 'Test content', got '%s'", got.Content)
	}
}
