package main

import (
	_ "github.com/joho/godotenv/autoload" // carga .env
	"go.uber.org/zap"

	//"time"

	"github.com/EliasEMC/rickpoke-poc/internal/adapter/api"
	"github.com/EliasEMC/rickpoke-poc/internal/config"
	"github.com/EliasEMC/rickpoke-poc/internal/infrastructure/circuitbreaker"
	httpinfra "github.com/EliasEMC/rickpoke-poc/internal/infrastructure/http"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
	"github.com/EliasEMC/rickpoke-poc/pkg/utils"
)

func main() {
	cfg := config.Load()

	// Logger
	logger := utils.NewLogger(cfg.LogLevel)
	defer logger.Sync()

	// Infra comunes
	httpClient := utils.NewHTTPClient(cfg.Timeout)

	// Breakers + adapters
	rickCB := circuitbreaker.New("rick")
	pokeCB := circuitbreaker.New("poke")

	rickClient := api.NewRickClient(cfg.RickURL, rickCB, cfg.Timeout)
	rickClient.Client = httpClient

	pokeClient := api.NewPokeClient(cfg.PokeURL, pokeCB, cfg.Timeout)
	pokeClient.Client = httpClient

	// Use‑cases
	combinedUC := usecase.CombinedUC{CharSvc: rickClient, PokeSvc: pokeClient}
	charUC := usecase.FetchCharacter{Svc: rickClient}
	pokeUC := usecase.FetchPokemon{Svc: pokeClient}

	// Router
	router := httpinfra.BuildRouter(logger, combinedUC, charUC, pokeUC)

	logger.Info("server up", zap.String("port", cfg.Port))
	logger.Fatal("server exited", zap.Error(router.Run(":"+cfg.Port)))
}
