package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
)

type PokemonHandler struct {
	uc usecase.FetchPokemon
}

func NewPokemonHandler(uc usecase.FetchPokemon) *PokemonHandler {
	return &PokemonHandler{uc: uc}
}

func (h *PokemonHandler) Get(c *gin.Context) {
	name := c.Param("name")
	res, err := h.uc.Get(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
