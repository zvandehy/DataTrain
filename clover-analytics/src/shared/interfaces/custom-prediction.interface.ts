export interface CustomCalculation {
  recency?: Factor[];
  recencyPct?: Factor[];
  similarPlayers?: Factor;
  similarTeams?: Factor;
  homeAwayWeight?: number;
  includePush: boolean;
  opponentWeight?: number;
}

export interface Factor {
  weight: number;
  count: number;
}
