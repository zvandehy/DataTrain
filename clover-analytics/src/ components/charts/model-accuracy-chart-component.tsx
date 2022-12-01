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
import { Bar, Chart } from "react-chartjs-2";
import { useGetPropositions } from "../../hooks/useGetPropositions";
import { DEFAULT_MODEL } from "../../shared/constants";
import { Proposition } from "../../shared/interfaces/graphql/proposition.interface";
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
  startDate: string;
  endDate: string;
  propositions: Proposition[];
}

const ModelAccuracyChart: React.FC<ModelAccuracyChartProps> = ({
  endDate,
  startDate,
  propositions,
}: ModelAccuracyChartProps) => {
  console.log("props", propositions.length);

  let days = [];
  //days between startDate and endDate
  for (let i = moment(startDate); i.isSameOrBefore(endDate); i.add(1, "days")) {
    if (
      propositions.find(
        (prop: Proposition) => prop.game.date === i.format("YYYY-MM-DD")
      )
    ) {
      days.push(i.format("YYYY-MM-DD"));
    }
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
        <Chart
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
              legend: {
                display: true,
                position: "top",
                labels: {
                  color: "#FFF",
                },
                title: {
                  color: "#FFF",
                },
              },
              tooltip: {
                enabled: true,
                mode: "index",
                caretPadding: 10,
              },
            },
            scales: {
              x: {
                ticks: {
                  color: "#FFF",
                },
              },
              y: {
                title: {
                  display: true,
                  text: "# PROPS",
                  color: "#FFF",
                },
                ticks: {
                  color: "#FFF",
                },
              },
              y2: {
                position: "right",
                min: 0,
                max: 100,
                ticks: {
                  color: "white",
                },
                title: {
                  display: true,
                  text: "ACCURACY (%)",
                  color: "white",
                },
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
                type: "bar",
              },
              {
                label: "Misses",
                data: days.map((day) => daysMap[day].misses),
                backgroundColor: COLORS.LOWER,
                type: "bar",
              },
              {
                label: "Pushed",
                data: days.map((day) => daysMap[day].pushes),
                backgroundColor: COLORS.PUSH,
                type: "bar",
              },
              {
                label: "Pending",
                data: days.map((day) => daysMap[day].pending),
                backgroundColor: "gray",
                type: "bar",
              },
              {
                label: "Accuracy",
                data: days.map((day) => daysMap[day].dayAccuracy),
                backgroundColor: COLORS.PRIMARY,
                pointBorderColor: (context: any) => {
                  const index = context.dataIndex;
                  const value = context.dataset.data[index];
                  return value > 60
                    ? COLORS.HIGHER_DARK
                    : value > 50
                    ? COLORS.PUSH_DARK
                    : COLORS.LOWER_DARK;
                },
                pointBackgroundColor: (context: any) => {
                  const index = context.dataIndex;
                  const value = context.dataset.data[index];
                  return value > 60
                    ? COLORS.HIGHER_DARK
                    : value > 50
                    ? COLORS.PUSH_DARK
                    : COLORS.LOWER_DARK;
                },
                borderColor: COLORS.PRIMARY,
                order: -1,
                type: "line",
                yAxisID: "y2",
              },
            ],
          }}
          type="bar"
          color={"#FFFFFF"}
        />
      </Box>
    </Box>
  );
};

export default ModelAccuracyChart;
