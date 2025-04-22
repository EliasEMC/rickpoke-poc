package main

import (
	_ "github.com/joho/godotenv/autoload" // carga .env
	"go.uber.org/zap"

	"github.com/EliasEMC/rickpoke-poc/internal/adapter/api"
	"github.com/EliasEMC/rickpoke-poc/internal/config"
	"github.com/EliasEMC/rickpoke-poc/internal/infrastructure/circuitbreaker"
	httpinfra "github.com/EliasEMC/rickpoke-poc/internal/infrastructure/http"
	"github.com/EliasEMC/rickpoke-poc/internal/usecase"
	"github.com/EliasEMC/rickpoke-poc/pkg/utils"
	"context"
	"github.com/EliasEMC/rickpoke-poc/internal/adapter/db/factory"
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

	// Database
	dbConfig := config.NewDatabaseConfig()
	repoFactory := factory.NewRepositoryFactory(dbConfig)
	repo, err := repoFactory.NewCharacterRepository(context.Background())
	if err != nil {
		logger.Fatal("failed to create repository", zap.Error(err))
	}

	// Use cases
	combinedUC := usecase.CombinedUC{CharSvc: rickClient, PokeSvc: pokeClient}
	charUC := usecase.FetchCharacter{Svc: rickClient}
	pokeUC := usecase.FetchPokemon{Svc: pokeClient}
	storeUC := usecase.StoreCharacter{Repo: repo}
	listUC := usecase.ListCharacters{Repo: repo}

	// Router
	router := httpinfra.BuildRouter(logger, combinedUC, charUC, pokeUC, storeUC, listUC, rickCB, pokeCB, repo)

	logger.Info("server up", zap.String("port", cfg.Port))
	logger.Fatal("server exited", zap.Error(router.Run(":"+cfg.Port)))
}
