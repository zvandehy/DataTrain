import { BETTING_CATEGORIES } from "../../shared/constants";
import { GetMaxConfidence } from "../../shared/functions/predictions.fn";
import { Projection } from "../../shared/interfaces/graphql/projection.interface";
import "./projections-summary.component.css";

interface ProjectionsSummaryProps {
  projections: Projection[];
}

const ProjectionsSummary: React.FC<ProjectionsSummaryProps> = ({
  projections,
}: ProjectionsSummaryProps) => {
  const uniqueTeams: string[] = [
    ...Array.from(
      new Set(projections.map((projection) => projection.opponent.abbreviation))
    ),
  ];
  let countCorrect = 0;
  let countIncorrect = 0;
  let countPush = 0;
  let countUnknown = 0;
  let countMaxUnknown = 0;
  let countTotal = 0;
  let countMaxTotal = 0;
  let countMaxCorrect = 0;
  let countMaxIncorrect = 0;
  let countMaxPush = 0;
  //TODO: if a statType has more than one proposition, use the most recent one
  projections.forEach((projection) => {
    if (projection.result) {
      BETTING_CATEGORIES.forEach((category) => {
        //TODO: get the most recent proposition
        const prop = projection.propositions.find(
          (p) => p.statType === category
        );
        if (prop) {
          if (prop.statType.score(projection.result) === prop.target) {
            countPush++;
            countTotal++;
          } else {
            if (prop.customPrediction.overUnderPrediction === "Over") {
              if (prop.statType.score(projection.result) > prop.target) {
                countCorrect++;
                countTotal++;
              } else {
                countIncorrect++;
                countTotal++;
              }
            }
            if (prop.customPrediction.overUnderPrediction === "Under") {
              if (prop.statType.score(projection.result) < prop.target) {
                countCorrect++;
                countTotal++;
              } else {
                countIncorrect++;
                countTotal++;
              }
            }
          }
        }
      });
      //get max proposition
      const maxProp = GetMaxConfidence(projection.propositions);
      if (maxProp.statType.score(projection.result) === maxProp.target) {
        countMaxPush++;
        countMaxTotal++;
      } else {
        if (maxProp.customPrediction.overUnderPrediction === "Over") {
          if (maxProp.statType.score(projection.result) > maxProp.target) {
            countMaxCorrect++;
            countMaxTotal++;
          } else {
            countMaxIncorrect++;
            countMaxTotal++;
          }
        }
        if (maxProp.customPrediction.overUnderPrediction === "Under") {
          if (maxProp.statType.score(projection.result) < maxProp.target) {
            countMaxCorrect++;
            countMaxTotal++;
          } else {
            countMaxIncorrect++;
            countMaxTotal++;
          }
        }
      }
    } else {
      countUnknown += projection.propositions.length;
      countMaxUnknown += 1;
    }
  });
  return (
    <div id="projection-summary">
      <span id="count-teams">{uniqueTeams.length} Teams</span>
      <span id="count-teams">
        {countMaxTotal} Players ({countMaxUnknown} TBD)
      </span>
      <span id="count-teams">
        {countTotal} Props ({countUnknown} TBD)
      </span>
      <span id="count-correct">
        Correct: {countCorrect} (
        {(
          (countCorrect / (countCorrect + countIncorrect + countPush)) *
          100
        ).toFixed(2)}
        %)
      </span>
      <span id="count-incorrect">
        Incorrect: {countIncorrect} (
        {(
          (countIncorrect / (countCorrect + countIncorrect + countPush)) *
          100
        ).toFixed(2)}
        %)
      </span>
      <span id="count-push">
        Push: {countPush} (
        {(
          (countPush / (countCorrect + countIncorrect + countPush)) *
          100
        ).toFixed(2)}
        %)
      </span>

      <span id="count-correct">
        Highest Correct: {countMaxCorrect} (
        {(
          (countMaxCorrect /
            (countMaxCorrect + countMaxIncorrect + countMaxPush)) *
          100
        ).toFixed(2)}
        %)
      </span>
      <span id="count-incorrect">
        Highest Incorrect: {countMaxIncorrect} (
        {(
          (countMaxIncorrect /
            (countMaxCorrect + countMaxIncorrect + countMaxPush)) *
          100
        ).toFixed(2)}
        %)
      </span>
      <span id="count-push">
        Highest Push: {countMaxPush} (
        {(
          (countMaxPush /
            (countMaxCorrect + countMaxIncorrect + countMaxPush)) *
          100
        ).toFixed(2)}
        % )
      </span>
    </div>
  );
};

export default ProjectionsSummary;
