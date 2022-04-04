import React from 'react'
import {round} from 'mathjs'
import {GetColor} from '../utils'

const LastNGames = (props) => {
    const {countObj, header} = props;
    const {avg, n, over, under} = countObj;
    const pct_over = round((over/n)*100,2);
  return (
    <>
    <span className="header last-n-games">Last {n}</span>
    <span className="player-stat">{avg}</span>
    <span className={`player-stat-result ${header === "PCT" ? GetColor("pct", pct_over) : ""}`}>{header === "PCT" ? pct_over : header === "OVER" ? over : under}%</span>
    </>
  )
}

export default LastNGames