import React from "react";
import LastNGames from "./LastNGames";
import { GetPropScore, AveragePropScore, GetColor } from "../utils";
import { RelevantStats } from "../utils";
import { SimilarPlayerRows } from "./SimilarPlayers";

export const PlayerStatsPreview = (props) => {
  const { predictions, selected, matchups } = props;
  const prediction = predictions.filter(
    (item) =>
      item.stat.recognize.toLowerCase() === selected.recognize.toLowerCase()
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
    </div>
  );
};

export const PlayerStatsTable = (props) => {
  const { predictions, selected, games, matchups, opponent, similar } = props;
  const prediction = predictions.filter(
    (item) =>
      item.stat.recognize.toLowerCase() === selected.recognize.toLowerCase()
  )[0];
  const playerCounts = prediction.counts;
  //TODO: Add state for cycling between PCT, OVER, UNDER, etc. (use in header)
  return (
    <div className="player-stats">
      <table>
        <thead>
          <tr>
            <th># Games</th>
            {RelevantStats[selected.recognize].map((item) => (
              <th key={item.label}>{item.label}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {playerCounts.map((item, i) => (
            <tr key={item.n + " " + i}>
              <td>{item.n}</td>
              {RelevantStats[selected.recognize].map((stat, i) => {
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
          {games.filter((game) => game.playoffs).length > 0 ? (
            <tr>
              <td>Playoffs ({games.filter((game) => game.playoffs).length})</td>
              {RelevantStats[selected.recognize].map((stat, i) => {
                const nGames = games.filter((game) => game.playoffs);

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
          ) : (
            <></>
          )}
          {matchups.length ? (
            <>
              <tr>
                <th>Matchup</th>
                {RelevantStats[selected.recognize].map((item, i) => (
                  <th key={item.label}>{item.label}</th>
                ))}
              </tr>
              {matchups.reverse().map((game) => {
                return (
                  <tr key={`'vs' ${game.opponent.abbreviation} ${game.date}`}>
                    <td>
                      vs {game.opponent.abbreviation} {game.date}
                    </td>
                    {RelevantStats[selected.recognize].map((stat, i) => {
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
          <SimilarPlayerRows
            similar={similar}
            opponent={opponent}
            prediction={prediction}
            selected={selected}
            average={AveragePropScore(games, prediction.stat.recognize)}
          />
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
