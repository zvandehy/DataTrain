import { KeyboardArrowUp, KeyboardArrowDown } from "@material-ui/icons";
import {
  IconButton,
  Paper,
  Table,
  TableBody,
  TableContainer,
  TableHead,
} from "@mui/material";
import { ColorCompare, ColorPct } from "../../shared/functions/color.fn";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
import { Proposition } from "../../shared/interfaces/graphql/game.interface";
import { SimilarCalculation } from "../../shared/interfaces/similarCalculation.interface";
import { Minutes, Stat } from "../../shared/interfaces/stat.interface";
import { StyledTableCell } from "../styled-table/styled-table-cell.component";
import { StyledTableRow } from "../styled-table/styled-table-row.component";
import "./player-stats-preview/similar-preview.component.css";

interface SimilarTableProps {
  selectedProp: Proposition;
  header: string;
  sim: SimilarCalculation;
  customModel: CustomCalculation;
  children: React.ReactNode;
  setOpen: (open: boolean) => void;
  open: boolean;
}

const SimilarTable: React.FC<SimilarTableProps> = ({
  selectedProp,
  header,
  sim,
  customModel,
  children,
  setOpen,
  open,
}: SimilarTableProps) => {
  return (
    <TableContainer className={"player-stats"} component={Paper}>
      <Table aria-label="similar teams table">
        <TableHead>
          <StyledTableRow>
            <StyledTableCell
              colSpan={11 + (selectedProp.statType.relatedStats?.length ?? 0)}
            >
              {header}
            </StyledTableCell>
          </StyledTableRow>
          <StyledTableRow>
            <StyledTableCell rowSpan={2}>
              <IconButton
                aria-label="expand row"
                size="small"
                onClick={() => setOpen(!open)}
                color={"secondary"}
              >
                {open ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
              </IconButton>
            </StyledTableCell>
            <StyledTableCell>GAMES</StyledTableCell>
            <StyledTableCell>
              {selectedProp.statType.abbreviation}
            </StyledTableCell>
            <StyledTableCell>
              {selectedProp.statType.abbreviation}/MIN
            </StyledTableCell>
            {selectedProp.statType.relatedStats?.map((related: Stat) => {
              return (
                <StyledTableCell key={related.label}>
                  {related.abbreviation}
                </StyledTableCell>
              );
            })}
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
            <StyledTableCell>ADJ.</StyledTableCell>
            <StyledTableCell>Weight</StyledTableCell>
          </StyledTableRow>
        </TableHead>
        <TableBody>
          <StyledTableRow>
            <StyledTableCell></StyledTableCell>
            <StyledTableCell>{`${sim.similarGames.length}`}</StyledTableCell>
            <StyledTableCell>{sim.similarAvg}</StyledTableCell>
            <StyledTableCell>{sim.similarAvgPerMin}</StyledTableCell>
            {selectedProp.statType.relatedStats?.map((related: Stat) => {
              return (
                <StyledTableCell key={related.label + "score"}>
                  {related.average(sim.similarGames)}
                </StyledTableCell>
              );
            })}
            <StyledTableCell>
              {Minutes.average(sim.similarGames)}
            </StyledTableCell>
            <StyledTableCell className={ColorCompare(sim.similarDiff, 0)}>
              {`${sim.similarDiff > 0 ? "+" : ""}${sim.similarDiff}`}
            </StyledTableCell>
            <StyledTableCell
              className={`${ColorPct(sim.simOverPct)}`}
            >{`${sim.countSimOver}-${sim.countSimUnder}-${sim.countSimPush}`}</StyledTableCell>
            {customModel.includePush ? (
              <StyledTableCell
                className={`${ColorPct(sim.simOverPct + sim.simPushPct)}`}
              >{`${sim.simOverPct + sim.simPushPct}%`}</StyledTableCell>
            ) : (
              <StyledTableCell
                className={`${ColorPct(sim.simOverPct)}`}
              >{`${sim.simOverPct}%`}</StyledTableCell>
            )}

            {customModel.includePush ? (
              <StyledTableCell
                className={`${ColorPct(sim.simUnderPct + sim.simPushPct)}`}
              >{`${sim.simUnderPct + sim.simPushPct}%`}</StyledTableCell>
            ) : (
              <StyledTableCell
                className={`${ColorPct(sim.simUnderPct)}`}
              >{`${sim.simUnderPct}%`}</StyledTableCell>
            )}

            <StyledTableCell
              className={ColorCompare(sim.playerAvgAdj, selectedProp.target)}
            >
              {sim.playerAvgAdj}
            </StyledTableCell>
            <StyledTableCell>{sim.weight.toFixed(0)}%</StyledTableCell>
          </StyledTableRow>
          {children}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default SimilarTable;
