import { Game } from "./graphql/game.interface";

export interface SimilarCalculation {
  similarCount: number;
  similarGames: Game[];
  similarAvg: number;
  similarAvgPerMin: number;
  similarAvgPerMinDiff: number;
  similarDiff: number;
  similarDiffPct: number;
  playerAvgAdj: number;
  countSimOver: number;
  simOverPct: number;
  simPushPct: number;
  simUnderPct: number;
  countSimUnder: number;
  countSimPush: number;
}
