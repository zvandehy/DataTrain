import { Stat } from "../stat.interface";
import { Player } from "./player.interface";
import { PlayersInGame } from "./playersInGame.interface";
import { Prediction, Projection } from "./projection.interface";
import { Team } from "./team.interface";
import { TeamGame } from "./teamgame.interface";

export interface Game {
  assist_percentage: number;
  assists: number;
  date: string;
  defensive_rebound_percentage: number;
  defensive_rebounds: number;
  effective_field_goal_percentage: number;
  field_goal_percentage: number;
  field_goals_attempted: number;
  field_goals_made: number;
  free_throws_attempted: number;
  free_throws_made: number;
  free_throws_percentage: number;
  gameID: string;
  home_or_away: string;
  outcome: string;
  margin: number;
  minutes: string;
  offensive_rebound_percentage: number;
  offensive_rebounds: number;
  opponent: Team;
  opponentStats: TeamGame;
  team: Team;
  teamStats: TeamGame;
  personal_fouls_drawn: number;
  personal_fouls: number;
  points: number;
  player: Player;
  playoffs: boolean;
  season: string;
  three_point_percentage: number;
  three_pointers_attempted: number;
  three_pointers_made: number;
  rebounds: number;
  true_shooting_percentage: number;
  turnovers: number;
  usage: number;
  blocks: number;
  steals: number;
  playersInGame: PlayersInGame;
  projections: Projection[];
  points_rebounds: number;
  points_assists: number;
  points_rebounds_assists: number;
  rebounds_assists: number;
  blocks_steals: number;
  fantasy_score: number;
  prediction: PredictionBreakdown;
}

export interface AverageStats {
  assists: number;
  defensive_rebounds: number;
  field_goals_attempted: number;
  field_goals_made: number;
  free_throws_attempted: number;
  free_throws_made: number;
  minutes: number;
  offensive_rebounds: number;
  personal_fouls_drawn: number;
  personal_fouls: number;
  points: number;
  three_pointers_attempted: number;
  three_pointers_made: number;
  rebounds: number;
  turnovers: number;
  usage: number;
  blocks: number;
  steals: number;
  points_rebounds: number;
  points_assists: number;
  points_rebounds_assists: number;
  rebounds_assists: number;
  blocks_steals: number;
  fantasy_score: number;
}

export interface PredictionBreakdown {
  weightedTotal: AverageStats;
  predictionAccuracy: AverageStats;
  fragments: BreakdownFragment[];
  // TODO: move these calculations to the backend
  propositions: PropositionA[];
}

export interface BreakdownFragment {
  name: string;
  weight: number;
  propositions: PropositionA[];
  derived: AverageStats;
  base: AverageStats;
  pctChange: AverageStats;
}

export interface PropositionA {
  target: number;
  type: string;
  statType?: Stat;
  sportsbook: string;
  lastModified: Date;
  analysis?: PropSummary;
  actual: number;
  actualPerMin: number;
  estimation: number;
  estimationPerMin: number;
  predictionTargetDiff: number; // Difference between prediction and target
  predictionTargetDiffPCT: number; // % Difference between prediction and target
  actualDiff: number;
  actualDiffPCT: number;
  actualDiffPerMin: number;
  actualDiffPerMinPCT: number;
  // TODO: add a significance field
  prediction: string; // OVER or UNDER
  predictionHit: string; // HIT or MISS or PUSH
}

export interface Proposition {
  target: number;
  type: string;
  statType: Stat;
  sportsbook: string;
  predictions: Prediction[];
  customPrediction: Prediction;
  lastModified: Date;
  analysis?: PropSummary;
}

export interface PropSummary {
  numOver: number;
  numUnder: number;
  numPush: number;
  pctOver: number;
  pctUnder: number;
  pctPush: number;
}
