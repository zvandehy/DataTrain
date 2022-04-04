import React from 'react'
import LastNGames from './LastNGames'

const PlayerStatsPreview = () => {
  return (
    <div className="player-stats">
            <span></span>
            <span className="player-stat-header">AVG</span>
            <span className="player-stat-result-header">PCT</span>

            <LastNGames count={70} over={25} under={40} avg={30} header={"PCT"}/>
            <LastNGames count={30} over={20} under={10} avg={34} header={"PCT"}/>
            <LastNGames count={10} over={5} under={4} avg={29} header={"PCT"}/>
            <LastNGames count={5} over={1} under={4} avg={27} header={"PCT"}/>
            
            {/* VERSUS STATS */}

            {/* Technically this could just be the average of all matchups with X opponent? */}

            <span className="header player-vs-opp-instance">vs BKN</span>
            <span className="player-vs-opp-stat">28</span>
            <span className="player-vs-opp-stat-result">Over</span>

            {/* <span className="header player-vs-opp-instance">@ BKN</span>
            <span className="player-vs-opp-stat">28</span>
            <span className="player-vs-opp-stat-result">Over</span> */}

            {/* SIMILAR STATS */}
            
            <span className="header similar-players-header">Similar Players</span>
            <span className="similar-players-stat">22.2</span>
            <span className="similar-players-stat-result">25%</span>

            <span className="header similar-opp-header">Similar Opp</span>
            <span className="similar-opp-stat">25.5</span>
            <span className="similar-opp-stat-result">50%</span>

        </div>
  )
}

export default PlayerStatsPreview