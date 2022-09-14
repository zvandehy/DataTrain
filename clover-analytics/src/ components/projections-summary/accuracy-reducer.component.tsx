import {
  Accuracy,
  DistributionItem,
} from "../../shared/interfaces/distribution.interface";
export const INITIALIZE_DISTRIBUTION = (): Accuracy => {
  let distributionItems: DistributionItem[] = [];

  const stepSize = 10;
  for (let x = 0; x <= 100 - stepSize; x += stepSize) {
    distributionItems.push({
      min: x,
      max: x + stepSize,
      correct: 0,
      incorrect: 0,
    });
  }
  return {
    total: 0,
    totalIncorrect: 0,
    totalCorrect: 0,
    distributions: distributionItems,
  };
};

export enum AccuracyActionType {
  SET_CORRECT = "SET_CORRECT",
  SET_INCORRECT = "SET_INCORRECT",
}

export interface AccuracyAction {
  type: AccuracyActionType;
  confidence: number;
}
export const AccuracyReducer = (
  state: Accuracy,
  action: AccuracyAction
): Accuracy => {
  switch (action.type) {
    case AccuracyActionType.SET_CORRECT:
      if (action.confidence > 70) {
        console.log("SET CORRECT");
      }
      state.distributions.forEach((dist) => {
        if (action.confidence >= dist.min && action.confidence < dist.max) {
          dist.correct++;
        }
      });
      return {
        ...state,
        totalCorrect: state.totalCorrect + 1,
        total: state.total + 1,
      };
    case AccuracyActionType.SET_INCORRECT:
      if (action.confidence > 70) {
        console.log("SET INCORRECT");
      }
      state.distributions.forEach((dist) => {
        if (action.confidence >= dist.min && action.confidence < dist.max) {
          dist.incorrect++;
        }
      });
      return {
        ...state,
        totalIncorrect: state.totalIncorrect + 1,
        total: state.total + 1,
      };
  }
  return state;
};
