import { ExpandLess } from "@material-ui/icons";
import { ExpandCircleDown } from "@mui/icons-material";
import {
  Avatar,
  Grid,
  IconButton,
  TableCell,
  TableRow,
  useTheme,
} from "@mui/material";
import moment from "moment";
import { useState } from "react";
import { Proposition } from "../../shared/interfaces/graphql/proposition.interface";
import { GetStatAbbreviation } from "../../shared/interfaces/stat.interface";
import { COLORS } from "../../shared/styles/constants";
import { HitMissIcon } from "../icons/hit-miss.icon";
import { OverUnderIcon } from "../icons/overUnderIcon.component";

interface PlayerRowProps {
  propositions: Proposition[];
  onClick: (prop: Proposition) => void;
}

const PlayerRow: React.FC<PlayerRowProps> = ({
  propositions,
  onClick,
}: PlayerRowProps) => {
  const theme = useTheme();
  const [expand, setExpand] = useState(false);
  const game = propositions[0].game;
  const player = propositions[0].game.player;
  return (
    <>
      {propositions.map((prop, i) => {
        const borderBottomColor = "white";
        //   (i === 0 && !expand) ||
        //   (i === game.prediction.propositions.length - 1 && expand)
        //     ? "white"
        //     : "inherit";
        return (
          <TableRow
            key={player.playerID + player.name + prop.type + prop.target}
            sx={{
              "& *": { verticalAlign: "middle" },
              "&:hover": {
                cursor: "pointer",
                backgroundColor: theme.palette.grey[900],
              },
            }}
            onClick={() => onClick(prop)}
          >
            {i === 0 ? (
              <>
                <TableCell sx={{ borderTop: "1px solid" }}>
                  <IconButton onClick={() => setExpand(!expand)}>
                    {!expand ? <ExpandCircleDown /> : <ExpandLess />}
                  </IconButton>
                </TableCell>
                <TableCell sx={{ borderTop: "1px solid" }}>
                  <Grid container alignItems={"center"} columnGap={1}>
                    {
                      <Avatar
                        sx={{
                          borderRadius: 5,
                          width: 48,
                          height: 48,
                          bgcolor: COLORS.AVATAR,
                        }}
                        src={player.image}
                        alt={"player" + { player }}
                      />
                    }
                    {player.name}
                  </Grid>
                </TableCell>
                <TableCell sx={{ borderTop: "1px solid" }}>
                  {player.team.abbreviation}
                </TableCell>
                <TableCell sx={{ borderTop: "1px solid" }}>
                  {(game.home_or_away === "home" ? "vs" : "@") +
                    " " +
                    game.opponent.abbreviation}
                </TableCell>
              </>
            ) : i === 1 && expand ? (
              <TableCell
                rowSpan={propositions.length - 1} // should be # of types (or number of props displayed)
                colSpan={4}
                sx={{ borderBottom: "1px solid white" }}
              />
            ) : (
              <></>
            )}
            {i === 0 || expand ? (
              <>
                <TableCell
                  sx={{
                    borderLeft: "1px solid",
                    borderBottom: `1px solid ${borderBottomColor}`,
                  }}
                >
                  {prop.sportsbook === "PrizePicks" ? (
                    <img
                      src="https://assets.website-files.com/5c777b25f39dfe05e88c3dda/6112a587690c49744507c59c_icon.png"
                      alt="PrizePicks"
                      width="16px"
                      style={{ verticalAlign: "middle" }}
                    />
                  ) : (
                    prop.sportsbook
                  )}
                  {`${moment(prop.lastModified).format("M/D h:mm a")}`}
                </TableCell>
                <TableCell
                  sx={{
                    borderBottom: `1px solid ${borderBottomColor}`,
                  }}
                >
                  {GetStatAbbreviation(prop.type)}
                </TableCell>
                <TableCell
                  sx={{
                    borderBottom: `1px solid ${borderBottomColor}`,
                  }}
                >
                  {prop.target}
                </TableCell>
                <TableCell
                  sx={{
                    borderBottom: `1px solid ${borderBottomColor}`,
                  }}
                >
                  {prop.prediction.estimation}
                </TableCell>
                <TableCell
                  sx={{
                    borderBottom: `1px solid ${borderBottomColor}`,
                  }}
                >
                  <OverUnderIcon overUnder={prop.prediction.wager} />
                  {prop.prediction.significance + "%"}
                </TableCell>
                <TableCell
                  sx={{
                    borderBottom: `1px solid ${borderBottomColor}`,
                  }}
                >
                  {prop.prediction.stdDev}
                </TableCell>
                <TableCell
                  sx={{
                    borderBottom: `1px solid ${borderBottomColor}`,
                  }}
                >
                  {prop.prediction.cumulativeOver +
                    prop.prediction.cumulativeUnder +
                    prop.prediction.cumulativePush}
                </TableCell>
                <TableCell
                  sx={{
                    borderLeft: "1px solid",
                    borderBottom: `1px solid ${borderBottomColor}`,
                  }}
                >
                  {game.outcome !== "PENDING" ? (
                    <>
                      {prop.actualResult}{" "}
                      <HitMissIcon outcome={prop.prediction.wagerOutcome} />
                    </>
                  ) : (
                    "TBD"
                  )}
                </TableCell>
              </>
            ) : (
              <></>
            )}
          </TableRow>
        );
      })}
    </>
  );
};

export default PlayerRow;
