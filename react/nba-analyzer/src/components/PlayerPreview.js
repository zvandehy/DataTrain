import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { mean, round } from "mathjs";
import { GetPropScore, GetShortType } from "../utils";

const PlayerPreview = (props) => {
  const { playerProp } = props;
  let player = playerProp.player;
  let seasonData = playerProp.player.games
    // .filter((game) => game.season === "2021-22")
    .sort(function (a, b) {
      var c = new Date(a.date);
      var d = new Date(b.date);
      return c - d;
    });
  const [count, setCount] = useState(0);

  if (!seasonData) {
    return <li></li>;
  }

  return (
    <>
      <div className="playercard">
        <img
          className="player-photo"
          alt={player.name}
          src={
            "https://ak-static.cms.nba.com/wp-content/uploads/headshots/nba/latest/260x190/" +
            player.playerID +
            ".png"
          }
        ></img>
        <div className="player-info">
          <h2 className="player-name">
            <Link
              to={`/players/${player.playerID}`}
              state={{ playerProp: playerProp }}
            >
              {player.name}
            </Link>
            <span className="player-position">
              ({player.position ? player.position : "X"})
            </span>
          </h2>
          <p className="player-team">{player.currentTeam.abbreviation}</p>
          <p>
            vs{" "}
            <span className="player-opponent">
              {playerProp.opponent.abbreviation}
            </span>
          </p>
          <ul className="player-since">
            <li>
              <button onClick={(event) => setCount(0)}>
                Season ({seasonData.length})
              </button>
            </li>
            <li>
              <button onClick={(event) => setCount(-30)}>30</button>
            </li>
            <li>
              <button onClick={(event) => setCount(-10)}>10</button>
            </li>
            <li>
              <button onClick={(event) => setCount(-5)}>5</button>
            </li>
          </ul>
        </div>
        <table className="player-stats">
          <thead>
            <tr>
              <th>{count ? count : seasonData.length}</th>
              {playerProp.propositions.map((projection) => {
                return <th>{GetShortType(projection.type)}</th>;
              })}
            </tr>
          </thead>
          <tbody>
            <tr className="player-avg-row">
              <th>AVG</th>
              {playerProp.propositions.map((projection) => {
                const data = seasonData
                  .slice(count)
                  .map((game) => GetPropScore(game, projection.type));
                return (
                  <td key={player.playerID + projection.type + "avg"}>
                    {data.length > 0 ? round(mean(data), 1) : 0}
                  </td>
                );
              })}
            </tr>
            <tr className="player-tar-row">
              <th>TAR</th>
              {playerProp.propositions.map((projection) => {
                return (
                  <td
                    className="player-over"
                    key={player.playerID + projection.type + "target"}
                  >
                    {projection.target ?? "-"}
                  </td>
                );
              })}
            </tr>
            <tr className="player-ovr-row">
              <th>OVR</th>
              {playerProp.propositions.map((projection) => {
                const data = seasonData
                  .slice(count)
                  .map((game) => GetPropScore(game, projection.type));
                const over = seasonData
                  .slice(count)
                  .filter(
                    (game) =>
                      GetPropScore(game, projection.type) > projection.target
                  ).length;
                return (
                  <td
                    className="player-over"
                    key={player.playerID + projection.type + "over"}
                  >
                    {data.length > 0 ? over : "-"}
                  </td>
                );
              })}
            </tr>
            <tr className="player-und-row">
              <th>PCT</th>
              {playerProp.propositions.map((projection) => {
                const data = seasonData
                  .slice(count)
                  .map((game) => GetPropScore(game, projection.type));
                const pct = round(
                  (seasonData
                    .slice(count)
                    .filter(
                      (game) =>
                        GetPropScore(game, projection.type) > projection.target
                    ).length /
                    data.length) *
                    100,
                  2
                );
                return (
                  <td key={player.playerID + projection.type + "pct"}>
                    <span
                      className={
                        "player-ovr-pct " +
                        (pct > 60 ? "high" : pct >= 50 ? "med" : "low")
                      }
                    >
                      {data.length > 0 ? pct + "%" : "-"}
                    </span>
                  </td>
                );
              })}
            </tr>
          </tbody>
        </table>
      </div>
    </>
  );
};

export default PlayerPreview;
