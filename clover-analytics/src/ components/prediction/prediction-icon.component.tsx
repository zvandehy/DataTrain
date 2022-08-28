import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faForward as up,
  faForward as down,
  faPlay as uncertain,
} from "@fortawesome/free-solid-svg-icons";

import "./prediction.component.css";
import { Prediction } from "../../shared/interfaces/graphql/projection.interface";

interface PredictionIconProps {
  prediction: Prediction;
  actual?: string;
}

const PredictionIcon: React.FC<PredictionIconProps> = ({
  prediction,
  actual,
}: PredictionIconProps) => {
  return (
    <div className="prediction-icon">
      <FontAwesomeIcon
        className={`arrow ${
          prediction.confidence >= 60
            ? "high"
            : prediction.confidence <= 40
            ? "low"
            : "med"
        }`}
        icon={getIcon(prediction.confidence, prediction.overUnderPrediction)}
        rotation={getRotation(prediction.overUnderPrediction)}
      />
      <div>
        <p className={`bold tall prediction-result`}>
          {prediction.overUnderPrediction.toUpperCase()}
        </p>
        <p
          className={`${
            prediction.confidence >= 60
              ? "high"
              : prediction.confidence <= 40
              ? "low"
              : "med"
          }`}
        >
          {prediction.confidence}%
        </p>
        {actual ? (
          <p
            className={
              actual.toUpperCase() === "PUSH"
                ? "med"
                : actual === prediction.overUnderPrediction
                ? "high"
                : "low"
            }
          >
            {actual.toUpperCase() === "PUSH"
              ? "PUSH"
              : actual === prediction.overUnderPrediction
              ? "CORRECT"
              : "INCORRECT"}
          </p>
        ) : (
          <></>
        )}
      </div>
    </div>
  );
};

export default PredictionIcon;

const PredictionIconSmall: React.FC<PredictionIconProps> = ({
  prediction,
}: PredictionIconProps) => {
  return (
    <div className="hide">
      <FontAwesomeIcon
        className={
          prediction.confidence >= 60
            ? "high"
            : prediction.confidence <= 40
            ? "low"
            : "med"
        }
        icon={getIcon(prediction.confidence, prediction.overUnderPrediction)}
        rotation={getRotation(prediction.overUnderPrediction)}
      />
      <p
        className={
          prediction.confidence >= 60
            ? "high"
            : prediction.confidence <= 40
            ? "low"
            : "med"
        }
      >
        {prediction.confidence}%
      </p>
    </div>
  );
};

function getIcon(confidence: number, overUnderPrediction: string) {
  return confidence < 60
    ? uncertain
    : overUnderPrediction.toUpperCase() === "OVER"
    ? up
    : down;
}

function getRotation(overUnderPrediction: string) {
  return overUnderPrediction.toUpperCase() === "OVER" ? 270 : 90;
}

export { PredictionIcon, PredictionIconSmall };
