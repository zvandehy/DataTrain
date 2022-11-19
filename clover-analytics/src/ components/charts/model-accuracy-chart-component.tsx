import { Box } from "@mui/material";
import {
  BarController,
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LinearScale,
  LineController,
  LineElement,
  PointElement,
  Title,
  Tooltip,
} from "chart.js";
import moment from "moment";
// import moment from "moment";
import React from "react";
import { Bar } from "react-chartjs-2";
import { useGetPropositions } from "../../hooks/useGetPropositions";
import { DEFAULT_MODEL } from "../../shared/constants";
import { COLORS } from "../../shared/styles/constants";
// import { ALL_STATS } from "../../../shared/constants";
// import { FilterGames } from "../../../shared/functions/filters.fn";
// import { GameFilter } from "../../../shared/interfaces/graphql/filters.interface";
// import { Game } from "../../../shared/interfaces/graphql/game.interface";
// import { Projection } from "../../../shared/interfaces/graphql/projection.interface";
// import { Stat } from "../../../shared/interfaces/stat.interface";

ChartJS.register(
  CategoryScale,
  LineController,
  BarController,
  LinearScale,
  BarElement,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

interface ModelAccuracyChartProps {
  endDate: string;
}

const ModelAccuracyChart: React.FC<ModelAccuracyChartProps> = ({
  endDate,
}: ModelAccuracyChartProps) => {
  const startDate = moment(endDate).subtract(2, "days").format("YYYY-MM-DD");
  const {
    loading,
    error,
    data: propositions,
  } = useGetPropositions({
    startDate: startDate,
    endDate: endDate,
    customModel: DEFAULT_MODEL,
  });

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    return <div>Error: {error.message}</div>;
  }
  console.log("props", propositions.length);

  let days = [];
  //days between startDate and endDate
  for (let i = moment(startDate); i.isSameOrBefore(endDate); i.add(1, "days")) {
    days.push(i.format("YYYY-MM-DD"));
  }

  let daysMap: {
    [key: string]: {
      hits: number;
      misses: number;
      pushes: number;
      pending: number;
      dayAccuracy: number;
    };
  } = {};
  days.forEach((day) => {
    daysMap[day] = {
      hits: 0,
      misses: 0,
      pushes: 0,
      pending: 0,
      dayAccuracy: 0,
    };
  });
  let totalHits = 0;
  let totalMisses = 0;
  let totalPushes = 0;
  let totalPending = 0;
  let totalAccuracy = 0;
  propositions.forEach((prop) => {
    if (prop.prediction?.wagerOutcome === "HIT") {
      daysMap[prop.game.date].hits++;
      totalHits++;
    } else if (prop.prediction?.wagerOutcome === "MISS") {
      daysMap[prop.game.date].misses++;
      totalMisses++;
    } else if (prop.prediction?.wagerOutcome === "PUSH") {
      daysMap[prop.game.date].pushes++;
      totalPushes++;
    } else if (prop.prediction?.wagerOutcome === "PENDING") {
      daysMap[prop.game.date].pending++;
      totalPending++;
    }
  });

  days.forEach((day) => {
    const dayAccuracy =
      (daysMap[day].hits / (daysMap[day].hits + daysMap[day].misses)) * 100;
    daysMap[day].dayAccuracy = dayAccuracy;
  });

  totalAccuracy = (totalHits / (totalHits + totalMisses)) * 100;

  return (
    <Box id="chart-wrapper">
      <Box id="chart" minHeight={"400px"} p={2}>
        <Bar
          datasetIdKey="id"
          // options={options}
          options={{
            responsive: true,
            maintainAspectRatio: false,
            interaction: { mode: "index", intersect: false },
            // stacked: true,
            plugins: {
              title: {
                display: true,
                text: `Model Accuracy between ${startDate} and ${endDate}`,
                color: "#FFF",
              },
              subtitle: {
                display: true,
                text: `Total Accuracy: ${totalAccuracy.toFixed(
                  2
                )}% (${totalHits} hits, ${totalMisses} misses)`,
                color: "#FFF",
              },
              tooltip: {
                enabled: true,
                mode: "index",
                caretPadding: 10,
              },
            },
          }}
          data={{
            labels: days,
            datasets: [
              {
                label: "Hits",
                data: days.map((day) => daysMap[day].hits),
                backgroundColor: COLORS.HIGHER,
              },
              {
                label: "Misses",
                data: days.map((day) => daysMap[day].misses),
                backgroundColor: COLORS.LOWER,
              },
              {
                label: "Pushed",
                data: days.map((day) => daysMap[day].pushes),
                backgroundColor: COLORS.PUSH,
              },
              {
                label: "Pending",
                data: days.map((day) => daysMap[day].pending),
                backgroundColor: COLORS.PRIMARY,
              },
            ],
          }}
          color={"#FFFFFF"}
        />
      </Box>
    </Box>
  );
};

export default ModelAccuracyChart;
