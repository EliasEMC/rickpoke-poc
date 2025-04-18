package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	mocks "github.com/EliasEMC/rickpoke-poc/internal/mocks/internal_/domain/service"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
)

func TestFetchPokemon_Get(t *testing.T) {
	pokeMock := &mocks.PokemonFetcher{}
	pokeMock.
		On("GetByName", mock.Anything, "bulbasaur").
		Return(&model.Pokemon{Name: "Bulbasaur", Type: "grass"}, nil).
		Once()

	uc := usecase.FetchPokemon{Svc: pokeMock}

	got, err := uc.Get(context.Background(), "bulbasaur")
	require.NoError(t, err)

	poke := got.(*model.Pokemon)
	require.Equal(t, "Bulbasaur", poke.Name)
	require.Equal(t, "grass", poke.Type)

	pokeMock.AssertExpectations(t)
}
