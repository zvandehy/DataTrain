import React from 'react'
import {Link} from 'react-router-dom';
import EventsModal from "./EventsModal"

const PlayerContext = (props) => {
    const {player, opponent, game} = props;
  return (
    <div className="player-context">
        <div className="player-info">
            <h2 className="player-name"><Link to={`/players/${player.playerID}`}>{player.name}</Link><span className="player-position tall titlecase">{player.position}</span></h2>
            <div className="player-matchup">
                <span>
                    <img className="team-logo" alt={`team-logo for ${player.currentTeam.abbreviation}`} src={`https://cdn.nba.com/logos/nba/${player.currentTeam.teamID}/primary/D/logo.svg`}></img>
                    <span className="team-abr">{player.currentTeam.abbreviation}</span>
                </span>
                <span className="versus">vs.</span>
                <span>
                    <img className="team-logo" alt={`team-logo for ${opponent.abbreviation}`} src={`https://cdn.nba.com/logos/nba/${opponent.teamID}/primary/D/logo.svg`}></img>
                    <span className="opp-abr">{opponent.abbreviation}</span>
                </span>
                <EventsModal/>
                {game && game.date ? <div className="matchup-date">
                    <span>{game.date}</span>
                </div> : <></>}
            </div>
        </div>
        <img className="player-photo" alt={player.name} src={`https://ak-static.cms.nba.com/wp-content/uploads/headshots/nba/latest/260x190/${player.playerID}.png`}></img>
    </div>
  )
}

export const PlayerPageContext = (props) => {
    const {player, opponent, game} = props;
    return (
      <div className="player-context">
          <div className="player-info">
              <h2 className="player-name">{player.name}<span className="player-position tall titlecase">{player.position}</span></h2>
              <div className="player-matchup">
              {game && game.date ? <div className="matchup-date">
                    <span>{game.date}</span>
                </div> : <></>}
                  <span>
                      <img className="team-logo" alt={`team-logo for ${player.currentTeam.abbreviation}`} src={`https://cdn.nba.com/logos/nba/${player.currentTeam.teamID}/primary/D/logo.svg`}></img>
                      <span className="team-abr">{player.currentTeam.abbreviation}</span>
                  </span>
                  <span className="versus">vs.</span>
                  <span>
                      <img className="team-logo" alt={`team-logo for ${opponent.abbreviation}`} src={`https://cdn.nba.com/logos/nba/${opponent.teamID}/primary/D/logo.svg`}></img>
                      <span className="opp-abr">{opponent.abbreviation}</span>
                  </span>
                  <EventsModal/>
              </div>
          </div>
          <div className="container">
              <img className="player-photo" alt={player.name} src={`https://ak-static.cms.nba.com/wp-content/uploads/headshots/nba/latest/260x190/${player.playerID}.png`}></img>
            </div>
      </div>
    )
}


export default PlayerContext