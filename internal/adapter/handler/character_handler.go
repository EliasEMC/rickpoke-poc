package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
)

func RegisterCharacterRoutes(r *gin.RouterGroup, uc usecase.FetchCharacter) {
	r.GET("/rick/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be int"})
			return
		}
		res, err := uc.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})
}
