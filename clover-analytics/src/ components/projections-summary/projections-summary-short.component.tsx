import moment from "moment";
// import React from "react";
// import { BETTING_CATEGORIES } from "../../shared/constants";
// import { ColorPct } from "../../shared/functions/color.fn";
// import { GetMaxConfidence } from "../../shared/functions/predictions.fn";
// import { DistributionItem } from "../../shared/interfaces/accuracy.interface";
// import { Projection } from "../../shared/interfaces/graphql/projection.interface";
// import { Stat } from "../../shared/interfaces/stat.interface";
// import {
//   AccuracyAction,
//   AccuracyActionType,
// } from "./accuracy-reducer.component";
// import "./projections-summary.component.css";

// interface ProjectionsSummaryProps {
//   projections: Projection[];
//   dispatch: (action: AccuracyAction) => void;
//   filteredStat?: Stat | undefined;
// }

// const ProjectionsSummaryShort: React.FC<ProjectionsSummaryProps> = ({
//   projections,
//   dispatch,
//   filteredStat,
// }: ProjectionsSummaryProps) => {
//   const uniqueTeams: string[] = [
//     ...Array.from(
//       new Set(projections.map((projection) => projection.opponent.abbreviation))
//     ),
//   ];
//   let countPush = 0;
//   let countTotal = 0;
//   let countMaxTotal = 0;
//   let countMaxCorrect = 0;
//   let countMaxIncorrect = 0;
//   let countMaxPush = 0;

//   let maxDistribution: DistributionItem[] = [];

//   const stepSize = 10;
//   for (let x = 0; x <= 100 - stepSize; x += stepSize) {
//     maxDistribution.push({
//       min: x,
//       max: x + stepSize,
//       correct: 0,
//       incorrect: 0,
//       push: 0,
//     });
//   }

//   BETTING_CATEGORIES.filter(
//     (stat) => (filteredStat && stat === filteredStat) || !filteredStat
//   ).forEach((stat) => {
//     projections.forEach((projection) => {
//       if (projection.result) {
//         let latestStatProps = projection.propositions
//           .filter((proposition) => proposition.statType === stat)
//           .sort((a, b) =>
//             a.customPrediction.confidence > b.customPrediction.confidence
//               ? -1
//               : 1
//           );
//         // .sort((a, b) => moment(a.lastModified).diff(b.lastModified));
//         if (latestStatProps.length > 0) {
//           let prop = latestStatProps[0];
//           if (prop.statType.score(projection.result) === prop.target) {
//             countPush++;
//             countTotal++;
//             //highest
//             if (
//               projection.propositions.sort((a, b) =>
//                 a.customPrediction.confidence > b.customPrediction.confidence
//                   ? -1
//                   : 1
//               )[0].statType === stat
//             ) {
//               countMaxPush++;
//               countMaxTotal++;
//               // TODO: dispatch push?
//             }
//           } else {
//             if (prop.customPrediction.overUnderPrediction === "Over") {
//               if (prop.statType.score(projection.result) > prop.target) {
//                 countTotal++;
//                 //highest
//                 if (
//                   projection.propositions.sort((a, b) =>
//                     a.customPrediction.confidence >
//                     b.customPrediction.confidence
//                       ? -1
//                       : 1
//                   )[0].statType === stat
//                 ) {
//                   maxDistribution.forEach((dist) => {
//                     if (
//                       prop.customPrediction.confidence >= dist.min &&
//                       prop.customPrediction.confidence < dist.max
//                     ) {
//                       dist.correct++;
//                     }
//                   });
//                   countMaxCorrect++;
//                   countMaxTotal++;
//                   dispatch({
//                     type: AccuracyActionType.SET_CORRECT,
//                     confidence: prop.customPrediction.confidence,
//                   });
//                 }
//               } else {
//                 countTotal++;
//                 //highest
//                 if (
//                   projection.propositions.sort((a, b) =>
//                     a.customPrediction.confidence >
//                     b.customPrediction.confidence
//                       ? -1
//                       : 1
//                   )[0].statType === stat
//                 ) {
//                   maxDistribution.forEach((dist) => {
//                     if (
//                       prop.customPrediction.confidence >= dist.min &&
//                       prop.customPrediction.confidence < dist.max
//                     ) {
//                       dist.incorrect++;
//                     }
//                   });
//                   countMaxIncorrect++;
//                   countMaxTotal++;
//                   dispatch({
//                     type: AccuracyActionType.SET_INCORRECT,
//                     confidence: prop.customPrediction.confidence,
//                   });
//                   if (prop.customPrediction.confidence > 70) {
//                   }
//                 }
//               }
//             }
//             if (prop.customPrediction.overUnderPrediction === "Under") {
//               if (prop.statType.score(projection.result) < prop.target) {
//                 countTotal++;
//                 //highest
//                 if (
//                   projection.propositions.sort((a, b) =>
//                     a.customPrediction.confidence >
//                     b.customPrediction.confidence
//                       ? -1
//                       : 1
//                   )[0].statType === stat
//                 ) {
//                   maxDistribution.forEach((dist) => {
//                     if (
//                       prop.customPrediction.confidence >= dist.min &&
//                       prop.customPrediction.confidence < dist.max
//                     ) {
//                       dist.correct++;
//                     }
//                   });
//                   countMaxCorrect++;
//                   countMaxTotal++;
//                   dispatch({
//                     type: AccuracyActionType.SET_CORRECT,
//                     confidence: prop.customPrediction.confidence,
//                   });
//                 }
//               } else {
//                 countTotal++;
//                 //highest
//                 if (
//                   projection.propositions.sort((a, b) =>
//                     a.customPrediction.confidence >
//                     b.customPrediction.confidence
//                       ? -1
//                       : 1
//                   )[0].statType === stat
//                 ) {
//                   maxDistribution.forEach((dist) => {
//                     if (
//                       prop.customPrediction.confidence >= dist.min &&
//                       prop.customPrediction.confidence < dist.max
//                     ) {
//                       dist.incorrect++;
//                     }
//                   });
//                   countMaxIncorrect++;
//                   countMaxTotal++;
//                   dispatch({
//                     type: AccuracyActionType.SET_INCORRECT,
//                     confidence: prop.customPrediction.confidence,
//                   });
//                   if (prop.customPrediction.confidence > 70) {
//                   }
//                 }
//               }
//             }
//           }
//         }
//       }
//     });
//   });

//   return (
//     <>
//       <div id="max-distribution-summary">
//         <>
//           <span style={{ paddingInline: "5px" }}>
//             {uniqueTeams.length} Teams,
//           </span>
//           <span style={{ paddingInline: "5px" }}>{countMaxTotal} Players,</span>
//           <span style={{ paddingInline: "5px" }}>{countTotal} Props,</span>
//           <span style={{ paddingInline: "5px" }}>
//             Highest Confidence Prop per Player:{" "}
//           </span>
//           {maxDistribution
//             .filter((dist) => dist.correct + dist.incorrect > 0)
//             .map((dist) => {
//               return (
//                 <span
//                   key={`highest-${dist.min}-${dist.max}`}
//                   className={ColorPct(
//                     dist.correct / (dist.correct + dist.incorrect)
//                   )}
//                   style={{ paddingInline: "5px" }}
//                 >
//                   {dist.min}-{dist.max}%: {dist.correct}-{dist.incorrect} (
//                   {(
//                     (dist.correct / (dist.correct + dist.incorrect)) *
//                     100
//                   ).toFixed(2)}
//                   %),{" "}
//                 </span>
//               );
//             })}
//         </>
//       </div>
//     </>
//   );
// };

// export default ProjectionsSummaryShort;
