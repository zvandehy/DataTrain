import React, { useState } from "react";
import { useLocation } from "react-router-dom";
import { FormatDate, RelevantStats } from "../utils.js";
import PlayerStatsChart from "./PlayerStatsChart";
import { gql, useQuery } from "@apollo/client";
import { PlayerPageContext } from "./Playercontext";
import PlayerProfileChart from "./PlayerProfileChart.js";

import "../styles/player.css";
import StatSelectBtns from "./Statselectbtns.js";
import Prediction from "./Prediction.js";
import { PlayerStatsTable } from "./PlayerStats.js";
import { CalculatePredictions } from "../predictions.js";
import { round, mean } from "mathjs";

const Player = () => {
  let location = useLocation();
  const date = FormatDate(new Date());
  const playerID = parseInt(
    location.pathname.split("/")[location.pathname.split("/").length - 1]
  );
  const query = gql`
    query Player($playerID: Int!, $date: String!) {
      player(input: { playerID: $playerID }) {
        name
        playerID
        position
        currentTeam {
          abbreviation
          teamID
        }
        projections(input: { sportsbook: "PrizePicks", startDate: $date }) {
          date
          opponent {
            abbreviation
            teamID
          }
          targets {
            target
            type
          }
        }
        games(input: { season: "2021-22" }) {
          points
          season
          assists
          assist_percentage
          rebounds
          offensive_rebounds
          offensive_rebound_percentage
          defensive_rebounds
          defensive_rebound_percentage
          personal_fouls_drawn
          steals
          blocks
          turnovers
          teamStats {
            points
            assists
            rebounds
            three_pointers_attempted
            blocks
            steals
          }
          opponent {
            abbreviation
            teamID
          }
          minutes
          date
          field_goals_attempted
          field_goal_percentage
          field_goals_made
          three_pointers_attempted
          three_pointers_made
          free_throws_attempted
          free_throws_made
          free_throws_percentage
          usage
          three_point_percentage
          effective_field_goal_percentage
        }
        similarPlayers(input: { season: "2021-22" }) {
          name
          playerID
          games(input: { season: "2021-22" }) {
            points
            season
            assists
            assist_percentage
            rebounds
            offensive_rebounds
            offensive_rebound_percentage
            defensive_rebounds
            defensive_rebound_percentage
            personal_fouls_drawn
            steals
            blocks
            turnovers
            opponent {
              abbreviation
              teamID
            }
            minutes
            date
            field_goals_attempted
            field_goal_percentage
            field_goals_made
            three_pointers_attempted
            three_pointers_made
            free_throws_attempted
            free_throws_made
            free_throws_percentage
            usage
            three_point_percentage
            effective_field_goal_percentage
          }
        }
      }
    }
  `;
  const { loading, error, data } = useQuery(query, {
    variables: { playerID: playerID, date: date },
  });

  const [stat, setStat] = useState("Points");
  function onStatSelect(stat) {
    setStat(stat);
  }

  if (loading) return "Loading...";
  if (error) {
    return `Error! ${error.message}. ${loading}. ${data}`;
  }
  console.log(data);
  const statData = data.player.games.filter(
    (game) => game.season === "2021-22"
  );
  const predictions = CalculatePredictions(
    data.player.projections[0],
    statData
  );

  const game = data.player.games.find((game) => game.date === date);
  const matchups = data.player.games.filter(
    (game) =>
      game.opponent.teamID === data.player.projections[0].opponent.teamID
  );

  const percentOfTeamStats = RelevantStats["Profile"].map((item) => {
    const avg = mean(data.player.games.map((game) => game[item.recognize]));
    const pct = round(
      avg /
        mean(data.player.games.map((game) => game.teamStats[item.recognize])),
      2
    );
    return {
      stat: item.label,
      avg: avg,
      pct: pct,
    };
  });

  return (
    <div className="player-page">
      <PlayerPageContext
        player={data.player}
        opponent={data.player.projections[0].opponent}
        game={game}
      />
      <PlayerProfileChart stats={percentOfTeamStats} />
      <StatSelectBtns
        predictions={predictions}
        playername={data.player.name}
        selected={stat}
        onStatSelect={onStatSelect}
      />
      <Prediction predictions={predictions} selected={stat} game={game} />
      <PlayerStatsTable
        predictions={predictions}
        selected={stat}
        matchups={matchups}
        games={statData}
        similar={data.player.similarPlayers}
        opponent={data.player.projections[0].opponent.teamID}
      />
      <PlayerStatsChart games={data.player.games} />
    </div>
  );
};

export default Player;
