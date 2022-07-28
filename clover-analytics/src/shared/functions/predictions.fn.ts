import { DEFAULT_WEIGHTS } from "../constants";
import { GameFilter } from "../interfaces/graphql/filters.interface";
import { Game } from "../interfaces/graphql/game.interface";
import {
  Prediction,
  PredictionFragment,
  Projection,
  Proposition,
} from "../interfaces/graphql/projection.interface";
import { GetStat } from "../interfaces/stat.interface";
import { FilterGames } from "./filters.fn";

export function CalculatePredictions(
  projections: Projection[],
  predictionFilter: GameFilter
) {
  return projections.map((projection) =>
    CalculatePrediction(
      projection,
      FilterGames(projection.player.games, predictionFilter)
    )
  );
}

// calculate the customPrediction (using UI settings) for all the propositions in a player's projection
export function CalculatePrediction(
  projection: Projection,
  games: Game[]
): Projection {
  let updatedProps: Proposition[] = [];
  projection.propositions.forEach((proposition) => {
    updatedProps.push(UpdatePropositionWithPrediction(proposition, games));
  });
  let updatedProjection = {
    ...projection,
    propositions: updatedProps,
    games: games,
  };
  return updatedProjection;
}

export function UpdatePropositionWithPrediction(
  proposition: Proposition,
  games: Game[]
): Proposition {
  let updatedProp;
  let prediction: Prediction = {
    model: "Custom",
    overUnderPrediction: "",
    confidence: 0,
    totalPrediction: 0,
    predictionFragments: [],
  };
  let skipped_weight_sum = 0;
  DEFAULT_WEIGHTS.forEach((item) => {
    const nGames = games.slice(item.count);
    //handle when player has played exactly item.count games
    if (nGames.length < item.count * -1) {
      skipped_weight_sum += item.weight;
      return;
    }
    const stat = proposition?.statType || GetStat(proposition);
    let numOver: number = 0,
      numUnder: number = 0,
      numPush: number = 0;
    nGames.forEach((game) => {
      const score = stat.score(game);
      if (score > proposition.target) numOver++;
      else if (score < proposition.target) numUnder++;
      else numPush++;
    });
    let fragment: PredictionFragment = {
      numGames: nGames.length,
      weight: item.weight,
      average: stat.average(nGames),
      median: stat.median(nGames),
      numOver: numOver,
      numUnder: numUnder,
      numPush: numPush,
      pctOver: numOver / nGames.length,
      pctUnder: numUnder / nGames.length,
    };
    prediction.predictionFragments.push(fragment);
  });
  const distributed_weight =
    skipped_weight_sum / prediction.predictionFragments.length;
  prediction.totalPrediction = +(
    prediction.predictionFragments.reduce(
      (acc, cur) => acc + (cur.weight + distributed_weight) * cur.average,
      0
    ) / DEFAULT_WEIGHTS.length
  ).toFixed(2);
  let overConfidence = +(
    prediction.predictionFragments.reduce(
      (acc, cur) => acc + (cur.weight + distributed_weight) * cur.pctOver,
      0
    ) * 100
  ).toFixed(2);
  let underConfidence = +(
    prediction.predictionFragments.reduce(
      (acc, cur) => acc + (cur.weight + distributed_weight) * cur.pctUnder,
      0
    ) * 100
  ).toFixed(2);
  prediction.confidence = +Math.max(overConfidence, underConfidence).toFixed(2);
  prediction.overUnderPrediction =
    overConfidence > underConfidence ? "Over" : "Under";
  updatedProp = {
    ...proposition,
    customPrediction: prediction,
    statType: proposition?.statType || GetStat(proposition),
  };
  return updatedProp;
}

export function PropositionScore(proposition: Proposition, game: Game): number {
  const stat = proposition?.statType || GetStat(proposition);
  return stat.score(game);
}

export function PropositionResult(
  proposition: Proposition,
  game: Game
): "Over" | "Under" | "Push" {
  const stat = proposition?.statType || GetStat(proposition);
  const score = stat.score(game);
  if (score > proposition.target) return "Over";
  if (score < proposition.target) return "Under";
  return "Push";
}

export function GetMaxConfidence(propositions: Proposition[]): Proposition {
  return propositions.sort(
    (a: Proposition, b: Proposition) =>
      b.customPrediction.confidence - a.customPrediction.confidence
  )[0];
}
