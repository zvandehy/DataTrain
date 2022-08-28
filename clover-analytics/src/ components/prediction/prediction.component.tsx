import { InputLabel, ListSubheader, MenuItem, Select } from "@mui/material";
import moment from "moment";
import React from "react";
import { PropositionResult } from "../../shared/functions/predictions.fn";
import {
  Projection,
  Proposition,
} from "../../shared/interfaces/graphql/projection.interface";
import { Stat } from "../../shared/interfaces/stat.interface";
import PredictionIcon from "./prediction-icon.component";
import "./prediction.component.css";

interface PredictionProps {
  projection: Projection;
  selectedProp: Proposition | undefined;
  selectedStat: Stat;
  onPredictionSelect: (proposition: Proposition) => void;
  customTarget?: number;
}

const Prediction: React.FC<PredictionProps> = ({
  projection,
  selectedProp,
  selectedStat,
  customTarget,
  onPredictionSelect,
}: PredictionProps) => {
  let actual: string = "";
  if (projection.result && selectedProp) {
    actual = PropositionResult(selectedProp, projection.result);
  }
  const uniqueSportsbooks: string[] = [
    ...Array.from(
      new Set(
        projection.propositions.map((proposition) => proposition.sportsbook)
      )
    ),
  ];
  return (
    <div className="prediction">
      <div className={"selected-stat-display"}>
        <div>{selectedStat.display}</div>
        <div>{selectedProp?.sportsbook}</div>
        <div>
          {" on "}
          {moment(selectedProp?.lastModified).format("MM-DD [at] hh:mm a")}
        </div>
      </div>
      <div>
        <div className="target-input">
          <InputLabel
            id="target-label"
            className={"target-label"}
            sx={{
              fontSize: "inherit",
              display: "inline",
              color: "inherit",
            }}
          >
            Target:
          </InputLabel>
          <Select
            className={"target-select"}
            labelId="target-label"
            id="target-select"
            sx={{
              // color: theme.palette.text.primary,
              // border: "1px solid var(--color-accent)",
              paddingRight: "0.25rem",
              paddingLeft: "0.25rem",
              "&:after, &:before": {
                borderColor: "var(--color-accent)",
              },
              input: {
                padding: "0.5rem",
              },
            }}
            inputProps={{
              classes: {
                input: {
                  fontSize: "20rem",
                },
              },
            }}
            value={selectedProp?.target || 0}
            onChange={(event) => {
              console.log(event);
              onPredictionSelect(
                projection.propositions.filter(
                  (prop) =>
                    prop.target === Number(event.target.value) &&
                    prop.statType === selectedStat
                )[0]
              );
            }}
          >
            {uniqueSportsbooks.map((groupedSportsbook) => {
              let items: Proposition[] = projection.propositions;
              if (
                selectedProp &&
                !items.some((prop) => prop.statType === selectedProp.statType)
              )
                items.push(selectedProp);
              const menuItems = items
                .filter(
                  (prop) =>
                    prop.statType === selectedStat &&
                    prop.sportsbook === groupedSportsbook
                )
                .map((prop) => {
                  return (
                    <MenuItem
                      key={`${prop.sportsbook} ${prop.type} ${prop.target}`}
                      value={prop.target}
                      onClick={() => onPredictionSelect(prop)}
                      sx={{
                        color: "var(--color-black)",
                      }}
                    >
                      {prop.target}
                    </MenuItem>
                  );
                });
              return [
                <ListSubheader key={groupedSportsbook}>
                  {groupedSportsbook}
                </ListSubheader>,
                menuItems,
              ];
            })}
          </Select>
        </div>
        {projection.result ? (
          <div className={"actual-result"}>
            ACTUAL:
            <span
              className={
                "actual-score " +
                (actual === "Over"
                  ? "high"
                  : actual === "Under"
                  ? "low"
                  : "med")
              }
            >
              {selectedStat.score(projection.result)}
            </span>
          </div>
        ) : (
          <></>
        )}
        {selectedProp ? (
          <PredictionIcon
            prediction={selectedProp.customPrediction}
            actual={actual}
          />
        ) : (
          <></>
        )}
      </div>
    </div>
  );
};

export default Prediction;

// const Prediction = (props) => {
//   const { propositions, predictions, selected, game } = props;

//   const projection = predictions.filter(
//     (item) =>
//       item.stat.recognize.toLowerCase() === selected.recognize.toLowerCase()
//   )[0];
//   const score =
//     game && projection.target
//       ? GetPropScore(game, projection.stat.recognize)
//       : "";
//   const actual =
//     score !== "" ? (score > projection.target ? "OVER" : "UNDER") : "";
//   const { overUnder, confidence } = projection;
//   return (

//   );
// };
