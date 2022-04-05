import React from 'react'
import {GetColor} from '../utils'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faAngleDoubleUp as up, faAngleDoubleDown as down, faEquals} from '@fortawesome/free-solid-svg-icons'

const Prediction = (props) => {
  const {projections,selected} = props
  const projection = projections.filter(item => item.stat.toLowerCase() === selected.toLowerCase())[0];
  //TODO: Use constant stat type mappings to get the selected projection
  const {prediction, confidence} = projection;
  // const [target, setTarget] = useState(projection.target);
  // function onChange(e) {
  //   setTarget(e.target.value)
  // }
  return (
    <div className="prediction">
      {/* TODO: Fix state of changing projection onChange={onChange}*/}
            <p>TARGET: <input type="number" disabled={true} min={0} max={100} step={0.5} value={projection.target} /></p>
            <PredictionIcon confidence={confidence} prediction={prediction}/>
        </div>
  )
}

const PredictionIcon = (props) => {
  const {confidence, prediction} = props;
  return (<div className="prediction-icon">
  <FontAwesomeIcon className={`arrow ${GetColor("pct", confidence)}`} icon={getIcon(confidence,prediction)}/>
  <p className={`bold tall prediction-result`}>{prediction}</p>
  <p className={`${GetColor("pct", confidence)}`}>{confidence}%</p>
</div>)
}

const PredictionIconSmall = (props) => {
  const {confidence, prediction} = props;
  return (
  <div className="hide">
      <FontAwesomeIcon className={GetColor("pct", confidence)} icon={getIcon(confidence,prediction)}/>
      <p className={GetColor("pct", confidence)}>{confidence}%</p>
  </div>)
}

function getIcon(confidence, prediction) {
  return confidence >= 50 && confidence < 60 ? faEquals : prediction === "OVER" ? up : down;
}

export {PredictionIcon, PredictionIconSmall}
export default Prediction