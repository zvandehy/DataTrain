import React from 'react'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faAngleDoubleUp as up } from '@fortawesome/free-solid-svg-icons'
import { GetColor } from '../utils';

const StatSelectButton = (props) => {
    const {stat, target, confidence} = props;
  return (
    <button className="stat-select-btn">
        <p className="bold titlecase">{stat}</p>
        <p className="hide">T: {target}</p>
        <div className="hide">
            <FontAwesomeIcon className={GetColor("pct", confidence)} icon={up}/>
            <p className={GetColor("pct", confidence)}>{confidence}%</p>
        </div>
    </button>
  )
}

export default StatSelectButton