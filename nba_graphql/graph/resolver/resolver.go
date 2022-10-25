package resolver

import "github.com/zvandehy/DataTrain/nba_graphql/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Db *database.NBADatabaseClient
}
