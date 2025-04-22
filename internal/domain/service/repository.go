package service

import (
	"context"

	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
)

// CharacterRepository define las operaciones que el dominio necesita.
type CharacterRepository interface {
	Save(ctx context.Context, c *model.Rick) error
	FindByID(ctx context.Context, id int) (*model.Rick, error)
}
