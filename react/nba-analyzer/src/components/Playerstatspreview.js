import React from 'react'
import LastNGames from './LastNGames'

const PlayerStatsPreview = (props) => {
  const {playerCounts, matchupData, similarData} = props;
  //TODO: Add state for cycling between PCT, OVER, UNDER, etc. (use in header)
  return (
    <div className="player-stats">
            <span></span>
            <span className="player-stat-header">AVG</span>
            <span className="player-stat-result-header">PCT</span>

            {playerCounts.map(item => <LastNGames countObj={item} header={"PCT"}/>)}
            
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