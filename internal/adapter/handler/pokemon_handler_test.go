package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/EliasEMC/rickpoke-poc/internal/adapter/handler"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	mocks "github.com/EliasEMC/rickpoke-poc/internal/mocks/internal_/domain/service"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
)

func TestPokemonRoute_OK(t *testing.T) {
	// 1. Mock del puerto
	pokeMock := &mocks.PokemonFetcher{}
	pokeMock.
		On("GetByName", mock.Anything, "pikachu").
		Return(&model.Pokemon{Name: "Pikachu", Type: "electric"}, nil).
		Once()

	uc := usecase.FetchPokemon{Svc: pokeMock}

	// 2. Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("")
	handler.RegisterPokemonRoutes(group, uc)

	// 3. Request
	req, _ := http.NewRequest("GET", "/pokemon/pikachu", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "Pikachu")
	pokeMock.AssertExpectations(t)
}
