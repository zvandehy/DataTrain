import React from "react";
import DataListInput from "react-datalist-input";
import { Link } from "react-router-dom";
import { CompareDates, FormatDate, GamesWithSeasonType } from "../utils";
import EventsModal from "./EventsModal";

const PlayerContext = (props) => {
  const { player, opponent, date } = props;
  return (
    <div className="player-context">
      <PlayerInfo player={player} opponent={opponent} link={true} date={date} />
      <img
        loading="lazy"
        className="player-photo"
        alt={player.name}
        src={`https://ak-static.cms.nba.com/wp-content/uploads/headshots/nba/latest/260x190/${player.playerID}.png`}
      ></img>
    </div>
  );
};

export const PlayerPageContext = (props) => {
  const { player, opponent, date, selectDate, seasonType, setSeasonType } =
    props;
  return (
    <div className="player-context">
      <PlayerInfo
        player={player}
        seasonType={seasonType}
        opponent={opponent}
        link={false}
        date={date}
        selectDate={selectDate}
        setSeasonType={setSeasonType}
      />
      <div className="container">
        <img
          loading="lazy"
          className="player-photo"
          alt={player.name}
          src={`https://ak-static.cms.nba.com/wp-content/uploads/headshots/nba/latest/260x190/${player.playerID}.png`}
        ></img>
      </div>
    </div>
  );
};

function GetGameDropdowns(player, seasonType) {
  const games = GamesWithSeasonType(player.games, seasonType);
  let dropdowns = games.map((game) => {
    return { date: game.date, opponent: game.opponent, team: game.team };
  });
  player.projections.forEach((projection) => {
    if (CompareDates(FormatDate(new Date()), projection.date) <= 0) {
      dropdowns.unshift({
        date: projection.date,
        opponent: projection.opponent,
        team: player.currentTeam,
      });
    }
  });
  return dropdowns.map((game, i, arr) => {
    return {
      key: game.date,
      label: `${game.opponent.abbreviation} ${arr.length - i} | ${game.date} | 
      ${game.team.abbreviation} vs ${game.opponent.abbreviation}`,
      value: game.date,
    };
  });
}

export const PlayerInfo = (props) => {
  const { player, opponent, link, selectDate, seasonType, setSeasonType } =
    props;
  return (
    <div className="player-info">
      <h2 className="player-name">
        {link ? (
          <Link to={`/players/${player.playerID}`}>{player.name}</Link>
        ) : (
          player.name
        )}
        <span className="player-position tall titlecase">
          {player.position}
        </span>
      </h2>
      {selectDate ? (
        <div id={"games-dropdown"}>
          <DataListInput
            placeholder="Select a game"
            items={GetGameDropdowns(player, seasonType)}
            onSelect={(date) => selectDate(date.value)}
            clearInputOnClick={true}
            suppressReselect={true}
          />
        </div>
      ) : (
        <></>
      )}
      <div className="player-matchup">
        <span>
          <img
            loading="lazy"
            className="team-logo"
            alt={`team-logo for ${player.currentTeam.abbreviation}`}
            src={`https://cdn.nba.com/logos/nba/${player.currentTeam.teamID}/primary/D/logo.svg`}
          ></img>
          <span className="team-abr">{player.currentTeam.abbreviation}</span>
        </span>
        <span className="versus">vs.</span>
        <span>
          <img
            loading="lazy"
            className="team-logo"
            alt={`team-logo for ${opponent.abbreviation}`}
            src={`https://cdn.nba.com/logos/nba/${opponent.teamID}/primary/D/logo.svg`}
          ></img>
          <span className="opp-abr">{opponent.abbreviation}</span>
        </span>
        <EventsModal player={player} opponent={opponent} />
      </div>
      <button
        onClick={() =>
          setSeasonType(
            seasonType === "REG"
              ? "PLAYOFFS"
              : seasonType === "PLAYOFFS"
              ? ""
              : "REG"
          )
        }
      >
        {seasonType === "" ? "ANY" : seasonType}
      </button>
    </div>
  );
};

export default PlayerContext;
