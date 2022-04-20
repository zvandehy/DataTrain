import React from "react";
import LastNGames from "./LastNGames";
import { GetPropScore, AveragePropScore, GetColor } from "../utils";
import { PredictionIconSmall } from "./Prediction";

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

  const relevantStats = {
    Points: [
      { recognize: "points", label: "PTS" },
      { recognize: "field_goals_made", label: "FGM" },
      { recognize: "field_goals_attempted", label: "FGA" },
      { recognize: "field_goal_percentage", label: "FG%" },
      { recognize: "three_pointers_made", label: "3PM" },
      { recognize: "three_pointers_attempted", label: "3PA" },
      { recognize: "three_point_percentage", label: "3P%" },
      { recognize: "free_throws_made", label: "FTM" },
      { recognize: "free_throws_attempted", label: "FTA" },
      // {recognize:,label:"FT RT"},
      { recognize: "effective_field_goal_percentage", label: "eFG%" },
      { recognize: "usage", label: "USG%" },
      { recognize: "minutes", label: "MIN" },
    ],
    Assists: [
      { recognize: "assists", label: "AST" },
      { recognize: "assist_percentage", label: "AST%" },
      // { recognize: "potential_assists", label: "POT. AST" },
      { recognize: "turnovers", label: "TOV" },
      { recognize: "field_goals_attempted", label: "FGA" },
      { recognize: "usage", label: "USG%" },
      { recognize: "minutes", label: "MIN" },
    ],
    "3 Pointers": [
      { recognize: "three_pointers_made", label: "3PM" },
      { recognize: "three_pointers_attempted", label: "3PA" },
      { recognize: "three_point_percentage", label: "3P%" },
      { recognize: "field_goals_attempted", label: "FGA" },
      { recognize: "percent_three_pointers", label: "%3P" },
      { recognize: "usage", label: "USG%" },
      { recognize: "minutes", label: "MIN" },
    ],
    "PTS + REB + AST": [
      { recognize: "pts+rebs+asts", label: "PRA" },
      { recognize: "points", label: "PTS" },
      { recognize: "total_rebounds", label: "REB" },
      { recognize: "assists", label: "AST" },
      { recognize: "field_goals_attempted", label: "FGA" },
      { recognize: "effective_field_goal_percentage", label: "eFG%" },
      { recognize: "offensive_rebounds", label: "OREB" },
      { recognize: "defensive_rebounds", label: "DREB" },
      { recognize: "usage", label: "USG%" },
      { recognize: "minutes", label: "MIN" },
    ],
    Rebounds: [
      { recognize: "total_rebounds", label: "REB" },
      { recognize: "offensive_rebounds", label: "OREB" },
      { recognize: "defensive_rebounds", label: "DREB" },
      { recognize: "offensive_rebound_percentage", label: "OREB%" },
      { recognize: "defensive_rebound_percentage", label: "DREB%" },
    ],
    "Free Throws": [
      { recognize: "free_throws_made", label: "FTM" },
      { recognize: "free_throws_attempted", label: "FTA" },
      { recognize: "free_throws_percentage", label: "FT%" },
      // { recognize: "free_throw_rate", label: "FT RT" },
      { recognize: "personal_fouls_drawn", label: "PFD" },
    ],
    Fantasy: [
      { recognize: "fantasy score", label: "FAN" },
      { recognize: "points", label: "PTS (1)" },
      { recognize: "assists", label: "AST (1.5)" },
      { recognize: "total_rebounds", label: "REB (1.2)" },
      { recognize: "blocks", label: "BLK (3)" },
      { recognize: "steals", label: "STL (3)" },
      { recognize: "turnovers", label: "TOV (-1)" },
    ],
    "Blocks + Steals": [
      { recognize: "blocks", label: "BLK" },
      { recognize: "steals", label: "STL" },
      { recognize: "personal_fouls", label: "PF" },
      { recognize: "minutes", label: "MIN" },
    ],
    "Double Double": [
      { recognize: "double-double", label: "DD" },
      { recognize: "points", label: "PTS" },
      { recognize: "assists", label: "AST" },
      { recognize: "total_rebounds", label: "REB" },
      { recognize: "minutes", label: "MIN" },
    ],
  };
  console.log(relevantStats, selected);
  //TODO: Add state for cycling between PCT, OVER, UNDER, etc. (use in header)
  return (
    <div className="player-stats">
      <table>
        <thead>
          <tr>
            <th># Games</th>
            {relevantStats[selected].map((item) => (
              <th>{item.label}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {playerCounts.map((item, i) => (
            <tr key={item.n + " " + i}>
              <td>{item.n}</td>
              {relevantStats[selected].map((stat, i) => {
                const nGames = games.slice(item.n * -1);

                const cellTarget =
                  i === 0 && prediction.target
                    ? prediction.target
                    : AveragePropScore(games, stat.recognize);
                console.log(stat, selected, cellTarget);
                return (
                  <AverageStatCell
                    games={nGames}
                    stat={stat}
                    target={cellTarget}
                  />
                );
              })}
            </tr>
          ))}
          {matchups.length ? (
            <>
              <tr>
                <th>Matchup</th>
                {relevantStats[selected].map((item) => (
                  <th>{item.label}</th>
                ))}
              </tr>
              {matchups.map((game) => {
                return (
                  <tr>
                    <td>
                      vs {game.opponent.abbreviation} {game.date}
                    </td>
                    {relevantStats[selected].map((stat, i) => {
                      const cellTarget =
                        i === 0 && prediction.target
                          ? prediction.target
                          : AveragePropScore(games, stat.recognize);
                      return (
                        <GameStatCell
                          game={game}
                          stat={stat}
                          target={cellTarget}
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
