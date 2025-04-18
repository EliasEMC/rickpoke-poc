package usecase

import (
	"context"

	"github.com/sony/gobreaker"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
)

type CombinedUC struct {
	CharSvc service.CharacterFetcher
	PokeSvc service.PokemonFetcher
}

func (uc CombinedUC) Get(ctx context.Context, charID int, pokeName string) (map[string]interface{}, error) {
	char, err := uc.CharSvc.GetByID(ctx, charID)
	if err != nil {
		return nil, err
	}
	poke, err := uc.PokeSvc.GetByName(ctx, pokeName)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"character": char,
		"pokemon":   poke,
	}, nil
}

// ‑‑ Helpers para /health
func (uc CombinedUC) CharSvcState() string { return breakerState(uc.CharSvc) }
func (uc CombinedUC) PokeSvcState() string { return breakerState(uc.PokeSvc) }

func breakerState(svc interface{}) string {
	if b, ok := svc.(interface{ State() gobreaker.State }); ok {
		return b.State().String()
	}
	return "unknown"
}
