package service

import (
	"context"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
)

type CharacterFetcher interface {
	GetByID(ctx context.Context, id int) (*model.Rick, error)
}

type PokemonFetcher interface {
	GetByName(ctx context.Context, name string) (*model.Pokemon, error)
}
