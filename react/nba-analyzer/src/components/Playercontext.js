import React from 'react'
import EventsModal from "./EventsModal"

const PlayerContext = (props) => {
    const {player, opponent} = props;
  return (
    <div className="player-context">
        <div className="player-info">
            <h2 className="player-name">{player.name} <span className="player-position tall titlecase">{player.position}</span></h2>
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
            </div>
        </div>
        <img className="player-photo" alt={player.name} src={`https://ak-static.cms.nba.com/wp-content/uploads/headshots/nba/latest/260x190/${player.playerID}.png`}></img>
    </div>
  )
}

export default PlayerContext