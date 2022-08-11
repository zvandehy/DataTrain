import { gql, useQuery } from "@apollo/client";
import {
  GameFilter,
  PlayerFilter,
} from "../shared/interfaces/graphql/filters.interface";
import { Player } from "../shared/interfaces/graphql/player.interface";

const GET_PLAYERS = gql`
  query GetPlayers($playerFilter: PlayerFilter!, $gameFilter: GameFilter!) {
    filterPlayers(input: $playerFilter) {
      name
      playerID
      position
      currentTeam {
        name
        location
        abbreviation
      }
      games(input: $gameFilter) {
        date
        opponent {
          name
          abbreviation
        }
        points
        rebounds
        assists
        steals
        blocks
        turnovers
        three_pointers_attempted
        three_pointers_made
        free_throws_attempted
        free_throws_made
      }
    }
  }
`;

export const useGetPlayers = (
  playerFilter: PlayerFilter,
  gameFilter: GameFilter
): Player[] | undefined => {
  const { data } = useQuery(GET_PLAYERS, {
    variables: { playerFilter: playerFilter, gameFilter: gameFilter },
  });
  if (data) {
    // console.log(data.filterPlayers);
    return data?.filterPlayers;
  }
};
