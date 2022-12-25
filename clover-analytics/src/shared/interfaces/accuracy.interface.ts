import moment from "moment";
// import { Game } from "./graphql/game.interface";
// import { Stat } from "./stat.interface";

// export enum HitCriteriaType {
//   ALL_PROPS = "ALL_PROPS",
//   PLAYER_TOP = "PLAYER_TOP",
// }

// export interface HitCriteria {
//   confidenceThreshold: number;
//   hitType: HitCriteriaType;
//   stepSize: number;
// }

// export interface HitRate {
//   correct: number;
//   incorrect: number;
//   push: number;
//   pct: number;
// }

// export enum PredictionResultEnum {
//   CORRECT = 1,
//   INCORRECT = -1,
//   PUSH = 0,
// }

// export interface PredictionResult {
//   result: PredictionResultEnum;
//   confidence: number;
//   date: string;
// }

// export class BreakdownFilter {
//   min: number = 0;
//   max: number = 100;
//   date?: string;

//   constructor(options?: { min?: number; max?: number; date?: string }) {
//     this.min = options?.min ?? 0;
//     this.max = options?.max ?? 100;
//     this.date = options?.date;
//   }
//   filter(prop: PredictionResult): boolean {
//     if (this.date) {
//       return (
//         prop.confidence >= this.min &&
//         prop.confidence <= this.max &&
//         moment(prop.date).isSame(this.date)
//       );
//     } else {
//       return prop.confidence >= this.min! && prop.confidence <= this.max!;
//     }
//   }
// }

// export class Accuracy {
//   allProps: PredictionResult[] = [];
//   topPlayerProps: PredictionResult[] = [];
//   statBreakdowns: Map<Stat, PredictionResult[]> = new Map<
//     Stat,
//     PredictionResult[]
//   >();
//   total: number = 0;
//   totalPush: number = 0;
//   totalCorrect: number = 0;
//   totalIncorrect: number = 0;

//   constructor(filteredProjections?: Projection[]) {
//     this.add(filteredProjections ?? []);
//   }

//   add(projections: Projection[]) {
//     projections.forEach((projection) => {
//       if (!projection.result) return;
//       // TODO: error handling / cleaning
//       projection.propositions
//         .sort((a, b) =>
//           a.customPrediction.confidence > b.customPrediction.confidence ? -1 : 1
//         )
//         .forEach((prop, index) => {
//           const result: PredictionResult = GetPredictionResult(
//             prop,
//             projection.result
//           );
//           this.total++;
//           if (result.result === PredictionResultEnum.CORRECT)
//             this.totalCorrect++;
//           if (result.result === PredictionResultEnum.INCORRECT)
//             this.totalIncorrect++;
//           if (result.result === PredictionResultEnum.PUSH) this.totalPush++;
//           // all Props
//           this.allProps.push(result);
//           // top Player prop
//           if (index === 0) {
//             this.topPlayerProps.push(result);
//           }
//           // stat props
//           const prev = this.statBreakdowns.get(prop.statType);
//           if (prev === undefined) {
//             this.statBreakdowns.set(prop.statType, [result]);
//           } else {
//             this.statBreakdowns.set(prop.statType, [...prev, result]);
//           }
//         });
//     });
//   }
//   filter(filter: BreakdownFilter): Accuracy {
//     let accuracy = new Accuracy();
//     if (filter.min === undefined) filter.min = 0;
//     if (filter.max === undefined) filter.max = 100;
//     accuracy.allProps = this.allProps.filter((prop) => filter.filter(prop));
//     accuracy.allProps.forEach((r) => {
//       if (r.result === PredictionResultEnum.CORRECT) accuracy.totalCorrect++;
//       if (r.result === PredictionResultEnum.INCORRECT)
//         accuracy.totalIncorrect++;
//       if (r.result === PredictionResultEnum.PUSH) accuracy.totalPush++;
//     });
//     accuracy.topPlayerProps = this.topPlayerProps.filter((prop) =>
//       filter.filter(prop)
//     );
//     this.statBreakdowns.forEach((breakdown, stat) => {
//       accuracy.statBreakdowns.set(
//         stat,
//         breakdown.filter((prop) => filter.filter(prop))
//       );
//     });
//     return accuracy;
//   }
// }

// function GetPredictionResult(
//   prop: Proposition,
//   result: Game
// ): PredictionResult {
//   const score = prop.statType.score(result);
//   const date = moment(result.date).format("YYYY-MM-DD");
//   if (score === prop.target) {
//     return {
//       result: PredictionResultEnum.PUSH,
//       confidence: prop.customPrediction.confidence,
//       date: date,
//     };
//   }
//   if (prop.customPrediction.overUnderPrediction === "Over") {
//     if (score > prop.target)
//       return {
//         result: PredictionResultEnum.CORRECT,
//         confidence: prop.customPrediction.confidence,
//         date: date,
//       };
//     return {
//       result: PredictionResultEnum.INCORRECT,
//       confidence: prop.customPrediction.confidence,
//       date: date,
//     };
//   } else {
//     if (score < prop.target)
//       return {
//         result: PredictionResultEnum.CORRECT,
//         confidence: prop.customPrediction.confidence,
//         date: date,
//       };
//     return {
//       result: PredictionResultEnum.INCORRECT,
//       confidence: prop.customPrediction.confidence,
//       date: date,
//     };
//   }
// }

// export function GetHitRate(
//   hitCriteria: HitCriteria,
//   results: PredictionResult[]
// ): HitRate {
//   let ret: HitRate = {
//     correct: 0,
//     incorrect: 0,
//     push: 0,
//     pct: 0,
//   };
//   ret.correct = results.filter(
//     (item) => item.result === PredictionResultEnum.CORRECT
//   ).length;
//   ret.incorrect = results.filter(
//     (item) => item.result === PredictionResultEnum.INCORRECT
//   ).length;
//   ret.push = results.filter(
//     (item) => item.result === PredictionResultEnum.PUSH
//   ).length;
//   ret.pct = +(
//     (ret.correct + ret.push) /
//     (ret.correct + ret.incorrect + ret.push)
//   ).toFixed(2);
//   return ret;
// }
