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

var _ service.CharacterRepository = (*CharacterRepo)(nil)
