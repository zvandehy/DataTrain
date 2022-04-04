import React from 'react'
import StatSelectButton from './StatSelectButton'


const StatSelectBtns = (props) => {
  const {projections,playername} = props;
  // console.log(projections)
  //targets = [{stat:"Points", target:20.5, prediction:"OVER", confidence=55%}]

  return (
    <div className="stat-select">
      {projections.map(item => {
        return <StatSelectButton key={`${playername} ${item.stat}`} projection={item}/>
      })}
{/* 
      <StatSelectButton stat="assists" target={30.5} confidence={67.5}/>
      <StatSelectButton stat="Fantasy" target={30.5} confidence={67.5}/>
      <StatSelectButton stat="PTS + REB + AST" target={30.5} confidence={67.5}/>
      <StatSelectButton stat="Free Throws" target={30.5} confidence={67.5}/>
      <StatSelectButton stat="3 Pointers" target={30.5} confidence={67.5}/>
      <StatSelectButton stat="Blocks + Steals" target={30.5} confidence={67.5}/>
      <StatSelectButton stat="Double Double" target={30.5} confidence={67.5}/> */}
      {/* <button className="stat-select-btn">
          <p className="bold">Other</p>
      </button> */}
    </div>
  )
}

export default StatSelectBtns