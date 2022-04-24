import React, { useState } from "react";
import PlayerContext from "./Playercontext";
import StatSelectBtns from "./Statselectbtns";
import Prediction from "./Prediction";
import { PlayerStatsPreview } from "./PlayerStats";
import { ParseDate, CompareDates } from "../utils";
import { CalculatePredictions, GetHighestConfidence } from "../predictions.js";

const Playercard = (props) => {
  const { projection, player, date, statPreference } = props;
  //TODO: move to utils (or a filters.js) as function
  let seasonData = player.games
    .filter(
      (game) =>
        game.season === "2021-22" &&
        CompareDates(date, ParseDate(game.date)) > 1
    )
    .sort(CompareDates);

  const game = player.games.find((game) => game.date === date);
  let predictions = CalculatePredictions(projection, seasonData);

  const matchups = seasonData.filter(
    (game) => game.opponent.teamID === projection.opponent.teamID
  );

  const initialStat =
    statPreference === ""
      ? GetHighestConfidence(predictions).stat
      : statPreference;

  const [stat, setStat] = useState(initialStat);
  function onStatSelect(stat) {
    setStat(stat);
  }
  React.useEffect(() => {
    if (statPreference !== "" && statPreference !== "ANY") {
      setStat(statPreference);
    } else {
      setStat(initialStat);
    }
  }, [initialStat, statPreference]);
  return (
    <div className="playercard">
      <PlayerContext
        player={player}
        opponent={projection.opponent}
      ></PlayerContext>
      <StatSelectBtns
        predictions={predictions}
        playername={player.name}
        selected={stat}
        onStatSelect={onStatSelect}
      ></StatSelectBtns>
      <Prediction
        predictions={predictions}
        selected={stat}
        game={game}
      ></Prediction>
      <PlayerStatsPreview
        predictions={predictions}
        selected={stat}
        matchups={matchups}
      ></PlayerStatsPreview>
    </div>
  );
};

export default Playercard;
