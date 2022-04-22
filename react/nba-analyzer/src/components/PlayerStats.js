import React from "react";
import LastNGames from "./LastNGames";
import { GetPropScore, AveragePropScore, GetColor } from "../utils";
import { RelevantStats } from "../utils";

export const PlayerStatsPreview = (props) => {
  const { predictions, selected, matchups } = props;
  const prediction = predictions.filter(
    (item) => item.stat.label.toLowerCase() === selected.toLowerCase()
  )[0];
  const playerCounts = prediction.counts;
  //TODO: Add state for cycling between PCT, OVER, UNDER, etc. (use in header)
  return (
    <div className="player-stats">
      <span></span>
      <span className="player-stat-header">AVG</span>
      <span className="player-stat-result-header">PCT</span>

      {playerCounts.map((item, i) => (
        <LastNGames key={item.n + " " + i} countObj={item} header={"PCT"} />
      ))}

      {/* VERSUS STATS */}

      {/* TODO: Investigate if should be average or all matchups */}

      {matchups.map((game) => {
        const score = GetPropScore(game, prediction.stat.recognize);
        const result = score > prediction.target ? "OVER" : "UNDER";
        // TODO: Make this its own component?
        return (
          <React.Fragment key={game.gameID}>
            {/* TODO: make "vs."" dynamic between 'vs' (home) &  '@' (away) */}
            <span className="header player-vs-opp-instance">
              vs {game.opponent.abbreviation}
            </span>
            <span className="player-vs-opp-stat">{score}</span>
            <span
              className={`player-vs-opp-stat-result ${GetColor(
                "over",
                result
              )}`}
            >
              {result}
            </span>
          </React.Fragment>
        );
      })}

      {/* SIMILAR STATS */}

      <span className="header similar-players-header">Similar Players</span>
      <span className="similar-players-stat">xx.x</span>
      <span className="similar-players-stat-result">XX%</span>

      <span className="header similar-opp-header">Similar Opp</span>
      <span className="similar-opp-stat">xx.x</span>
      <span className="similar-opp-stat-result">XX%</span>
    </div>
  );
};

export const PlayerStatsTable = (props) => {
  const { predictions, selected, games, matchups } = props;
  const prediction = predictions.filter(
    (item) => item.stat.label.toLowerCase() === selected.toLowerCase()
  )[0];
  const playerCounts = prediction.counts;
  //TODO: Add state for cycling between PCT, OVER, UNDER, etc. (use in header)
  return (
    <div className="player-stats">
      <table>
        <thead>
          <tr>
            <th># Games</th>
            {RelevantStats[selected].map((item) => (
              <th key={item.label}>{item.label}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {playerCounts.map((item, i) => (
            <tr key={item.n + " " + i}>
              <td>{item.n}</td>
              {RelevantStats[selected].map((stat, i) => {
                const nGames = games.slice(item.n * -1);

                const cellTarget =
                  i === 0 && prediction.target
                    ? prediction.target
                    : AveragePropScore(games, stat.recognize);
                return (
                  <AverageStatCell
                    games={nGames}
                    stat={stat}
                    target={cellTarget}
                    key={`${nGames.length} ${stat.label} ${i}`}
                  />
                );
              })}
            </tr>
          ))}
          {matchups.length ? (
            <>
              <tr>
                <th>Matchup</th>
                {RelevantStats[selected].map((item, i) => (
                  <th key={item.label}>{item.label}</th>
                ))}
              </tr>
              {matchups.reverse().map((game) => {
                return (
                  <tr key={`'vs' ${game.opponent.abbreviation} ${game.date}`}>
                    <td>
                      vs {game.opponent.abbreviation} {game.date}
                    </td>
                    {RelevantStats[selected].map((stat, i) => {
                      const cellTarget =
                        i === 0 && prediction.target
                          ? prediction.target
                          : AveragePropScore(games, stat.recognize);
                      return (
                        <GameStatCell
                          game={game}
                          stat={stat}
                          target={cellTarget}
                          key={game.date + " " + stat.label}
                        />
                      );
                    })}
                  </tr>
                );
              })}
            </>
          ) : (
            <></>
          )}
        </tbody>
      </table>
    </div>
  );
};

export const AverageStatCell = (props) => {
  const { games, stat, target } = props;
  const avg = AveragePropScore(games, stat.recognize);
  return (
    <td className={avg > target ? "high" : avg < target ? "low" : "med"}>
      {avg}
    </td>
  );
};

export const GameStatCell = (props) => {
  const { game, stat, target } = props;
  const score = GetPropScore(game, stat.recognize);
  return (
    <td className={score > target ? "high" : score < target ? "low" : "med"}>
      {score}
    </td>
  );
};
