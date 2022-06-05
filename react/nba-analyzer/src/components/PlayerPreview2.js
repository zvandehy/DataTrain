import React, { useState } from "react";
import { Link } from "react-router-dom";
import { mean, round } from "mathjs";
import { GetPropScore, GetShortType } from "../utils";

const PlayerPreview2 = (props) => {
  const { playerProp } = props;
  let player = playerProp.player;
  let seasonData = playerProp.player.games
    .filter((game) => game.season === "2021-22")
    .sort(function (a, b) {
      var c = new Date(a.date);
      var d = new Date(b.date);
      return c - d;
    });

  const counts = [0, -30, -10, -5];
  const weights = [0.4, 0.25, 0.2, 0.15]; //TODO: determine best weights to use

  let statData = new Map();
  let maxType = playerProp.propositions[0].type;
  let max = 0.0;
  playerProp.propositions.forEach((projection) => {
    statData.set(projection.type, new Map());
    counts.forEach((count) => {
      const data = seasonData
        .slice(count)
        .map((game) => GetPropScore(game, projection.type));
      const over = data.filter((score) => score > projection.target).length;
      const under = data.filter((score) => score < projection.target).length;
      const pct = round((over / data.length) * 100, 2);
      const pct_under = round((under / data.length) * 100, 2);
      statData.get(projection.type).set(count, {
        over: over,
        under: under,
        pct: pct,
        pct_u: pct_under,
        mean: data.length ? round(mean(data), 1) : 0,
      });
    });
    let conf = 0.0;
    let conf_u = 0.0;
    for (let i = 0; i < counts.length; i++) {
      conf += statData.get(projection.type).get(counts[i]).pct * weights[i];
      conf_u += statData.get(projection.type).get(counts[i]).pct_u * weights[i];
    }
    let confStr = "";
    conf = round(conf, 2);
    conf_u = round(conf_u, 2);
    if (conf > conf_u) {
      confStr = "O" + conf;
      statData.get(projection.type).set("conf", conf);
      if (conf > max) {
        maxType = projection.type;
        max = conf;
      }
    } else {
      confStr = "U" + conf_u;
      statData.get(projection.type).set("conf", conf_u);
      if (conf_u > max) {
        maxType = projection.type;
        max = conf_u;
      }
    }
    statData.get(projection.type).set("confStr", confStr);
  });

  const [type, setType] = useState(maxType);

  const displayConf = statData.get(type).get("conf");
  const displayConfStr = statData.get(type).get("confStr");

  if (!seasonData) {
    return <div></div>;
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
            {playerProp.propositions.map((projection) => {
              return (
                <li
                  key={"select_" + playerProp.player.playerID + projection.type}
                >
                  <button
                    className={
                      projection.type === type ? "selected" : "unselected"
                    }
                    onClick={() => setType(projection.type)}
                  >
                    {GetShortType(projection.type)}
                  </button>
                </li>
              );
            })}
          </ul>
          <p>
            Target:{" "}
            {
              playerProp.propositions.find(
                (projection) => projection.type === type
              ).target
            }
          </p>
          <div>
            <p>Key Events:</p>
            <ul className="events">
              <li>BKN Playername traded</li>
              <li>BKN Playername injured</li>
              <li>CHI Playername returned</li>
              <li>CHI Playername returns</li>
            </ul>
          </div>
        </div>
        <table className="player-stats">
          <thead>
            <tr>
              <th
                colspan="2"
                className={
                  displayConf > 60 ? "high" : displayConf >= 50 ? "med" : "low"
                }
              >
                {displayConfStr}
              </th>
              <th>AVG</th>
              <th>OVR</th>
              <th>PCT</th>
            </tr>
          </thead>
          <tbody>
            {counts.map((count, i) => {
              const target = playerProp.propositions.find(
                (projection) => projection.type === type
              ).target;
              return (
                <>
                  {
                    <tr key={player.playerID + type + count}>
                      {i === 0 ? (
                        <th rowspan={counts.length}>Last N Games</th>
                      ) : (
                        <></>
                      )}
                      <th>{count ? count * -1 : seasonData.length}</th>
                      <td>{statData.get(type).get(count).mean}</td>
                      <td>{statData.get(type).get(count).over}</td>
                      <td
                        className={
                          "player-over-pct " +
                          (target === "-"
                            ? ""
                            : statData.get(type).get(count).pct > 60
                            ? "high"
                            : statData.get(type).get(count).pct >= 50
                            ? "med"
                            : "low")
                        }
                      >
                        {statData.get(type).get(count).pct + "%"}
                      </td>
                    </tr>
                  }
                </>
              );
            })}
            {seasonData
              .filter(
                (game) =>
                  game.opponent.abbreviation ===
                  playerProp.opponent.abbreviation
              )
              .map((game, i, arr) => {
                const target = playerProp.propositions.find(
                  (projection) => projection.type === type
                ).target;
                const pct = round(
                  (arr.filter((game) => GetPropScore(game, type) > target)
                    .length /
                    arr.length) *
                    100,
                  2
                );
                return (
                  <>
                    {
                      <tr key={player.playerID + type + game.date}>
                        {i === 0 ? (
                          <th rowSpan={arr.length}>
                            {playerProp.player.name.split(" ", 2)[1]} vs{" "}
                            {playerProp.opponent.abbreviation}
                          </th>
                        ) : (
                          <></>
                        )}
                        <th>{game.date}</th>
                        <td>{GetPropScore(game, type)}</td>
                        <td>{GetPropScore(game, type) > target ? "O" : "U"}</td>
                        <td
                          className={
                            "player-over-pct " +
                            (pct > 60 ? "high" : pct >= 50 ? "med" : "low")
                          }
                        >
                          {pct + "%"}
                        </td>
                      </tr>
                    }
                  </>
                );
              })}
            <tr>
              <th colspan="2">
                Similar Players vs {playerProp.opponent.abbreviation}
              </th>
              <td>xx</td>
              <td>##/xx</td>
              <td className={"player-over-pct high"}>{"xx.x%"}</td>
            </tr>
            <tr>
              <th colspan="2">
                {playerProp.player.name.split(" ", 2)[1]} vs Similar Opp
              </th>
              <td>xx</td>
              <td>##/xx</td>
              <td className={"player-over-pct low"}>{"xx.x%"}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </>
  );
};

export default PlayerPreview2;
