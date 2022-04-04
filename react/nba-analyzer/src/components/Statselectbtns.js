import React from 'react'
import StatSelectButton from './StatSelectButton'


const StatSelectBtns = () => {
  return (
    <div className="stat-select">
            <StatSelectButton stat="Points" target={30.5} confidence={67.5}/>
            <StatSelectButton stat="assists" target={30.5} confidence={67.5}></StatSelectButton>
            <StatSelectButton stat="Fantasy" target={30.5} confidence={67.5}></StatSelectButton>
            <StatSelectButton stat="PTS + REB + AST" target={30.5} confidence={67.5}></StatSelectButton>
            <StatSelectButton stat="Free Throws" target={30.5} confidence={67.5}></StatSelectButton>
            <StatSelectButton stat="3 Pointers" target={30.5} confidence={67.5}></StatSelectButton>
            <StatSelectButton stat="Blocks + Steals" target={30.5} confidence={67.5}></StatSelectButton>
            <StatSelectButton stat="Double Double" target={30.5} confidence={67.5}></StatSelectButton>
            <button className="stat-select-btn">
                <p className="bold">Other</p>
            </button>
        </div>
  )
}

export default StatSelectBtns