import {
  GameFilter,
  Period,
  PlayerFilter,
  TeamFilter,
} from "./graphql/filters.interface";

export interface CustomCalculation {
  includePush: boolean;
  includeOnDifferentTeam: boolean;
  recency?: Factor[];
  recencyPct?: Factor[];
  similarPlayers?: Factor;
  similarTeams?: Factor;
  homeAwayWeight?: number;
  opponentWeight?: number;
  playoffs?: Factor;
  // hitCriteria: HitCriteria;
}

export interface Factor {
  weight: number;
  count?: number;
  seasons?: string[];
}

export interface ModelInput {
  model: string;
  gameBreakdowns: GameBreakdownInput[];
  similarPlayerInput?: SimilarPlayerInput;
  similarTeamInput?: SimilarTeamInput;
}

export interface GameBreakdownInput {
  name: string;
  filter: GameFilter;
  weight: number;
}

export interface SimilarPlayerInput {
  limit: number;
  statsOfInterest: string[]; // TODO: Maybe this should be an enum?
  playerPoolFilter?: PlayerFilter;
  weight: number;
}

export interface SimilarTeamInput {
  limit: number;
  statsOfInterest: string[]; //TODO: Maybe enum?
  teamPoolFilter?: TeamFilter[];
  period: Period;
  weight: number;
}
