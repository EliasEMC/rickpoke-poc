package service

import (
	"context"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
)

type CharacterRepository interface {
	Save(ctx context.Context, c *model.Rick) error
	FindByID(ctx context.Context, id int) (*model.Rick, error)
	FindAll(ctx context.Context) ([]*model.Rick, error)
	Ping(ctx context.Context) error
} 