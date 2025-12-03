package mongo_db

import (
	"context"

	"github.com/tqhuy-dev/xgen/utilities"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	*mongo.Client
	option Option
}

func NewMongoDB(option Option) (*MongoDB, error) {
	client, err := mongo.Connect(options.Client().
		ApplyURI(option.URI()).
		SetMaxPoolSize(option.MaxPoolSize).
		SetMinPoolSize(option.MinPoolSize).
		SetMaxConnIdleTime(option.MaxIdleTime),
	)
	if err != nil {
		return nil, err
	}
	ctx := utilities.Ternary(option.Ctx == nil, context.Background(), option.Ctx)
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &MongoDB{
		Client: client,
		option: option,
	}, nil
}

func (db *MongoDB) Collection(collection string) *mongo.Collection {
	coll := db.Database(db.option.DB).Collection(collection)
	return coll
}

func (db *MongoDB) Close() error {
	return db.Client.Disconnect(context.Background())
}
