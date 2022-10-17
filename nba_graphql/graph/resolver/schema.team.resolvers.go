package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	"github.com/zvandehy/DataTrain/nba_graphql/graph/generated"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
)

// Players is the resolver for the players field.
func (r *teamResolver) Players(ctx context.Context, obj *model.Team) ([]*model.Player, error) {
	panic(fmt.Errorf("(r *teamResolver) Players not implemented"))
}

// LogoImage is the resolver for the logoImage field.
func (r *teamResolver) LogoImage(ctx context.Context, obj *model.Team) (string, error) {
	if obj.League == "NBA" {
		return fmt.Sprintf("https://cdn.nba.com/logos/nba/%d/primary/D/logo.svg", obj.TeamID), nil
	}
	return fmt.Sprintf("https://%s.wnba.com/wp-content/themes/wnba-parent/img/logos/%s-primary-logo.svg", strings.ToLower(obj.Name), strings.ToLower(obj.Name)), nil
}

// Opponent is the resolver for the opponent field.
func (r *teamGameResolver) Opponent(ctx context.Context, obj *model.TeamGame) (*model.Team, error) {
	panic(fmt.Errorf("(r *teamResolver) Players not implemented"))
}

// Team returns generated.TeamResolver implementation.
func (r *Resolver) Team() generated.TeamResolver { return &teamResolver{r} }

// TeamGame returns generated.TeamGameResolver implementation.
func (r *Resolver) TeamGame() generated.TeamGameResolver { return &teamGameResolver{r} }

type teamResolver struct{ *Resolver }
type teamGameResolver struct{ *Resolver }
