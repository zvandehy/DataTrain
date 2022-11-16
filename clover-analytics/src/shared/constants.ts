import { gql } from "@apollo/client";
import { ModelInput } from "./interfaces/custom-prediction.interface";
import { SeasonOption } from "./interfaces/graphql/filters.interface";
import {
  Assists,
  Blocks,
  BlocksSteals,
  Fantasy,
  FreeThrowsAttempted,
  FreeThrowsMade,
  Minutes,
  Points,
  PointsAssists,
  PointsRebounds,
  PointsReboundsAssists,
  Rebounds,
  ReboundsAssists,
  Steals,
  ThreeFGA,
  ThreeFGM,
  Turnovers,
} from "./interfaces/stat.interface";

//array of all stats found in stat.interface
export const ALL_STATS = [
  Points,
  Rebounds,
  Assists,
  Steals,
  Blocks,
  BlocksSteals,
  Turnovers,
  PointsAssists,
  PointsRebounds,
  ReboundsAssists,
  PointsReboundsAssists,
  Fantasy,
  ThreeFGM,
  ThreeFGA,
  FreeThrowsAttempted,
  FreeThrowsMade,
  Minutes,
];

export const BETTING_CATEGORIES = [
  Points,
  Rebounds,
  Assists,
  PointsAssists,
  PointsRebounds,
  ReboundsAssists,
  PointsReboundsAssists,
  Fantasy,
  BlocksSteals,
  Steals,
  Blocks,
  Turnovers,
  ThreeFGM,
  FreeThrowsMade,
];

export const DEFAULT_WEIGHTS = [
  { count: 0, weight: 0.25 },
  { count: -20, weight: 0.25 },
  { count: -10, weight: 0.25 },
  { count: -5, weight: 0.25 },
  // { count: 0, weight: 1 },
  // { count: -20, weight: 0 },
  // { count: -10, weight: 0 },
  // { count: -5, weight: 0 },
];

export const DEFAULT_MODEL: ModelInput = {
  model: "DEFAULT",
  gameBreakdowns: [
    {
      name: "2021-22 Season",
      filter: {
        seasons: [SeasonOption.SEASON_2021_22],
      },
      weight: 25,
    },
    {
      name: "2021-22 Matchup",
      filter: {
        seasons: [SeasonOption.SEASON_2021_22],
        opponentMatch: true,
      },
      weight: 20,
    },
    {
      name: "2022-23 Last 10",
      filter: {
        seasons: [SeasonOption.SEASON_2022_23],
        lastX: 10,
      },
      weight: 20,
    },
  ],
  similarPlayerInput: {
    limit: 10,
    statsOfInterest: [
      "Points",
      "Rebounds",
      "Assists",
      "Steals",
      "Blocks",
      "Turnovers",
      "ThreePointersMade",
    ],
    weight: 20,
    playerPoolFilter: {
      positionLooseMatch: true,
      seasons: [SeasonOption.SEASON_2022_23, SeasonOption.SEASON_2021_22],
    },
  },
  similarTeamInput: {
    limit: 3,
    statsOfInterest: ["OppPoints", "OppRebounds", "OppAssists"],
    period: {
      seasons: [SeasonOption.SEASON_2022_23, SeasonOption.SEASON_2021_22],
    },
    weight: 15,
  },
};

// OLD PROJECTIONS QUERIES
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
        date
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
        date
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
        date
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
        date
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

// let QUERY =
//     customModel.similarPlayers &&
//     customModel.similarPlayers.weight > 0 &&
//     customModel.similarTeams &&
//     customModel.similarTeams.weight > 0
//       ? GET_PROJECTIONS_AND_SIMILAR_PLAYERS_AND_TEAMS
//       : customModel.similarPlayers && customModel.similarPlayers.weight > 0
//       ? GET_PROJECTIONS_AND_SIMILAR_PLAYERS
//       : customModel.similarTeams && customModel.similarTeams.weight > 0
//       ? GET_PROJECTIONS_AND_SIMILAR_TEAMS
//       : GET_PROJECTIONS;
