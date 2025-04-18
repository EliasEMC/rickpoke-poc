package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
)

func RegisterPokemonRoutes(r *gin.RouterGroup, uc usecase.FetchPokemon) {
	r.GET("/pokemon/:name", func(c *gin.Context) {
		name := strings.ToLower(c.Param("name"))
		res, err := uc.Get(c.Request.Context(), name)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
