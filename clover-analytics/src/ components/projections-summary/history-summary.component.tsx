import { Box } from "@mui/material";
import moment from "moment";
import { useMemo, useReducer, useState } from "react";
import { ColorPct } from "../../shared/functions/color.fn";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
import { Stat } from "../../shared/interfaces/stat.interface";
import {
  AccuracyReducer,
  INITIALIZE_DISTRIBUTION,
} from "./accuracy-reducer.component";
import DayAccuracy from "./day-accuracy.component";

interface ModelAccuracyHistoryProps {
  customModel: CustomCalculation;
  filteredStat: Stat | undefined;
  lookup: string;
  date: string;
}

const ModelAccuracyHistory: React.FC<ModelAccuracyHistoryProps> = ({
  customModel,
  filteredStat,
  lookup,
  date,
}: ModelAccuracyHistoryProps) => {
  console.log("Render ModelAccuracyHistory");
  const [totalDist, dispatch] = useReducer(
    AccuracyReducer,
    INITIALIZE_DISTRIBUTION()
  );
  // TODO: Doing many games is slow on the client
  const [dateRange, setDateRange] = useState({
    start: moment(date),
    numDays: 1,
  });
  let dates: string[] = [];
  for (let i = 1; i <= dateRange.numDays; i++) {
    dates.push(
      moment(dateRange.start)
        .add(i * -1, "days")
        .format("YYYY-MM-DD")
    );
  }
  const daySummaries = useMemo(() => {
    console.log("Get all summaries");
    return dates.map((date) => (
      <DayAccuracy
        key={moment(date).format("YYYY-MM-DD")}
        date={date}
        customModel={customModel}
        statType={filteredStat}
        lookup={lookup}
        dispatch={dispatch}
      />
    ));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [date, JSON.stringify(customModel)]);
  return (
    <Box
      className="model-accuracy"
      sx={{ display: "block", gridArea: "accuracy", marginBottom: "2rem" }}
    >
      <h1>ALL</h1>
      <div>
        {totalDist.distributions
          .filter((dist) => dist.correct + dist.incorrect > 0)
          .map((dist) => {
            return (
              <span
                key={`highest-${dist.min}-${dist.max}`}
                className={ColorPct(
                  dist.correct / (dist.correct + dist.incorrect)
                )}
                style={{ paddingInline: "5px" }}
              >
                {dist.min}-{dist.max}%: {dist.correct}-{dist.incorrect} (
                {(
                  (dist.correct / (dist.correct + dist.incorrect)) *
                  100
                ).toFixed(2)}
                %),{" "}
              </span>
            );
          })}
      </div>
      <ul>{daySummaries.map((day) => day)}</ul>
    </Box>
  );
};

export default ModelAccuracyHistory;
