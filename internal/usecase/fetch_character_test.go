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

func TestFetchCharacter_Get(t *testing.T) {
	charMock := &mocks.CharacterFetcher{}
	charMock.
		On("GetByID", mock.Anything, 7).
		Return(&model.Rick{ID: 7, Name: "Abradolf Lincler"}, nil).
		Once()

	uc := usecase.FetchCharacter{Svc: charMock}

	got, err := uc.Get(context.Background(), 7)
	require.NoError(t, err)

	char := got.(*model.Rick)
	require.Equal(t, 7, char.ID)
	require.Equal(t, "Abradolf Lincler", char.Name)

	charMock.AssertExpectations(t)
}
