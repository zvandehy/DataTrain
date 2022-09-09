import { Collapse, Table, TableBody, TableHead } from "@mui/material";
import { ColorCompare } from "../../shared/functions/color.fn";
import { Game } from "../../shared/interfaces/graphql/game.interface";
import { Proposition } from "../../shared/interfaces/graphql/projection.interface";
import { ScoreType } from "../../shared/interfaces/score-type.enum";
import { Stat, Minutes } from "../../shared/interfaces/stat.interface";
import { StyledTableCell } from "../styled-table/styled-table-cell.component";
import { StyledTableRow } from "../styled-table/styled-table-row.component";

interface ExpandTeamTableProps {
  selectedProp: Proposition;
  vsTeams: Game[];
  open: boolean;
}

const ExpandTeamTable: React.FC<ExpandTeamTableProps> = ({
  selectedProp,
  vsTeams,
  open,
}: ExpandTeamTableProps) => {
  return (
    <StyledTableRow>
      <StyledTableCell colSpan={7} sx={{ borderBottom: "none" }}>
        <Collapse in={open} unmountOnExit>
          <Table>
            <TableHead>
              <StyledTableCell>DATE</StyledTableCell>
              <StyledTableCell>OPP</StyledTableCell>
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
            <TableBody
              sx={{
                "& tr:nth-child(even)": {
                  backgroundColor: "rgba(0, 0, 0, 0.2)",
                },
              }}
            >
              {vsTeams.map((game) => {
                return (
                  <StyledTableRow
                    key={game.gameID + " " + game.player.playerID}
                  >
                    <StyledTableCell>{game.date}</StyledTableCell>
                    <StyledTableCell>
                      {game.opponent.abbreviation}
                    </StyledTableCell>
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
                        selectedProp.statType.score(game) - selectedProp.target,
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
                        selectedProp.statType.score(game) - selectedProp.target
                      ).toFixed(2)}`}
                    </StyledTableCell>
                    <StyledTableCell
                      className={ColorCompare(
                        selectedProp.statType.score(game) - selectedProp.target,
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
  );
};

export default ExpandTeamTable;
