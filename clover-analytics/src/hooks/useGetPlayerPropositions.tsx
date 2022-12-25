import { gql, useQuery } from "@apollo/client";
import { ModelInput } from "../shared/interfaces/custom-prediction.interface";
import { Proposition } from "../shared/interfaces/graphql/proposition.interface";

export const USE_GET_PLAYER_PROPOSITIONS = gql`
  query GetPropositions(
    $playerID: Int!
    $startDate: String!
    $endDate: String!
    $customModel: ModelInput!
  ) {
    propositions(
      input: { startDate: $startDate, endDate: $endDate, PlayerID: $playerID }
    ) {
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
        breakdowns {
          name
          over
          under
          push
          overPct
          underPct
          pushPct
          derivedAverage
          weight
          stdDev
          pctChange
          contribution
          base
          derivedGames {
            date
            opponent {
              abbreviation
              teamID
            }
            player {
              name
              playerID
            }
            gameID
            home_or_away
            minutes
          }
        }
      }
    }
  }
`;

export interface PlayerPropositionsQuery {
  playerID: number;
  startDate: string;
  endDate: string;
  customModel: ModelInput;
}

export const useGetPlayerPropositions = ({
  playerID,
  startDate,
  endDate,
  customModel,
}: PlayerPropositionsQuery) => {
  const { loading, error, data, fetchMore } = useQuery(
    USE_GET_PLAYER_PROPOSITIONS,
    {
      variables: { playerID, startDate, endDate, customModel },
    }
  );
  return {
    loading,
    error,
    data: data?.propositions.sort(
      (a: Proposition, b: Proposition) =>
        a.prediction.significance - b.prediction.significance
    ) as Proposition[],
    fetchMore,
  };
};
