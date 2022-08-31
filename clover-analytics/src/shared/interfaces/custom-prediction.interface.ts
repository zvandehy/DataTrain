export interface CustomCalculation {
  recency?: Factor[];
  recencyPct?: Factor[];
  similarPlayers?: Factor;
  similarTeams?: Factor;
  homeAwayWeight?: number;
  includePush: boolean;
}

export interface Factor {
  weight: number;
  count: number;
}
