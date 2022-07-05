import React from "react";
import { GetColor, GetPropScore } from "../utils";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faForward as up,
  faForward as down,
  faPlay as uncertain,
} from "@fortawesome/free-solid-svg-icons";

const Prediction = (props) => {
  const { propositions, predictions, selected, game } = props;

  const projection = predictions.filter(
    (item) =>
      item.stat.recognize.toLowerCase() === selected.recognize.toLowerCase()
  )[0];
  const score =
    game && projection.target
      ? GetPropScore(game, projection.stat.recognize)
      : "";
  const actual =
    score !== "" ? (score > projection.target ? "OVER" : "UNDER") : "";
  const { overUnder, confidence } = projection;
  return (
    <div className="prediction">
      <p>{projection.stat.label}</p>
      {propositions.map((prop) => {
        return prop.type.toLowerCase() === projection.stat.recognize ? (
          <p>
            {prop.sportsbook} {prop.target}
          </p>
        ) : (
          <></>
        );
      })}
      <p>
        TARGET:{" "}
        <input
          type="number"
          disabled={true}
          min={0}
          max={100}
          step={0.5}
          value={projection.target}
        />
      </p>
      {game ? (
        <p className={"actual-result"}>
          ACTUAL:{" "}
          <span className={score !== "" ? GetColor("OVER", actual) : ""}>
            {GetPropScore(game, projection.stat.recognize)}
          </span>
        </p>
      ) : (
        <></>
      )}
      <PredictionIcon
        confidence={confidence}
        overUnder={overUnder}
        actual={actual}
      />
    </div>
  );
};

const PredictionIcon = (props) => {
  const { confidence, overUnder, actual } = props;
  return (
    <div className="prediction-icon">
      <FontAwesomeIcon
        className={`arrow ${GetColor("pct", confidence)}`}
        icon={getIcon(confidence, overUnder)}
        rotation={getRotation(overUnder)}
      />
      <p className={`bold tall prediction-result`}>{overUnder}</p>
      <div>
        <p className={`${GetColor("pct", confidence)}`}>{confidence}%</p>
        {actual ? (
          <p className={actual === overUnder ? "high" : "low"}>
            {actual === overUnder ? "CORRECT" : "INCORRECT"}
          </p>
        ) : (
          <></>
        )}
      </div>
    </div>
  );
};

const PredictionIconSmall = (props) => {
  const { confidence, overUnder } = props;
  return (
    <div className="hide">
      <FontAwesomeIcon
        className={GetColor("pct", confidence)}
        icon={getIcon(confidence, overUnder)}
        rotation={getRotation(overUnder)}
      />
      <p className={GetColor("pct", confidence)}>{confidence}%</p>
    </div>
  );
};

function getIcon(confidence, overUnder) {
  return confidence < 60 ? uncertain : overUnder === "OVER" ? up : down;
}

function getRotation(overUnder) {
  return overUnder === "OVER" ? 270 : 90;
}

export { PredictionIcon, PredictionIconSmall };
export default Prediction;
