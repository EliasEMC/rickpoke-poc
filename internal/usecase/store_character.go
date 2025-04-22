package usecase

import (
	"context"
	"fmt"

	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
)

type StoreCharacter struct {
	Repo service.CharacterRepository
}

func (uc StoreCharacter) Save(ctx context.Context, c *model.Rick) error {
	if c.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return uc.Repo.Save(ctx, c)
}
