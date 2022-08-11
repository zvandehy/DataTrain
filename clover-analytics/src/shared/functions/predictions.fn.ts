import { DEFAULT_WEIGHTS } from "../constants";
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
      })
    );
  });
}

// calculate the customPrediction (using UI settings) for all the propositions in a player's projection
export function CalculatePrediction(
  projection: Projection,
  games: Game[]
): Projection {
  let updatedProps: Proposition[] = [];
  projection.propositions.forEach((proposition) => {
    let updatedProp = UpdatePropositionWithStatType(proposition);
    updatedProp = UpdatePropositionWithPrediction(
      updatedProp,
      games,
      projection
    );
    // if (projection.player.similarPlayers?.length > 0) {
    //   updatedProp = UpdatePropositionWithSimilarCalculation(
    //     SimilarPlayerCalculation(projection, updatedProp),
    //     updatedProp
    //   );
    // }
    // if (projection.opponent.similarTeams?.length > 0) {
    //   updatedProp = UpdatePropositionWithSimilarCalculation(
    //     SimilarTeamCalculation(projection, updatedProp),
    //     updatedProp
    //   );
    // }
    updatedProps.push(updatedProp);
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

// EXPERIMENTAL COUNT PREDICTION
interface CountPrediction {
  games: number[];
  opponent: number[];
  similarPlayerAdj: number;
  similarTeamAdj: number;
  similarPlayerPct: number;
  similarTeamPct: number;
}

export function UpdatePropositionWithPrediction(
  proposition: Proposition,
  games: Game[],
  projection: Projection
): Proposition {
  let countPrediction: CountPrediction = {
    games: [],
    opponent: [],
    similarPlayerAdj: 0,
    similarTeamAdj: 0,
    similarPlayerPct: 0,
    similarTeamPct: 0,
  };
  let updatedProp;
  let prediction: Prediction = {
    model: "Custom",
    overUnderPrediction: "",
    confidence: 0,
    totalPrediction: 0,
    predictionFragments: [],
    // predictionFragments: {
    //   [ScoreType.PerGame]: [],
    //   [ScoreType.PerMin]: [],
    //   [ScoreType.Per36Min]: [],
    // },
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

      //TODO: add option to include numPush as either over, under, or both
      pctOver: numOver / nGames.length,
      pctUnder: numUnder / nGames.length,
    };
    prediction.predictionFragments.push(fragment);

    if (fragment.pctOver > fragment.pctUnder) {
      countPrediction.games.push(1);
    } else {
      countPrediction.games.push(-1);
    }
  });
  // COUNT PREDICTION EXPERIEMENT
  games
    .filter((game) => game.opponent.teamID === projection.opponent.teamID)
    .forEach((game) => {
      countPrediction.opponent.push(
        proposition.statType.score(game) > proposition.target ? 1 : -1
      );
    });

  const distributed_weight =
    skipped_weight_sum / prediction.predictionFragments.length;
  prediction.totalPrediction = +(
    prediction.predictionFragments.reduce(
      (acc, cur) => acc + (cur.weight + distributed_weight) * cur.average,
      0
    ) / DEFAULT_WEIGHTS.length
  ).toFixed(2);
  let overConfidence = 0;
  let underConfidence = 0;
  prediction.predictionFragments.forEach((fragment) => {
    overConfidence += (fragment.weight + distributed_weight) * fragment.pctOver;
    underConfidence +=
      (fragment.weight + distributed_weight) * fragment.pctUnder;
  });
  overConfidence = +(overConfidence * 100).toFixed(2);
  underConfidence = +(underConfidence * 100).toFixed(2);

  //TODO: may have to use a different condition to calculate confidence from similarity
  if (projection.player.similarPlayers?.length > 0) {
    const simPlayerCalc = SimilarPlayerCalculation(
      projection,
      proposition,
      proposition.statType.average(games ?? projection.player.games)
    );
    if (projection.opponent.similarTeams?.length > 0) {
      const simTeamCalc = SimilarTeamCalculation(
        projection,
        proposition,
        games ?? projection.player.games
      );
      // if (projection.player.name === "Diana Taurasi") {
      //   console.log(
      //     `${overConfidence}*.06 + ${simTeamCalc.simOverPct}*.02 + ${
      //       simPlayerCalc.simOverPct
      //     }*0.2 = ${
      //       overConfidence * 0.6 +
      //       simTeamCalc.simOverPct * 0.2 +
      //       simPlayerCalc.simOverPct * 0.2
      //     }`
      //   );
      //   console.log(
      //     `${underConfidence}*.06 + ${simTeamCalc.simUnderPct}*.02 + ${
      //       simPlayerCalc.simUnderPct
      //     }*0.2 = ${
      //       underConfidence * 0.6 +
      //       simTeamCalc.simUnderPct * 0.2 +
      //       simPlayerCalc.simUnderPct * 0.2
      //     }`
      //   );
      // }
      overConfidence =
        overConfidence * 0.6 +
        simTeamCalc.simOverPct * 0.2 +
        simPlayerCalc.simOverPct * 0.2;
      underConfidence =
        underConfidence * 0.6 +
        simTeamCalc.simUnderPct * 0.2 +
        simPlayerCalc.simUnderPct * 0.2;

      // COUNT PREDICTION EXPERIMENT
      countPrediction.similarPlayerAdj =
        simPlayerCalc.playerAvgAdj > proposition.target ? 1 : -1;
      countPrediction.similarTeamAdj =
        simTeamCalc.playerAvgAdj > proposition.target ? 1 : -1;
      countPrediction.similarPlayerPct = simPlayerCalc.simOverPct > 50 ? 1 : -1;
      countPrediction.similarTeamPct = simTeamCalc.simOverPct > 50 ? 1 : -1;
    } else {
      overConfidence = overConfidence * 0.6 + simPlayerCalc.simOverPct * 0.4;
      underConfidence = underConfidence * 0.6 + simPlayerCalc.simUnderPct * 0.4;

      // COUNT PREDICTION EXPERIMENT
      countPrediction.similarPlayerAdj =
        simPlayerCalc.playerAvgAdj > proposition.target ? 1 : -1;
      countPrediction.similarPlayerPct = simPlayerCalc.simOverPct > 50 ? 1 : -1;
    }
  } else if (projection.opponent.similarTeams?.length > 0) {
    const simTeamCalc = SimilarTeamCalculation(
      projection,
      proposition,
      games ?? projection.player.games
    );
    overConfidence = overConfidence * 0.6 + simTeamCalc.simOverPct * 0.6;
    underConfidence = underConfidence * 0.6 + simTeamCalc.simUnderPct * 0.6;

    // COUNT PREDICTION EXPERIMENT
    countPrediction.similarTeamAdj =
      simTeamCalc.playerAvgAdj > proposition.target ? 1 : -1;
    countPrediction.similarTeamPct = simTeamCalc.simOverPct > 50 ? 1 : -1;
  }

  prediction.confidence = +Math.max(overConfidence, underConfidence).toFixed(2);
  prediction.overUnderPrediction =
    overConfidence > underConfidence ? "Over" : "Under";

  let predictionCountOver = countPrediction.games.filter(
    (val) => val > 0
  ).length;
  predictionCountOver += countPrediction.opponent.filter(
    (val) => val > 0
  ).length;
  predictionCountOver += countPrediction.similarPlayerAdj > 0 ? 1 : 0;
  predictionCountOver += countPrediction.similarTeamAdj > 0 ? 1 : 0;
  predictionCountOver += countPrediction.similarTeamPct > 0 ? 1 : 0;
  predictionCountOver += countPrediction.similarPlayerPct > 0 ? 1 : 0;

  let predictionCountUnder = countPrediction.games.filter(
    (val) => val < 0
  ).length;
  predictionCountUnder += countPrediction.opponent.filter(
    (val) => val < 0
  ).length;
  predictionCountUnder += countPrediction.similarPlayerAdj < 0 ? 1 : 0;
  predictionCountUnder += countPrediction.similarTeamAdj < 0 ? 1 : 0;
  predictionCountUnder += countPrediction.similarTeamPct < 0 ? 1 : 0;
  predictionCountUnder += countPrediction.similarPlayerPct < 0 ? 1 : 0;

  const predictionCountTotal =
    countPrediction.games.length +
    countPrediction.opponent.length +
    Math.abs(countPrediction.similarPlayerAdj) +
    Math.abs(countPrediction.similarTeamPct) +
    Math.abs(countPrediction.similarPlayerPct) +
    Math.abs(countPrediction.similarTeamAdj);

  // DISABLE COUNT PREDICTION
  // if (predictionCountOver > predictionCountUnder) {
  //   prediction.overUnderPrediction = "Over";
  //   prediction.confidence = +(
  //     (predictionCountOver / predictionCountTotal) *
  //     100
  //   ).toFixed(2);
  // } else {
  //   prediction.overUnderPrediction = "Under";
  //   prediction.confidence = +(
  //     (predictionCountUnder / predictionCountTotal) *
  //     100
  //   ).toFixed(2);
  // }

  // console.group(projection.player.name);
  // console.log(
  //   proposition.type,
  //   predictionCountOver,
  //   predictionCountUnder,
  //   predictionCountTotal,
  //   prediction.confidence,
  //   prediction.overUnderPrediction,
  //   countPrediction
  // );
  // console.groupEnd();
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
    // if (selectedProp.statType.scorePer(game, DEFAULT_SCORE_TYPE) > selectedProp.target) {
    //   countSimPlayerOverVsOpp++;
    // } else if (selectedProp.statType.scorePer(game, DEFAULT_SCORE_TYPE) < selectedProp.target) {
    //   countSimPlayerUnderVsOpp++;
    // } else {
    //   countSimPlayerPushVsOpp++;
    // }
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
    countSimUnder: countsimTeamUnderVsOpp,
    countSimPush: countsimTeamPushVsOpp,
  };
}
