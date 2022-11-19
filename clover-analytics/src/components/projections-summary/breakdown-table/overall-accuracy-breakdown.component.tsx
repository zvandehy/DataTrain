import { Button, Table, TableBody, TableHead } from "@mui/material";
// import moment from "moment";
// import { useGetProjections } from "../../../hooks/useGetProjections";
// import apolloClient from "../../../shared/apollo-client";
// import { BETTING_CATEGORIES } from "../../../shared/constants";
// import { ColorCompare } from "../../../shared/functions/color.fn";
// import { Match } from "../../../shared/functions/filters.fn";
// import { CalculatePredictions } from "../../../shared/functions/predictions.fn";
// import {
//   Accuracy,
//   BreakdownFilter,
//   GetHitRate,
//   HitCriteria,
//   PredictionResultEnum,
// } from "../../../shared/interfaces/accuracy.interface";
// import { CustomCalculation } from "../../../shared/interfaces/custom-prediction.interface";
// import { DateRange } from "../../../shared/interfaces/dateRange.interface";
// import {
//   GameFilter,
//   ProjectionFilter,
// } from "../../../shared/interfaces/graphql/filters.interface";
// import { Projection } from "../../../shared/interfaces/graphql/proposition.interface";
// import { Stat } from "../../../shared/interfaces/stat.interface";
// import { StyledTableCell } from "../../styled-table/styled-table-cell.component";
// import { StyledTableRow } from "../../styled-table/styled-table-row.component";
// import "./overall-accuracy-breakdown.component.css";

// interface OverallAccuracyBreakdownTableProps {
//   // totalAccuracy: Accuracy;
//   dateRange: DateRange;
//   customModel: CustomCalculation;
//   lookup?: string;
//   hitCriteria: HitCriteria;
//   setDates: (dateRange: DateRange) => void;
// }

// const OverallAccuracyBreakdownTable: React.FC<
//   OverallAccuracyBreakdownTableProps
// > = ({
//   dateRange,
//   customModel,
//   lookup,
//   hitCriteria,
//   setDates,
// }: OverallAccuracyBreakdownTableProps) => {
//   // const [allProjections, setAllProjections] = useState([] as Projection[]);
//   // const [totalAccuracy, setTotalAccuracy] = useState(new Accuracy());
//   // const [allProjections, setAllProjections] = useState(
//   //   new Map<string, Projection[]>()
//   // );
//   // console.group("RENDER TABLE COMPONENT");
//   // console.log(dateRange.start, dateRange.end);
//   // console.log(customModel);
//   // console.log(lookup);
//   // console.log(hitCriteria);
//   // console.groupEnd();

//   // const { loading, error, data } = useGetProjections({
//   //   customModel: customModel,
//   //   gameFilter: { endDate: dateRange.end },
//   //   predictionFilter: { endDate: dateRange.end },
//   //   projectionFilter: { startDate: dateRange.start, endDate: dateRange.end },
//   // });

//   // if (loading) return loading;
//   // if (error || !data) return error;
//   // let allProjections: Map<string, Projection[]> = new Map();
//   // data.forEach((projection) => {
//   //   const date = moment(projection.startTime).format("YYYY-MM-DD");
//   //   const foundProjections = allProjections.get(date);
//   //   const calculatedProjection = CalculatePredictions(
//   //     [projection],
//   //     { endDate: date },
//   //     customModel
//   //   )[0];
//   //   if (foundProjections) {
//   //     allProjections.set(date, [...foundProjections, calculatedProjection]);
//   //   } else {
//   //     allProjections.set(date, [calculatedProjection]);
//   //   }
//   // });

//   // useEffect(() => {
//   //   let ignore = false;

//   //   async function startFetching() {
//   //     const date = dateRange.start;
//   //     if (!ignore) {
//   //       GetDay(date, customModel)
//   //         .then((projections: Projection[]) => {
//   //           console.group(`Got ${projections.length} from ${date}`);
//   //           setAllProjections((prev) => {
//   //             //   console.log("PREV: ", prev);
//   //             //   console.log("NEW: ", projections);
//   //             //   console.log("COMBINED: ", [...prev, ...projections]);
//   //             prev.set(date, projections);
//   //             return prev;
//   //           });
//   //           console.log(allProjections);
//   //           console.groupEnd();
//   //           // setTotalAccuracy(new Accuracy(allProjections));
//   //           // console.log(dayAccuracies, allProjections, totalAccuracy);
//   //         })
//   //         .catch(() => {});
//   //     }
//   //   }

//   //   startFetching();

//   //   return () => {
//   //     ignore = true;
//   //   };
//   // }, [dateRange.start]);

//   // // REFRESH ALL DATES
//   // useMemo(() => {
//   //   for (
//   //     let cur = moment(dateRange.start);
//   //     cur.isSameOrBefore(dateRange.end, "day");
//   //     cur.add(1, "day")
//   //   ) {
//   //     const date = cur.format("YYYY-MM-DD");
//   //     GetDay(date, customModel)
//   //       .then((projections: Projection[]) => {
//   //         // console.log("Got: ", date, projections);
//   //         setDayAccuracies((prev) => {
//   //           return prev.set(date, projections);
//   //         });
//   //         // setAllProjections((prev) => {
//   //         //   console.log("PREV: ", prev);
//   //         //   console.log("NEW: ", projections);
//   //         //   console.log("COMBINED: ", [...prev, ...projections]);
//   //         //   return [...prev, ...projections];
//   //         // });
//   //         // setTotalAccuracy(new Accuracy(allProjections));
//   //       })
//   //       .catch(() => {});
//   //   }
//   //   // console.log(dayAccuracies);
//   // }, [customModel]);
//   // // ADD SINGLE DAY
//   // useMemo(() => {
//   //   const date = dateRange.start;
//   //   console.log("Get day and update total: ", date);
//   //   GetDay(date, customModel)
//   //     .then((projections: Projection[]) => {
//   //       // console.log("Got: ", date, projections);
//   //       setDayAccuracies((prev) => {
//   //         return prev.set(date, projections);
//   //       });
//   //       // setAllProjections((prev) => {
//   //       //   console.log("PREV: ", prev);
//   //       //   console.log("NEW: ", projections);
//   //       //   console.log("COMBINED: ", [...prev, ...projections]);
//   //       //   return [...prev, ...projections];
//   //       // });
//   //       // setTotalAccuracy(new Accuracy(allProjections));
//   //       // console.log(dayAccuracies, allProjections, totalAccuracy);
//   //     })
//   //     .catch(() => {});
//   // }, [dateRange.start]);
//   // let p: Projection[] = Array.from(dayAccuracies.values()).flat();
//   // console.log(dayAccuracies, p);
//   // const totalAccuracy: Accuracy = new Accuracy();
//   return (
//     <></>
//     // <Table>
//     //   <TableHead>
//     //     <StyledTableRow>
//     //       <StyledTableCell>
//     //         <Button
//     //           onClick={() =>
//     //             setDates({
//     //               end: dateRange.end,
//     //               start: moment(dateRange.start)
//     //                 .subtract(7, "day")
//     //                 .format("YYYY-MM-DD"),
//     //             })
//     //           }
//     //         >
//     //           Load{" "}
//     //           {moment(dateRange.start).subtract(7, "day").format("YYYY-MM-DD")}
//     //         </Button>
//     //       </StyledTableCell>
//     //       <StyledTableCell>ALL</StyledTableCell>
//     //       {BETTING_CATEGORIES.map((stat: Stat) => (
//     //         <StyledTableCell key={stat.label}>
//     //           {stat.abbreviation}
//     //         </StyledTableCell>
//     //       ))}
//     //     </StyledTableRow>
//     //   </TableHead>
//     //   <TableBody>
//     //     <AllSummary
//     //       allProjections={allProjections}
//     //       hitCriteria={hitCriteria}
//     //       key={allProjections.size}
//     //     />
//     //     {Array.from(allProjections, (mapEntry, k) => {
//     //       const date = mapEntry[0];
//     //       const accuracy = new Accuracy(mapEntry[1]).filter(
//     //         new BreakdownFilter({ min: 60 })
//     //       );
//     //       return { accuracy, date };
//     //     })
//     //       .sort((itemA, itemB) => moment(itemB.date).diff(itemA.date))
//     //       .map(({ accuracy, date }) => (
//     //         // console.log(mapEntry[1], accuracy);
//     //         <StyledTableRow key={date}>
//     //           <StyledTableCell>{date}</StyledTableCell>
//     //           <StyledTableCell
//     //             className={ColorCompare(
//     //               (accuracy.totalCorrect + accuracy.totalPush) /
//     //                 (accuracy.totalCorrect +
//     //                   accuracy.totalIncorrect +
//     //                   accuracy.totalPush),
//     //               0.6
//     //             )}
//     //           >
//     //             {accuracy.totalCorrect}-{accuracy.totalIncorrect}-
//     //             {accuracy.totalPush}
//     //           </StyledTableCell>
//     //           {BETTING_CATEGORIES.map((stat) => {
//     //             const breakdown = accuracy.statBreakdowns.get(stat);
//     //             if (!breakdown)
//     //               return (
//     //                 <StyledTableCell
//     //                   key={stat.label + " " + date}
//     //                 ></StyledTableCell>
//     //               );
//     //             const result = breakdown.filter((r) => r.confidence > 60);
//     //             let correct = 0;
//     //             let incorrect = 0;
//     //             let push = 0;
//     //             result.forEach((r) => {
//     //               if (r.result === PredictionResultEnum.CORRECT) correct++;
//     //               if (r.result === PredictionResultEnum.INCORRECT) incorrect++;
//     //               if (r.result === PredictionResultEnum.PUSH) push++;
//     //             });
//     //             return (
//     //               <StyledTableCell
//     //                 key={stat.label + " " + date}
//     //                 className={ColorCompare(
//     //                   (correct + push) / (correct + incorrect + push),
//     //                   0.6
//     //                 )}
//     //               >
//     //                 {correct}-{incorrect}-{push}
//     //               </StyledTableCell>
//     //             );
//     //           })}
//     //         </StyledTableRow>
//     //       ))}
//     //     {/* {dayAccuracies.forEach((accuracy, date) => {
//     //       return (
//     //         <StyledTableRow>
//     //           <StyledTableCell>{date}</StyledTableCell>
//     //           <StyledTableCell>{accuracy.totalCorrect}</StyledTableCell>
//     //         </StyledTableRow>
//     //       );
//     //     })} */}
//     //   </TableBody>
//     // </Table>
//   );
// };

// export default OverallAccuracyBreakdownTable;

// const AllSummary: React.FC<{
//   allProjections: Map<string, Projection[]>;
//   hitCriteria: HitCriteria;
// }> = ({ allProjections, hitCriteria }) => {
//   console.log("Rendered All Summary");
//   let combinedProjections: Projection[] = [];
//   allProjections.forEach((dayProjections, date) => {
//     combinedProjections.push(...dayProjections);
//   });
//   const totalAccuracy = new Accuracy(combinedProjections);
//   return (
//     <StyledTableRow key={"all"}>
//       <StyledTableCell>ALL</StyledTableCell>
//       <StyledTableCell
//         className={ColorCompare(
//           (totalAccuracy.totalCorrect + totalAccuracy.totalPush) /
//             (totalAccuracy.totalCorrect +
//               totalAccuracy.totalIncorrect +
//               totalAccuracy.totalPush),
//           0.6
//         )}
//       >
//         {totalAccuracy.totalCorrect}-{totalAccuracy.totalIncorrect}-
//         {totalAccuracy.total -
//           totalAccuracy.totalCorrect -
//           totalAccuracy.totalIncorrect}
//       </StyledTableCell>
//       {BETTING_CATEGORIES.map((stat) => {
//         const breakdown = GetHitRate(
//           hitCriteria,
//           totalAccuracy.statBreakdowns
//             .get(stat)
//             ?.filter(
//               (item) => item.confidence > hitCriteria.confidenceThreshold
//             ) ?? []
//         );
//         return (
//           <StyledTableCell
//             key={stat.label}
//             className={ColorCompare(breakdown?.pct ?? 0, 0.6)}
//           >
//             {breakdown?.correct}-{breakdown?.incorrect}-{breakdown?.push}
//           </StyledTableCell>
//         );
//       })}
//     </StyledTableRow>
//   );
// };

// async function GetDay(
//   date: string,
//   customModel: CustomCalculation,
//   lookup?: string
// ): Promise<Projection[]> {
//   // console.log("Get summary for: ", date);
//   // await new Promise((r) => setTimeout(r, 3000));
//   // let projectionFilter: ProjectionFilter = {
//   //   startDate: moment(date).format("YYYY-MM-DD"),
//   //   endDate: moment(date).format("YYYY-MM-DD"),
//   // };
//   // let gameFilter: GameFilter = {
//   //   endDate: moment(date).format("YYYY-MM-DD"),
//   // };
//   // const predictionFilter: GameFilter = {
//   //   season: "2022-23",
//   //   endDate: moment(date).format("YYYY-MM-DD"),
//   // };
//   // console.log("wait for response: ", date);
//   // const response = await apolloClient.query({
//   //   query: getQuery({ customModel }),
//   //   variables: {
//   //     playerFilter: projectionFilter,
//   //     gameFilter: predictionFilter,
//   //   },
//   // });
//   // console.log("got response for ", date, response);
//   // if (response.data && !response.error) {
//   //   // console.log("Calculate predictions for: ", response.data);
//   //   const projections: Projection[] = response.data?.projections.map(
//   //     (projection: Projection) => {
//   //       let player = { ...projection.player };
//   //       let games = player.games.map((game) => {
//   //         return { ...game };
//   //       });
//   //       games.sort((a, b) => {
//   //         return moment(a.date).diff(b.date);
//   //       });
//   //       player.games = games;
//   //       return { ...projection, player: player };
//   //     }
//   //   );
//   // const filteredProjections = CalculatePredictions(
//   //   projections,
//   //   gameFilter,
//   //   customModel
//   // );
//   // if (lookup) {
//   //   filteredProjections.filter((projection) => {
//   //     // if (statType !== undefined) {
//   //     //   return Match(projection, {
//   //     //     lookup: lookup,
//   //     //     statType: statType as Stat,
//   //     //   });
//   //     // }
//   //     return Match(projection, { lookup: lookup });
//   //   });
//   // }
//   // // console.log("Got projections: ", filteredProjections);
//   // return filteredProjections;
//   // }

//   // if (statType) {
//   //   filteredProjections.forEach((projection, i) => {
//   //     filteredProjections[i].propositions = projection.propositions.filter(
//   //       (prop) => prop.statType.abbreviation === statType.abbreviation
//   //     );
//   //   });
//   // }
//   return [];
// }
