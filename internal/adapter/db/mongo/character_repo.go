package mongo

import (
	"context"
	"time"

	"github.com/EliasEMC/rickpoke-poc/internal/domain/model"
	"github.com/EliasEMC/rickpoke-poc/internal/domain/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type CharacterRepo struct {
	Coll *mongo.Collection
}

func NewCharacterRepo(coll *mongo.Collection) *CharacterRepo {
	return &CharacterRepo{Coll: coll}
}

func (r *CharacterRepo) Save(ctx context.Context, c *model.Rick) error {
	_, err := r.Coll.UpdateOne(
		ctx,
		map[string]interface{}{"id": c.ID},
		map[string]interface{}{"$set": c},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *CharacterRepo) FindByID(ctx context.Context, id int) (*model.Rick, error) {
	var character model.Rick
	err := r.Coll.FindOne(ctx, map[string]interface{}{"id": id}).Decode(&character)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &character, nil
}

func (r *CharacterRepo) FindAll(ctx context.Context) ([]*model.Rick, error) {
	cursor, err := r.Coll.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var characters []*model.Rick
	if err := cursor.All(ctx, &characters); err != nil {
		return nil, err
	}
	return characters, nil
}

func (r *CharacterRepo) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return r.Coll.Database().Client().Ping(ctx, readpref.Primary())
}

var _ service.CharacterRepository = (*CharacterRepo)(nil) 