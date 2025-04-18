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

func TestCharacterRoute_OK(t *testing.T) {
	charMock := &mocks.CharacterFetcher{}
	charMock.
		On("GetByID", mock.Anything, 1).
		Return(&model.Rick{ID: 1, Name: "Rick"}, nil).Once()

	uc := usecase.FetchCharacter{Svc: charMock}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("")
	handler.RegisterCharacterRoutes(group, uc)

	for _, ri := range r.Routes() {
		t.Logf("route registered: %s %s", ri.Method, ri.Path)
	}

	req, _ := http.NewRequest("GET", "/rick/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "Rick")
	charMock.AssertExpectations(t)
}
