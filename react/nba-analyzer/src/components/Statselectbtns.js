import React from 'react'
import {PredictionIconSmall} from './Prediction'


const StatSelectBtns = (props) => {
  const {projections,playername, onStatSelect, selected} = props;
  return (
    <div className="stat-select">
      {projections.map(item => {
        return <StatSelectButton key={`${playername} ${item.stat.label}`} projection={item} selected={item.stat.label === selected} onStatSelect={onStatSelect}/>
      })}
      <button className="stat-select-btn span-3">
          <p className="bold">Other</p>
      </button>
    </div>
  )
}
const StatSelectButton = (props) => {
  const {projection, onStatSelect, selected} = props;
  const {stat, target, prediction, confidence} = projection;
  return (
    <button className={`stat-select-btn ${selected ? 'selected' : ''}`} onClick={() => onStatSelect(stat.label)}>
        <p className="bold titlecase">{stat.label}</p>
        {target && confidence ? <>
          <p className="hide">T: {target}</p>
          <PredictionIconSmall confidence={confidence} prediction={prediction}/>
        </>: <></>}
    </button>
  )
}

export default StatSelectBtns