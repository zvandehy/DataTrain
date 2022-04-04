import React from 'react'
import {PredictionIconSmall} from './Prediction'

const StatSelectButton = (props) => {
  const {stat, target, prediction, confidence} = props.projection;
  return (
    <button className="stat-select-btn">
        <p className="bold titlecase">{stat}</p>
        {target && confidence ? <>
          <p className="hide">T: {target}</p>
          <PredictionIconSmall confidence={confidence} prediction={prediction}/>
        </>: <></>}
    </button>
  )
}

export default StatSelectButton