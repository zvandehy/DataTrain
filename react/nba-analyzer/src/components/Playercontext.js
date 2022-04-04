import React from 'react'
import EventsModal from "./EventsModal"

const PlayerContext = () => {
  return (
    <div className="player-context">
        <div className="player-info">
            <h2 className="player-name">Giannis Antetokounmpo <span className="player-position tall">C</span></h2>
            <div className="player-matchup">
                <span>
                    <img className="team-logo" alt="team-logo" src="https://cdn.nba.com/logos/nba/1610612755/primary/D/logo.svg"></img>
                    <span className="team-abr">PHI</span>
                </span>
                <span className="versus">vs.</span>
                <span>
                    <img className="team-logo" alt="team-logo" src="https://cdn.nba.com/logos/nba/1610612751/primary/D/logo.svg"></img>
                    <span className="opp-abr">BKN</span>
                </span>
                
                <EventsModal/>
            </div>
        </div>
        <img className="player-photo" alt="playername" src={"https://ak-static.cms.nba.com/wp-content/uploads/headshots/nba/latest/260x190/201988.png"}></img>
    </div>
  )
}

export default PlayerContext