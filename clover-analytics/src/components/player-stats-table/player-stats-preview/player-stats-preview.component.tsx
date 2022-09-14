import {
  Paper,
  Table,
  TableBody,
  TableContainer,
  TableHead,
} from "@mui/material";
import React from "react";
import { ColorCompare, ColorPct } from "../../../shared/functions/color.fn";
import {
  Projection,
  Proposition,
} from "../../../shared/interfaces/graphql/projection.interface";
import { StyledTableCell } from "../../styled-table/styled-table-cell.component";
import { StyledTableRow } from "../../styled-table/styled-table-row.component";
import SimilarPreview from "./similar-preview.component";
import "./player-stats-preview.component.css";
import { CustomCalculation } from "../../../shared/interfaces/custom-prediction.interface";

interface PlayerStatsPreviewProps {
  selectedProp: Proposition;
  projection: Projection;
  customModel: CustomCalculation;
}

const PlayerStatsPreview: React.FC<PlayerStatsPreviewProps> = ({
  selectedProp,
  projection,
  customModel,
}: PlayerStatsPreviewProps) => {
  return (
    <div className={"player-stats"}>
      <TableContainer component={Paper}>
        <Table aria-label="simple table">
          <TableHead>
            <StyledTableRow>
              <StyledTableCell></StyledTableCell>
              <StyledTableCell>AVG</StyledTableCell>
              <StyledTableCell>AVG/MIN</StyledTableCell>
              <StyledTableCell>MINS</StyledTableCell>
              <StyledTableCell>DIFF</StyledTableCell>
              <StyledTableCell>O-U-P</StyledTableCell>
              {customModel.includePush ? (
                <StyledTableCell>OVER / PUSH %</StyledTableCell>
              ) : (
                <StyledTableCell>OVER %</StyledTableCell>
              )}

              {customModel.includePush ? (
                <StyledTableCell>UNDER / PUSH %</StyledTableCell>
              ) : (
                <StyledTableCell>UNDER %</StyledTableCell>
              )}
              <StyledTableCell>Weight</StyledTableCell>
            </StyledTableRow>
          </TableHead>
          <TableBody>
            {selectedProp.customPrediction.recencyFragments.map((fragment) => (
              <StyledTableRow
                className={"player-row"}
                key={`${projection.player.name} ${fragment.games.length}`}
                sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
              >
                <StyledTableCell>Last {fragment.games.length}</StyledTableCell>
                <StyledTableCell>{fragment.average}</StyledTableCell>
                <StyledTableCell>{fragment.avgPerMin}</StyledTableCell>
                <StyledTableCell>{fragment.minutes}</StyledTableCell>
                <StyledTableCell
                  className={ColorCompare(
                    fragment.average - selectedProp.target,
                    0
                  )}
                >
                  {`${
                    +(fragment.average - selectedProp.target).toFixed(2) > 0
                      ? "+"
                      : ""
                  }${(fragment.average - selectedProp.target).toFixed(2)}`}
                </StyledTableCell>
                <StyledTableCell
                  className={`${ColorPct(fragment.pctOver)}`}
                >{`${fragment.numOver}-${fragment.numUnder}-${fragment.numPush}`}</StyledTableCell>

                {customModel.includePush ? (
                  <StyledTableCell
                    className={`${ColorPct(fragment.pctPushOrMore)}`}
                  >
                    {(fragment.pctPushOrMore * 100).toFixed(2)}%
                  </StyledTableCell>
                ) : (
                  <StyledTableCell className={`${ColorPct(fragment.pctOver)}`}>
                    {(fragment.pctOver * 100).toFixed(2)}%
                  </StyledTableCell>
                )}

                {customModel.includePush ? (
                  <StyledTableCell
                    className={`${ColorPct(fragment.pctPushOrLess)}`}
                  >
                    {(fragment.pctPushOrLess * 100).toFixed(2)}%
                  </StyledTableCell>
                ) : (
                  <StyledTableCell className={`${ColorPct(fragment.pctUnder)}`}>
                    {(fragment.pctUnder * 100).toFixed(2)}%
                  </StyledTableCell>
                )}
                <StyledTableCell>{fragment.weight.toFixed(0)}%</StyledTableCell>
              </StyledTableRow>
            ))}
            {/* {projection.player.games.map((game) => {
              if (game.opponent.teamID === projection.opponent.teamID) {
                return (
                  <StyledTableRow>
                    <StyledTableCell className={"flex-cell"}>
                      <span>vs {projection.opponent.abbreviation}</span>
                      <span>{moment(game.date).format("MM/DD/YY")}</span>
                    </StyledTableCell>
                    <StyledTableCell>
                      {selectedProp.statType.score(game)}
                    </StyledTableCell>
                    <StyledTableCell>
                      {selectedProp.statType.scorePer(game, ScoreType.PerMin)}
                    </StyledTableCell>
                    <StyledTableCell>
                      {ConvertMinutes(game.minutes)}
                    </StyledTableCell>
                    <StyledTableCell
                      className={ColorCompare(
                        selectedProp.statType.score(game),
                        selectedProp.target
                      )}
                    >
                      {`${
                        +(
                          selectedProp.statType.score(game) -
                          selectedProp.target
                        ).toFixed(2) > 0
                          ? "+"
                          : ""
                      }${(
                        selectedProp.statType.score(game) - selectedProp.target
                      ).toFixed(2)}`}
                    </StyledTableCell>
                  </StyledTableRow>
                );
              }
              return <></>;
            })} */}
          </TableBody>
        </Table>
      </TableContainer>
      {customModel.opponentWeight &&
      selectedProp.customPrediction.vsOpponent &&
      projection.player.games.some(
        (game) =>
          game.opponent.abbreviation === projection.opponent.abbreviation
      ) ? (
        <SimilarPreview
          projection={projection}
          selectedProp={selectedProp}
          header={`${projection.player.name} vs ${projection.opponent.abbreviation}`}
          sim={selectedProp.customPrediction.vsOpponent}
          customModel={customModel}
        />
      ) : (
        <></>
      )}
      {projection.opponent.similarTeams?.length > 0 &&
      selectedProp.customPrediction.vsSimilarTeams ? (
        <SimilarPreview
          projection={projection}
          selectedProp={selectedProp}
          header={`${projection.player.name} vs ${selectedProp.customPrediction.vsSimilarTeams.similarCount} Similar Teams`}
          sim={selectedProp.customPrediction.vsSimilarTeams}
          customModel={customModel}
        />
      ) : (
        <></>
      )}
      {projection.player.similarPlayers?.length > 0 &&
      selectedProp.customPrediction.similarPlayersVsOpponent ? (
        <SimilarPreview
          projection={projection}
          selectedProp={selectedProp}
          header={`${selectedProp.customPrediction.similarPlayersVsOpponent.similarCount} Similar Players vs ${projection.opponent.abbreviation}`}
          sim={selectedProp.customPrediction.similarPlayersVsOpponent}
          customModel={customModel}
        />
      ) : (
        <></>
      )}
    </div>
  );
};

export default PlayerStatsPreview;
