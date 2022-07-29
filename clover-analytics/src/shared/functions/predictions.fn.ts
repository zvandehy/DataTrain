import ProjectionsSummary from "../../ components/projections-summary/projections-summary.component";
import { DEFAULT_WEIGHTS } from "../constants";
import { GameFilter } from "../interfaces/graphql/filters.interface";
import { Game } from "../interfaces/graphql/game.interface";
import {
  Prediction,
  PredictionFragment,
  Projection,
  Proposition,
} from "../interfaces/graphql/projection.interface";
import { SimilarCalculation } from "../interfaces/similarCalculation.interface";
import { GetStat } from "../interfaces/stat.interface";
import { FilterGames } from "./filters.fn";

export function CalculatePredictions(
  projections: Projection[],
  predictionFilter: GameFilter
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
      FilterGames(projection.player.games, predictionFilter)
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

export function UpdatePropositionWithPrediction(
  proposition: Proposition,
  games: Game[],
  projection: Projection
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

  //TODO: may have to use a different condition to calculate confidence from similarity
  if (projection.player.similarPlayers?.length > 0) {
    const simPlayerCalc = SimilarPlayerCalculation(projection, proposition);
    if (projection.opponent.similarTeams?.length > 0) {
      const simTeamCalc = SimilarTeamCalculation(projection, proposition);
      overConfidence =
        overConfidence * 0.5 +
        simTeamCalc.simOverPct * 0.25 +
        simPlayerCalc.simOverPct * 0.25;
      underConfidence =
        underConfidence * 0.5 +
        simTeamCalc.simUnderPct * 0.25 +
        simPlayerCalc.simUnderPct * 0.25;
    } else {
      overConfidence = overConfidence * 0.5 + simPlayerCalc.simOverPct * 0.5;
      underConfidence = underConfidence * 0.5 + simPlayerCalc.simUnderPct * 0.5;
    }
  } else if (projection.opponent.similarTeams?.length > 0) {
    const simTeamCalc = SimilarTeamCalculation(projection, proposition);
    overConfidence = overConfidence * 0.5 + simTeamCalc.simOverPct * 0.5;
    underConfidence = underConfidence * 0.5 + simTeamCalc.simUnderPct * 0.5;
  }

  prediction.confidence = +Math.max(overConfidence, underConfidence).toFixed(2);
  prediction.overUnderPrediction =
    overConfidence > underConfidence ? "Over" : "Under";
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

// export function UpdatePropositionWithSimilarCalculation(
//   sim: SimilarCalculation,
//   proposition: Proposition
// ): Proposition {
//   //TODO: revisit the best way to do this
//   let overConfidence = proposition.customPrediction.
//   if (sim.simOverPct > sim.simUnderPct) {
//     if(proposition.customPrediction.overUnderPrediction === "Over") {
//       updatedConfidence = +(
//         proposition.customPrediction.confidence * 0.5 +
//         sim.simOverPct * 0.5
//       ).toFixed(2);
//   }
//   let updatedConfidence = proposition.customPrediction.confidence;
//   if (
//     (sim.similarDiff > 0 &&
//       proposition.customPrediction.overUnderPrediction === "Over") ||
//     (sim.similarDiff < 0 &&
//       proposition.customPrediction.overUnderPrediction === "Under")
//   ) {

//   } else {
//     updatedConfidence = +(
//       proposition.customPrediction.confidence * 0.5 +
//       (1 - sim.simOverPct) * 0.5
//     ).toFixed(2);
//   }
//   //TODO: need to reverse the confidence if the prediction is opposite
//   let updatedProp: Proposition = {
//     ...proposition,
//     customPrediction: {
//       ...proposition.customPrediction,
//       confidence: updatedConfidence,
//     },
//   };
//   return updatedProp;
// }

export function SimilarPlayerCalculation(
  projection: Projection,
  selectedProp: Proposition
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
  let simPlayerDifference: number = +(
    simPlayerAverageVsOpp - simPlayerAverage
  ).toFixed(2);
  let simPlayerDiffPct: number = +(
    simPlayerDifference / simPlayerAverage
  ).toFixed(2);

  let playerAvgAdj: number = +(
    selectedProp.statType.average(projection.player.games) +
    selectedProp.statType.average(projection.player.games) * simPlayerDiffPct
  ).toFixed(2);

  let countSimPlayerOverVsOpp: number = 0;
  let countSimPlayerUnderVsOpp: number = 0;
  let countSimPlayerPushVsOpp: number = 0;
  simPlayerGamesVsOpp.forEach((game) => {
    // if (selectedProp.statType.score(game) > selectedProp.target) {
    //   countSimPlayerOverVsOpp++;
    // } else if (selectedProp.statType.score(game) < selectedProp.target) {
    //   countSimPlayerUnderVsOpp++;
    // } else {
    //   countSimPlayerPushVsOpp++;
    // }
    if (selectedProp.statType.score(game) > playerAvgAdj) {
      countSimPlayerOverVsOpp++;
    } else if (selectedProp.statType.score(game) < playerAvgAdj) {
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
  selectedProp: Proposition
): SimilarCalculation {
  let gamesVsSimTeams: Game[] = [];
  projection.opponent.similarTeams.forEach((team) => {
    gamesVsSimTeams.push(
      ...projection.player.games.filter(
        (game) => game.opponent.teamID === team.teamID
      )
    );
  });
  let avgVsSimTeams: number = selectedProp.statType.average(gamesVsSimTeams);
  let playerAvg: number = selectedProp.statType.average(
    projection.player.games
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
