package usecase

import (
	"context"

	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
)

type FetchCharacter struct {
	Svc service.CharacterFetcher
}

func (uc FetchCharacter) Get(ctx context.Context, id int) (interface{}, error) {
	return uc.Svc.GetByID(ctx, id)
}
