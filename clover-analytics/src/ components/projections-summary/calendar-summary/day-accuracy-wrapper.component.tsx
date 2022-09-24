import { CircularProgress } from "@mui/material";
import { Box } from "@mui/system";
import moment from "moment";
import { useGetProjections } from "../../../hooks/useGetProjections";
import { Match } from "../../../shared/functions/filters.fn";
import { CalculatePredictions } from "../../../shared/functions/predictions.fn";
import { CustomCalculation } from "../../../shared/interfaces/custom-prediction.interface";
import {
  Accuracy,
  HitCriteria,
} from "../../../shared/interfaces/accuracy.interface";
import {
  GameFilter,
  ProjectionFilter,
} from "../../../shared/interfaces/graphql/filters.interface";
import { Stat } from "../../../shared/interfaces/stat.interface";
import DayAccuracy from "./day-accuracy.component";

interface DayAccuracyProps {
  date: string;
  customModel: CustomCalculation;
  statType: Stat | undefined;
  lookup: string;
  setAccuracy: React.Dispatch<React.SetStateAction<Accuracy>>;
  hitCriteria: HitCriteria;
}

const DayAccuracyWrapper: React.FC<DayAccuracyProps> = ({
  date,
  customModel,
  statType,
  lookup,
  setAccuracy,
  hitCriteria,
}: DayAccuracyProps) => {
  let projectionFilter: ProjectionFilter = {
    startDate: moment(date).format("YYYY-MM-DD"),
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  let gameFilter: GameFilter = {
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  const predictionFilter: GameFilter = {
    season: "2022-23",
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  const { loading, error, data } = useGetProjections({
    projectionFilter,
    gameFilter,
    predictionFilter,
    customModel,
  });

  if (loading) {
    return (
      <Box>
        <CircularProgress size={15} />
      </Box>
    );
  }
  if (error || !data) {
    return <></>;
  }
  let filteredProjections = CalculatePredictions(
    data,
    gameFilter,
    customModel
  ).filter((projection) => {
    if (statType !== undefined) {
      return Match(projection, { lookup: lookup, statType: statType as Stat });
    }
    return Match(projection, { lookup: lookup });
  });
  if (statType) {
    filteredProjections.forEach((projection, i) => {
      filteredProjections[i].propositions = projection.propositions.filter(
        (prop) => prop.statType.abbreviation === statType.abbreviation
      );
    });
  }
  if (filteredProjections.length === 0) {
    return <Box>0-0</Box>;
  }
  return filteredProjections.length > 0 ? (
    <DayAccuracy
      accuracy={new Accuracy(filteredProjections)}
      setAccuracy={setAccuracy}
      projections={filteredProjections}
      hitCriteria={hitCriteria}
    />
  ) : (
    <></>
  );
};

export default DayAccuracyWrapper;
