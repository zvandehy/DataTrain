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

ChartJS.register(
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend
);

const PlayerProfileChart = (props) => {
  const { stats } = props;

  const data = {
    labels: stats.map((item) => item.stat),
    datasets: [
      {
        label: "Percent of Team",
        data: stats.map((item) => item.pct),
        backgroundColor: "rgba(255, 99, 132, .7)",
        borderColor: "rgba(255, 99, 132, 1)",
        borderWidth: 1,
        fill: true,
      },
    ],
  };
  const options = {
    responsive: true,
    maintainAspectRatio: false,
    color: "rgb(255,255,255)",
    scales: {
      r: {
        grid: {
          color: ["rgba(255,255,255,1)"],
        },
        angleLines: {
          color: ["rgba(255,255,255,1)"],
        },
        max: 0.3,
        min: 0,
        ticks: {
          display: false,
        },
        pointLabels: {
          color: ["rgba(255,255,255,1)"],
          font: 40,
        },
      },
    },
  };
  return (
    <div id="profile-chart">
      <Radar data={data} options={options} />
    </div>
  );
};

export default PlayerProfileChart;
