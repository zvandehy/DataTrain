import React, { useState } from "react";
import { useLocation } from "react-router-dom";
import {
  CompareDates,
  FormatDate,
  GamesWithSeasonType,
  RelevantStats,
} from "../utils.js";
import PlayerStatsChart from "./PlayerStatsChart";
import { gql, useQuery } from "@apollo/client";
import { PlayerPageContext } from "./Playercontext";
import PlayerProfileChart from "./PlayerProfileChart.js";

import "../styles/player.css";
import StatSelectBtns from "./Statselectbtns.js";
import Prediction from "./Prediction.js";
import { PlayerStatsTable } from "./PlayerStats.js";
import { CalculatePredictions, StatObjects } from "../predictions.js";
import { round, mean } from "mathjs";

const Player = (prop) => {
  const { client, league } = prop;
  const season = league === "nba" ? "2021-22" : "2022-23";
  let location = useLocation();
  const playerID = parseInt(
    location.pathname.split("/")[location.pathname.split("/").length - 1]
  );
  const [date, setDate] = useState(FormatDate(new Date()));
  //TODO: Handle error when pick game that isn't in seasonType
  const [seasonType, setSeasonType] = useState("");
  const query = gql`
    query Player($playerID: Int!, $season: String!) {
      player(input: { playerID: $playerID }) {
        name
        playerID
        position
        currentTeam {
          abbreviation
          teamID
          name
          injuries {
            startDate
            returnDate
            status
            player {
              name
            }
          }
        }
        projections(input: {}) {
          date
          opponent {
            abbreviation
            teamID
            name
            injuries {
              startDate
              returnDate
              status
              player {
                name
              }
            }
          }
          propositions {
            target
            type
            sportsbook
          }
        }
        games(input: { season: $season }) {
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
          team {
            abbreviation
            teamID
          }
          opponent {
            abbreviation
            teamID
            injuries {
              startDate
              returnDate
              status
              player {
                name
              }
            }
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
          playoffs
        }
        similarPlayers(input: { season: $season }) {
          name
          playerID
          games(input: { season: $season }) {
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
            playoffs
          }
        }
      }
    }
  `;
  const { loading, error, data } = useQuery(query, {
    variables: { playerID: playerID, season: season },
    client: client,
  });

  const [stat, setStat] = useState(
    StatObjects.find((stat) => stat.label === "Points")
  );
  function onStatSelect(stat) {
    setStat(stat);
  }

  if (loading) return "Loading...";
  if (error) {
    return `Error! ${error.message}. ${error}. ${loading}. ${data}`;
  }
  let games = GamesWithSeasonType(data.player.games, seasonType);
  games = games.filter((game) => CompareDates(game.date, date) < 0);
  games = games.sort((a, b) => CompareDates(a.date, b.date));
  // const statData = games.filter((game) => game.season === "2021-22");
  // get the projection with the current date or the most recent date
  const projection =
    data.player.projections.find((p) => p.date === date) ||
    data.player.projections.find(
      (p) => p.date === games[games.length - 1].date
    ) ||
    data.player.projections[0];
  // TODO: when automatically getting last game, the result doesn't generate
  //TODO: handle if no projection is found
  const statData = games;
  const predictions = CalculatePredictions(projection, statData);
  const game = data.player.games.find((game) => game.date === date);
  let matchups = games.filter(
    (matchup) =>
      matchup.opponent.teamID ===
      (projection?.opponent.teamID ?? game?.opponent.teamID)
  );
  if (game !== undefined) matchups.push(game);

  const percentOfTeamStats = RelevantStats["Profile"].map((item) => {
    const avg = mean(games.map((game) => game[item.recognize]));
    const pct = round(
      avg / mean(games.map((game) => game.teamStats[item.recognize])),
      2
    );
    return {
      stat: item.label,
      avg: avg,
      pct: pct,
    };
  });

  function onDateSelect(newDate) {
    setDate(newDate);
  }

  // TODO: handle injuries for previous games
  return (
    <div className="player-page">
      <PlayerPageContext
        player={data.player}
        opponent={projection?.opponent ?? game?.opponent}
        game={game}
        date={projection?.date ?? game?.date}
        selectDate={onDateSelect}
        seasonType={seasonType}
        setSeasonType={setSeasonType}
        league={league}
      />
      <PlayerProfileChart stats={percentOfTeamStats} />
      <StatSelectBtns
        predictions={predictions}
        playername={data.player.name}
        selected={stat}
        onStatSelect={onStatSelect}
      />
      <Prediction
        propositions={projection.propositions}
        predictions={predictions}
        selected={stat}
        game={game}
      />
      <PlayerStatsTable
        predictions={predictions}
        selected={stat}
        matchups={matchups}
        games={statData}
        similar={data.player.similarPlayers}
        opponent={
          projection?.opponent ??
          game?.opponent ??
          games[games.length - 1].opponent
        }
      />
      <PlayerStatsChart
        games={games}
        predictions={predictions}
        selected={stat}
      />
    </div>
  );
};

export default Player;
