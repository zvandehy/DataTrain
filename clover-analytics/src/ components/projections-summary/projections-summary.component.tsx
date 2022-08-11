import { BETTING_CATEGORIES } from "../../shared/constants";
import { ColorPct } from "../../shared/functions/color.fn";
import { GetMaxConfidence } from "../../shared/functions/predictions.fn";
import { Projection } from "../../shared/interfaces/graphql/projection.interface";
import { Stat } from "../../shared/interfaces/stat.interface";
import "./projections-summary.component.css";

interface ProjectionsSummaryProps {
  projections: Projection[];
  filteredStat?: Stat | undefined;
}

const ProjectionsSummary: React.FC<ProjectionsSummaryProps> = ({
  projections,
  filteredStat,
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

  interface DistributionItem {
    min: number;
    max: number;
    correct: number;
    incorrect: number;
  }

  let distribution: DistributionItem[] = [];
  let maxDistribution: DistributionItem[] = [];

  const stepSize = 10;
  for (let x = 0; x <= 100 - stepSize; x += stepSize) {
    distribution.push({
      min: x,
      max: x + stepSize,
      correct: 0,
      incorrect: 0,
    });
    maxDistribution.push({
      min: x,
      max: x + stepSize,
      correct: 0,
      incorrect: 0,
    });
  }

  //TODO: if a statType has more than one proposition, use the most recent one
  projections.forEach((projection) => {
    if (projection.result) {
      BETTING_CATEGORIES.filter(
        (category) => category === filteredStat || !filteredStat
      ).forEach((category) => {
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
                distribution.forEach((dist) => {
                  if (
                    prop.customPrediction.confidence >= dist.min &&
                    prop.customPrediction.confidence < dist.max
                  ) {
                    dist.correct++;
                  }
                });
                countCorrect++;
                countTotal++;
              } else {
                distribution.forEach((dist) => {
                  if (
                    prop.customPrediction.confidence >= dist.min &&
                    prop.customPrediction.confidence < dist.max
                  ) {
                    dist.incorrect++;
                  }
                });
                countIncorrect++;
                countTotal++;
              }
            }
            if (prop.customPrediction.overUnderPrediction === "Under") {
              if (prop.statType.score(projection.result) < prop.target) {
                distribution.forEach((dist) => {
                  if (
                    prop.customPrediction.confidence >= dist.min &&
                    prop.customPrediction.confidence < dist.max
                  ) {
                    dist.correct++;
                  }
                });
                countCorrect++;
                countTotal++;
              } else {
                distribution.forEach((dist) => {
                  if (
                    prop.customPrediction.confidence >= dist.min &&
                    prop.customPrediction.confidence < dist.max
                  ) {
                    dist.incorrect++;
                  }
                });
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
            maxDistribution.forEach((dist) => {
              if (
                maxProp.customPrediction.confidence > dist.min &&
                maxProp.customPrediction.confidence <= dist.max
              ) {
                dist.correct++;
              }
            });
          } else {
            maxDistribution.forEach((dist) => {
              if (
                maxProp.customPrediction.confidence > dist.min &&
                maxProp.customPrediction.confidence <= dist.max
              ) {
                dist.incorrect++;
              }
            });
            countMaxIncorrect++;
            countMaxTotal++;
          }
        }
        if (maxProp.customPrediction.overUnderPrediction === "Under") {
          if (maxProp.statType.score(projection.result) < maxProp.target) {
            maxDistribution.forEach((dist) => {
              if (
                maxProp.customPrediction.confidence > dist.min &&
                maxProp.customPrediction.confidence <= dist.max
              ) {
                dist.correct++;
              }
            });
            countMaxCorrect++;
            countMaxTotal++;
          } else {
            maxDistribution.forEach((dist) => {
              if (
                maxProp.customPrediction.confidence > dist.min &&
                maxProp.customPrediction.confidence <= dist.max
              ) {
                dist.incorrect++;
              }
            });
            countMaxIncorrect++;
            countMaxTotal++;
          }
        }
      }
    } else {
      countUnknown += projection.propositions.filter(
        (prop) => prop.statType === filteredStat || !filteredStat
      ).length;
      countMaxUnknown += 1;
    }
  });
  return (
    <>
      <div id="projection-summary">
        <span id="count-teams">{uniqueTeams.length} Teams</span>
        <span id="count-teams">
          {countMaxTotal} Players ({countMaxUnknown} TBD)
        </span>
        <span id="count-teams">
          {countTotal} Props ({countUnknown} TBD)
        </span>
        {countCorrect + countIncorrect ? (
          <>
            <span
              id="count-correct"
              className={ColorPct(
                countCorrect / (countCorrect + countIncorrect + countPush)
              )}
            >
              Correct: {countCorrect} (
              {(
                (countCorrect / (countCorrect + countIncorrect + countPush)) *
                100
              ).toFixed(2)}
              %)
            </span>
            <span
              id="count-incorrect"
              className={ColorPct(
                countIncorrect / (countCorrect + countIncorrect + countPush)
              )}
            >
              Incorrect: {countIncorrect} (
              {(
                (countIncorrect / (countCorrect + countIncorrect + countPush)) *
                100
              ).toFixed(2)}
              %)
            </span>
          </>
        ) : (
          <></>
        )}
        {countPush ? (
          <span
            id="count-push"
            className={ColorPct(
              countPush / (countCorrect + countIncorrect + countPush)
            )}
          >
            Push: {countPush} (
            {(
              (countPush / (countCorrect + countIncorrect + countPush)) *
              100
            ).toFixed(2)}
            %)
          </span>
        ) : (
          <></>
        )}

        {countMaxCorrect + countMaxIncorrect ? (
          <>
            <span
              id="count-correct"
              className={ColorPct(
                countMaxCorrect /
                  (countMaxCorrect + countMaxIncorrect + countMaxPush)
              )}
            >
              Highest Correct: {countMaxCorrect} (
              {(
                (countMaxCorrect /
                  (countMaxCorrect + countMaxIncorrect + countMaxPush)) *
                100
              ).toFixed(2)}
              %)
            </span>
            <span
              id="count-incorrect"
              className={ColorPct(
                countMaxIncorrect /
                  (countMaxCorrect + countMaxIncorrect + countMaxPush)
              )}
            >
              Highest Incorrect: {countMaxIncorrect} (
              {(
                (countMaxIncorrect /
                  (countMaxCorrect + countMaxIncorrect + countMaxPush)) *
                100
              ).toFixed(2)}
              %)
            </span>
          </>
        ) : (
          <></>
        )}
        {countMaxPush ? (
          <span
            id="count-push"
            className={ColorPct(
              countMaxPush /
                (countMaxCorrect + countMaxIncorrect + countMaxPush)
            )}
          >
            Highest Push: {countMaxPush} (
            {(
              (countMaxPush /
                (countMaxCorrect + countMaxIncorrect + countMaxPush)) *
              100
            ).toFixed(2)}
            % )
          </span>
        ) : (
          <></>
        )}
      </div>
      {countTotal > 0 ? (
        <div id="distribution-summary">
          <span>All</span>
          {distribution
            .filter((dist) => dist.correct + dist.incorrect > 0)
            .map((dist) => {
              return (
                <span
                  className={ColorPct(
                    dist.correct / (dist.correct + dist.incorrect)
                  )}
                >
                  {dist.min}-{dist.max}%: {dist.correct}-{dist.incorrect} (
                  {(
                    (dist.correct / (dist.correct + dist.incorrect)) *
                    100
                  ).toFixed(2)}
                  %)
                </span>
              );
            })}
        </div>
      ) : (
        <></>
      )}
      {countMaxTotal > 0 ? (
        <div id="max-distribution-summary">
          <span>Highest</span>
          {maxDistribution
            .filter((dist) => dist.correct + dist.incorrect > 0)
            .map((dist) => {
              return (
                <span
                  className={ColorPct(
                    dist.correct / (dist.correct + dist.incorrect)
                  )}
                >
                  {dist.min}-{dist.max}%: {dist.correct}-{dist.incorrect} (
                  {(
                    (dist.correct / (dist.correct + dist.incorrect)) *
                    100
                  ).toFixed(2)}
                  %)
                </span>
              );
            })}
        </div>
      ) : (
        <></>
      )}
    </>
  );
};

export default ProjectionsSummary;
