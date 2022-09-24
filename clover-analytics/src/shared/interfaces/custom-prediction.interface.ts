import { HitCriteria } from "./accuracy.interface";

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
  hitCriteria: HitCriteria;
}

export interface Factor {
  weight: number;
  count?: number;
  seasons?: string[];
}
