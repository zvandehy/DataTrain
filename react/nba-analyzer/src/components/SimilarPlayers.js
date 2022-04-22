import React from "react";
import { AveragePropScore, GetColor, GetPropScore } from "../utils";
import { round } from "mathjs";

export const SimilarPlayersPreview = (props) => {
  const { similar, opponent, prediction } = props;
  console.log(prediction);
  let avg = 0;
  let avgVsOpp = 0;
  let gamesOver = 0;
  let gamesPlayed = 0;
  similar.forEach((player) => {
    const playerAvg = AveragePropScore(player.games, prediction.stat.recognize);
    avg += playerAvg;

    const simPlayerVSOpp = player.games.filter(
      (game) => game.opponent.teamID === opponent
    );
    gamesPlayed += simPlayerVSOpp.length;

    if (simPlayerVSOpp.length !== 0) {
      const playerAvgVsOpp = AveragePropScore(
        player.games.filter((game) => game.opponent.teamID === opponent),
        prediction.stat.recognize
      );
      avgVsOpp += playerAvgVsOpp;
    }
    const over = simPlayerVSOpp.filter(
      (game) =>
        GetPropScore(game, prediction.stat.recognize) > prediction.target
    ).length;

    gamesOver += over;
  });
  avgVsOpp = round(avgVsOpp / similar.length, 2);
  avg = round(avg / similar.length, 2);
  const pct = round((gamesOver / gamesPlayed) * 100, 2);
  return (
    <>
      <span className="header similar-players-header">
        Similar Players ({avg})
      </span>
      <span className="similar-players-stat">{avgVsOpp}</span>
      <span className="similar-players-stat-result">{pct}%</span>
    </>
  );
};

export const SimilarPlayerRows = (props) => {
  const { similar, opponent, prediction, selected } = props;
  let similarPlayers = similar.filter(
    (player) =>
      player.games.filter((game) => game.opponent.teamID === opponent)
        .length !== 0
  );

  let avg = 0;
  let avgVsOpp = 0;
  let gamesOver = 0;
  let gamesPlayed = 0;
  let avgDiff = 0;
  let simRows = [
    <tr>
      <th>Similar Players (AVG)</th>
      <th>{selected}</th>
      <th>DIFF</th>
    </tr>,
  ];
  similarPlayers.forEach((player) => {
    const simPlayerVSOpp = player.games.filter(
      (game) => game.opponent.teamID === opponent
    );

    if (simPlayerVSOpp.length === 0) {
      return;
    }

    const playerAvg = AveragePropScore(player.games, prediction.stat.recognize);
    avg += playerAvg;

    gamesPlayed += simPlayerVSOpp.length;

    const playerAvgVsOpp = AveragePropScore(
      player.games.filter((game) => game.opponent.teamID === opponent),
      prediction.stat.recognize
    );
    const diff = round(playerAvgVsOpp - playerAvg, 2);
    avgVsOpp += playerAvgVsOpp;
    avgDiff += diff;

    const over = simPlayerVSOpp.filter(
      (game) =>
        GetPropScore(game, prediction.stat.recognize) > prediction.target
    ).length;

    gamesOver += over;

    simRows.push(
      <tr>
        <td>
          {player.name} ({playerAvg})
        </td>
        <td>
          {player.games
            .filter((game) => game.opponent.teamID === opponent)
            .map((game, i) => {
              const score = GetPropScore(game, prediction.stat.recognize);
              return (
                <span
                  className={
                    score > prediction.target
                      ? "high"
                      : score < prediction.target
                      ? "low"
                      : "med"
                  }
                >
                  {score}
                  {i <
                  player.games.filter(
                    (game) => game.opponent.teamID === opponent
                  ).length -
                    1
                    ? ","
                    : ""}{" "}
                </span>
              );
            })}
        </td>
        <td className={diff > 0 ? "high" : "low"}>
          {diff > 0 ? "+" : ""}
          {diff}
        </td>
      </tr>
    );
  });

  avgVsOpp = round(avgVsOpp / similarPlayers.length, 2);
  avg = round(avg / similarPlayers.length, 2);
  avgDiff = round(avgDiff / similarPlayers.length, 2);

  const pctDiff = round((avgDiff / avg) * 100, 2);
  console.log(avgDiff, avg, avgVsOpp, pctDiff);
  const pct = round((gamesOver / gamesPlayed) * 100, 2);
  simRows.push(
    <tr>
      <th>Average ({avg})</th>
      <td>
        <span
          className={
            avgVsOpp > prediction.target
              ? "high"
              : avgVsOpp < prediction.target
              ? "low"
              : "med"
          }
        >
          {avgVsOpp},
        </span>
        <span className={GetColor("pct", pct)}> {pct}%</span>
      </td>
      <td className={avgDiff > 0 ? "high" : "low"}>
        {avgDiff}, {avgDiff > 0 ? "+" : ""}
        {pctDiff}%
      </td>
    </tr>
  );
  return simRows;
};
