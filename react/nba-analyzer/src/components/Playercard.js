import React, {useState, useEffect} from 'react'
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
    target:10.5,
    prediction:"OVER",
    confidence:85.5, 
    //TODO: generate counts data using counts option & seasonData
    counts: [{n:70, avg:11, over:40, under:25},{n:30, avg:12, over:18, under:12},{n:10, avg:8, over:3, under:6},{n:5, avg:9, over:3, under:2}]
  }})

  //TODO: Get the type from the selected button (default to highest confidence or preference from filters)
  const [stat, setStat] = useState("points");
  const [selectedProjection, setProjection] = useState(projections.find(item => item.stat.toLowerCase() === stat.toLowerCase()));
  useEffect(() => {
    const projection = projections.find(item => item.stat.toLowerCase() === stat.toLowerCase());
    setProjection(projection)
  }, [stat] //TODO: projections in the dependencies causes a loop that crashes the site for too many re-renders
  );

  return (
    <div className="playercard">
        <PlayerContext player={player} opponent={playerProp.opponent}></PlayerContext>
        <StatSelectBtns projections={projections} playername={player.name}></StatSelectBtns>
        <Prediction projection={selectedProjection}></Prediction>
        <PlayerStatsPreview playerCounts={selectedProjection.counts} matchupData={""} similarData={""}></PlayerStatsPreview>
    </div>
  )
}

export default Playercard