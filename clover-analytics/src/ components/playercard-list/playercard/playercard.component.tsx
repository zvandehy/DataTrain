import { useEffect, useState } from "react";
import {
  GetMaxConfidence,
  UpdatePropositionWithPrediction,
} from "../../../shared/functions/predictions.fn";
import {
  Projection,
  Proposition,
} from "../../../shared/interfaces/graphql/projection.interface";
import { Stat } from "../../../shared/interfaces/stat.interface";
import PlayerStatsPreview from "../../player-stats/player-stats-preview.component";
import Prediction from "../../prediction/prediction.component";
import PlayercardContext from "./playercard-context/playercard-context.component";
import "./playercard.component.css";
import StatSelectButtons from "./stat-select-buttons/stat-select-buttons.component";

interface PlayerCardProps {
  projection: Projection;
  filteredStatType: Stat | undefined;
}

function PreselectProp(
  filteredStatType: Stat | undefined,
  projection: Projection
): Proposition {
  if (filteredStatType !== undefined) {
    return projection.propositions.find(
      (p) => p.statType === filteredStatType
    )!;
  }
  return GetMaxConfidence(projection.propositions);
}

const PlayerCard: React.FC<PlayerCardProps> = ({
  projection,
  filteredStatType,
}: PlayerCardProps) => {
  useEffect(() => {
    onPropSelect(PreselectProp(filteredStatType, projection));
  }, [projection, filteredStatType]);

  const [selectedProp, selectProp] = useState(
    PreselectProp(filteredStatType, projection)
  );
  const [stat, statSelect] = useState(selectedProp.statType);
  const onPropSelect = (prop: Proposition) => {
    selectProp(prop);
    statSelect(prop.statType);
  };

  const onStatSelect = (stat: Stat) => {
    statSelect(stat);
    let customTarget = 0;
    let customProp: Proposition = {
      target: customTarget || 0,
      statType: stat,
      type: stat.label,
      sportsbook: "None",
      predictions: [],
      customPrediction: {
        model: "Custom",
        overUnderPrediction: "",
        confidence: 0,
        totalPrediction: 0,
        predictionFragments: [],
      },
    };
    const foundProp =
      projection.propositions.find((p) => p.statType === stat) ||
      UpdatePropositionWithPrediction(customProp, projection.player.games);
    selectProp(foundProp);
  };
  return (
    <div className="playercard">
      <PlayercardContext projection={projection} />
      <StatSelectButtons
        propositions={projection.propositions}
        selectedStat={stat}
        onStatSelect={onStatSelect}
      />
      <Prediction
        projection={projection}
        selectedProp={selectedProp}
        selectedStat={stat}
        onPredictionSelect={onPropSelect}
      ></Prediction>
      <PlayerStatsPreview
        selectedProp={selectedProp}
        projection={projection}
      ></PlayerStatsPreview>
    </div>
  );
};

export default PlayerCard;
