import {
  Assists,
  Blocks,
  BlocksSteals,
  Fantasy,
  FreeThrowsAttempted,
  FreeThrowsMade,
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
];
