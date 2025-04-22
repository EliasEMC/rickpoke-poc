package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
	"github.com/sony/gobreaker"
)

type HealthHandler struct {
	rickCB *gobreaker.CircuitBreaker
	pokeCB *gobreaker.CircuitBreaker
	repo   service.CharacterRepository
}

func NewHealthHandler(rickCB, pokeCB *gobreaker.CircuitBreaker, repo service.CharacterRepository) *HealthHandler {
	return &HealthHandler{
		rickCB: rickCB,
		pokeCB: pokeCB,
		repo:   repo,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	// Verificar estado de la base de datos
	dbStatus := "ok"
	if err := h.repo.Ping(c.Request.Context()); err != nil {
		dbStatus = "error"
	} else {
		dbStatus = "ok"
	}

	// Verificar estado de las APIs
	rickStatus := "ok"
	if h.rickCB.State() == gobreaker.StateOpen {
		rickStatus = "open"
	} else {
		rickStatus = "closed"
	}

	pokeStatus := "ok"
	if h.pokeCB.State() == gobreaker.StateOpen {
		pokeStatus = "open"
	} else {
		pokeStatus = "closed"
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"services": gin.H{
			"database": dbStatus,
			"rick_api": rickStatus,
			"poke_api": pokeStatus,
		},
	})
} 