import moment from "moment";
import { Confidence } from "../interfaces/confidence.interface";
import {
  CustomCalculation,
  Factor,
} from "../interfaces/custom-prediction.interface";
import { GameFilter } from "../interfaces/graphql/filters.interface";
import { Game, Proposition } from "../interfaces/graphql/game.interface";
import {
  Prediction,
  PredictionFragment,
  Projection,
} from "../interfaces/graphql/projection.interface";
import { ScoreType } from "../interfaces/score-type.enum";
import { SimilarCalculation } from "../interfaces/similarCalculation.interface";
import { ConvertMinutes, GetStat } from "../interfaces/stat.interface";
import { FilterGames } from "./filters.fn";
import { AdjustConfidence, CalculateSimilarity } from "./similarity.fn";

export function CalculatePredictions(
  projections: Projection[],
  gameFilter: GameFilter,
  customModel: CustomCalculation
) {
  return projections.map((projection) => {
    let updatedProjection = projection;
    if (projection.opponent.similarTeams?.length > 0) {
      updatedProjection = {
        ...projection,
        opponent: {
          ...projection.opponent,
          similarTeams: [...projection.opponent.similarTeams],
        },
      };
    }
    let c = CalculatePrediction(updatedProjection, gameFilter, customModel);
    c.propositions = c.propositions.sort(
      (prop) => prop.customPrediction.confidence
    );
    return c;
  });
}

// calculate the customPrediction (using UI settings) for all the propositions in a player's projection
export function CalculatePrediction(
  projection: Projection,
  gameFilter: GameFilter,
  customModel: CustomCalculation
): Projection {
  let updatedProps: Proposition[] = [];
  projection.propositions.forEach((proposition) => {
    let updatedProp = UpdatePropositionWithStatType(proposition);
    updatedProp = UpdatePropositionWithPrediction(
      updatedProp,
      gameFilter,
      projection,
      customModel
    );
    updatedProps.push(updatedProp);
    updatedProps = updatedProps.sort((a, b) => {
      return moment(a.lastModified).diff(b.lastModified);
    });
  });
  let updatedProjection: Projection = {
    ...projection,
    propositions: updatedProps,
  };
  return updatedProjection;
}

export function UpdatePropositionWithStatType(
  proposition: Proposition
): Proposition {
  let updatedProp = {
    ...proposition,
    statType: proposition?.statType || GetStat(proposition),
  };
  return updatedProp;
}

export function UpdatePropositionWithPrediction(
  proposition: Proposition,
  gameFilter: GameFilter,
  projection: Projection,
  customModel: CustomCalculation
): Proposition {
  let updatedProp;
  let prediction: Prediction = {
    model: "Custom",
    overUnderPrediction: "",
    confidence: 0,
    totalPrediction: 0,
    recencyFragments: [],
  };
  let skipped_weight_sum = 0;
  let filteredGames = FilterGames(projection.player.games, gameFilter).sort(
    (a, b) => {
      return moment(a.date).diff(b.date);
    }
  );
  const gamesVsOpponent = filteredGames.filter(
    (game) => game.opponent.abbreviation === projection.opponent.abbreviation
  );
  //setup recency
  customModel.recency?.forEach((item) => {
    const nGames = FilterGames(
      projection.player.games.slice(item.count! * -1),
      gameFilter
    );

    if (
      // TODO: This is skipping count of 0 (all games)
      projection.player.games.slice(item.count! * -1).length < item.count! ||
      filteredGames.length === item.count
    ) {
      skipped_weight_sum += item.weight;
      return;
    }
    let fragment = CalculateFragment(nGames, proposition, item);
    prediction.recencyFragments.push(fragment);
  });
  if (customModel.opponentWeight && gamesVsOpponent.length === 0) {
    skipped_weight_sum += customModel.opponentWeight;
  }
  const distribute_between_num =
    prediction.recencyFragments.length +
    (customModel.similarPlayers ? 1 : 0) +
    (customModel.opponentWeight && gamesVsOpponent.length > 0 ? 1 : 0) +
    (customModel.similarTeams ? 1 : 0);

  // evenly distribute any weight that couldn't be calculated
  const distributed_weight = skipped_weight_sum / distribute_between_num;

  let confidence: Confidence = {
    over: 0,
    under: 0,
    overOrPush: 0,
    underOrPush: 0,
  };
  prediction.recencyFragments.forEach((fragment, i) => {
    confidence.over +=
      (fragment.weight + distributed_weight) * fragment.pctOver;
    confidence.under +=
      (fragment.weight + distributed_weight) * fragment.pctUnder;
    confidence.overOrPush +=
      (fragment.weight + distributed_weight) * fragment.pctPushOrMore;
    confidence.underOrPush +=
      (fragment.weight + distributed_weight) * fragment.pctPushOrLess;
    prediction.recencyFragments[i].weight += distributed_weight;
  });
  //opponent
  if (
    customModel.opponentWeight &&
    customModel.opponentWeight > 0 &&
    FilterGames(gamesVsOpponent, gameFilter).length > 0
  ) {
    const opponentCalculation = CalculateSimilarity(
      1,
      filteredGames,
      gamesVsOpponent,
      proposition.statType,
      customModel.opponentWeight,
      proposition.target
    );
    confidence = AdjustConfidence(
      opponentCalculation,
      confidence,
      customModel.opponentWeight,
      distributed_weight
    );
    opponentCalculation.weight =
      customModel.opponentWeight + distributed_weight;
    prediction.vsOpponent = opponentCalculation;
  }

  //similar teams
  if (
    customModel.similarTeams &&
    customModel.similarTeams.weight > 0 &&
    projection.opponent.similarTeams &&
    FilterGames(
      filteredGames.filter((game) =>
        projection.opponent.similarTeams
          .map((team) => team.abbreviation)
          .includes(game.opponent.abbreviation)
      ),
      gameFilter
    ).length > 0
  ) {
    const simTeamCalc = SimilarTeamCalculation(
      projection,
      proposition,
      customModel.similarTeams.weight,
      filteredGames
    );
    confidence = AdjustConfidence(
      simTeamCalc,
      confidence,
      customModel.similarTeams.weight,
      distributed_weight
    );
    simTeamCalc.weight = customModel.similarTeams.weight + distributed_weight;
    prediction.vsSimilarTeams = simTeamCalc;
  }

  if (
    customModel.similarPlayers &&
    customModel.similarPlayers.weight > 0 &&
    FilterGames(
      projection.player.similarPlayers?.map((player) => player.games).flat(),
      gameFilter
    ).length > 0
  ) {
    const simPlayerCalc = SimilarPlayerCalculation(
      projection,
      gameFilter,
      proposition,
      customModel.similarPlayers.weight
    );
    confidence = AdjustConfidence(
      simPlayerCalc,
      confidence,
      customModel.similarPlayers.weight,
      distributed_weight
    );
    simPlayerCalc.weight =
      customModel.similarPlayers.weight + distributed_weight;
    prediction.similarPlayersVsOpponent = simPlayerCalc;
  }

  confidence.over = +confidence.over.toFixed(2);
  confidence.under = +confidence.under.toFixed(2);
  confidence.overOrPush = +confidence.overOrPush.toFixed(2);
  confidence.underOrPush = +confidence.underOrPush.toFixed(2);

  if (customModel.includePush) {
    prediction.confidence = +Math.max(
      confidence.overOrPush,
      confidence.underOrPush
    ).toFixed(2);
    prediction.overUnderPrediction =
      confidence.overOrPush > confidence.underOrPush ? "Over" : "Under";
  } else {
    prediction.confidence = +Math.max(
      confidence.over,
      confidence.under
    ).toFixed(2);
    prediction.overUnderPrediction =
      confidence.over > confidence.under ? "Over" : "Under";
  }
  updatedProp = {
    ...proposition,
    customPrediction: prediction,
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

export function SimilarPlayerCalculation(
  projection: Projection,
  gameFilter: GameFilter,
  selectedProp: Proposition,
  weight: number
): SimilarCalculation {
  let calculation: SimilarCalculation = {
    similarCount: 0,
    similarGames: [],
    similarAvg: 0,
    similarAvgPerMin: 0,
    similarAvgPerMinDiff: 0,
    similarDiff: 0,
    similarDiffPct: 0,
    countSimOver: 0,
    countSimUnder: 0,
    countSimPush: 0,
    simOverPct: 0,
    simPushPct: 0,
    simUnderPct: 0,
    playerAvgAdj: 0,
    weight: weight,
  };
  let allGames: Game[] = [];
  projection.player.similarPlayers.forEach((similarPlayer) => {
    const filteredGames = FilterGames(similarPlayer.games, gameFilter);
    const targetFromGames = selectedProp.statType.average(filteredGames);
    const simVsGames = filteredGames.filter(
      (game) => game.opponent.abbreviation === projection.opponent.abbreviation
    );
    if (simVsGames.length > 0) {
      calculation.similarCount++;
      allGames.push(...filteredGames);
    }
    simVsGames.forEach((game) => {
      const gameScore = selectedProp.statType.score(game);
      if (gameScore > targetFromGames) {
        calculation.countSimOver++;
      } else if (gameScore < targetFromGames) {
        calculation.countSimUnder++;
      } else {
        calculation.countSimPush++;
      }
      calculation.similarDiff += +(gameScore - targetFromGames).toFixed(1);
      calculation.similarAvgPerMinDiff +=
        selectedProp.statType.scorePer(game, ScoreType.PerMin) -
        selectedProp.statType.averagePer(filteredGames, ScoreType.PerMin);
    });
    calculation.similarGames.push(...simVsGames);
  });
  calculation.simOverPct = +(
    (calculation.countSimOver / calculation.similarGames.length) *
    100
  ).toFixed(2);
  calculation.simUnderPct = +(
    (calculation.countSimUnder / calculation.similarGames.length) *
    100
  ).toFixed(2);
  calculation.simPushPct = +(
    (calculation.countSimPush / calculation.similarGames.length) *
    100
  ).toFixed(2);
  calculation.similarDiff = +(
    calculation.similarDiff / calculation.similarGames.length
  ).toFixed(2);
  calculation.similarAvg = selectedProp.statType.average(allGames);
  calculation.similarAvgPerMin = selectedProp.statType.averagePer(
    allGames,
    ScoreType.PerMin
  );
  calculation.similarDiffPct = +(
    calculation.similarDiff / calculation.similarAvg
  ).toFixed(2);
  const similarDiffPerMinPct = +(
    calculation.similarAvgPerMinDiff / calculation.similarAvgPerMin
  ).toFixed(2);
  const playerAvg = selectedProp.statType.average(
    FilterGames(projection.player.games, gameFilter)
  );
  calculation.playerAvgAdj = +(
    playerAvg +
    calculation.similarDiffPct * playerAvg
  ).toFixed(2);
  return calculation;
}

export function SimilarTeamCalculation(
  projection: Projection,
  selectedProp: Proposition,
  weight: number,
  filteredGames: Game[]
): SimilarCalculation {
  const vsGames: Game[] = filteredGames.filter((game) =>
    projection.opponent.similarTeams
      .map((team) => team.abbreviation)
      .includes(game.opponent.abbreviation)
  );
  return CalculateSimilarity(
    projection.opponent.similarTeams.length,
    filteredGames,
    vsGames,
    selectedProp.statType,
    weight,
    selectedProp.target
  );
}

export function CalculateFragment(
  nGames: Game[],
  proposition: Proposition,
  item: Factor
): PredictionFragment {
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
    games: nGames,
    weight: item.weight,
    average: stat.average(nGames),
    avgPerMin: stat.averagePer(nGames, ScoreType.PerMin),
    minutes: +(
      nGames
        .map((game) => ConvertMinutes(game.minutes))
        .reduce((a, b) => a + b, 0) / nGames.length
    ).toFixed(2),
    median: stat.median(nGames),
    numOver: numOver,
    numUnder: numUnder,
    numPush: numPush,

    pctOver: numOver / nGames.length,
    pctUnder: numUnder / nGames.length,
    pctPushOrMore: (numOver + numPush) / nGames.length,
    pctPushOrLess: (numUnder + numPush) / nGames.length,
  };
  return fragment;
}
