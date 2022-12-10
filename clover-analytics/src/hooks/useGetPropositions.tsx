import { gql, useQuery } from "@apollo/client";
import moment from "moment";
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
  // const sortedPropositions: Proposition[] = data?.propositions.sort(
  //   (a: Proposition, b: Proposition) => {
  //     const aDate = moment(a.game.date);
  //     const bDate = moment(b.game.date);
  //     if (aDate.isBefore(bDate)) {
  //       return -1;
  //     }
  //     if (aDate.isAfter(bDate)) {
  //       return 1;
  //     }
  //     return 0;
  //   }
  // );
  // console.log(startDate);
  return {
    loading,
    error,
    data: data?.propositions as Proposition[],
    fetchMore,
    minDate: startDate,
    // maxDate: sortedPropositions?.[sortedPropositions.length - 1]?.game.date,
  };
};
