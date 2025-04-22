package factory

import (
    "context"
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    mongodriver "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/EliasEMC/rickpoke-poc/internal/config"
    "github.com/EliasEMC/rickpoke-poc/internal/domain/service"
    mongorepo "github.com/EliasEMC/rickpoke-poc/internal/adapter/db/mongo"
    "github.com/EliasEMC/rickpoke-poc/internal/adapter/db/postgres"
    "go.uber.org/zap"
)

type RepositoryFactory struct {
    config *config.DatabaseConfig
    logger *zap.Logger
}

func NewRepositoryFactory(config *config.DatabaseConfig) *RepositoryFactory {
    logger, _ := zap.NewProduction()
    return &RepositoryFactory{
        config: config,
        logger: logger,
    }
}

func (f *RepositoryFactory) NewCharacterRepository(ctx context.Context) (service.CharacterRepository, error) {
    switch f.config.Type {
    case "mongodb":
        return f.newMongoRepo(ctx)
    case "postgres":
        return f.newPostgresRepo(ctx)
    default:
        return nil, fmt.Errorf("unsupported database type: %s", f.config.Type)
    }
}

func (f *RepositoryFactory) newMongoRepo(ctx context.Context) (service.CharacterRepository, error) {
    uri := fmt.Sprintf("mongodb://%s:%d", f.config.Host, f.config.Port)
    client, err := mongodriver.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }

    db := client.Database(f.config.Name)
    coll := db.Collection("characters")
    return &mongorepo.CharacterRepo{Coll: coll}, nil
}

func (f *RepositoryFactory) newPostgresRepo(ctx context.Context) (service.CharacterRepository, error) {
    connStr := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=disable",
        f.config.User,
        f.config.Password,
        f.config.Host,
        f.config.Port,
        f.config.Name,
    )

    f.logger.Info("connecting to postgres",
        zap.String("host", f.config.Host),
        zap.Int("port", f.config.Port),
        zap.String("user", f.config.User),
        zap.String("dbname", f.config.Name),
    )

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        f.logger.Error("failed to open postgres connection", zap.Error(err))
        return nil, err
    }

    if err := db.PingContext(ctx); err != nil {
        f.logger.Error("failed to ping postgres", zap.Error(err))
        return nil, err
    }

    f.logger.Info("successfully connected to postgres")
    return postgres.NewCharacterRepo(db), nil
} 