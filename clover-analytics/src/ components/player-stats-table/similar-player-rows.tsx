import { useState } from "react";
import { ColorCompare } from "../../shared/functions/color.fn";
import {
  Projection,
  Proposition,
} from "../../shared/interfaces/graphql/projection.interface";
import { ScoreType } from "../../shared/interfaces/score-type.enum";
import { Stat, Minutes } from "../../shared/interfaces/stat.interface";
import { StyledTableCell } from "../styled-table/styled-table-cell.component";
import { StyledTableRow } from "../styled-table/styled-table-row.component";

interface ExpandSimilarPlayerRowsProps {
  selectedProp: Proposition;
  projection: Projection;
}

const ExpandSimilarPlayerRows: React.FC<ExpandSimilarPlayerRowsProps> = ({
  selectedProp,
  projection,
}: ExpandSimilarPlayerRowsProps) => {
  if (selectedProp.customPrediction.similarPlayersVsOpponent === undefined) {
    return <StyledTableRow></StyledTableRow>;
  }
  return (
    <>
      {selectedProp.customPrediction.similarPlayersVsOpponent?.similarGames.map(
        (game) => {
          const seasonAvg = selectedProp.statType.average(
            projection.player.similarPlayers!.find(
              (player) => player.name === game.player.name
            )!.games
          );
          return (
            <StyledTableRow key={game.gameID + " " + game.player.playerID}>
              <StyledTableCell>{game.player.name}</StyledTableCell>
              <StyledTableCell>{seasonAvg}</StyledTableCell>
              <StyledTableCell>{game.date}</StyledTableCell>
              <StyledTableCell>{game.opponent.abbreviation}</StyledTableCell>
              <StyledTableCell>
                {selectedProp.statType.score(game)}
              </StyledTableCell>
              <StyledTableCell>
                {" "}
                {selectedProp.statType.scorePer(game, ScoreType.PerMin)}
              </StyledTableCell>
              {selectedProp.statType.relatedStats?.map((related: Stat) => {
                return (
                  <StyledTableCell key={related.label + "score"}>
                    {related.score(game)}
                  </StyledTableCell>
                );
              })}
              <StyledTableCell>{Minutes.score(game)}</StyledTableCell>
              <StyledTableCell
                className={ColorCompare(
                  selectedProp.statType.score(game) - seasonAvg,
                  0
                )}
              >
                {`${
                  +(selectedProp.statType.score(game) - seasonAvg).toFixed(2) >
                  0
                    ? "+"
                    : ""
                }${(selectedProp.statType.score(game) - seasonAvg).toFixed(2)}`}
              </StyledTableCell>
              <StyledTableCell
                className={ColorCompare(
                  selectedProp.statType.score(game) - seasonAvg,
                  0
                )}
              >{`${
                selectedProp.statType.score(game) > seasonAvg
                  ? "OVER"
                  : selectedProp.statType.score(game) < seasonAvg
                  ? "UNDER"
                  : "PUSH"
              }`}</StyledTableCell>
            </StyledTableRow>
          );
        }
      )}
    </>
  );
};

export default ExpandSimilarPlayerRows;
