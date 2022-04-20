import React from "react";
import "chart.js/auto";
import { Chart } from "react-chartjs-2";
import zoomPlugin from "chartjs-plugin-zoom";
import { GetPropScore } from "../utils";

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
import { LineController } from "chart.js";
import { BarController } from "chart.js";
import { index } from "mathjs";

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
  Legend,
  zoomPlugin
);

const PlayerStatsChart = (props) => {
  let { games } = props;
  const options = {
    responsive: true,
    interaction: {
      mode: "index",
      intersect: false,
    },
    stacked: false,
    plugins: {
      zoom: {
        pan: {
          enabled: true,
          mode: "x",
        },
        zoom: {
          // zoom options and/or events
          wheel: { enabled: true, mode: "x" },
          drag: { enabled: true, mode: "x", modifierKey: "shift" },
          pinch: { enabled: true, mode: "x" },
          mode: "x",
        },
      },
      title: {
        display: true,
        text: "Chart.js Line Chart - Multi Axis",
      },
      tooltip: {
        enabled: true,
        mode: "index",
        // callbacks: {
        //     title: () => "Title",
        //     label: (context) => {},
        //     }
      },
    },
    scales: {
      y: {
        minutes: { display: false, min: 0, max: 48, type: "linear" },
        percent: { display: false, min: 0, max: 1, type: "linear" },
        display: false,
        min: 0,
        type: "linear",
      },
    },
  };
  const dates = games.map((game) => game.date);
  let type = "points";
  const lines = ["minutes", "field_goal_percentage"];
  let datasets = Object.getOwnPropertyNames(games[0])
    .filter((item) => item !== "__typename" && item !== "date")
    .map((prop, i) => {
      const r = Math.floor(Math.random() * 255);
      const g = Math.floor(Math.random() * 255);
      const b = Math.floor(Math.random() * 255);
      return {
        id: prop,
        label: prop,
        type: lines.find((item) => item === prop) ? "line" : "line",
        fill: true,
        data: games.map((game) => GetPropScore(game, prop)),
        yAxisID:
          prop.indexOf("percent") > -1 || prop === "usage"
            ? "percent"
            : prop === "minutes"
            ? prop
            : "y",
        hidden: prop !== type,
        order: prop !== type ? (prop === "minutes" ? 100 : i) : 0,
        backgroundColor: `rgba(${r}, ${g}, ${b}, 0.5)`,
        color: `rgba(${r}, ${g}, ${b}, 0.5)`,
        fillColor: `rgba(${r}, ${g}, ${b}, 0.5)`,
        borderWidth: 2,
        pointBorderColor: function (context) {
          const index = context.dataIndex;
          const value = context.dataset.data[index];
          if (type !== prop) {
            return `rgba(${r}, ${g}, ${b}, 0.5)`;
          }
          return value < target ? "rgba(255,0,0,0.5)" : "rgba(0,255,0,0.5)";
        },
      };
    });

  let target = 9.5;
  // let propDataset = datasets[datasets.findIndex(item => item.id === prop)]
  // for (let i=0; i<datasets[datasets.findIndex(item => item.id === prop)].data.length; i++) {
  //     if (GetPropScore({points:datasets[datasets.findIndex(item => item.id === prop)].data[i]},prop) > target) {
  //         datasets[datasets.findIndex(item => item.id === prop)].pointBackgroundColors.push("rgba(0,255,0,0.5)")
  //     } else {
  //         datasets[datasets.findIndex(item => item.id === prop)].pointBackgroundColors.push("rgba(255,0,0,0.5)")
  //     }
  // datasets[datasets.findIndex(item => item.id === prop)].pointBackgroundColor = GetPropScore({points:datasets[datasets.findIndex(item => item.id === prop)].data[i]},prop) > 20 ? 'green' : 'red'
  // }
  return (
    <div id="chart">
      <Chart
        datasetIdKey="id"
        options={options}
        data={{
          labels: dates.map((item, i) => i),
          datasets: datasets,
        }}
      />
    </div>
  );
};

export default PlayerStatsChart;
