import {
  Paper,
  Table,
  TableBody,
  TableContainer,
  TableHead,
} from "@mui/material";
// import { ColorCompare, ColorPct } from "../../../shared/functions/color.fn";
// import { CustomCalculation } from "../../../shared/interfaces/custom-prediction.interface";
// import { Proposition } from "../../../shared/interfaces/graphql/game.interface";
// import { Projection } from "../../../shared/interfaces/graphql/proposition.interface";
// import { SimilarCalculation } from "../../../shared/interfaces/similarCalculation.interface";
// import { Minutes } from "../../../shared/interfaces/stat.interface";
// import { StyledTableCell } from "../../styled-table/styled-table-cell.component";
// import { StyledTableRow } from "../../styled-table/styled-table-row.component";
// import "./similar-preview.component.css";

// interface SimilarPreviewProps {
//   projection: Projection;
//   selectedProp: Proposition;
//   header: string;
//   sim: SimilarCalculation;
//   customModel: CustomCalculation;
// }

// const SimilarPreview: React.FC<SimilarPreviewProps> = ({
//   projection,
//   selectedProp,
//   header,
//   sim,
//   customModel,
// }: SimilarPreviewProps) => {
//   return (
//     <TableContainer className={"player-stats"} component={Paper}>
//       <Table aria-label="similar teams table">
//         <TableHead>
//           <StyledTableRow>
//             <StyledTableCell colSpan={10}>{header}</StyledTableCell>
//           </StyledTableRow>
//           <StyledTableRow>
//             <StyledTableCell>GAMES</StyledTableCell>
//             <StyledTableCell>AVG</StyledTableCell>
//             <StyledTableCell>AVG/MIN</StyledTableCell>
//             <StyledTableCell>MINS</StyledTableCell>
//             <StyledTableCell>DIFF</StyledTableCell>
//             <StyledTableCell>O-U-P</StyledTableCell>
//             {customModel.includePush ? (
//               <StyledTableCell>OVER / PUSH %</StyledTableCell>
//             ) : (
//               <StyledTableCell>OVER %</StyledTableCell>
//             )}

//             {customModel.includePush ? (
//               <StyledTableCell>UNDER / PUSH %</StyledTableCell>
//             ) : (
//               <StyledTableCell>UNDER %</StyledTableCell>
//             )}
//             <StyledTableCell>ADJ.</StyledTableCell>
//             <StyledTableCell>Weight</StyledTableCell>
//           </StyledTableRow>
//         </TableHead>
//         <TableBody>
//           <StyledTableRow>
//             <StyledTableCell>{`${sim.similarGames.length}`}</StyledTableCell>
//             <StyledTableCell>{sim.similarAvg}</StyledTableCell>
//             <StyledTableCell>{sim.similarAvgPerMin}</StyledTableCell>
//             <StyledTableCell>
//               {Minutes.average(sim.similarGames)}
//             </StyledTableCell>
//             <StyledTableCell className={ColorCompare(sim.similarDiff, 0)}>
//               {`${sim.similarDiff > 0 ? "+" : ""}${sim.similarDiff}`}
//             </StyledTableCell>
//             <StyledTableCell
//               className={`${ColorPct(sim.simOverPct)}`}
//             >{`${sim.countSimOver}-${sim.countSimUnder}-${sim.countSimPush}`}</StyledTableCell>
//             {customModel.includePush ? (
//               <StyledTableCell
//                 className={`${ColorPct(sim.simOverPct + sim.simPushPct)}`}
//               >{`${(sim.simOverPct + sim.simPushPct).toFixed(
//                 2
//               )}%`}</StyledTableCell>
//             ) : (
//               <StyledTableCell
//                 className={`${ColorPct(sim.simOverPct)}`}
//               >{`${sim.simOverPct}%`}</StyledTableCell>
//             )}

//             {customModel.includePush ? (
//               <StyledTableCell
//                 className={`${ColorPct(sim.simUnderPct + sim.simPushPct)}`}
//               >{`${(sim.simUnderPct + sim.simPushPct).toFixed(
//                 2
//               )}%`}</StyledTableCell>
//             ) : (
//               <StyledTableCell
//                 className={`${ColorPct(sim.simUnderPct)}`}
//               >{`${sim.simUnderPct}%`}</StyledTableCell>
//             )}

//             <StyledTableCell
//               className={ColorCompare(sim.playerAvgAdj, selectedProp.target)}
//             >
//               {sim.playerAvgAdj}
//             </StyledTableCell>
//             <StyledTableCell>{sim.weight.toFixed(0)}%</StyledTableCell>
//           </StyledTableRow>
//         </TableBody>
//       </Table>
//     </TableContainer>
//   );
// };

// export default SimilarPreview;
