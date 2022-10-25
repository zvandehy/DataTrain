import moment from "moment";
// import {
//   Accuracy,
//   Breakdown,
//   DistributionItem,
// } from "../../shared/interfaces/accuracy.interface";
// import { Stat } from "../../shared/interfaces/stat.interface";
// export const INITIALIZE_DISTRIBUTION = (): Accuracy => {
//   let distributionItems: DistributionItem[] = [];

//   const stepSize = 10;
//   for (let x = 0; x <= 100 - stepSize; x += stepSize) {
//     distributionItems.push({
//       min: x,
//       max: x + stepSize,
//       correct: 0,
//       incorrect: 0,
//       push: 0,
//     });
//   }
//   return {
//     total: 0,
//     totalIncorrect: 0,
//     totalCorrect: 0,
//     distributions: distributionItems,
//     topDistributions: [],
//     statTypeBreakdowns: new Map<Stat, Breakdown>(),
//   };
// };

// export enum AccuracyActionType {
//   SET_CORRECT = "SET_CORRECT",
//   SET_INCORRECT = "SET_INCORRECT",
// }

// export interface AccuracyAction {
//   type: AccuracyActionType;
//   confidence: number;
// }
// export const AccuracyReducer = (
//   state: Accuracy,
//   action: AccuracyAction
// ): Accuracy => {
//   switch (action.type) {
//     case AccuracyActionType.SET_CORRECT:
//       if (action.confidence > 70) {
//       }
//       state.distributions.forEach((dist) => {
//         if (action.confidence >= dist.min && action.confidence < dist.max) {
//           dist.correct++;
//         }
//       });
//       return {
//         ...state,
//         totalCorrect: state.totalCorrect + 1,
//         total: state.total + 1,
//       };
//     case AccuracyActionType.SET_INCORRECT:
//       if (action.confidence > 70) {
//       }
//       state.distributions.forEach((dist) => {
//         if (action.confidence >= dist.min && action.confidence < dist.max) {
//           dist.incorrect++;
//         }
//       });
//       return {
//         ...state,
//         totalIncorrect: state.totalIncorrect + 1,
//         total: state.total + 1,
//       };
//   }
//   return state;
// };

// export interface CalendarAccuracies {
//   startDate: string;
//   endDate: string;
//   loadedAccuracies: Accuracy[];
// }

// export interface CalendarAccuracyAction {
//   type: CalendarAccuracyActionType;
//   date: string;
//   loadedAccuracy?: Accuracy;
// }

// export enum CalendarAccuracyActionType {
//   ADD_DATE = "ADD_DATE",
//   ADD_ACCURACY = "ADD_ACCURACY",
// }

// export const CalendarAccuracyReducer = (
//   state: CalendarAccuracies,
//   action: CalendarAccuracyAction
// ): CalendarAccuracies => {
//   switch (action.type) {
//     case CalendarAccuracyActionType.ADD_DATE:
//       if (moment(action.date).isBefore(state.startDate)) {
//         return {
//           ...state,
//           startDate: action.date,
//         };
//       }
//       if (moment(action.date).isAfter(state.endDate)) {
//         return {
//           ...state,
//           endDate: action.date,
//         };
//       }
//       return state;
//   }
//   return state;
// };
