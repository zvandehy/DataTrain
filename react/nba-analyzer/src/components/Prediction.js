import React, {useState} from 'react'
import {GetColor} from '../utils'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faAngleDoubleUp as up } from '@fortawesome/free-solid-svg-icons'

const Prediction = (props) => {
  const {targets, prediction, confidence} = props
  const [target, setTarget] = useState(targets[0]);
  return (
    <div className="prediction">
            <p>TARGET: <input type="number" disabled={true} min={0} max={100} step={0.5} value={target} onChange={(e) => setTarget(e.target.value)}/></p>
            <div className="prediction-icon">
                <FontAwesomeIcon className={`arrow ${GetColor("pct", confidence)}`} icon={up}/>
                <p className={`bold tall prediction-result`}>{prediction}</p>
                <p className={`${GetColor("pct", confidence)}`}>{confidence}%</p>
            </div>
        </div>
  )
}

export default Prediction