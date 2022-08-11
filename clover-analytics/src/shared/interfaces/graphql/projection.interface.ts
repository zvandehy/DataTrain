import { ScoreType } from "../score-type.enum";
import { Stat } from "../stat.interface";
import { Game } from "./game.interface";
import { Player } from "./player.interface";
import { Team } from "./team.interface";

export interface Projection {
  player: Player;
  opponent: Team;
  propositions: Proposition[];
  result: Game;
  startTime: string;
  date: string;
}

export interface Proposition {
  target: number;
  type: string;
  statType: Stat;
  sportsbook: string;
  predictions: Prediction[];
  customPrediction: Prediction;
}

export interface Prediction {
  model: string;
  overUnderPrediction: string;
  confidence: number;
  totalPrediction: number;
  predictionFragments: PredictionFragment[];
  // predictionFragments: {
  //   [key in ScoreType]: PredictionFragment[];
  // };
}

export interface PredictionFragment {
  numGames: number;
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
}

export interface CountWeight {
  count: number;
  weight: number;
}
