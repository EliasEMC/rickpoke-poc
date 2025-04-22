package httpinfra

import (
	"strconv"
	"net/http"
	"github.com/EliasEMC/rickpoke-poc/internal/adapter/handler"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/EliasEMC/rickpoke-poc/internal/middleware"
	"github.com/sony/gobreaker"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
)

func BuildRouter(
	logger *zap.Logger,
	combinedUC usecase.CombinedUC,
	charUC usecase.FetchCharacter,
	pokeUC usecase.FetchPokemon,
	storeUC usecase.StoreCharacter,
	listUC usecase.ListCharacters,
	rickCB *gobreaker.CircuitBreaker,
	pokeCB *gobreaker.CircuitBreaker,
	repo service.CharacterRepository,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))

	// Handlers
	characterHandler := handler.NewCharacterHandler(storeUC, listUC)
	pokemonHandler := handler.NewPokemonHandler(pokeUC)
	healthHandler := handler.NewHealthHandler(rickCB, pokeCB, repo)

	// Routes
	router.GET("/character/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be int"})
			return
		}
		res, err := charUC.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})
	router.GET("/pokemon/:name", pokemonHandler.Get)
	router.GET("/combined", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Query("char_id"))
		name := c.Query("pokemon")
		resp, err := combinedUC.Get(c.Request.Context(), id, name)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})
	router.GET("/health", healthHandler.Check)

	// Nuevas rutas para el repositorio
	router.POST("/characters", characterHandler.Save)
	router.GET("/characters", characterHandler.List)

	return router
}
