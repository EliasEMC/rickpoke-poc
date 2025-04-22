package mongo

import (
	"context"

	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CharacterRepo struct {
	Coll *mongo.Collection
}

func (r CharacterRepo) Save(ctx context.Context, c *model.Rick) error {
	_, err := r.Coll.InsertOne(ctx, c)
	return err
}

func (r CharacterRepo) FindByID(ctx context.Context, id int) (*model.Rick, error) {
	var out model.Rick
	err := r.Coll.FindOne(ctx, bson.M{"id": id}).Decode(&out)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &out, err
}

// FindAll recupera todos los personajes de la base de datos
func (r *CharacterRepo) FindAll(ctx context.Context) ([]*model.Rick, error) {
	var characters []model.Rick
	cursor, err := r.Coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &characters); err != nil {
		return nil, err
	}
	
	// Convertir a slice de punteros
	result := make([]*model.Rick, len(characters))
	for i := range characters {
		result[i] = &characters[i]
	}
	return result, nil
}

// Ping verifica la conexión con la base de datos
func (r *CharacterRepo) Ping(ctx context.Context) error {
	return r.Coll.Database().Client().Ping(ctx, nil)
}

var _ service.CharacterRepository = (*CharacterRepo)(nil)
