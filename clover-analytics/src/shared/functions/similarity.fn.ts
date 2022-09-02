import { Confidence } from "../interfaces/confidence.interface";
import { Game } from "../interfaces/graphql/game.interface";
import { ScoreType } from "../interfaces/score-type.enum";
import { SimilarCalculation } from "../interfaces/similarCalculation.interface";
import { Stat } from "../interfaces/stat.interface";

export function CalculateSimilarity(
  fromCount: number,
  allGames: Game[],
  vsGames: Game[],
  stat: Stat,
  weight: number,
  target: number // either line prop or player average
): SimilarCalculation {
  const overallAverage = stat.average(allGames);
  const similarAvg = stat.average(vsGames);
  const similarAvgPerMin = stat.averagePer(vsGames, ScoreType.PerMin);
  const similarDiff = +(similarAvg - overallAverage).toFixed(2);
  const simDiffPct = similarDiff / overallAverage;
  const adjustedAverage = +(
    simDiffPct * overallAverage +
    overallAverage
  ).toFixed(2);
  let countOver: number = 0;
  let countUnder: number = 0;
  let countPush: number = 0;
  vsGames.forEach((game) => {
    const score = stat.score(game);
    // base over/under on the outcome vs player's average
    if (score > target) {
      countOver++;
    } else if (score < target) {
      countUnder++;
    } else {
      countPush++;
    }
  });
  let pctOver: number = +(countOver / vsGames.length).toFixed(2);
  let pctUnder: number = +(countUnder / vsGames.length).toFixed(2);
  let pctPush: number = +(countPush / vsGames.length).toFixed(2);
  return {
    similarCount: fromCount,
    similarGames: vsGames,
    similarAvg: similarAvg,
    similarAvgPerMin: similarAvgPerMin,
    similarDiff: similarDiff,
    similarAvgPerMinDiff: +(
      similarAvgPerMin - stat.averagePer(allGames, ScoreType.PerMin)
    ).toFixed(2),
    similarDiffPct: +(simDiffPct * 100).toFixed(2),
    playerAvgAdj: adjustedAverage,
    countSimOver: countOver,
    countSimUnder: countUnder,
    countSimPush: countPush,
    simOverPct: +(pctOver * 100).toFixed(2),
    simUnderPct: +(pctUnder * 100).toFixed(2),
    simPushPct: +(pctPush * 100).toFixed(2),
    weight: weight,
  };
}

export function AdjustConfidence(
  similarityCalculation: SimilarCalculation,
  confidence: Confidence,
  startWeight: number,
  distributedWeight: number
): Confidence {
  const adjWeight = startWeight + distributedWeight;

  confidence.over += (similarityCalculation.simOverPct / 100) * adjWeight;
  confidence.under += (similarityCalculation.simUnderPct / 100) * adjWeight;
  confidence.overOrPush +=
    (similarityCalculation.simOverPct / 100 +
      similarityCalculation.simPushPct / 100) *
    adjWeight;
  confidence.underOrPush +=
    (similarityCalculation.simUnderPct / 100 +
      similarityCalculation.simPushPct / 100) *
    adjWeight;
  return confidence;
}
