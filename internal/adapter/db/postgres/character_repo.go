package postgres

import (
	"context"
	"database/sql"

	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
	"github.com/jmoiron/sqlx"
)

type CharacterRepo struct{ DB *sqlx.DB }

func (r CharacterRepo) Save(ctx context.Context, c *model.Rick) error {
	_, err := r.DB.NamedExecContext(ctx,
		`INSERT INTO characters (id, name, origin) VALUES (:id, :name, :origin)`,
		c)
	return err
}

func (r CharacterRepo) FindByID(ctx context.Context, id int) (*model.Rick, error) {
	var out model.Rick
	err := r.DB.GetContext(ctx, &out, `SELECT id, name, origin FROM characters WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &out, err
}

var _ service.CharacterRepository = (*CharacterRepo)(nil)
