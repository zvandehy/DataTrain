package resolver

import "github.com/zvandehy/DataTrain/nba_graphql/database"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	database.BasketballRepository
}

// NewResolver returns a new instance of the resolver.
func NewResolver(db database.BasketballRepository) *Resolver {
	return &Resolver{BasketballRepository: db}
}
