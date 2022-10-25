import { Collapse, Table, TableBody, TableHead } from "@mui/material";
import { ColorCompare } from "../../shared/functions/color.fn";
import { Proposition } from "../../shared/interfaces/graphql/game.interface";
import { Projection } from "../../shared/interfaces/graphql/projection.interface";
import { ScoreType } from "../../shared/interfaces/score-type.enum";
import { Stat, Minutes } from "../../shared/interfaces/stat.interface";
import { StyledTableCell } from "../styled-table/styled-table-cell.component";
import { StyledTableRow } from "../styled-table/styled-table-row.component";

interface ExpandSimilarPlayersProps {
  selectedProp: Proposition;
  projection: Projection;
  open: boolean;
}

const ExpandSimilarPlayers: React.FC<ExpandSimilarPlayersProps> = ({
  selectedProp,
  projection,
  open,
}: ExpandSimilarPlayersProps) => {
  return (
    <StyledTableRow>
      <StyledTableCell colSpan={100} sx={{ borderBottom: "none" }}>
        <Collapse in={open} unmountOnExit>
          <Table>
            <TableHead>
              <StyledTableCell>PLAYER (AVG)</StyledTableCell>
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
              {selectedProp.customPrediction.similarPlayersVsOpponent?.similarGames.map(
                (game, i) => {
                  const seasonAvg = selectedProp.statType.average(
                    projection.player.similarPlayers!.find(
                      (player) => player.name === game.player.name
                    )!.games
                  );

                  return (
                    <>
                      {i === 0 ||
                      selectedProp.customPrediction.similarPlayersVsOpponent
                        ?.similarGames[i - 1].player.name !==
                        game.player.name ? (
                        <StyledTableRow
                          key={
                            game.gameID + " " + game.player.playerID + " head"
                          }
                          sx={{ fontWeight: "bold" }}
                        >
                          <StyledTableCell
                            colSpan={0}
                            sx={{ fontWeight: "bold" }}
                          >
                            {game.player.name} ({seasonAvg})
                          </StyledTableCell>
                          <StyledTableCell colSpan={100}></StyledTableCell>
                        </StyledTableRow>
                      ) : (
                        <></>
                      )}
                      <StyledTableRow
                        key={game.gameID + " " + game.player.playerID}
                      >
                        <StyledTableCell>
                          {game.team.abbreviation}{" "}
                          {game.home_or_away === "home" ? "vs" : "@"}{" "}
                          {game.opponent.abbreviation} {game.date}
                        </StyledTableCell>
                        <StyledTableCell>
                          {selectedProp.statType.score(game)}
                        </StyledTableCell>
                        <StyledTableCell>
                          {" "}
                          {selectedProp.statType.scorePer(
                            game,
                            ScoreType.PerMin
                          )}
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
                            selectedProp.statType.score(game) - seasonAvg,
                            0
                          )}
                        >
                          {`${
                            +(
                              selectedProp.statType.score(game) - seasonAvg
                            ).toFixed(2) > 0
                              ? "+"
                              : ""
                          }${(
                            selectedProp.statType.score(game) - seasonAvg
                          ).toFixed(2)}`}
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
                    </>
                  );
                }
              )}
            </TableBody>
          </Table>
        </Collapse>
      </StyledTableCell>
    </StyledTableRow>
  );
};

export default ExpandSimilarPlayers;
