import React from "react";
import {
  Chart as ChartJS,
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend,
} from "chart.js";
import { Radar } from "react-chartjs-2";
import { RelevantStats, AveragePropScore } from "../utils";

ChartJS.register(
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend
);

const PlayerProfileChart = (props) => {
  const { games } = props;
  const profileStats = RelevantStats["Profile"];

  const data = {
    labels: profileStats.map((item) => item.label),
    datasets: [
      {
        label: "Player Profile",
        data: profileStats.map((item) =>
          AveragePropScore(games, item.recognize)
        ),
        backgroundColor: "rgba(255, 99, 132, 0.2)",
        borderColor: "rgba(255, 99, 132, 1)",
        borderWidth: 1,
      },
    ],
  };
  const options = {
    responsive: true,
  };
  return (
    <div id="profile-chart">
      <Radar data={data} options={options} />
    </div>
  );
};

export default PlayerProfileChart;
