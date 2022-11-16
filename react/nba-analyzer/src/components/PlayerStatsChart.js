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
  let { games, predictions, selected } = props;
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
        display: false,
        // text: `${}`,
        // color: "#FFF",
      },
      tooltip: {
        enabled: true,
        mode: "index",
        callbacks: {
          // label: function (context) {
          //   return "Hello world";
          // },
          title: function (context) {
            const index = context[0].dataIndex;
            const game = games[index];
            return `${game.team.abbreviation} ${game.margin > 0 ? "+" : ""}${
              game.margin
            } ${game.home_or_away === "home" ? "vs" : "@"} ${
              game.opponent.abbreviation
            } ${game.date}`;
          },
        },
      },
    },
    // scales: {
    //   y: {
    //     minutes: { display: false, min: 0, max: 48, type: "linear" },
    //     percent: { display: false, min: 0, max: 1, type: "linear" },
    //     display: false,
    //     min: 0,
    //     type: "linear",
    //   },
    // },
  };
  const colors = [
    "#5D2E8C",
    "#1B998B",
    "#3185FC",
    "#E94900",
    "#16DB93",
    "#83E377",
    "#FF7B9C",
    "#EFEA5A",
    "#F29E4C",
    "#F29E4C",
  ];
  const labels = games.map((game) => `${game.opponent.abbreviation}`);
  let datasets = [
    "Fantasy Score",
    "Points",
    "Rebounds",
    "Assists",
    "Steals",
    "Blocks",
    "Turnovers",
    "Minutes",
    "Free Throws Made",
    "3-PT Made",
    "Three_Pointers_Attempted",
  ].map((prop, i) => {
    return {
      id: prop,
      label: prop,
      type:
        prop.toLowerCase() === selected.recognize?.toLowerCase()
          ? "bar"
          : "line",
      data: games.map((game) => GetPropScore(game, prop)),
      yAxisID:
        prop.indexOf("percent") > -1 || prop === "usage"
          ? "percent"
          : prop === "minutes"
          ? prop
          : "y",
      hidden: prop !== selected.recorgnize,
      order: prop.toLowerCase() !== selected.recognize.toLowerCase() ? 0 : 1,
      backgroundColor: `${colors[i % colors.length]}`,
      borderColor: function (context) {
        if (prop.toLowerCase() !== selected.recognize?.toLowerCase()) {
          return `${colors[i % colors.length]}`;
        }
        return graphColor(context, predictions, prop, colors, i);
      },
      borderWidth: 2,
      pointBorderColor: function (context) {
        return graphColor(context, predictions, prop, colors, i);
      },
    };
  });
  return (
    <div id="chart-wrapper">
      <div id="chart">
        <Chart
          datasetIdKey="id"
          options={options}
          data={{
            labels: labels,
            datasets: datasets,
          }}
          color={"#FFFFFF"}
        />
      </div>
    </div>
  );
};

function graphColor(context, predictions, prop, colors, i) {
  const index = context.dataIndex;
  const value = context.dataset.data[index];
  const prediction = predictions.find(
    (prediction) =>
      prediction.stat.recognize.toLowerCase() === prop.toLowerCase()
  );
  if (prediction === undefined) {
    return `${colors[i % colors.length]}`;
  }
  return value <= prediction.target ? "rgba(255,0,0)" : "rgba(0,255,0, .8)";
}

export default PlayerStatsChart;
