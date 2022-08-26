import moment from "moment";
import React from "react";
import { BETTING_CATEGORIES } from "../../../../shared/constants";
import { Proposition } from "../../../../shared/interfaces/graphql/projection.interface";
import { Stat } from "../../../../shared/interfaces/stat.interface";
import { PredictionIconSmall } from "../../../prediction/prediction-icon.component";
import "./stat-select-buttons.component.css";

interface StatSelectButtonsProps {
  propositions: Proposition[];
  selectedStat: Stat;
  onStatSelect: (stat: Stat) => void;
}

const StatSelectButtons: React.FC<StatSelectButtonsProps> = ({
  propositions,
  selectedStat: selected,
  onStatSelect,
}: StatSelectButtonsProps) => {
  return (
    <div className="stat-select">
      {BETTING_CATEGORIES.map((stat) => {
        const proposition = propositions.find((prop) => prop.statType === stat);
        return (
          <StatSelectButton
            key={`${stat.label}`}
            proposition={proposition}
            stat={stat}
            isSelected={stat.label === selected.label}
            onStatSelect={onStatSelect}
          />
        );
      })}
    </div>
  );
};

export default StatSelectButtons;

interface StatSelectButtonProps {
  onStatSelect: (stat: Stat) => void;
  proposition: Proposition | undefined;
  stat: Stat;
  isSelected: boolean;
}

const StatSelectButton: React.FC<StatSelectButtonProps> = ({
  onStatSelect,
  proposition,
  stat,
  isSelected,
}: StatSelectButtonProps) => {
  return (
    <button
      className={`stat-select-btn ${isSelected ? "selected" : ""}`}
      onClick={() => onStatSelect(stat)}
    >
      <p className="bold titlecase">{stat.abbreviation}</p>
      {proposition?.target && proposition.customPrediction.confidence ? (
        <>
          <p className="hide">T: {proposition.target}</p>
          <PredictionIconSmall prediction={proposition.customPrediction} />
        </>
      ) : (
        <></>
      )}
    </button>
  );
};
