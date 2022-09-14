import { Player } from "../../shared/interfaces/graphql/player.interface";
import TableCell from "@mui/material/TableCell";
import getPlayerPhotoUrl from "../../shared/functions/photos.fn";

interface PlayerCellProps {
  player: Player;
}

const PlayerCell: React.FC<PlayerCellProps> = ({ player }: PlayerCellProps) => {
  const league = "WNBA";
  return (
    <TableCell className={"player-cell"} component="th" scope="row">
      <img
        loading="lazy"
        className="player-photo"
        alt={player.name}
        src={getPlayerPhotoUrl(player.playerID, league)}
      ></img>
      {player.name}
    </TableCell>
  );
};

export default PlayerCell;
