import { SimilarCalculation } from "../similarCalculation.interface";
import { Game, Proposition } from "./game.interface";
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

export interface CountWeight {
  count: number;
  weight: number;
}
