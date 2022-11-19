// import { Projection } from "../interfaces/graphql/proposition.interface";
import {
  Assists,
  Points,
  PointsAssists,
  PointsRebounds,
  PointsReboundsAssists,
  Rebounds,
  ReboundsAssists,
  Stat,
} from "../interfaces/stat.interface";

// export function GetImpliedTarget(projection: Projection, stat: Stat): number {
//   let target = 0;

//   if (stat === Rebounds) {
//     const pointsReboundsTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === PointsRebounds
//       )?.target ?? 0;
//     const pointsTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === Points
//       )?.target ?? 0;
//     const assistsTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === Assists
//       )?.target ?? 0;
//     const astRebTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === ReboundsAssists
//       )?.target ?? 0;
//     const praTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === PointsReboundsAssists
//       )?.target ?? 0;
//     if (pointsReboundsTarget && pointsTarget) {
//       target = pointsReboundsTarget - pointsTarget;
//     } else if (assistsTarget && astRebTarget) {
//       target = astRebTarget - assistsTarget;
//     } else if (praTarget && pointsTarget && assistsTarget) {
//       target = praTarget - (pointsTarget + assistsTarget);
//     }
//   } else if (stat === Assists) {
//     const pointsAssistsTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === PointsAssists
//       )?.target ?? 0;
//     const pointsTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === Points
//       )?.target ?? 0;
//     const reboundsTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === Rebounds
//       )?.target ?? 0;
//     const rebAstTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === ReboundsAssists
//       )?.target ?? 0;
//     const praTarget =
//       projection.propositions.find(
//         (proposition) => proposition.statType === PointsReboundsAssists
//       )?.target ?? 0;
//     if (pointsAssistsTarget && pointsTarget) {
//       target = pointsAssistsTarget - pointsTarget;
//     } else if (reboundsTarget && rebAstTarget) {
//       target = rebAstTarget - reboundsTarget;
//     } else if (praTarget && pointsTarget && reboundsTarget) {
//       target = praTarget - (pointsTarget + reboundsTarget);
//     }
//   }
//   return target;
// }
