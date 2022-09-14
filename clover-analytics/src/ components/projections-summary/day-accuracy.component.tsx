import { Box } from "@mui/system";
import moment from "moment";
import { useMemo } from "react";
import { useGetProjections } from "../../hooks/useGetProjections";
import { Match } from "../../shared/functions/filters.fn";
import { CalculatePredictions } from "../../shared/functions/predictions.fn";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
import {
  GameFilter,
  ProjectionFilter,
} from "../../shared/interfaces/graphql/filters.interface";
import { Stat } from "../../shared/interfaces/stat.interface";
import { AccuracyAction } from "./accuracy-reducer.component";
import ProjectionsSummaryShort from "./projections-summary-short.component";

interface DayAccuracyProps {
  date: string;
  customModel: CustomCalculation;
  statType: Stat | undefined;
  lookup: string;
  dispatch: (action: AccuracyAction) => void;
}

const DayAccuracy: React.FC<DayAccuracyProps> = ({
  date,
  customModel,
  statType,
  lookup,
  dispatch,
}: DayAccuracyProps) => {
  // console.log("Get summary for: ", date);
  let projectionFilter: ProjectionFilter = {
    startDate: moment(date).format("YYYY-MM-DD"),
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  let gameFilter: GameFilter = {
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  const predictionFilter: GameFilter = {
    season: "2022-23", // TODO
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  const { loading, error, data } = useGetProjections({
    projectionFilter,
    gameFilter,
    predictionFilter,
    customModel,
  });

  if (loading) {
    return <p>Loading...</p>;
  }
  if (error || !data) {
    return <>{error}</>;
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
  return filteredProjections.length > 0 ? (
    <Box sx={{ display: "flex", alignItems: "start" }}>
      <div>{moment(date).format("YYYY-MM-DD")}</div>
      <ProjectionsSummaryShort
        projections={filteredProjections}
        filteredStat={statType}
        dispatch={dispatch}
      />
    </Box>
  ) : (
    <></>
  );
};

export default DayAccuracy;
