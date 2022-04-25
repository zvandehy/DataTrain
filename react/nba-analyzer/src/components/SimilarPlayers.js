import React from "react";
import {
  AveragePropScore,
  GetColor,
  GetPropScore,
  RelevantStats,
} from "../utils";
import { round } from "mathjs";

export const SimilarPlayersPreview = (props) => {
  const { similar, opponent, prediction } = props;
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
  const { similar, opponent, prediction, selected, average } = props;
  //don't include players without games vs opponent
  let similarPlayers = similar.filter(
    (player) =>
      player.games.filter((game) => game.opponent.teamID === opponent.teamID)
        .length !== 0
  );

  let avg = 0;
  let avgVsOpp = 0;
  let gamesOver = 0;
  let gamesPlayed = 0;
  let avgDiff = 0;
  let simRows = [
    <tr>
      <th>Similar Players (AVG) vs {opponent.abbreviation}</th>
      {RelevantStats[selected.recognize].map((item, i) => (
        <th key={item.label}>{item.label}</th>
      ))}
    </tr>,
  ];
  similarPlayers.forEach((player) => {
    const simPlayerVSOpp = player.games.filter(
      (game) => game.opponent.teamID === opponent.teamID
    );

    if (simPlayerVSOpp.length === 0) {
      return;
    }

    const playerAvg = AveragePropScore(player.games, prediction.stat.recognize);
    avg += playerAvg;

    gamesPlayed += simPlayerVSOpp.length;

    const playerAvgVsOpp = AveragePropScore(
      player.games.filter((game) => game.opponent.teamID === opponent.teamID),
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
        {RelevantStats[selected.recognize].map((stat, i) => {
          const cellTarget = AveragePropScore(player.games, stat.recognize);
          const games = player.games.filter(
            (game) => game.opponent.teamID === opponent.teamID
          );
          const score = AveragePropScore(games, stat.recognize);
          return (
            <td
              className={
                score > cellTarget ? "high" : score < cellTarget ? "low" : "med"
              }
            >
              {score}
            </td>
          );
        })}
      </tr>
    );
  });

  avgVsOpp = round(avgVsOpp / similarPlayers.length, 2);
  avg = round(avg / similarPlayers.length, 2);
  avgDiff = round(avgDiff / similarPlayers.length, 2);

  const pctDiff = round((avgDiff / avg) * 100, 2);
  const pct = round((gamesOver / gamesPlayed) * 100, 2);
  simRows.push(
    <>
      <tr>
        <th className={"right"}>Average ({avg})</th>
        <td
          className={avgVsOpp > avg ? "high" : avgVsOpp < avg ? "low" : "med"}
        >
          {avgVsOpp}
        </td>
      </tr>
      <tr>
        <th className="right">
          % Over ({gamesOver} / {gamesPlayed})
        </th>
        <td className={GetColor("pct", pct)}>{pct}%</td>
      </tr>
      <tr>
        <th className={"right"}>Difference</th>
        <td className={avgDiff > 0 ? "high" : avgDiff < 0 ? "low" : "med"}>
          {avgDiff > 0 ? "+" : ""}
          {avgDiff}, {avgDiff > 0 ? "+" : ""}
          {pctDiff}%
        </td>
      </tr>
      <tr>
        <th className={"right"}>
          {average}*{pctDiff}%
        </th>
        <td
          className={
            average * (1 + pctDiff / 100) > prediction.target ? "high" : "low"
          }
        >
          {round(average * (1 + pctDiff / 100), 2)}
        </td>
      </tr>
    </>
  );
  return simRows;
};
