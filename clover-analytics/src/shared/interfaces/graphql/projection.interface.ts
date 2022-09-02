import { ScoreType } from "../score-type.enum";
import { SimilarCalculation } from "../similarCalculation.interface";
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
  lastModified: Date;
}

export interface Prediction {
  model: string;
  overUnderPrediction: string;
  confidence: number;
  totalPrediction: number;
  recencyFragments: PredictionFragment[];
  vsOpponent?: SimilarCalculation;
  vsSimilarTeams?: SimilarCalculation;
  similarPlayersVsOpponent?: SimilarCalculation;
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
  pctPushOrMore: number;
  pctPushOrLess: number;
}

export interface CountWeight {
  count: number;
  weight: number;
}
