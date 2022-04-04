import React, {useState} from 'react'
import {GetColor} from '../utils'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faAngleDoubleUp as up, faAngleDoubleDown as down, faEquals} from '@fortawesome/free-solid-svg-icons'

const Prediction = (props) => {
  const {projection} = props
  //TODO: Use constant stat type mappings to get the selected projection
  const {prediction, confidence} = projection;
  const [target, setTarget] = useState(projection.target);
  function onChange(e) {
    setTarget(e.target.value)
  }
  return (
    <div className="prediction">
            <p>TARGET: <input type="number" disabled={true} min={0} max={100} step={0.5} value={target} onChange={onChange}/></p>
            <PredictionIcon confidence={confidence} prediction={prediction}/>
        </div>
  )
}

const PredictionIcon = (props) => {
  const {confidence, prediction} = props;

  return (<div className="prediction-icon">
  <FontAwesomeIcon className={`arrow ${GetColor("pct", confidence)}`} icon={up}/>
  <p className={`bold tall prediction-result`}>{prediction}</p>
  <p className={`${GetColor("pct", confidence)}`}>{confidence}%</p>
</div>)
}

const PredictionIconSmall = (props) => {
  const {confidence, prediction} = props;
  const icon = confidence > 40 && confidence < 60 ? faEquals : prediction === "OVER" ? up : down;
  return (
  <div className="hide">
      <FontAwesomeIcon className={GetColor("pct", confidence)} icon={icon}/>
      <p className={GetColor("pct", confidence)}>{confidence}%</p>
  </div>)
}

export {PredictionIcon, PredictionIconSmall}
export default Prediction