import { gql, useQuery } from "@apollo/client";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import { CalculatePredictions } from "../shared/functions/predictions.fn";
import { GameFilter } from "../shared/interfaces/graphql/filters.interface";
import { Player } from "../shared/interfaces/graphql/player.interface";

export const GET_PLAYER = gql`
  query GetPlayer($playerID: Int!, $season: String!) {
    player(input: { playerID: $playerID }) {
      name
      playerID
      position
      currentTeam {
        abbreviation
        teamID
        name
        injuries {
          startDate
          returnDate
          status
          player {
            name
          }
        }
      }
      projections(input: {}) {
        date
        player {
          currentTeam {
            abbreviation
            name
            teamID
          }
          similarPlayers(input: { season: $season }) {
            name
            playerID
            games(input: { season: $season }) {
              points
              season
              assists
              assist_percentage
              rebounds
              offensive_rebounds
              offensive_rebound_percentage
              defensive_rebounds
              defensive_rebound_percentage
              personal_fouls_drawn
              steals
              blocks
              turnovers
              opponent {
                abbreviation
                teamID
              }
              team {
                abbreviation
                teamID
              }
              minutes
              date
              field_goals_attempted
              field_goal_percentage
              field_goals_made
              three_pointers_attempted
              three_pointers_made
              free_throws_attempted
              free_throws_made
              free_throws_percentage
              usage
              three_point_percentage
              effective_field_goal_percentage
              playoffs
            }
          }
        }
        opponent {
          abbreviation
          teamID
          name
          injuries {
            startDate
            returnDate
            status
            player {
              name
            }
          }
          similarTeams(input: { season: $season }) {
            teamID
            name
            abbreviation
          }
        }
        propositions {
          target
          type
          sportsbook
        }
        result {
          points
          season
          assists
          assist_percentage
          rebounds
          offensive_rebounds
          offensive_rebound_percentage
          defensive_rebounds
          defensive_rebound_percentage
          personal_fouls_drawn
          steals
          blocks
          turnovers
          opponent {
            abbreviation
            teamID
          }
          team {
            abbreviation
            teamID
          }
          minutes
          date
          field_goals_attempted
          field_goal_percentage
          field_goals_made
          three_pointers_attempted
          three_pointers_made
          free_throws_attempted
          free_throws_made
          free_throws_percentage
          usage
          three_point_percentage
          effective_field_goal_percentage
          playoffs
        }
      }
      games(input: { season: $season }) {
        points
        season
        assists
        assist_percentage
        rebounds
        offensive_rebounds
        offensive_rebound_percentage
        defensive_rebounds
        defensive_rebound_percentage
        personal_fouls_drawn
        steals
        blocks
        turnovers
        teamStats {
          points
          assists
          rebounds
          three_pointers_attempted
          blocks
          steals
        }
        team {
          abbreviation
          teamID
          name
        }
        opponent {
          abbreviation
          teamID
          name
          injuries {
            startDate
            returnDate
            status
            player {
              name
            }
          }
        }
        minutes
        margin
        home_or_away
        date
        field_goals_attempted
        field_goal_percentage
        field_goals_made
        three_pointers_attempted
        three_pointers_made
        free_throws_attempted
        free_throws_made
        free_throws_percentage
        usage
        three_point_percentage
        effective_field_goal_percentage
        playoffs
      }
    }
  }
`;
// # const { loading, error, data } = useQuery(query, {
// #   variables: { playerID: playerID, season: season },
// #   client: client,
// # });
// # `;

export interface QueryResult {
  loading?: any;
  error?: any;
  data?: any;
}

export interface ProjectionQueryResult extends QueryResult {
  data?: Player;
}

export const useGetPlayerDetails = ({
  playerID,
  predictionFilter,
}: {
  playerID: number;
  predictionFilter: GameFilter;
}): ProjectionQueryResult => {
  const { loading, error, data } = useQuery(GET_PLAYER, {
    //TODO: refactor query to use predictionFilter
    variables: { playerID: playerID, season: predictionFilter.season },
  });
  let loadingComponent;
  if (loading) {
    loadingComponent = (
      <Box className={"loading-player"}>
        <h1>Loading </h1>
        <CircularProgress />
      </Box>
    );
  }
  let errorComponent;
  if (error) {
    errorComponent = <Box>{JSON.stringify(error) + error.message}</Box>;
  }
  if (data && data?.player?.projections) {
    const projections = CalculatePredictions(
      data.player.projections,
      predictionFilter,
      data.player.games
    );
    const player: Player = {
      ...data.player,
      projections: projections,
    };
    console.log(player);
    return {
      loading: loadingComponent,
      error: errorComponent,
      data: player,
    };
  }
  return {
    loading: loadingComponent,
    error: errorComponent,
  };
};
