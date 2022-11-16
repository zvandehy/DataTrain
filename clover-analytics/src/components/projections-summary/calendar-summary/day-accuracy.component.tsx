import { Box } from "@mui/system";
import moment from "moment";
import { useEffect } from "react";
import {
  Accuracy,
  BreakdownFilter,
  GetHitRate,
  HitCriteria,
  HitCriteriaType,
  HitRate,
} from "../../../shared/interfaces/accuracy.interface";
import { Projection } from "../../../shared/interfaces/graphql/projection.interface";

interface DayAccuracyProps {
  accuracy: Accuracy;
  setAccuracy?: React.Dispatch<React.SetStateAction<Accuracy>>;
  projections?: Projection[];
  hitCriteria: HitCriteria;
}

const DayAccuracy: React.FC<DayAccuracyProps> = ({
  accuracy,
  setAccuracy,
  projections,
  hitCriteria,
}: DayAccuracyProps) => {
  useEffect(() => {
    if (
      setAccuracy &&
      projections &&
      !accuracy.allProps.find((prop) =>
        moment(prop.date).isSame(projections[0].startTime)
      )
    ) {
      setAccuracy((prev) => {
        let newAccuracy = prev;
        newAccuracy.add(projections);
        return newAccuracy;
      });
    }
  }, [accuracy, projections, setAccuracy]);
  const withinConfidence = accuracy.filter(
    new BreakdownFilter({
      min: hitCriteria.confidenceThreshold,
    })
  );
  const results =
    hitCriteria.hitType === HitCriteriaType.ALL_PROPS
      ? withinConfidence.allProps
      : withinConfidence.topPlayerProps;
  const hit = GetHitRate(hitCriteria, results);
  return <Box>{`${hit.correct}-${hit.incorrect}`}</Box>;
};

export default DayAccuracy;
