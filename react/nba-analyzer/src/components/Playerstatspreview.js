import React from 'react'
import LastNGames from './LastNGames'
import {GetPropScore, GetColor} from '../utils'

const PlayerStatsPreview = (props) => {
  const {projections, selected, matchups} = props;
  const projection = projections.filter(item => item.stat.label.toLowerCase() === selected.toLowerCase())[0]
  const playerCounts = projection.counts;
  //TODO: Add state for cycling between PCT, OVER, UNDER, etc. (use in header)
  return (
    <div className="player-stats">
            <span></span>
            <span className="player-stat-header">AVG</span>
            <span className="player-stat-result-header">PCT</span>

            {playerCounts.map((item,i)=> <LastNGames key={item.n+" "+i} countObj={item} header={"PCT"}/>)}
            
            {/* VERSUS STATS */}

            {/* TODO: Investigate if should be average or all matchups */}

            {matchups.map(game => {
              const score = GetPropScore(game, projection.stat.recognize);
              const result = score > projection.target ? "OVER" : "UNDER";
              // TODO: Make this its own component?
              return <React.Fragment key={game.gameID}>
              {/* TODO: make "vs."" dynamic between 'vs' (home) &  '@' (away) */}
              <span className="header player-vs-opp-instance">vs {game.opponent.abbreviation}</span>
              <span className="player-vs-opp-stat">{score}</span>
              <span className={`player-vs-opp-stat-result ${GetColor("over", result)}`}>{result}</span>
            </React.Fragment>})}
            

            {/* <span className="header player-vs-opp-instance">@ BKN</span>
            <span className="player-vs-opp-stat">28</span>
            <span className="player-vs-opp-stat-result">Over</span> */}

            {/* SIMILAR STATS */}
            
            <span className="header similar-players-header">Similar Players</span>
            <span className="similar-players-stat">xx.x</span>
            <span className="similar-players-stat-result">XX%</span>

            <span className="header similar-opp-header">Similar Opp</span>
            <span className="similar-opp-stat">xx.x</span>
            <span className="similar-opp-stat-result">XX%</span>

        </div>
  )
}

export default PlayerStatsPreview