import { KeyboardArrowUp, KeyboardArrowDown } from "@material-ui/icons";
import {
  Collapse,
  IconButton,
  Table,
  TableBody,
  TableHead,
} from "@mui/material";
import { ColorCompare, ColorPct } from "../../shared/functions/color.fn";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
import { Proposition } from "../../shared/interfaces/graphql/game.interface";
import {
  PredictionFragment,
  Projection,
} from "../../shared/interfaces/graphql/projection.interface";
import { ScoreType } from "../../shared/interfaces/score-type.enum";
import { Stat, Minutes } from "../../shared/interfaces/stat.interface";
import { StyledTableCell } from "../styled-table/styled-table-cell.component";
import { StyledTableRow } from "../styled-table/styled-table-row.component";

interface ExpandFragmentTableProps {
  selectedProp: Proposition;
  open: boolean;
  fragment: PredictionFragment;
  projection: Projection;
  customModel: CustomCalculation;
  setOpen: (open: boolean) => void;
}

const ExpandFragmentTable: React.FC<ExpandFragmentTableProps> = ({
  selectedProp,
  open,
  fragment,
  projection,
  customModel,
  setOpen,
}: ExpandFragmentTableProps) => {
  return (
    <>
      <StyledTableRow
        className={"player-row"}
        key={`${projection.player.name} ${fragment.games.length}`}
        sx={{
          "&:last-child td, &:last-child th": { border: 0 },
          "& tr:nth-child(even)": {
            backgroundColor: "rgba(0, 0, 0, 0.2)",
          },
        }}
      >
        <StyledTableCell>
          <IconButton
            color="secondary"
            aria-label="expand row"
            size="small"
            onClick={() => setOpen(!open)}
          >
            {open ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
          </IconButton>
        </StyledTableCell>
        <StyledTableCell>Last {fragment.games.length}</StyledTableCell>
        <StyledTableCell>{fragment.average}</StyledTableCell>
        <StyledTableCell>{fragment.avgPerMin}</StyledTableCell>
        {selectedProp.statType.relatedStats?.map((related: Stat) => {
          return (
            <StyledTableCell key={related.label + "score"}>
              {related.average(fragment.games)}
            </StyledTableCell>
          );
        })}
        <StyledTableCell>{fragment.minutes}</StyledTableCell>
        <StyledTableCell
          className={ColorCompare(fragment.average - selectedProp.target, 0)}
        >
          {`${
            +(fragment.average - selectedProp.target).toFixed(2) > 0 ? "+" : ""
          }${(fragment.average - selectedProp.target).toFixed(2)}`}
        </StyledTableCell>
        <StyledTableCell
          className={`${ColorPct(fragment.pctOver)}`}
        >{`${fragment.numOver}-${fragment.numUnder}-${fragment.numPush}`}</StyledTableCell>

        {customModel.includePush ? (
          <StyledTableCell className={`${ColorPct(fragment.pctPushOrMore)}`}>
            {(fragment.pctPushOrMore * 100).toFixed(2)}%
          </StyledTableCell>
        ) : (
          <StyledTableCell className={`${ColorPct(fragment.pctOver)}`}>
            {(fragment.pctOver * 100).toFixed(2)}%
          </StyledTableCell>
        )}

        {customModel.includePush ? (
          <StyledTableCell className={`${ColorPct(fragment.pctPushOrLess)}`}>
            {(fragment.pctPushOrLess * 100).toFixed(2)}%
          </StyledTableCell>
        ) : (
          <StyledTableCell className={`${ColorPct(fragment.pctUnder)}`}>
            {(fragment.pctUnder * 100).toFixed(2)}%
          </StyledTableCell>
        )}
        <StyledTableCell>{fragment.weight.toFixed(0)}%</StyledTableCell>
      </StyledTableRow>
      <StyledTableRow>
        <StyledTableCell colSpan={100} sx={{ borderBottom: "none" }}>
          <Collapse in={open} unmountOnExit>
            <Table>
              <TableHead>
                <StyledTableCell></StyledTableCell>
                <StyledTableCell></StyledTableCell>
                <StyledTableCell></StyledTableCell>
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
                <StyledTableCell>OVER/UNDER</StyledTableCell>
              </TableHead>
              <TableBody>
                {fragment.games.map((game) => {
                  return (
                    <StyledTableRow>
                      <StyledTableCell colSpan={2}></StyledTableCell>
                      <StyledTableCell>{game.date}</StyledTableCell>
                      <StyledTableCell>
                        {selectedProp.statType.score(game)}
                      </StyledTableCell>
                      <StyledTableCell>
                        {" "}
                        {selectedProp.statType.scorePer(game, ScoreType.PerMin)}
                      </StyledTableCell>
                      {selectedProp.statType.relatedStats?.map(
                        (related: Stat) => {
                          return (
                            <StyledTableCell key={related.label + "score"}>
                              {related.score(game)}
                            </StyledTableCell>
                          );
                        }
                      )}
                      <StyledTableCell>{Minutes.score(game)}</StyledTableCell>
                      <StyledTableCell
                        className={ColorCompare(
                          selectedProp.statType.score(game) -
                            selectedProp.target,
                          0
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
                          selectedProp.statType.score(game) -
                          selectedProp.target
                        ).toFixed(2)}`}
                      </StyledTableCell>
                      <StyledTableCell
                        className={ColorCompare(
                          selectedProp.statType.score(game) -
                            selectedProp.target,
                          0
                        )}
                      >{`${
                        selectedProp.statType.score(game) > selectedProp.target
                          ? "OVER"
                          : selectedProp.statType.score(game) <
                            selectedProp.target
                          ? "UNDER"
                          : "PUSH"
                      }`}</StyledTableCell>
                    </StyledTableRow>
                  );
                })}
              </TableBody>
            </Table>
          </Collapse>
        </StyledTableCell>
      </StyledTableRow>
    </>
  );
};

export default ExpandFragmentTable;
