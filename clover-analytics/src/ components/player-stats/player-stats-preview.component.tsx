import {
  Paper,
  styled,
  Table,
  TableBody,
  TableCell,
  tableCellClasses,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";
import React from "react";
import { ColorPct } from "../../shared/functions/color.fn";
import {
  Projection,
  Proposition,
} from "../../shared/interfaces/graphql/projection.interface";
import "./player-stats-preview.component.css";

interface PlayerStatsPreviewProps {
  selectedProp: Proposition;
  projection: Projection;
}

const PlayerStatsPreview: React.FC<PlayerStatsPreviewProps> = ({
  selectedProp,
  projection,
}: PlayerStatsPreviewProps) => {
  return (
    <TableContainer className={"player-stats"} component={Paper}>
      <Table aria-label="simple table">
        <TableHead>
          <StyledTableRow>
            <StyledTableCell
              sx={{ "*": { color: "inherit" } }}
            ></StyledTableCell>
            <StyledTableCell sx={{ "*": { color: "inherit" } }}>
              AVG
            </StyledTableCell>
            <StyledTableCell sx={{ "*": { color: "inherit" } }}>
              O-U-P
            </StyledTableCell>
            <StyledTableCell sx={{ "*": { color: "inherit" } }}>
              OVER %
            </StyledTableCell>
          </StyledTableRow>
        </TableHead>
        <TableBody>
          {selectedProp.customPrediction.predictionFragments.map((fragment) => (
            <StyledTableRow
              className={"player-row"}
              key={`${projection.player.name} ${fragment.numGames}`}
              sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
            >
              <StyledTableCell>Last {fragment.numGames}</StyledTableCell>
              <StyledTableCell>{fragment.average}</StyledTableCell>
              <StyledTableCell
                className={`${ColorPct(fragment.pctOver)}`}
              >{`${fragment.numOver}-${fragment.numUnder}-${fragment.numPush}`}</StyledTableCell>
              <StyledTableCell className={`${ColorPct(fragment.pctOver)}`}>
                {(fragment.pctOver * 100).toFixed(2)}%
              </StyledTableCell>
            </StyledTableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default PlayerStatsPreview;
//TODO: Implement theming
const StyledTableCell = styled(TableCell)(({ theme }) => ({
  [`&.${tableCellClasses.head}`]: {
    backgroundColor: "inherit",
    padding: "15px min(1.5vw, 1.5rem);",
    paddingBottom: "0px",
    color: "inherit",
    fontSize: "unset",
  },
  [`&.${tableCellClasses.body}`]: {
    padding: "min(1vw, 1.5rem);",
    fontSize: "1rem",
    color: "inherit",
  },
}));

const StyledTableRow = styled(TableRow)(({ theme }) => ({
  "&:nth-of-type(odd)": {
    backgroundColor: "inherit",
  },
  // hide last border
  "&:last-child td, &:last-child th": {
    border: 0,
  },
}));
