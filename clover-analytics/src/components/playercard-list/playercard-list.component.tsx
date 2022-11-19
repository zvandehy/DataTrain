import moment from "moment";
// import { useEffect, useMemo, useState } from "react";
// import { Match, SortProjections } from "../../shared/functions/filters.fn";
// import { CalculatePredictions } from "../../shared/functions/predictions.fn";
// import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
// import { DateRange } from "../../shared/interfaces/dateRange.interface";
// import { Accuracy } from "../../shared/interfaces/accuracy.interface";
// import { GameFilter } from "../../shared/interfaces/graphql/filters.interface";
// import { Projection } from "../../shared/interfaces/graphql/proposition.interface";
// import { Stat } from "../../shared/interfaces/stat.interface";
// import CalendarSummary from "../projections-summary/calendar-summary/calendar-summary.component";
// import PlayerListFilters from "./list-filters/list-filters.component";
// import "./playercard-list.component.css";
// import PlayerCard from "./playercard/playercard.component";
// import OverallAccuracyBreakdownTable from "../projections-summary/breakdown-table/overall-accuracy-breakdown.component";
// import { Player } from "../../shared/interfaces/graphql/player.interface";
// import Table from "@mui/material/Table";
// import TableRow from "@mui/material/TableRow";
// import { TableBody, TableHead } from "@mui/material";
// import TableCell from "@material-ui/core/TableCell";
// import { AverageStats } from "../../shared/interfaces/graphql/game.interface";

// interface PlayerCardListProps {
//   // projections: Projection[];
//   // customModel: CustomCalculation;
//   // gameFilter: GameFilter;
//   players: Player[];
// }

// const PlayerCardList: React.FC<PlayerCardListProps> = ({
//   players,
// }: PlayerCardListProps) => {
//   // const calculatedProjections = useMemo(() => {
//   //   return CalculatePredictions(projections, gameFilter, customModel);
//   //   // eslint-disable-next-line react-hooks/exhaustive-deps
//   // }, [JSON.stringify(customModel), JSON.stringify(gameFilter)]);
//   // const [lookup, setLookup] = useState("");
//   // const [sortType, setSortType] = useState("");
//   // const [statType, setStatType] = useState(undefined as Stat | undefined);

//   // let filteredProjections = calculatedProjections.filter((projection) => {
//   //   if (statType !== undefined) {
//   //     return Match(projection, { lookup: lookup, statType: statType as Stat });
//   //   }
//   //   return Match(projection, { lookup: lookup });
//   // });
//   // filteredProjections = SortProjections(filteredProjections, {
//   //   sortBy: sortType,
//   //   statType: statType,
//   // });

//   // const [dateRange, setDateRange] = useState(
//   //   filteredProjections.length > 0
//   //     ? ({
//   //         start: moment(filteredProjections[0].startTime)
//   //           .subtract(7, "days")
//   //           .format("YYYY-MM-DD"),
//   //         end: moment(filteredProjections[0].startTime).format("YYYY-MM-DD"),
//   //       } as DateRange)
//   //     : {
//   //         start: moment(new Date()).subtract(7, "days").format("YYYY-MM-DD"),
//   //         end: moment(new Date()).format("YYYY-MM-DD"),
//   //       }
//   // );

//   // const onSetDates = (date: string) => {
//   //   if (moment(date).isBefore(dateRange.start)) {
//   //     setDateRange((prev) => ({
//   //       ...prev,
//   //       start: moment(date).format("YYYY-MM-DD"),
//   //     }));
//   //   }
//   //   if (moment(date).isAfter(dateRange.end)) {
//   //     setDateRange((prev) => ({
//   //       ...prev,
//   //       end: moment(date).format("YYYY-MM-DD"),
//   //     }));
//   //   }
//   // };

//   // const [accuracies, setAccuracies] = useState([] as Accuracy[]);
//   // const [totalAccuracy, setTotalAccuracy] = useState(new Accuracy());

//   // useEffect(() => {
//   //   setTotalAccuracy(new Accuracy());
//   //   setDateRange(
//   //     filteredProjections.length > 0
//   //       ? ({
//   //           start: moment(filteredProjections[0].startTime)
//   //             .add(-1, "days")
//   //             .format("YYYY-MM-DD"),
//   //           end: moment(filteredProjections[0].startTime).format("YYYY-MM-DD"),
//   //         } as DateRange)
//   //       : {
//   //           start: moment(new Date()).add(-1, "days").format("YYYY-MM-DD"),
//   //           end: moment(new Date()).format("YYYY-MM-DD"),
//   //         }
//   //   );
//   // }, [customModel, statType]);

//   // const calendarSummary = useMemo(() => {
//   //   console.log("reload calendar summary");
//   //   return (
//   //     // <CalendarSummary
//   //     //   accuracy={totalAccuracy}
//   //     //   dateRange={dateRange}
//   //     //   setDates={onSetDates}
//   //     //   setAccuracy={setTotalAccuracy}
//   //     //   customModel={customModel}
//   //     //   statType={statType}
//   //     //   lookup={lookup}
//   //     // ></CalendarSummary>
//   //   );
//   // }, [totalAccuracy.allProps, dateRange.start]);

//   // const totalTable = useMemo(() => {
//   //   console.log("Reload table");
//   //   return (
//   //     <OverallAccuracyBreakdownTable
//   //       customModel={customModel}
//   //       lookup={lookup}
//   //       dateRange={dateRange}
//   //       setDates={setDateRange}
//   //       hitCriteria={customModel.hitCriteria}
//   //     />
//   //   );
//   // }, [customModel, dateRange]);

//   return (
//     <>
//       {/* <PlayerListFilters
//         onSearchChange={setLookup}
//         onSortSelect={setSortType}
//         onStatSelect={setStatType}
//       /> */}
//       {/* {calendarSummary} */}
//       {/* {totalTable} */}
//       <div id="player-list">
//         {players.length > 0 ? (
//           // <PlayerCard
//           //   key={`${projection.player.playerID} ${projection.startTime}`}
//           //   projection={projection}
//           //   filteredStatType={statType}
//           //   gameFilter={gameFilter}
//           //   customModel={customModel}
//           // />
//           <Table>
//             <TableHead>
//               <TableRow>
//                 <TableCell>Player</TableCell>
//                 <TableCell>Position</TableCell>
//                 <TableCell>Team</TableCell>
//                 <TableCell>Opponent</TableCell>
//                 <TableCell>Date</TableCell>
//                 <TableCell>Sportsbook</TableCell>
//                 <TableCell>Time</TableCell>
//                 <TableCell>Stat</TableCell>
//                 <TableCell>Target</TableCell>
//                 <TableCell>Prediction</TableCell>
//                 <TableCell>Difference</TableCell>
//                 <TableCell>Wager</TableCell>
//                 <TableCell>Actual</TableCell>
//                 <TableCell>Prop Difference</TableCell>
//                 <TableCell>Prediction Accuracy</TableCell>
//                 <TableCell>Outcome</TableCell>
//               </TableRow>
//             </TableHead>
//             <TableBody>
//               {players.map((player) => {
//                 return player.games[0].prediction.fragments[0].propositions.map(
//                   (prop) => {
//                     const pending: boolean =
//                       player.games[0].outcome[0].toLowerCase() === "p";
//                     return (
//                       <TableRow
//                         key={
//                           player.playerID +
//                           " " +
//                           prop.type +
//                           " " +
//                           prop.target +
//                           " " +
//                           prop.sportsbook
//                         }
//                       >
//                         <TableCell>{player.name}</TableCell>
//                         <TableCell>{player.position}</TableCell>
//                         <TableCell>{player.team.abbreviation}</TableCell>
//                         <TableCell>
//                           {player.games[0].opponent.abbreviation}
//                         </TableCell>
//                         <TableCell>{player.games[0].date}</TableCell>
//                         <TableCell>{prop.sportsbook}</TableCell>
//                         <TableCell>
//                           {moment(prop.lastModified).format(
//                             "MM-DD [at] hh:mm a"
//                           )}
//                         </TableCell>
//                         <TableCell>{prop.type}</TableCell>
//                         <TableCell>{prop.target}</TableCell>
//                         <TableCell>{prop.estimation}</TableCell>
//                         <TableCell>
//                           {(prop.predictionTargetDiff > 0 ? "+" : "") +
//                             prop.predictionTargetDiff}
//                         </TableCell>
//                         <TableCell>{prop.prediction}</TableCell>
//                         <TableCell>
//                           {pending ? "PENDING" : prop.actual}
//                         </TableCell>
//                         <TableCell>
//                           {pending
//                             ? "PENDING"
//                             : prop.actual - prop.target > 0
//                             ? "+" + (prop.actual - prop.target).toFixed(2)
//                             : (prop.actual - prop.target).toFixed(2)}
//                         </TableCell>
//                         <TableCell>
//                           {pending ? "PENDING" : prop.predictionTargetDiff}
//                         </TableCell>
//                         <TableCell>
//                           {pending ? "PENDING" : prop.predictionHit}
//                         </TableCell>
//                       </TableRow>
//                     );
//                   }
//                 );
//               })}
//             </TableBody>
//           </Table>
//         ) : (
//           <div className={"no-results"}>
//             <h1>No Projections Found</h1>
//             <p>
//               Try another date! If you think this is an error then please
//               contact support.
//             </p>
//           </div>
//         )}
//       </div>
//     </>
//   );
// };

// export default PlayerCardList;
