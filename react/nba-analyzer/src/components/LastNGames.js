import React from 'react'
import {GetColor} from '../utils'

const LastNGames = (props) => {
    const {countObj, header} = props;
    const {avg, n, over, under, pct_o} = countObj;
  return (
    <>
    <span className="header last-n-games">Last {n}</span>
    <span className="player-stat">{avg}</span>
    <span className={`player-stat-result ${header === "PCT" ? GetColor("pct", pct_o) : ""}`}>{header === "PCT" ? pct_o : header === "OVER" ? over : under}%</span>
    </>
  )
}

export default LastNGames