import { ApolloError, gql, useQuery } from "@apollo/client";
import { ModelInput } from "../shared/interfaces/custom-prediction.interface";
import { Proposition } from "../shared/interfaces/graphql/proposition.interface";

export const USE_GET_PROPOSITIONS = gql`
  query GetPropositions(
    $startDate: String!
    $endDate: String! # $customModel: ModelInput!
  ) {
    propositions(input: { startDate: $startDate, endDate: $endDate }) {
      sportsbook
      game {
        date
        gameID
        home_or_away
        opponent {
          abbreviation
          teamID
          name
          logoImage
        }
        player {
          name
          playerID
          image
          team {
            abbreviation
            logoImage
          }
        }
      }
      type
      target
      outcome
      actualResult
      prediction(
        # input: $customModel
        input: {
          model: "SEASON"
          gameBreakdowns: [
            { name: "2022-23", weight: 45, filter: { seasonMatch: true } }
            {
              name: "2021-22"
              weight: 25
              filter: { previousSeasonMatch: true }
            }
            { name: "Opponent", weight: 20, filter: { opponentMatch: true } }
          ]
          similarPlayerInput: { weight: 30, limit: 5 }
        }
      ) {
        estimation
        significance
        wager
        wagerOutcome
        stdDev
        cumulativeOver
        cumulativeUnder
        cumulativePush
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
  data: Proposition[];
}

export const useGetPropositions = ({
  startDate,
  endDate,
  customModel,
}: {
  startDate: string;
  endDate: string;
  customModel: ModelInput;
}): ProjectionQueryResult => {
  const { loading, error, data } = useQuery(USE_GET_PROPOSITIONS, {
    variables: { startDate, endDate },
  });
  if (loading) {
    return { loading, error, data: [] };
  }
  if (error) {
    return { loading, error, data: [] };
  }
  if (data && data?.propositions) {
    return {
      loading: false,
      error: undefined,
      data: data.propositions,
    };
  }
  return { loading, error, data: [] };
};

// function GetMaxPctDiff(player: Player) {
//   let maxPctDiff = 0;
//   player.games.forEach((game) => {
//     game.prediction.propositions.forEach((prop) => {
//       if (Math.abs(prop.predictionTargetDiffPCT) > Math.abs(maxPctDiff)) {
//         maxPctDiff = Math.abs(prop.predictionTargetDiffPCT);
//       }
//     });
//   });
//   return maxPctDiff;
// }
