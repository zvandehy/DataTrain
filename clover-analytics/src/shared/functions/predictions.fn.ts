import moment from "moment";
import { Confidence } from "../interfaces/confidence.interface";
import {
  CustomCalculation,
  Factor,
} from "../interfaces/custom-prediction.interface";
import { GameFilter } from "../interfaces/graphql/filters.interface";
import { Game } from "../interfaces/graphql/game.interface";
import {
  Prediction,
  PredictionFragment,
  Projection,
  Proposition,
} from "../interfaces/graphql/projection.interface";
import { ListFilterOptions } from "../interfaces/listFilter.interface";
import { ScoreType } from "../interfaces/score-type.enum";
import { SimilarCalculation } from "../interfaces/similarCalculation.interface";
import {
  ConvertMinutes,
  GetStat,
  Minutes,
  Points,
  ReboundsAssists,
} from "../interfaces/stat.interface";
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
    return CalculatePrediction(updatedProjection, gameFilter, customModel);
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
      projection.player.games.slice(item.count),
      gameFilter
    );

    if (
      projection.player.games.slice(item.count).length < item.count * -1 ||
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

  confidence.over = +(confidence.over * 100).toFixed(2);
  confidence.under = +(confidence.under * 100).toFixed(2);
  confidence.overOrPush = +(confidence.overOrPush * 100).toFixed(2);
  confidence.underOrPush = +(confidence.underOrPush * 100).toFixed(2);

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
  // TODO: Verify that results / future games are not used in the calculation
  let overallOver = 0;
  let overallUnder = 0;
  let overallPush = 0;
  let overallDiff = 0;
  let overallDiffPerMin = 0;
  let similarCount = 0;
  let allGames: Game[] = [];
  let vsGames: Game[] = [];
  projection.player.similarPlayers.forEach((similarPlayer) => {
    const filteredGames = FilterGames(similarPlayer.games, gameFilter);
    const simVsGames = filteredGames.filter(
      (game) => game.opponent.abbreviation === projection.opponent.abbreviation
    );
    if (simVsGames.length > 0) {
      const similarity = CalculateSimilarity(
        1,
        filteredGames,
        simVsGames,
        selectedProp.statType,
        weight,
        selectedProp.statType.average(filteredGames)
      );
      overallOver += similarity.countSimOver;
      overallUnder += similarity.countSimUnder;
      overallPush += similarity.countSimPush;
      overallDiff += similarity.similarDiff;
      overallDiffPerMin += similarity.similarAvgPerMinDiff;
      allGames.push(...filteredGames);
      vsGames.push(...simVsGames);
      similarCount++;
    }
  });
  let average = selectedProp.statType.average(
    FilterGames(projection.player.games, gameFilter)
  );
  let compoundedSimilarity = CalculateSimilarity(
    similarCount,
    allGames,
    vsGames,
    selectedProp.statType,
    weight,
    selectedProp.statType.average(allGames)
  );
  compoundedSimilarity.countSimOver = overallOver;
  compoundedSimilarity.countSimUnder = overallUnder;
  compoundedSimilarity.countSimPush = overallPush;
  compoundedSimilarity.simOverPct = +(
    (overallOver / vsGames.length) *
    100
  ).toFixed(2);
  compoundedSimilarity.simUnderPct = +(
    (overallUnder / vsGames.length) *
    100
  ).toFixed(2);
  compoundedSimilarity.simPushPct = +(
    (overallPush / vsGames.length) *
    100
  ).toFixed(2);
  compoundedSimilarity.similarDiff = +overallDiff.toFixed(2);
  compoundedSimilarity.similarAvgPerMinDiff = +overallDiffPerMin.toFixed(2);
  compoundedSimilarity.similarDiffPct = +(
    (overallDiff / average) *
    100
  ).toFixed(2);
  compoundedSimilarity.playerAvgAdj = +(
    (compoundedSimilarity.similarDiffPct / 100) * average +
    average
  ).toFixed(2);
  return compoundedSimilarity;
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
    numGames: nGames.length,
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
