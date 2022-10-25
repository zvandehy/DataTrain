import { ApolloError, gql, useQuery } from "@apollo/client";
import moment from "moment";
import { ModelInput } from "../shared/interfaces/custom-prediction.interface";
import {
  AverageStats,
  Game,
  PropositionA,
} from "../shared/interfaces/graphql/game.interface";
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
        let newGames: Game[] = [];
        newPlayer.games.forEach((game) => {
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
                newProp.estimationPerMin = +(
                  newProp.estimation / newGame.prediction.weightedTotal.minutes
                ).toFixed(2);
                newProp.predictionTargetDiff = +(
                  prop.estimation - prop.target
                ).toFixed(2);
                newProp.predictionTargetDiffPCT = +(
                  (prop.predictionTargetDiff / prop.target) *
                  100
                ).toFixed();
                newProp.prediction =
                  prop.predictionTargetDiff > 0 ? "OVER" : "UNDER";
                newProp.actual = +(
                  newGame[prop.type as keyof AverageStats] as number
                ).toFixed();
                newProp.actualPerMin = +(
                  newProp.actual / newGame.prediction.weightedTotal.minutes
                ).toFixed(2);
                newProp.predictionHit =
                  newGame.outcome.toUpperCase() === "PENDING"
                    ? "PENDING"
                    : prop.actual === prop.target
                    ? "PUSH"
                    : prop.predictionTargetDiff > 0 && prop.actual > prop.target
                    ? "HIT"
                    : prop.predictionTargetDiff < 0 && prop.actual < prop.target
                    ? "HIT"
                    : "MISS";
                newProp.actualDiff = +(prop.actual - prop.target).toFixed(2);
                newProp.actualDiffPCT = +(
                  (prop.actualDiff / prop.target) *
                  100
                ).toFixed();
                newProp.actualDiffPerMin = +(
                  newProp.actualPerMin - newProp.estimationPerMin
                ).toFixed(2);
                newProp.actualDiffPerMinPCT = +(
                  (prop.actualDiffPerMin / prop.estimationPerMin) *
                  100
                ).toFixed();
                return newProp;
              });
            newGame.prediction.propositions =
              newGame.prediction.propositions.sort((a, b) =>
                Math.abs(a.predictionTargetDiffPCT) >
                Math.abs(b.predictionTargetDiffPCT)
                  ? -1
                  : 1
              );
          }
          //       const propositions = game.prediction.fragments[0].propositions;
          //       for (let i = 0; i < propositions.length; i++) {
          //         const prop = newGame.prediction.propositions[i];

          //         newGame.prediction.propositions.push(prop);
          //       }

          // game date is before today and they played (not pending) or game is today or later
          if (
            (moment(game.date).isBefore(startDate, "day") &&
              game.outcome.toUpperCase() !== "PENDING") ||
            moment(game.date).isSameOrAfter(startDate, "day")
          ) {
            newGames.push(newGame);
          } else {
            console.log(
              game.date,
              game.outcome,
              moment(game.date).isBefore(),
              moment(game.date).isSameOrAfter(),
              game.outcome.toUpperCase() !== "PENDING"
            );
          }
        });
        newPlayer.games = newGames;
        if (newPlayer.games.length > 0) {
          players.push(newPlayer);
        }
      }
      players = players.sort((a, b) => GetMaxPctDiff(b) - GetMaxPctDiff(a));
    });
    return {
      loading: false,
      error: undefined,
      data: players,
    };
  }
  return { loading, error, data: [] };
};

function GetMaxPctDiff(player: Player) {
  let maxPctDiff = 0;
  player.games.forEach((game) => {
    game.prediction.propositions.forEach((prop) => {
      if (Math.abs(prop.predictionTargetDiffPCT) > Math.abs(maxPctDiff)) {
        maxPctDiff = Math.abs(prop.predictionTargetDiffPCT);
      }
    });
  });
  return maxPctDiff;
}
