package httpinfra

import (
	"time"
	"strconv"
	"github.com/EliasEMC/rickpoke-poc/internal/adapter/handler"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func BuildRouter(
	logger *zap.Logger,
	combined usecase.CombinedUC,
	charUC usecase.FetchCharacter,
	pokeUC usecase.FetchPokemon,
) *gin.Engine {

	r := gin.New()
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// Agrupamos rutas
	api := r.Group("/")

	// Endpoints individuales
	handler.RegisterCharacterRoutes(api, charUC)
	handler.RegisterPokemonRoutes(api, pokeUC)

	// Endpoint combinado
	r.GET("/combined", func(c *gin.Context) {
        id, _ := strconv.Atoi(c.Query("char_id"))
        name := c.Query("pokemon")
        resp, err := combined.Get(c.Request.Context(), id, name)
        if err != nil {
            c.JSON(502, gin.H{"error": err.Error()})
            return
        }
        c.JSON(200, resp)
    })

	// Health
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"rick_cb": combined.CharSvcState(),
			"poke_cb": combined.PokeSvcState(),
		})
	})

	return r
}
