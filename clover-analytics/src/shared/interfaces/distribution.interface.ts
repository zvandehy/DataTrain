export interface DistributionItem {
  min: number;
  max: number;
  correct: number;
  incorrect: number;
}

export interface Accuracy {
  distributions: DistributionItem[];
  total: number;
  totalCorrect: number;
  totalIncorrect: number;
}
