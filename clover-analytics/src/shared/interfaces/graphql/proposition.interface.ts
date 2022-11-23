import { abs } from "mathjs";
import { Game } from "./game.interface";

export interface Proposition {
  target: number;
  type: string;
  sportsbook: string;
  prediction: PropPrediction;
  lastModified: Date;
  game: Game;
  outcome: "OVER" | "UNDER" | "PUSH";
  actualResult: number;
  // # how close the proposition was to the actual outcome
  accuracy: number;
}

export interface PropPrediction {
  model: string;
  estimation: number;
  estimationAccuracy?: number;
  significance: number;
  stdDev: number;
  cumulativeOver: number;
  cumulativeUnder: number;
  cumulativePush: number;
  cumulativeOverPct: number;
  cumulativeUnderPct: number;
  cumulativePushPct: number;
  wager: "OVER" | "UNDER";
  wagerOutcome: "HIT" | "MISS" | "PUSH" | "PENDING";
  breakdowns: PropBreakdown[];
}

export function ComparePropByPredictionDeviation(
  a: Proposition,
  b: Proposition
) {
  return GetPropPredictionDeviation(a) - GetPropPredictionDeviation(b);
}

export function GetPropPredictionDeviation(prop: Proposition): number {
  return +prop.prediction?.significance.toFixed(2);
  // return (
  //   +(
  //     abs((prop.prediction.estimation - prop.target) / prop.target) * 100
  //   ).toFixed(2) || 0
  // );
}

export interface PropBreakdown {
  name?: string;
  over?: number;
  under?: number;
  push?: number;
  overPct?: number;
  underPct?: number;
  pushPct?: number;
  derivedAverage?: number;
  weight?: number;
  pctChange?: number;
  base?: number;
  derivedGames?: Game[];
}

export interface PredictionFragment {
  games: Game[];
  minutes: number;
  avgPerMin: number;
  weight: number;
  average: number;
  median: number;
  numOver: number;
  numUnder: number;
  numPush: number;
  pctOver: number;
  pctUnder: number;
  pctPushOrMore: number;
  pctPushOrLess: number;
}
