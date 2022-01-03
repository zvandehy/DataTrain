package database

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var once sync.Once
var instance *NBADatabaseClient

type NBADatabaseClient struct {
	*mongo.Client
	conn string
}

func ConnectDB(ctx context.Context) (*NBADatabaseClient, error) {
	var connErr error
	once.Do(func() {
		instance = &NBADatabaseClient{conn: "mongodb+srv://datatrain:nbawinners@datatrain.i5rgk.mongodb.net/nba?retryWrites=true&w=majority"}
		client, connErr := mongo.NewClient(options.Client().ApplyURI(instance.conn))
		if connErr != nil {
			return
		}
		connErr = client.Connect(ctx)
		if connErr != nil {
			return
		}
		instance.Client = client
	})
	if connErr != nil {
		return nil, connErr
	}
	return instance, nil
}
