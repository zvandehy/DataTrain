import { Box, CircularProgress, Typography } from "@mui/material";
// import moment from "moment";
// import Calendar from "react-calendar";
// import "react-calendar/dist/Calendar.css";
// import "./calendar-summary.component.css";
// import { ChangeEvent } from "react";
// import {
//   Accuracy,
//   BreakdownFilter,
//   GetHitRate,
//   HitCriteria,
//   HitCriteriaType,
//   HitRate,
// } from "../../../shared/interfaces/accuracy.interface";
// import { ColorCompare } from "../../../shared/functions/color.fn";
// import { DateRange } from "../../../shared/interfaces/dateRange.interface";
// import DayAccuracy from "./day-accuracy.component";
// import { CustomCalculation } from "../../../shared/interfaces/custom-prediction.interface";
// import { Stat } from "../../../shared/interfaces/stat.interface";
// // import DayAccuracyWrapper from "./day-accuracy-wrapper.component";

// interface CalendarSummaryProps {
//   accuracy: Accuracy;
//   setAccuracy: React.Dispatch<React.SetStateAction<Accuracy>>;
//   dateRange: DateRange;
//   setDates: (date: string) => void;
//   customModel: CustomCalculation;
//   statType: Stat | undefined;
//   lookup: string;
// }

// function GetAccuracyContentForDate(
//   date: string,
//   accuracy: Accuracy,
//   setAccuracy: React.Dispatch<React.SetStateAction<Accuracy>>,
//   dateRange: DateRange,
//   customModel: CustomCalculation,
//   statType: Stat | undefined,
//   lookup: string
// ) {
//   //   let found = accuracies.find((accuracy) => moment(accuracy.date).isSame(date));
//   //   if (found) return <Box sx={{ textAlign: "end" }}>{HitRate(found, 0.6)}%</Box>;
//   if (
//     moment(date).isBetween(
//       moment(dateRange.start),
//       moment(dateRange.end),
//       undefined,
//       "[]"
//     )
//   ) {
//     const found = accuracy.filter(new BreakdownFilter({ date: date }));
//     if (found.total > 0) {
//       return (
//         <DayAccuracy accuracy={found} hitCriteria={customModel.hitCriteria} />
//       );
//     }
//     return (
//       <></>
//       // <DayAccuracyWrapper
//       //   date={date}
//       //   customModel={customModel}
//       //   statType={statType}
//       //   lookup={lookup}
//       //   setAccuracy={setAccuracy}
//       //   hitCriteria={customModel.hitCriteria}
//       // />
//     );
//   }
//   return <></>;
// }

// function GetAccuracyClassForDate(
//   date: string,
//   accuracy: Accuracy,
//   dateRange: DateRange,
//   hitCriteria: HitCriteria
// ): string {
//   return "";
//   // let found = accuracies.find((accuracy) => moment(accuracy.date).isSame(date));
//   // let classes = "";
//   // if (moment(date).isSame(dateRange.end)) {
//   //   classes = "active ";
//   // }
//   // if (!found) {
//   //   if (
//   //     moment(date).isBetween(
//   //       moment(dateRange.start),
//   //       moment(dateRange.end),
//   //       undefined,
//   //       "[]"
//   //     )
//   //   ) {
//   //     classes += "na";
//   //   }
//   //   return classes;
//   // }
//   // let hit = GetHitRate(found, hitCriteria);
//   // let pctCorrect = +(
//   //   ((hit.correct + hit.push) / (hit.correct + hit.push + hit.incorrect)) *
//   //   100
//   // ).toFixed(2);
//   // classes += ColorCompare(pctCorrect, 60);
//   // return classes;
// }

// function GetProjectionsForDates(
//   setDates: (dateRange: DateRange) => void
// ): (
//   values: [Date] | [Date, Date],
//   event: ChangeEvent<HTMLInputElement>
// ) => void {
//   return (
//     values: [Date] | [Date, Date],
//     event: ChangeEvent<HTMLInputElement>
//   ) => {
//     if (values.length >= 2) {
//       setDates({
//         start: moment(values[0]).format("YYYY-MM-DD"),
//         end: moment(values[1]).format("YYYY-MM-DD"),
//       });
//     }
//   };
// }

// function BeyondRange(
//   date: string,
//   dateRange: DateRange,
//   allow: number
// ): boolean {
//   if (
//     moment(date).isBefore(moment(dateRange.start).subtract(allow, "days")) ||
//     moment(date).isAfter(moment(dateRange.end).add(0, "days"))
//   )
//     return true;
//   return false;
// }

// const CalendarSummary: React.FC<CalendarSummaryProps> = ({
//   accuracy,
//   setAccuracy,
//   dateRange,
//   setDates,
//   customModel,
//   statType,
//   lookup,
// }: CalendarSummaryProps) => {
//   return (
//     <Box>
//       <Typography variant="h5" align="center">
//         Model Performance
//       </Typography>
//       <Typography variant="h6" align="center">
//         Hit Rate on{" "}
//         {customModel.hitCriteria.hitType === HitCriteriaType.ALL_PROPS
//           ? "All Props"
//           : "the Top Play for Each Player"}{" "}
//         Over {customModel.hitCriteria.confidenceThreshold}% Confidence
//       </Typography>
//       <Calendar
//         view="month"
//         allowPartialRange={true}
//         tileDisabled={(props) =>
//           BeyondRange(moment(props.date).format("YYYY-MM-DD"), dateRange, 3)
//         }
//         defaultValue={new Date(dateRange.end)}
//         // tileDisabled={(props) =>
//         //   NoDateExists(moment(props.date).format("YYYY-MM-DD"), accuracies)
//         // }
//         maxDate={new Date()}
//         tileContent={(props) =>
//           GetAccuracyContentForDate(
//             moment(props.date).format("YYYY-MM-DD"),
//             accuracy,
//             setAccuracy,
//             dateRange,
//             customModel,
//             statType,
//             lookup
//           )
//         }
//         tileClassName={(props) =>
//           GetAccuracyClassForDate(
//             moment(props.date).format("YYYY-MM-DD"),
//             accuracy,
//             dateRange,
//             customModel.hitCriteria
//           )
//         }
//         selectRange={false}
//         onClickDay={(value: Date, event: any) => {
//           setDates(moment(value).format("YYYY-MM-DD"));
//         }}
//         // onChange={GetProjectionsForDates(setDates)}
//       />
//     </Box>
//   );
// };

// export default CalendarSummary;
