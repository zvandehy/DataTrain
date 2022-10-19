import { ApolloError, gql, useQuery } from "@apollo/client";
import { ModelInput } from "../shared/interfaces/custom-prediction.interface";
import { AverageStats } from "../shared/interfaces/graphql/game.interface";
import { Player } from "../shared/interfaces/graphql/player.interface";

export const GET_UPCOMING_PROJECTIONS = gql`
  query GetUpcomingProjections(
    $startDate: String!
    $endDate: String!
    $customModel: ModelInput!
  ) {
    players(
      input: {
        withGames: { startDate: $startDate, endDate: $endDate }
        withPropositions: {
          period: { startDate: $startDate, endDate: $endDate }
        }
      }
    ) {
      playerID
      name
      image
      position
      team {
        abbreviation
      }
      games(input: { startDate: $startDate, endDate: $endDate }) {
        date
        outcome
        opponent {
          abbreviation
        }
        points
        rebounds
        assists
        steals
        blocks
        turnovers
        minutes
        fantasy_score
        points_rebounds_assists
        points_rebounds
        points_assists
        rebounds_assists
        blocks_steals
        free_throws_percentage
        free_throws_made
        field_goal_percentage
        field_goals_made
        three_point_percentage
        three_pointers_made
        three_pointers_attempted
        field_goals_attempted
        free_throws_attempted
        home_or_away
        prediction(input: $customModel) {
          weightedTotal {
            points
            rebounds
            assists
            points_rebounds
            points_assists
            rebounds_assists
            points_rebounds_assists
            blocks
            steals
            blocks_steals
            turnovers
            fantasy_score
            three_pointers_made
            free_throws_made
            minutes
          }
          predictionAccuracy {
            points
            rebounds
            assists
            points_rebounds
            points_assists
            rebounds_assists
            points_rebounds_assists
            blocks
            steals
            blocks_steals
            turnovers
            fantasy_score
            three_pointers_made
            free_throws_made
            minutes
          }
          fragments {
            name
            weight
            propositions {
              target
              type
              sportsbook
              lastModified
              analysis {
                numOver
                numUnder
                numPush
                pctOver
                pctUnder
                pctPush
              }
            }
            derived {
              points
              rebounds
              assists
              points_rebounds
              points_assists
              rebounds_assists
              points_rebounds_assists
              blocks
              steals
              blocks_steals
              turnovers
              fantasy_score
              three_pointers_made
              free_throws_made
              minutes
              games_played
            }
            base {
              points
              rebounds
              assists
              points_rebounds
              points_assists
              rebounds_assists
              points_rebounds_assists
              blocks
              steals
              blocks_steals
              turnovers
              fantasy_score
              three_pointers_made
              free_throws_made
              minutes
            }
            pctChange {
              points
              rebounds
              assists
              points_rebounds
              points_assists
              rebounds_assists
              points_rebounds_assists
              blocks
              steals
              blocks_steals
              turnovers
              fantasy_score
              three_pointers_made
              free_throws_made
              minutes
            }
            # derived games?
          }
        }
      }
    }
  }
`;
export interface QueryResult {
  loading?: boolean;
  error?: ApolloError;
  data?: any;
}

export interface ProjectionQueryResult extends QueryResult {
  data: Player[];
}

export const useGetProjections = ({
  startDate,
  endDate,
  customModel,
}: {
  startDate: string;
  endDate: string;
  customModel: ModelInput;
}): ProjectionQueryResult => {
  const { loading, error, data } = useQuery(GET_UPCOMING_PROJECTIONS, {
    variables: { startDate, endDate, customModel },
  });
  if (loading) {
    return { loading, error, data: [] };
  }
  if (error) {
    return { loading, error, data: [] };
  }
  if (data && data?.players) {
    let players: Player[] = [];
    data.players.forEach((player: Player) => {
      console.log(player.name);
      if (player.games.length > 0) {
        const newPlayer = player;
        newPlayer.games = newPlayer.games.map((game) => {
          const newGame = game;
          if (game.prediction.fragments.length > 0) {
            newGame.prediction.propositions =
              game.prediction.fragments[0].propositions.map((prop) => {
                const newProp = prop;
                newProp.analysis = undefined;
                newProp.estimation =
                  newGame.prediction.weightedTotal[
                    prop.type as keyof AverageStats
                  ];
                newProp.predictionDiff = +(
                  prop.estimation - prop.target
                ).toFixed(2);
                newProp.prediction = prop.predictionDiff > 0 ? "OVER" : "UNDER";
                newProp.actual = newGame[
                  prop.type as keyof AverageStats
                ] as number;
                newProp.predictionHit =
                  prop.actual === prop.target
                    ? "PUSH"
                    : prop.predictionDiff > 0 && prop.actual > prop.target
                    ? "HIT"
                    : prop.predictionDiff < 0 && prop.actual < prop.target
                    ? "HIT"
                    : "MISS";
                return newProp;
              });
          }
          //       const propositions = game.prediction.fragments[0].propositions;
          //       for (let i = 0; i < propositions.length; i++) {
          //         const prop = newGame.prediction.propositions[i];

          //         newGame.prediction.propositions.push(prop);
          //       }
          return newGame;
        });
        players.push(newPlayer);
      }
    });
    return {
      loading: false,
      error: undefined,
      data: players,
    };
  }
  return { loading, error, data: [] };
};
