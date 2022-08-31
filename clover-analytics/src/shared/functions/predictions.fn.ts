import moment from "moment";
import { DEFAULT_WEIGHTS } from "../constants";
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
import { ScoreType } from "../interfaces/score-type.enum";
import { SimilarCalculation } from "../interfaces/similarCalculation.interface";
import { ConvertMinutes, GetStat } from "../interfaces/stat.interface";
import { FilterGames } from "./filters.fn";

export function CalculatePredictions(
  projections: Projection[],
  predictionFilter: GameFilter,
  customModel: CustomCalculation,
  games?: Game[]
) {
  return projections.map((projection) => {
    let updatedProjection = projection;
    if (projection.opponent.similarTeams?.length > 0) {
      updatedProjection = {
        ...projection,
        opponent: {
          ...projection.opponent,
          similarTeams: [
            ...projection.opponent.similarTeams,
            projection.opponent,
          ],
        },
      };
    }
    return CalculatePrediction(
      updatedProjection,
      FilterGames(games ?? projection.player.games, {
        ...predictionFilter,
        endDate: projection.date,
      }),
      customModel
    );
  });
}

// calculate the customPrediction (using UI settings) for all the propositions in a player's projection
export function CalculatePrediction(
  projection: Projection,
  games: Game[],
  customModel: CustomCalculation
): Projection {
  let updatedProps: Proposition[] = [];
  projection.propositions.forEach((proposition) => {
    let updatedProp = UpdatePropositionWithStatType(proposition);
    updatedProp = UpdatePropositionWithPrediction(
      updatedProp,
      games,
      projection,
      customModel
    );
    updatedProps.push(updatedProp);
    updatedProps = updatedProps.sort((a, b) => {
      return moment(a.lastModified).diff(b.lastModified);
    });
  });
  let updatedProjection = {
    ...projection,
    propositions: updatedProps,
    games: games,
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
  games: Game[],
  projection: Projection,
  customModel: CustomCalculation
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

  //setup recency
  customModel.recency?.forEach((item) => {
    const nGames = games.slice(item.count);

    if (nGames.length < item.count * -1 || games.length === item.count) {
      skipped_weight_sum += item.weight;
      return;
    }
    let fragment = CalculateFragment(nGames, proposition, item);
    prediction.predictionFragments.push(fragment);
  });

  const distribute_between_num =
    prediction.predictionFragments.length +
    (customModel.similarPlayers ? 1 : 0) +
    (customModel.similarTeams ? 1 : 0);
  // evenly distribute any weight that couldn't be calculated
  const distributed_weight = skipped_weight_sum / distribute_between_num;

  // calculate recency
  let overConfidence = 0;
  let underConfidence = 0;
  let overOrPushConfidence = 0;
  let underOrPushConfidence = 0;
  prediction.predictionFragments.forEach((fragment) => {
    overConfidence += (fragment.weight + distributed_weight) * fragment.pctOver;
    underConfidence +=
      (fragment.weight + distributed_weight) * fragment.pctUnder;
    overOrPushConfidence +=
      (fragment.weight + distributed_weight) * fragment.pctPushOrMore;
    underOrPushConfidence +=
      (fragment.weight + distributed_weight) * fragment.pctPushOrLess;
  });

  //similar teams
  if (customModel.similarTeams && customModel.similarTeams.weight > 0) {
    const simTeamCalc = SimilarTeamCalculation(
      projection,
      proposition,
      games ?? projection.player.games
    );

    const teamWeight = customModel.similarTeams.weight + distributed_weight;
    if (projection.player.name === "Kelsey Plum") {
      console.log(
        proposition.statType.label,
        underOrPushConfidence,
        underOrPushConfidence * (1 - teamWeight) +
          (simTeamCalc.simUnderPct / 100 + simTeamCalc.simPushPct / 100) *
            teamWeight
      );
      console.log(
        proposition.statType.label,
        overOrPushConfidence,
        overOrPushConfidence * (1 - teamWeight) +
          (simTeamCalc.simOverPct / 100 + simTeamCalc.simPushPct / 100) *
            teamWeight
      );
    }
    overConfidence += (simTeamCalc.simOverPct / 100) * teamWeight;
    underConfidence += (simTeamCalc.simUnderPct / 100) * teamWeight;
    overOrPushConfidence +=
      (simTeamCalc.simOverPct / 100 + simTeamCalc.simPushPct / 100) *
      teamWeight;
    underOrPushConfidence +=
      (simTeamCalc.simUnderPct / 100 + simTeamCalc.simPushPct / 100) *
      teamWeight;
  }

  if (customModel.similarPlayers && customModel.similarPlayers.weight > 0) {
    const simPlayerCalc = SimilarPlayerCalculation(
      projection,
      proposition,
      proposition.statType.average(games ?? projection.player.games)
    );
    const playerWeight = customModel.similarPlayers.weight + distributed_weight;

    overConfidence += (simPlayerCalc.simOverPct / 100) * playerWeight;
    underConfidence += (simPlayerCalc.simUnderPct / 100) * playerWeight;
    overOrPushConfidence +=
      (simPlayerCalc.simOverPct / 100 + simPlayerCalc.simPushPct / 100) *
      playerWeight;
    underOrPushConfidence +=
      (simPlayerCalc.simUnderPct / 100 + simPlayerCalc.simPushPct / 100) *
      playerWeight;
  }

  overConfidence = +(overConfidence * 100).toFixed(2);
  underConfidence = +(underConfidence * 100).toFixed(2);
  overOrPushConfidence = +(overOrPushConfidence * 100).toFixed(2);
  underOrPushConfidence = +(underOrPushConfidence * 100).toFixed(2);

  if (customModel.includePush) {
    prediction.confidence = +Math.max(
      overOrPushConfidence,
      underOrPushConfidence
    ).toFixed(2);
    prediction.overUnderPrediction =
      overOrPushConfidence > underOrPushConfidence ? "Over" : "Under";
  } else {
    prediction.confidence = +Math.max(overConfidence, underConfidence).toFixed(
      2
    );
    prediction.overUnderPrediction =
      overConfidence > underConfidence ? "Over" : "Under";
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
  selectedProp: Proposition,
  playerAvg: number
): SimilarCalculation {
  let simPlayerGamesVsOpp: Game[] = [];
  let simPlayerAllGames: Game[] = [];
  projection.player.similarPlayers.forEach((player) => {
    simPlayerGamesVsOpp.push(
      ...player.games.filter(
        (game) => game.opponent.teamID === projection.opponent.teamID
      )
    );
    simPlayerAllGames.push(...player.games);
  });
  let simPlayerAverageVsOpp: number =
    selectedProp.statType.average(simPlayerGamesVsOpp);
  let simPlayerAverage: number =
    selectedProp.statType.average(simPlayerAllGames);
  let simPlayerAveragePerMin: number = selectedProp.statType.averagePer(
    simPlayerAllGames,
    ScoreType.PerMin
  );
  let simPlayerVsOppAveragePerMin: number = selectedProp.statType.averagePer(
    simPlayerGamesVsOpp,
    ScoreType.PerMin
  );
  let simPlayerDifference: number = +(
    simPlayerAverageVsOpp - simPlayerAverage
  ).toFixed(5);
  let simPlayerDiffPct: number = +(
    simPlayerDifference / simPlayerAverage
  ).toFixed(2);

  let playerAvgAdj: number = +(
    playerAvg +
    playerAvg * simPlayerDiffPct
  ).toFixed(2);

  let countSimPlayerOverVsOpp: number = 0;
  let countSimPlayerUnderVsOpp: number = 0;
  let countSimPlayerPushVsOpp: number = 0;
  simPlayerGamesVsOpp.forEach((game) => {
    // base over/under on the outcome vs player's average
    if (selectedProp.statType.score(game) > simPlayerAverage) {
      countSimPlayerOverVsOpp++;
    } else if (selectedProp.statType.score(game) < simPlayerAverage) {
      countSimPlayerUnderVsOpp++;
    } else {
      countSimPlayerPushVsOpp++;
    }
  });
  let simPlayerPctOverVsOpp: number = +(
    countSimPlayerOverVsOpp / simPlayerGamesVsOpp.length
  ).toFixed(2);
  let simPlayerPctUnderVsOpp: number = +(
    countSimPlayerUnderVsOpp / simPlayerGamesVsOpp.length
  ).toFixed(2);
  let simPlayerPctPushVsOpp: number = +(
    countSimPlayerPushVsOpp / simPlayerGamesVsOpp.length
  ).toFixed(2);
  return {
    similarCount: projection.player.similarPlayers.length,
    similarGames: simPlayerGamesVsOpp,
    similarAvg: simPlayerAverage,
    similarAvgPerMin: simPlayerAveragePerMin,
    similarAvgPerMinDiff: +(
      simPlayerVsOppAveragePerMin - simPlayerAveragePerMin
    ).toFixed(2),
    similarDiff: simPlayerDifference,
    similarDiffPct: +(simPlayerDiffPct * 100).toFixed(2),
    playerAvgAdj: playerAvgAdj,
    countSimOver: countSimPlayerOverVsOpp,
    simOverPct: +(simPlayerPctOverVsOpp * 100).toFixed(2),
    simPushPct: +(simPlayerPctPushVsOpp * 100).toFixed(2),
    simUnderPct: +(simPlayerPctUnderVsOpp * 100).toFixed(2),
    countSimUnder: countSimPlayerUnderVsOpp,
    countSimPush: countSimPlayerPushVsOpp,
  };
}

export function SimilarTeamCalculation(
  projection: Projection,
  selectedProp: Proposition,
  games?: Game[]
): SimilarCalculation {
  let gamesVsSimTeams: Game[] = [];
  if (!games) {
    games = projection.player.games;
  }
  projection.opponent.similarTeams.forEach((team) => {
    gamesVsSimTeams.push(
      ...games!.filter((game) => game.opponent.teamID === team.teamID)
    );
  });
  const playerAvg = selectedProp.statType.average(games);
  let avgVsSimTeams: number = selectedProp.statType.average(gamesVsSimTeams);
  const playerAvgPerMin = selectedProp.statType.averagePer(
    gamesVsSimTeams,
    ScoreType.PerMin
  );
  let avgVsSimTeamsPerMin: number = selectedProp.statType.averagePer(
    gamesVsSimTeams,
    ScoreType.PerMin
  );
  let simTeamDifference: number = +(avgVsSimTeams - playerAvg).toFixed(2);
  let simTeamDiffPct: number = +(simTeamDifference / playerAvg).toFixed(2);

  let playerAvgAdj: number = +(playerAvg + playerAvg * simTeamDiffPct).toFixed(
    2
  );

  let countsimTeamOverVsOpp: number = 0;
  let countsimTeamUnderVsOpp: number = 0;
  let countsimTeamPushVsOpp: number = 0;
  gamesVsSimTeams.forEach((game) => {
    if (selectedProp.statType.score(game) > selectedProp.target) {
      countsimTeamOverVsOpp++;
    } else if (selectedProp.statType.score(game) < selectedProp.target) {
      countsimTeamUnderVsOpp++;
    } else {
      countsimTeamPushVsOpp++;
    }
  });
  let simTeamPctOverVsOpp: number = +(
    countsimTeamOverVsOpp / gamesVsSimTeams.length
  ).toFixed(2);
  let simTeamPctUnderVsOpp: number = +(
    countsimTeamUnderVsOpp / gamesVsSimTeams.length
  ).toFixed(2);
  let simTeamPctPushVsOpp: number = +(
    countsimTeamPushVsOpp / gamesVsSimTeams.length
  ).toFixed(2);
  return {
    similarCount: projection.opponent.similarTeams.length,
    similarGames: gamesVsSimTeams,
    similarAvg: avgVsSimTeams,
    similarAvgPerMin: playerAvgPerMin,
    similarAvgPerMinDiff: +(avgVsSimTeamsPerMin - playerAvgPerMin).toFixed(2),
    similarDiff: simTeamDifference,
    similarDiffPct: +(simTeamDiffPct * 100).toFixed(2),
    playerAvgAdj: playerAvgAdj,
    countSimOver: countsimTeamOverVsOpp,
    simOverPct: +(simTeamPctOverVsOpp * 100).toFixed(2),
    simUnderPct: +(simTeamPctUnderVsOpp * 100).toFixed(2),
    simPushPct: +(simTeamPctPushVsOpp * 100).toFixed(2),
    countSimUnder: countsimTeamUnderVsOpp,
    countSimPush: countsimTeamPushVsOpp,
  };
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
        .reduce((a, b) => a + b) / nGames.length
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
