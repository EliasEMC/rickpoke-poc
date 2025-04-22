package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/EliasEMC/rickpoke-poc/internal/adapter/handler"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	mocks "github.com/EliasEMC/rickpoke-poc/internal/mocks/internal_/domain/service"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
)

func TestPokemonHandler_Get(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Create mock service
	pokeMock := &mocks.PokemonFetcher{}
	pokeMock.On("GetByName", mock.Anything, "pikachu").
		Return(&model.Pokemon{Name: "pikachu", Type: "electric"}, nil)
	pokeMock.On("GetByName", mock.Anything, "nonexistent").
		Return(nil, assert.AnError)
	
	// Create use case with mock
	uc := usecase.FetchPokemon{Svc: pokeMock}
	
	// Create handler
	pokemonHandler := handler.NewPokemonHandler(uc)
	
	// Register routes
	handler.RegisterPokemonRoutes(router, pokemonHandler)

	// Test cases
	tests := []struct {
		name       string
		pokemon    string
		wantStatus int
	}{
		{
			name:       "valid pokemon",
			pokemon:    "pikachu",
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid pokemon",
			pokemon:    "nonexistent",
			wantStatus: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/pokemon/"+tt.pokemon, nil)
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
	
	// Verify mock expectations
	pokeMock.AssertExpectations(t)
}
