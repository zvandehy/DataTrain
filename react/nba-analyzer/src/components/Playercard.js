import React, {useState, useEffect, useRef} from 'react'
import PlayerContext from "./Playercontext"
import StatSelectBtns from "./Statselectbtns"
import Prediction from "./Prediction"
import PlayerStatsPreview from "./Playerstatspreview"
import { GetPropScore } from '../utils'
import {round, mean} from "mathjs"

const Playercard = (props) => {
  let {confidence, playerProp} = props;
  const {player} = playerProp;
  //TODO: move to utils (or a filters.js) as function
  let seasonData = player.games.filter((game) => game.season === "2021-22").sort(function(a, b) {
      var c = new Date(a.date);
      var d = new Date(b.date);
      return c-d;
  });
  //TODO: add counts option to custom filters
  const counts = [0, -30, -10, -5];
  const weights = [.30,.27,.25,.18]; //TODO: determine best weights to use
  //TODO: create a constant list of stat types, shorthand, and backend score mapping
  const stats = ["Points","Assists","3 Pointers","Free Throws","PTS + REB + AST","Rebounds","Fantasy", "Double Double", "Blocks + Steals"];
  //TODO: Get the targets (if it exists) and the prediction & confidence for each stat type
  let projections = stats.map(item => {
    const target = getTarget(playerProp.targets, item);
    const playerStats = getStats(seasonData, counts, item, target);
    const predictionAndConfidence = getPredictionAndConfidence(playerStats, weights);
    return {
    stat:item, 
    target: target,
    prediction:predictionAndConfidence[0],
    confidence:predictionAndConfidence[1], 
    //TODO: generate counts data using counts option & seasonData
    counts: playerStats,
  }});

  //TODO: Make this more sophisticated
  function getPredictionAndConfidence(stats, weights) {
    let conf_o = 0;
    let conf_u = 0;
    stats.forEach((item,i)=> {
      conf_o += item.pct_o * weights[i];
      conf_u += item.pct_u * weights[i];
    })
    conf_o = round(conf_o,2);
    conf_u = round(conf_u,2);
    if (conf_o > conf_u) {
      return ["OVER", conf_o];
    } return ["UNDER", conf_u]
  }

  function getTarget(targets, stat) {
    //TODO: Update this using the constant list of types and score mappings to get all the types
    const exists = targets.filter(item => item.type.toLowerCase() === stat.toLowerCase());
    if (exists.length === 1) {
      return exists[0].target
    } 
    return 0;
  }

  function getStats(games, counts, stat, target) {
    let stats = [];
    const scores = games.map(game => GetPropScore(game, stat));
    counts.forEach(count => {
      //TODO: Apply game stat filters
      const data = scores.slice(count);
      const avg = round(mean(data),2);
      const over = data.filter(score => score > target).length;
      const under = data.filter(score => score < target).length;
      const pct_o = round((over / data.length)*100, 2);
      const pct_u = round((under / data.length)*100, 2);
      stats.push({n:data.length, avg:avg, over:over, under:under, pct_o:pct_o, pct_u:pct_u})
    })
    return stats
  }

  const matchups = seasonData.filter((game) => game.opponent.teamID === playerProp.opponent.teamID);

  //TODO: Default to highest confidence or preference from filters
  const [stat, setStat] = useState("Points");
  function onStatSelect(stat) {
    setStat(stat);
  }

  return (
    <div className="playercard">
        <PlayerContext player={player} opponent={playerProp.opponent}></PlayerContext>
        <StatSelectBtns projections={projections} playername={player.name} selected={stat} onStatSelect={onStatSelect}></StatSelectBtns>
        <Prediction projections={projections} selected={stat}></Prediction>
        <PlayerStatsPreview projections={projections} selected={stat} matchups={matchups} similarData={""}></PlayerStatsPreview>
    </div>
  )
}

export default Playercard