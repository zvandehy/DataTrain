import { Box } from "@mui/material";
import {
  BarController,
  ChartDataset,
  ChartOptions,
  LineController,
} from "chart.js";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import moment from "moment";
// import moment from "moment";
import React from "react";
import { Bar } from "react-chartjs-2";
import { useGetProjections } from "../../hooks/useGetProjections";
import { ALL_STATS, DEFAULT_MODEL } from "../../shared/constants";
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
  const startDate = moment(endDate).subtract(3, "days").format("YYYY-MM-DD");
  const {
    loading,
    error,
    data: players,
  } = useGetProjections({
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
  console.log("players", players.length);

  let days = [];
  //days between startDate and endDate
  for (let i = moment(startDate); i.isBefore(endDate); i.add(1, "days")) {
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
  players.forEach((player) => {
    player.games.forEach((game) => {
      game.prediction.propositions.forEach((prop) => {
        if (prop.predictionHit === "HIT") {
          daysMap[game.date].hits++;
          totalHits++;
        } else if (prop.predictionHit === "MISS") {
          daysMap[game.date].misses++;
          totalMisses++;
        } else if (prop.predictionHit === "PUSH") {
          daysMap[game.date].pushes++;
          totalPushes++;
        } else if (prop.predictionHit === "PENDING") {
          daysMap[game.date].pending++;
          totalPending++;
        }
      });
    });
  });

  days.forEach((day) => {
    const dayAccuracy =
      (daysMap[day].hits / (daysMap[day].hits + daysMap[day].misses)) * 100;
    daysMap[day].dayAccuracy = dayAccuracy;
  });

  totalAccuracy = (totalHits / (totalHits + totalMisses)) * 100;

  const options: ChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      mode: "index",
      intersect: false,
    },
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
  };

  return (
    <Box id="chart-wrapper">
      <Box id="chart" minHeight={"400px"} p={2}>
        <Bar
          datasetIdKey="id"
          options={options}
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
            ],
          }}
          color={"#FFFFFF"}
        />
      </Box>
    </Box>
  );
};

export default ModelAccuracyChart;
