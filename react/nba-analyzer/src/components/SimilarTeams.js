import React from "react";
import {
  AveragePropScore,
  GetColor,
  GetPropScore,
  RelevantStats,
} from "../utils";
import { round } from "mathjs";

//TODO: Not implmented
export const SimilarTeamsPreview = (props) => {
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

export const SimilarTeamRows = (props) => {
  const { similar, player, opponent, prediction, selected, average } = props;
  //don't include players without games vs opponent

  let similarTeams = similar.filter((team) =>
    player.games.some((game) => game.opponent.teamID === team.teamID)
  );

  //if player played any games vs opponent then add opponent to similarTeams
  if (player.games.some((game) => game.opponent.teamID === opponent.teamID)) {
    similarTeams.push(opponent);
  }

  let avg = 0;
  let avgVsSim = 0;
  let gamesOver = 0;
  let gamesPlayed = 0;
  let avgDiff = 0;
  let simRows = [
    <tr>
      <th>{player.name} vs Similar Teams</th>
      {RelevantStats[selected.recognize].map((item, i) => (
        <th key={item.label}>{item.label}</th>
      ))}
    </tr>,
  ];
  similarTeams.forEach((team) => {
    const gamesVsTeam = player.games.filter(
      (game) => game.opponent.teamID === team.teamID
    );

    if (gamesVsTeam.length === 0) {
      return;
    }

    const playerAvg = AveragePropScore(player.games, prediction.stat.recognize);
    avg += playerAvg;
    gamesPlayed += gamesVsTeam.length;

    const playerAvgVsOpp = AveragePropScore(
      gamesVsTeam,
      prediction.stat.recognize
    );
    const diff = round(playerAvgVsOpp - playerAvg, 2);
    avgVsSim += playerAvgVsOpp;
    avgDiff += diff;

    const over = gamesVsTeam.filter(
      (game) =>
        GetPropScore(game, prediction.stat.recognize) > prediction.target
    ).length;

    gamesOver += over;

    simRows.push(
      <tr>
        <td>
          {team.abbreviation} {team.name} ({gamesVsTeam.length} games)
        </td>
        {RelevantStats[selected.recognize].map((stat) => {
          const cellTarget = AveragePropScore(player.games, stat.recognize);
          const score = AveragePropScore(gamesVsTeam, stat.recognize);
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

  avgVsSim = round(avgVsSim / similarTeams.length, 2);
  avg = round(avg / similarTeams.length, 2);
  avgDiff = round(avgDiff / similarTeams.length, 2);

  const pctDiff = round((avgDiff / avg) * 100, 2);
  const pct = round((gamesOver / gamesPlayed) * 100, 2);
  simRows.push(
    <>
      <tr>
        <th className={"right"}>Average ({avg})</th>
        <td
          className={avgVsSim > avg ? "high" : avgVsSim < avg ? "low" : "med"}
        >
          {avgVsSim}
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
