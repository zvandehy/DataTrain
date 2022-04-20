import React from 'react'
import {PredictionIconSmall} from './Prediction'


const StatSelectBtns = (props) => {
  const {predictions,playername, onStatSelect, selected} = props;
  return (
    <div className="stat-select">
      {predictions.map(item => {
        return <StatSelectButton key={`${playername} ${item.stat.label}`} prediction={item} selected={item.stat.label === selected} onStatSelect={onStatSelect}/>
      })}
      {/* <button className="stat-select-btn span-3">
          <p className="bold">Other</p>
      </button> */}
    </div>
  )
}
const StatSelectButton = (props) => {
  const {prediction, onStatSelect, selected} = props;
  const {stat, target, overUnder, confidence} = prediction;
  return (
    <button className={`stat-select-btn ${selected ? 'selected' : ''}`} onClick={() => onStatSelect(stat.label)}>
        <p className="bold titlecase">{stat.label}</p>
        {target && confidence ? <>
          <p className="hide">T: {target}</p>
          <PredictionIconSmall confidence={confidence} overUnder={overUnder}/>
        </>: <></>}
    </button>
  )
}

export default StatSelectBtns