import { gql, useQuery } from "@apollo/client";
import { ModelInput } from "../shared/interfaces/custom-prediction.interface";
import { Proposition } from "../shared/interfaces/graphql/proposition.interface";

export const USE_GET_PROPOSITIONS = gql`
  query GetPropositions(
    $startDate: String!
    $endDate: String!
    $customModel: ModelInput!
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
      lastModified
      actualResult
      prediction(input: $customModel) {
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

export interface ProjectionQuery {
  startDate: string;
  endDate: string;
  customModel: ModelInput;
}

export const useGetPropositions = ({
  startDate,
  endDate,
  customModel,
}: ProjectionQuery) => {
  const { loading, error, data, fetchMore } = useQuery(USE_GET_PROPOSITIONS, {
    variables: { startDate, endDate, customModel },
  });
  return {
    loading,
    error,
    data: data?.propositions as Proposition[],
    fetchMore,
    minDate: startDate,
  };
};
