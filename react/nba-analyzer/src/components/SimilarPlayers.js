import React from "react";
import { AveragePropScore, GetPropScore } from "../utils";
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
  let avg = 0;
  let avgVsOpp = 0;
  let gamesOver = 0;
  let gamesPlayed = 0;
  similar.map((player) => {
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

    return (
      <tr>
        <td>{player.name}</td>
        <td>{playerAvg}</td>
        <td>
          {player.games
            .filter((game) => game.opponent.teamID === opponent)
            .map((game) => {
              const score = GetPropScore(game, prediction.stat.recognize);
              return <span></span>;
            })}
        </td>
      </tr>
    );
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
