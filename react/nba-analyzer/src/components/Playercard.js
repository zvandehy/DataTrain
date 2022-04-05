import React, {useState, useEffect, useRef} from 'react'
import PlayerContext from "./Playercontext"
import StatSelectBtns from "./Statselectbtns"
import Prediction from "./Prediction"
import PlayerStatsPreview from "./Playerstatspreview"

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
  const weights = [.40,.25,.2,.15]; //TODO: determine best weights to use
  //TODO: create a constant list of stat types, shorthand, and backend score mapping
  const stats = ["Points","Assists","3 Pointers","Free Throws","PTS + REB + AST","Rebounds","Fantasy", "Double Double", "Blocks + Steals"];
  //TODO: Get the targets (if it exists) and the prediction & confidence for each stat type
  let projections = stats.map(item => {return {
    stat:item, 
    target: getTarget(playerProp.targets, item),
    prediction:"OVER",
    confidence:85.5, 
    //TODO: generate counts data using counts option & seasonData
    counts: [{n:70, avg:11, over:40, under:25},{n:30, avg:12, over:18, under:12},{n:10, avg:8, over:3, under:6},{n:5, avg:9, over:3, under:2}],
  }})

  function getTarget(targets, stat) {
    //TODO: Update this using the constnt list of types and score mappings to get all the types
    const exists = targets.filter(item => item.type.toLowerCase() === stat.toLowerCase());
    if (exists.length === 1) {
      return exists[0].target
    } 
    return 0;
  }

  const matchups = seasonData.filter((game) => game.opponent.teamID === playerProp.opponent.teamID);

  //TODO: Get the type from the selected button (default to highest confidence or preference from filters)
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