import React from 'react'
import PlayerContext from "./Playercontext"
import StatSelectBtns from "./Statselectbtns"
import Prediction from "./Prediction"
import PlayerStatsPreview from "./Playerstatspreview"

const Playercard = (props) => {
  let {confidence} = props;
  if (!confidence) {
    confidence = 67.5;
  }
  return (
    <div className="playercard">
        <PlayerContext></PlayerContext>
        <StatSelectBtns></StatSelectBtns>
        <Prediction targets={[30.5]} prediction="OVER" confidence={confidence}></Prediction>
        <PlayerStatsPreview></PlayerStatsPreview>
    </div>
  )
}

export default Playercard