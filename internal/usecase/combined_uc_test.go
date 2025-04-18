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

func TestCombinedUC_Success(t *testing.T) {
	// 1. Crear mocks
	charMock := &mocks.CharacterFetcher{}
	pokeMock := &mocks.PokemonFetcher{}

	// 2. Configurar expectativas
	charMock.
		On("GetByID", mock.Anything, 1).
		Return(&model.Rick{ID: 1, Name: "Rick"}, nil).Once()

	pokeMock.
		On("GetByName", mock.Anything, "pikachu").
		Return(&model.Pokemon{Name: "Pikachu", Type: "electric"}, nil).Once()

	// 3. Inyectar mocks en el caso de uso
	uc := usecase.CombinedUC{CharSvc: charMock, PokeSvc: pokeMock}

	out, err := uc.Get(context.Background(), 1, "pikachu")
	require.NoError(t, err)
	require.Equal(t, "Rick", out["character"].(*model.Rick).Name)
	require.Equal(t, "electric", out["pokemon"].(*model.Pokemon).Type)

	// 4. Verificar que se llamaron los mocks
	charMock.AssertExpectations(t)
	pokeMock.AssertExpectations(t)
}
