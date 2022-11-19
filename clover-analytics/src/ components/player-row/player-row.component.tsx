import { ExpandLess } from "@material-ui/icons";
// import { ExpandCircleDown } from "@mui/icons-material";
// import {
//   TableRow,
//   TableCell,
//   IconButton,
//   Avatar,
//   Grid,
//   useTheme,
// } from "@mui/material";
// import { borderBottomColor } from "@mui/system";
// import { index } from "mathjs";
// import moment from "moment";
// import { useState } from "react";
// import {
//   Game,
//   PropositionA,
// } from "../../shared/interfaces/graphql/game.interface";
// import { Player } from "../../shared/interfaces/graphql/player.interface";
// import {
//   GetStatAbbreviation,
//   ConvertMinutes,
// } from "../../shared/interfaces/stat.interface";
// import { COLORS } from "../../shared/styles/constants";
// import { HitMissIcon } from "../icons/hit-miss.icon";
// import { OverUnderIcon } from "../icons/overUnderIcon.component";

// interface PlayerRowProps {
//   player: Player;
//   game: Game;
// }

// const PlayerRow: React.FC<PlayerRowProps> = ({
//   player,
//   game,
// }: PlayerRowProps) => {
//   const theme = useTheme();
//   const [expand, setExpand] = useState(false);
//   return (
//     <>
//       {" "}
//       {game.prediction.propositions.map((prop, i) => {
//         const borderBottomColor =
//           (i === 0 && !expand) ||
//           (i === game.prediction.propositions.length - 1 && expand)
//             ? "white"
//             : "inherit";
//         return (
//           <TableRow
//             key={player.playerID + player.name}
//             sx={{
//               "& *": { verticalAlign: "middle" },
//               "&:hover": {
//                 cursor: "pointer",
//                 backgroundColor: theme.palette.grey[900],
//               },
//             }}
//             onClick={() => console.log("click", player.name)}
//           >
//             {i === 0 ? (
//               <>
//                 <TableCell sx={{ borderTop: "1px solid" }}>
//                   <IconButton onClick={() => setExpand(!expand)}>
//                     {!expand ? <ExpandCircleDown /> : <ExpandLess />}
//                   </IconButton>
//                 </TableCell>
//                 <TableCell sx={{ borderTop: "1px solid" }}>
//                   <Grid container alignItems={"center"} columnGap={1}>
//                     {
//                       <Avatar
//                         sx={{
//                           borderRadius: 5,
//                           width: 48,
//                           height: 48,
//                           bgcolor: COLORS.AVATAR,
//                         }}
//                         src={player.image}
//                         alt={"player" + { player }}
//                       />
//                     }
//                     {player.name}
//                   </Grid>
//                 </TableCell>
//                 <TableCell sx={{ borderTop: "1px solid" }}>
//                   {player.team.abbreviation}
//                 </TableCell>
//                 <TableCell sx={{ borderTop: "1px solid" }}>
//                   {(game.home_or_away === "home" ? "vs" : "@") +
//                     " " +
//                     game.opponent.abbreviation}
//                 </TableCell>
//               </>
//             ) : i === 1 && expand ? (
//               <TableCell
//                 rowSpan={game.prediction.propositions.length - 1} // should be # of types (or number of props displayed)
//                 colSpan={4}
//                 sx={{ borderBottom: "1px solid white" }}
//               />
//             ) : (
//               <></>
//             )}
//             {i === 0 || expand ? (
//               <>
//                 <TableCell
//                   sx={{
//                     borderLeft: "1px solid",
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   {prop.sportsbook === "PrizePicks" ? (
//                     <img
//                       src="https://assets.website-files.com/5c777b25f39dfe05e88c3dda/6112a587690c49744507c59c_icon.png"
//                       alt="PrizePicks"
//                       width="16px"
//                       style={{ verticalAlign: "middle" }}
//                     />
//                   ) : (
//                     prop.sportsbook
//                   )}
//                   {`${moment(prop.lastModified).format("M/D h:mm a")}`}
//                 </TableCell>
//                 <TableCell
//                   sx={{
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   {GetStatAbbreviation(prop.type)}
//                 </TableCell>
//                 <TableCell
//                   sx={{
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   {prop.target}
//                 </TableCell>
//                 <TableCell
//                   sx={{
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   {prop.estimation}
//                 </TableCell>
//                 <TableCell
//                   sx={{
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   <OverUnderIcon overUnder={prop.prediction} />
//                   {(prop.predictionTargetDiffPCT > 0 ? "+" : "") +
//                     prop.predictionTargetDiffPCT +
//                     "%"}
//                 </TableCell>
//                 <TableCell
//                   sx={{
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   {ConvertMinutes(game.prediction.weightedTotal.minutes)}
//                 </TableCell>
//                 <TableCell
//                   sx={{
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   {prop.estimationPerMin}
//                 </TableCell>
//                 <TableCell
//                   sx={{
//                     borderLeft: "1px solid",
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   {game.outcome !== "PENDING" ? (
//                     <>
//                       {prop.actual} <HitMissIcon outcome={prop.predictionHit} />
//                     </>
//                   ) : (
//                     "TBD"
//                   )}
//                 </TableCell>
//                 <TableCell
//                   sx={{
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   {game.outcome !== "PENDING" ? game.minutes : "TBD"}
//                 </TableCell>
//                 <TableCell
//                   sx={{
//                     borderBottom: `1px solid ${borderBottomColor}`,
//                   }}
//                 >
//                   {game.outcome !== "PENDING"
//                     ? prop.actualPerMin
//                     : // +
//                       //   " (" +
//                       //   (prop.actualDiffPerMin > 0
//                       //     ? "+" + prop.actualDiffPerMinPCT
//                       //     : prop.actualDiffPerMinPCT) +
//                       //   "%)"
//                       "TBD"}
//                 </TableCell>
//               </>
//             ) : (
//               <></>
//             )}
//           </TableRow>
//         );
//       })}
//     </>
//   );
// };

// export default PlayerRow;
