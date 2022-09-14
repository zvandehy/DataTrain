import { gql, useQuery } from "@apollo/client";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import moment from "moment";
import { CustomCalculation } from "../shared/interfaces/custom-prediction.interface";
import {
  GameFilter,
  ProjectionFilter,
} from "../shared/interfaces/graphql/filters.interface";
import { Projection } from "../shared/interfaces/graphql/projection.interface";

export const GET_PROJECTIONS = gql`
  query GetProjections(
    $playerFilter: ProjectionFilter!
    $gameFilter: GameFilter!
  ) {
    projections(input: $playerFilter) {
      player {
        name
        position
        playerID
        currentTeam {
          abbreviation
          teamID
          name
        }
        games(input: $gameFilter) {
          season
          date
          gameID
          playoffs
          opponent {
            name
            teamID
            abbreviation
          }
          points
          assists
          rebounds
          defensive_rebounds
          offensive_rebounds
          three_pointers_attempted
          three_pointers_made
          free_throws_attempted
          free_throws_made
          minutes
          blocks
          turnovers
          steals
        }
      }
      opponent {
        abbreviation
        teamID
        name
      }
      propositions {
        target
        type
        sportsbook
        lastModified
      }
      startTime
      result {
        points
        assists
        rebounds
        defensive_rebounds
        offensive_rebounds
        three_pointers_attempted
        three_pointers_made
        free_throws_attempted
        free_throws_made
        minutes
        blocks
        turnovers
        steals
      }
    }
  }
`;

export const GET_PROJECTIONS_AND_SIMILAR_PLAYERS = gql`
  query GetProjections(
    $playerFilter: ProjectionFilter!
    $gameFilter: GameFilter!
  ) {
    projections(input: $playerFilter) {
      player {
        name
        position
        playerID
        currentTeam {
          abbreviation
          teamID
          name
        }
        games(input: $gameFilter) {
          season
          date
          gameID
          playoffs
          opponent {
            name
            teamID
            abbreviation
          }
          points
          assists
          rebounds
          defensive_rebounds
          offensive_rebounds
          three_pointers_attempted
          three_pointers_made
          free_throws_attempted
          free_throws_made
          minutes
          blocks
          turnovers
          steals
        }
        similarPlayers(input: $gameFilter) {
          name
          currentTeam {
            name
            abbreviation
          }
          position
          playerID
          games(input: $gameFilter) {
            season
            date
            gameID
            playoffs
            opponent {
              name
              teamID
              abbreviation
            }
            points
            assists
            rebounds
            defensive_rebounds
            offensive_rebounds
            three_pointers_attempted
            three_pointers_made
            free_throws_attempted
            free_throws_made
            minutes
            blocks
            turnovers
            steals
          }
        }
      }
      opponent {
        abbreviation
        teamID
        name
      }
      propositions {
        target
        type
        sportsbook
      }
      startTime
      result {
        points
        assists
        rebounds
        defensive_rebounds
        offensive_rebounds
        three_pointers_attempted
        three_pointers_made
        free_throws_attempted
        free_throws_made
        minutes
        blocks
        turnovers
        steals
      }
    }
  }
`;

export const GET_PROJECTIONS_AND_SIMILAR_TEAMS = gql`
  query GetProjections(
    $playerFilter: ProjectionFilter!
    $gameFilter: GameFilter!
  ) {
    projections(input: $playerFilter) {
      player {
        name
        position
        playerID
        currentTeam {
          abbreviation
          teamID
          name
        }
        games(input: $gameFilter) {
          season
          date
          gameID
          playoffs
          opponent {
            name
            teamID
            abbreviation
          }
          points
          assists
          rebounds
          defensive_rebounds
          offensive_rebounds
          three_pointers_attempted
          three_pointers_made
          free_throws_attempted
          free_throws_made
          minutes
          blocks
          turnovers
          steals
        }
      }
      opponent {
        abbreviation
        teamID
        name
        similarTeams(input: $gameFilter) {
          name
          abbreviation
          teamID
        }
      }
      propositions {
        target
        type
        sportsbook
      }
      startTime
      result {
        points
        assists
        rebounds
        defensive_rebounds
        offensive_rebounds
        three_pointers_attempted
        three_pointers_made
        free_throws_attempted
        free_throws_made
        minutes
        blocks
        turnovers
        steals
      }
    }
  }
`;

export const GET_PROJECTIONS_AND_SIMILAR_PLAYERS_AND_TEAMS = gql`
  query GetProjections(
    $playerFilter: ProjectionFilter!
    $gameFilter: GameFilter!
  ) {
    projections(input: $playerFilter) {
      player {
        name
        position
        playerID
        currentTeam {
          abbreviation
          teamID
          name
        }
        games(input: $gameFilter) {
          season
          date
          gameID
          playoffs
          opponent {
            name
            teamID
            abbreviation
          }
          points
          assists
          rebounds
          defensive_rebounds
          offensive_rebounds
          three_pointers_attempted
          three_pointers_made
          free_throws_attempted
          free_throws_made
          minutes
          blocks
          turnovers
          steals
        }
        similarPlayers(input: $gameFilter) {
          name
          currentTeam {
            name
            abbreviation
          }
          position
          playerID
          games(input: $gameFilter) {
            season
            date
            gameID
            playoffs
            opponent {
              name
              teamID
              abbreviation
            }
            points
            assists
            rebounds
            defensive_rebounds
            offensive_rebounds
            three_pointers_attempted
            three_pointers_made
            free_throws_attempted
            free_throws_made
            minutes
            blocks
            turnovers
            steals
          }
        }
      }
      opponent {
        abbreviation
        teamID
        name
        similarTeams(input: $gameFilter) {
          name
          abbreviation
          teamID
        }
      }
      propositions {
        target
        type
        sportsbook
        lastModified
      }
      startTime
      result {
        points
        assists
        rebounds
        defensive_rebounds
        offensive_rebounds
        three_pointers_attempted
        three_pointers_made
        free_throws_attempted
        free_throws_made
        minutes
        blocks
        turnovers
        steals
      }
    }
  }
`;

export interface QueryResult {
  loading?: any;
  error?: any;
  data?: any;
}

export interface ProjectionQueryResult extends QueryResult {
  data: Projection[];
}

export const useGetProjections = ({
  projectionFilter,
  gameFilter,
  predictionFilter,
  customModel,
}: {
  projectionFilter: ProjectionFilter;
  gameFilter: GameFilter;
  predictionFilter: GameFilter;
  customModel: CustomCalculation;
}): ProjectionQueryResult => {
  // console.log("Get projections for: ", projectionFilter, gameFilter);
  let QUERY =
    customModel.similarPlayers &&
    customModel.similarPlayers.weight > 0 &&
    customModel.similarTeams &&
    customModel.similarTeams.weight > 0
      ? GET_PROJECTIONS_AND_SIMILAR_PLAYERS_AND_TEAMS
      : customModel.similarPlayers && customModel.similarPlayers.weight > 0
      ? GET_PROJECTIONS_AND_SIMILAR_PLAYERS
      : customModel.similarTeams && customModel.similarTeams.weight > 0
      ? GET_PROJECTIONS_AND_SIMILAR_TEAMS
      : GET_PROJECTIONS;
  // QUERY = GET_PROJECTIONS;
  const { loading, error, data } = useQuery(QUERY, {
    variables: { playerFilter: projectionFilter, gameFilter: predictionFilter },
  });
  let loadingComponent;
  if (loading) {
    loadingComponent = (
      <Box className={"loading-results"}>
        <h1>Loading </h1>
        <CircularProgress />
      </Box>
    );
  }
  let errorComponent;
  if (error) {
    errorComponent = <Box>{JSON.stringify(error) + error.message}</Box>;
  }
  if (data && data?.projections) {
    let projections = data?.projections.map((projection: Projection) => {
      let player = { ...projection.player };
      let games = player.games.map((game) => {
        return { ...game };
      });
      games.sort((a, b) => {
        return moment(a.date).diff(b.date);
      });
      player.games = games;
      return { ...projection, player: player };
    });

    return {
      loading: loadingComponent,
      error: errorComponent,
      data: projections,
    };
  }
  return {
    loading: loadingComponent,
    error: errorComponent,
    data: [],
  };
};
