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
  SimilarPlayerCalculation,
  SimilarTeamCalculation,
} from "../../../shared/functions/predictions.fn";
import {
  Projection,
  Proposition,
} from "../../../shared/interfaces/graphql/projection.interface";
import { StyledTableCell } from "../../styled-table/styled-table-cell.component";
import { StyledTableRow } from "../../styled-table/styled-table-row.component";
import SimilarPreview from "./similar-preview.component";
import "./player-stats-preview.component.css";
import { ConvertMinutes } from "../../../shared/interfaces/stat.interface";
import moment from "moment";
import { ScoreType } from "../../../shared/interfaces/score-type.enum";

interface PlayerStatsPreviewProps {
  selectedProp: Proposition;
  projection: Projection;
}

const PlayerStatsPreview: React.FC<PlayerStatsPreviewProps> = ({
  selectedProp,
  projection,
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
              <StyledTableCell>OVER %</StyledTableCell>
            </StyledTableRow>
          </TableHead>
          <TableBody>
            {selectedProp.customPrediction.predictionFragments.map(
              (fragment) => (
                <StyledTableRow
                  className={"player-row"}
                  key={`${projection.player.name} ${fragment.numGames}`}
                  sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
                >
                  <StyledTableCell>Last {fragment.numGames}</StyledTableCell>
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
                  <StyledTableCell className={`${ColorPct(fragment.pctOver)}`}>
                    {(fragment.pctOver * 100).toFixed(2)}%
                  </StyledTableCell>
                </StyledTableRow>
              )
            )}
            {projection.player.games.map((game) => {
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
            })}
          </TableBody>
        </Table>
      </TableContainer>
      {projection.opponent.similarTeams?.length > 0 ? (
        <SimilarPreview
          projection={projection}
          selectedProp={selectedProp}
          header={`${projection.player.name} vs ${projection.opponent.similarTeams?.length} Similar Teams`}
          sim={SimilarTeamCalculation(projection, selectedProp)}
        />
      ) : (
        <></>
      )}
      {projection.player.similarPlayers?.length > 0 ? (
        <SimilarPreview
          projection={projection}
          selectedProp={selectedProp}
          header={`${projection.player.similarPlayers?.length} Similar Players vs ${projection.opponent.abbreviation}`}
          sim={SimilarPlayerCalculation(
            projection,
            selectedProp,
            selectedProp.statType.average(projection.player.games)
          )}
        />
      ) : (
        <></>
      )}
    </div>
  );
};

export default PlayerStatsPreview;
