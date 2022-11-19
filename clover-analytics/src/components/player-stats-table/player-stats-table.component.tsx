import {
  Paper,
  Table,
  TableBody,
  TableContainer,
  TableHead,
} from "@mui/material";
// import React, { useState } from "react";
// import { Projection } from "../../shared/interfaces/graphql/proposition.interface";
// import { StyledTableCell } from "../styled-table/styled-table-cell.component";
// import { StyledTableRow } from "../styled-table/styled-table-row.component";
// import "./player-stats-table.component.css";
// import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
// import { Stat } from "../../shared/interfaces/stat.interface";
// import SimilarTable from "./similar-table.component";
// import ExpandSimilarPlayers from "./expand-similar-player-table";
// import ExpandTeamTable from "./expand-team-table";
// import ExpandFragmentTable from "./expand-fragment-table";
// import { Proposition } from "../../shared/interfaces/graphql/game.interface";

// interface PlayerStatsTableProps {
//   selectedProp: Proposition;
//   projection: Projection;
//   customModel: CustomCalculation;
// }

// const PlayerStatsTable: React.FC<PlayerStatsTableProps> = ({
//   selectedProp,
//   projection,
//   customModel,
// }: PlayerStatsTableProps) => {
//   const [openRecency, setOpenRecency] = useState(
//     customModel.recency?.map(() => false) ??
//       customModel.recencyPct?.map(() => false) ?? [false]
//   );
//   function setOpenRecencyFunc(index: number): () => void {
//     const newOpenRecency = [...openRecency];
//     newOpenRecency[index] = !newOpenRecency[index];
//     return () => setOpenRecency(newOpenRecency);
//   }
//   const [openOpponent, setOpenOpponent] = useState(false);
//   const [openSimilarTeams, setOpenSimilarTeams] = useState(false);
//   const [openSimilarPlayers, setOpenSimilarPlayers] = useState(false);

//   return (
//     <div className={"player-stats"}>
//       <TableContainer component={Paper}>
//         <Table aria-label="simple table">
//           <TableHead>
//             <StyledTableRow>
//               <StyledTableCell></StyledTableCell>
//               <StyledTableCell></StyledTableCell>
//               <StyledTableCell>
//                 {selectedProp.statType.abbreviation}
//               </StyledTableCell>
//               <StyledTableCell>
//                 {selectedProp.statType.abbreviation}/MIN
//               </StyledTableCell>
//               {selectedProp.statType.relatedStats?.map((related: Stat) => {
//                 return (
//                   <StyledTableCell key={related.label}>
//                     {related.abbreviation}
//                   </StyledTableCell>
//                 );
//               })}
//               <StyledTableCell>MINS</StyledTableCell>
//               <StyledTableCell>DIFF</StyledTableCell>
//               <StyledTableCell>O-U-P</StyledTableCell>
//               {customModel.includePush ? (
//                 <StyledTableCell>OVER / PUSH %</StyledTableCell>
//               ) : (
//                 <StyledTableCell>OVER %</StyledTableCell>
//               )}

//               {customModel.includePush ? (
//                 <StyledTableCell>UNDER / PUSH %</StyledTableCell>
//               ) : (
//                 <StyledTableCell>UNDER %</StyledTableCell>
//               )}
//               <StyledTableCell>Weight</StyledTableCell>
//             </StyledTableRow>
//           </TableHead>
//           <TableBody>
//             {selectedProp.customPrediction.recencyFragments.map(
//               (fragment, i) => (
//                 <ExpandFragmentTable
//                   open={openRecency[i]}
//                   selectedProp={selectedProp}
//                   fragment={fragment}
//                   projection={projection}
//                   customModel={customModel}
//                   setOpen={setOpenRecencyFunc(i)}
//                 />
//               )
//             )}
//           </TableBody>
//         </Table>
//       </TableContainer>
//       {/* VS Opponent */}
//       {customModel.opponentWeight &&
//       selectedProp.customPrediction.vsOpponent &&
//       projection.player.games.some(
//         (game) =>
//           game.opponent.abbreviation === projection.opponent.abbreviation
//       ) ? (
//         <SimilarTable
//           selectedProp={selectedProp}
//           header={`${projection.player.name} vs ${projection.opponent.abbreviation}`}
//           sim={selectedProp.customPrediction.vsOpponent}
//           customModel={customModel}
//           setOpen={setOpenOpponent}
//           open={openOpponent}
//         >
//           <ExpandTeamTable
//             open={openOpponent}
//             selectedProp={selectedProp}
//             vsTeams={selectedProp.customPrediction.vsOpponent.similarGames}
//           />
//         </SimilarTable>
//       ) : (
//         <></>
//       )}

//       {/* Similar Teams */}
//       {projection.opponent.similarTeams?.length > 0 &&
//       selectedProp.customPrediction.vsSimilarTeams ? (
//         <SimilarTable
//           selectedProp={selectedProp}
//           header={`${projection.player.name} vs ${projection.opponent.similarTeams?.length} Similar Teams`}
//           sim={selectedProp.customPrediction.vsSimilarTeams}
//           customModel={customModel}
//           setOpen={setOpenSimilarTeams}
//           open={openSimilarTeams}
//         >
//           <ExpandTeamTable
//             open={openSimilarTeams}
//             selectedProp={selectedProp}
//             vsTeams={selectedProp.customPrediction.vsSimilarTeams.similarGames}
//           />
//         </SimilarTable>
//       ) : (
//         <></>
//       )}
//       {projection.player.similarPlayers?.length > 0 &&
//       selectedProp.customPrediction.similarPlayersVsOpponent ? (
//         <SimilarTable
//           selectedProp={selectedProp}
//           header={`${projection.player.similarPlayers?.length} Similar Players vs ${projection.opponent.abbreviation}`}
//           sim={selectedProp.customPrediction.similarPlayersVsOpponent}
//           customModel={customModel}
//           setOpen={setOpenSimilarPlayers}
//           open={openSimilarPlayers}
//         >
//           <ExpandSimilarPlayers
//             selectedProp={selectedProp}
//             projection={projection}
//             open={openSimilarPlayers}
//           />
//         </SimilarTable>
//       ) : (
//         <></>
//       )}
//     </div>
//   );
// };

// export default PlayerStatsTable;
