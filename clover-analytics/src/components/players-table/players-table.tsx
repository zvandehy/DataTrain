import * as React from "react";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Paper from "@mui/material/Paper";
import { Player } from "../../shared/interfaces/graphql/player.interface";
import "./players-table.css";
import { ALL_STATS } from "../../shared/constants";
import PlayerCell from "./player-cell";

interface PlayersTableProps {
  players: Player[];
}

const PlayersTable: React.FC<PlayersTableProps> = ({
  players,
}: PlayersTableProps) => {
  return (
    <TableContainer component={Paper}>
      <Table sx={{ minWidth: 650, maxWidth: "100%" }} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>Player</TableCell>
            {ALL_STATS.map((stat) => (
              <TableCell key={`${stat.abbreviation}-header`}>
                {stat.abbreviation}
              </TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {players.map((player) => {
            if (player.games?.length === 0) {
              return <></>;
            }
            return (
              <TableRow
                className={"player-row"}
                key={player.name}
                sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
              >
                <PlayerCell player={player} />
                {ALL_STATS.map((stat) => {
                  return (
                    <TableCell key={`${player.name}-${stat.abbreviation}`}>
                      {`${stat.score(player.games[0])}, ${stat.average(
                        player.games
                      )}`}
                    </TableCell>
                  );
                })}
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default PlayersTable;
