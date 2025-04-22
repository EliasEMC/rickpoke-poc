package postgres

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
)

type CharacterRepo struct {
	db *sql.DB
}

func NewCharacterRepo(db *sql.DB) *CharacterRepo {
	return &CharacterRepo{db: db}
}

func (r *CharacterRepo) Save(ctx context.Context, c *model.Rick) error {
	query := `
		INSERT INTO characters (id, name)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET name = $2
	`
	_, err := r.db.ExecContext(ctx, query, c.ID, c.Name)
	return err
}

func (r *CharacterRepo) FindByID(ctx context.Context, id int) (*model.Rick, error) {
	query := `SELECT id, name FROM characters WHERE id = $1`
	var character model.Rick
	err := r.db.QueryRowContext(ctx, query, id).Scan(&character.ID, &character.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &character, nil
}

func (r *CharacterRepo) FindAll(ctx context.Context) ([]*model.Rick, error) {
	query := `SELECT id, name FROM characters ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var characters []*model.Rick
	for rows.Next() {
		var c model.Rick
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		characters = append(characters, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return characters, nil
}

func (r *CharacterRepo) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

var _ service.CharacterRepository = (*CharacterRepo)(nil)
