package postgres

import (
    "context"
    "database/sql"
    "testing"
    _ "github.com/lib/pq"
    "github.com/EliasEMC/rickpoke-poc/internal/domain/model"
)

func setupTestDB(t *testing.T) *sql.DB {
    connStr := "host=localhost port=5432 user=postgres password=postgres dbname=rickpoke_test sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }

    // Limpiar la tabla antes de cada test
    _, err = db.Exec("DROP TABLE IF EXISTS characters")
    if err != nil {
        t.Fatalf("Failed to drop table: %v", err)
    }

    // Crear la tabla
    _, err = db.Exec(`
        CREATE TABLE characters (
            id INTEGER PRIMARY KEY,
            name VARCHAR(255) NOT NULL
        )
    `)
    if err != nil {
        t.Fatalf("Failed to create table: %v", err)
    }

    return db
}

func TestCharacterRepo_Save(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := NewCharacterRepo(db)
    ctx := context.Background()

    character := &model.Rick{
        ID:   1,
        Name: "Rick Sanchez",
    }

    // Test Save
    err := repo.Save(ctx, character)
    if err != nil {
        t.Errorf("Failed to save character: %v", err)
    }

    // Test FindByID
    found, err := repo.FindByID(ctx, character.ID)
    if err != nil {
        t.Errorf("Failed to find character: %v", err)
    }
    if found == nil {
        t.Error("Character not found")
    }
    if found.ID != character.ID || found.Name != character.Name {
        t.Errorf("Expected %v, got %v", character, found)
    }
} 