package usecase

import (
	"context"

	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
)

type ListCharacters struct {
	Repo service.CharacterRepository
}

func (uc ListCharacters) GetAll(ctx context.Context) ([]*model.Rick, error) {
	return uc.Repo.FindAll(ctx)
} 