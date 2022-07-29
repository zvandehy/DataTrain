import { ApolloQueryResult, gql, useQuery } from "@apollo/client";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import { CalculatePredictions } from "../shared/functions/predictions.fn";
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
  similarPlayers,
  similarTeams,
}: {
  projectionFilter: ProjectionFilter;
  gameFilter: GameFilter;
  predictionFilter: GameFilter;
  similarPlayers: boolean;
  similarTeams: boolean;
}): ProjectionQueryResult => {
  let QUERY =
    similarPlayers && similarTeams
      ? GET_PROJECTIONS_AND_SIMILAR_PLAYERS_AND_TEAMS
      : similarPlayers
      ? GET_PROJECTIONS_AND_SIMILAR_PLAYERS
      : similarTeams
      ? GET_PROJECTIONS_AND_SIMILAR_TEAMS
      : GET_PROJECTIONS;
  console.log("similarTeams: ", similarTeams);
  console.log("similarPlayers: ", similarPlayers);
  // QUERY = GET_PROJECTIONS;
  const { loading, error, data } = useQuery(QUERY, {
    variables: { playerFilter: projectionFilter, gameFilter: gameFilter },
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
    const projections = CalculatePredictions(
      data.projections,
      predictionFilter
    );
    console.log(projections);
    // setSkip(true);
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
