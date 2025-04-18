package usecase

import (
	"context"

	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
)

type FetchPokemon struct {
	Svc service.PokemonFetcher
}

func (uc FetchPokemon) Get(ctx context.Context, name string) (interface{}, error) {
	return uc.Svc.GetByName(ctx, name)
}
