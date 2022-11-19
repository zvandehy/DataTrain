import {
  BarController,
  ChartDataset,
  ChartOptions,
  LineController,
} from "chart.js";
// import {
//   Chart as ChartJS,
//   CategoryScale,
//   LinearScale,
//   BarElement,
//   PointElement,
//   LineElement,
//   Title,
//   Tooltip,
//   Legend,
// } from "chart.js";
// import moment from "moment";
// import React from "react";
// import { Chart } from "react-chartjs-2";
// import { ALL_STATS } from "../../../shared/constants";
// import { FilterGames } from "../../../shared/functions/filters.fn";
// import { GameFilter } from "../../../shared/interfaces/graphql/filters.interface";
// import { Game } from "../../../shared/interfaces/graphql/game.interface";
// import { Projection } from "../../../shared/interfaces/graphql/proposition.interface";
// import { Stat } from "../../../shared/interfaces/stat.interface";
// import "./player-stats-chart.component.css";

// ChartJS.register(
//   CategoryScale,
//   LineController,
//   BarController,
//   LinearScale,
//   BarElement,
//   PointElement,
//   LineElement,
//   Title,
//   Tooltip,
//   Legend
//   //   zoomPlugin
// );

// interface PlayerStatsChartProps {
//   games: Game[];
//   selectedStat: Stat;
//   selectedProjection: Projection;
//   gameFilter: GameFilter;
// }

// const PlayerStatsChart: React.FC<PlayerStatsChartProps> = ({
//   games,
//   selectedProjection,
//   selectedStat,
//   gameFilter,
// }: PlayerStatsChartProps) => {
//   games = FilterGames(games, gameFilter).sort((a, b) =>
//     a.date > b.date ? 1 : -1
//   );
//   const options: ChartOptions = {
//     responsive: true,
//     maintainAspectRatio: false,
//     interaction: {
//       mode: "index",
//       intersect: false,
//     },
//     // stacked: false,
//     plugins: {
//       //   zoom: {
//       //     pan: {
//       //       enabled: true,
//       //       mode: "x",
//       //     },
//       //     zoom: {
//       //       // zoom options and/or events
//       //       wheel: { enabled: true, mode: "x" },
//       //       drag: { enabled: true, mode: "x", modifierKey: "shift" },
//       //       pinch: { enabled: true, mode: "x" },
//       //       mode: "x",
//       //     },
//       //   },
//       title: {
//         display: false,
//         // text: `${}`,
//         // color: "#FFF",
//       },
//       tooltip: {
//         enabled: true,
//         mode: "index",
//         caretPadding: 10,
//         callbacks: {
//           // label: function (context) {
//           //   return "Hello world";
//           // },
//           title: function (context) {
//             const index = context[0].dataIndex;
//             const game = games[index];
//             return `${game.team.abbreviation} ${game.margin > 0 ? "+" : ""}${
//               game.margin
//             } ${game.home_or_away === "home" ? "vs" : "@"} ${
//               game.opponent.abbreviation
//             } ${game.date}`;
//           },
//         },
//       },
//     },
//     // scales: {
//     //   y: {
//     //     minutes: { display: false, min: 0, max: 48, type: "linear" },
//     //     percent: { display: false, min: 0, max: 1, type: "linear" },
//     //     display: false,
//     //     min: 0,
//     //     type: "linear",
//     //   },
//     // },
//   };
//   const colors = [
//     "#5D2E8C",
//     "#1B998B",
//     "#3185FC",
//     "#E94900",
//     "#16DB93",
//     "#83E377",
//     "#FF7B9C",
//     "#EFEA5A",
//     "#F29E4C",
//     "#F29E4C",
//   ];
//   const labels = games.map(
//     (game) => `${game.opponent.abbreviation} ${moment(game.date).format("M/D")}`
//   );
//   let datasets: ChartDataset[] = [...ALL_STATS].map((stat, i) => {
//     let data: ChartDataset = {
//       //   id: stat.label,
//       label: stat.display,
//       type: stat === selectedStat ? "bar" : "line",
//       data: games.map((game) => stat.score(game)),
//       // TODO: put axis type as it's own property on Stat object
//       //   yAxisID:
//       //     stat.label.indexOf("percent") > -1
//       //       ? "percent"
//       //       : stat.display === "Minutes"
//       //       ? stat
//       //       : "y",
//       hidden: stat !== selectedStat,
//       order: stat !== selectedStat ? 0 : 1,
//       backgroundColor: function (context: any) {
//         if (stat !== selectedStat) {
//           return `${colors[i % colors.length]}`;
//         }
//         let color = graphColor(context, selectedProjection, stat, colors, i);
//         color = color.replace("rgb", "rgba");
//         color = color.replace(")", ",.5)");
//         return color;
//       },
//       borderColor: function (context: any) {
//         if (stat !== selectedStat) {
//           return `${colors[i % colors.length]}`;
//         }
//         return graphColor(context, selectedProjection, stat, colors, i);
//       },
//       borderWidth: 2,
//       pointBorderColor: function (context: any) {
//         return graphColor(context, selectedProjection, stat, colors, i);
//       },
//       hoverRadius: 8,
//     };
//     return data;
//   });
//   return (
//     <div id="chart-wrapper">
//       <div id="chart">
//         <Chart
//           datasetIdKey="id"
//           options={options}
//           type="line"
//           data={{
//             labels: labels,
//             datasets: datasets,
//           }}
//           color={"#FFFFFF"}
//         />
//       </div>
//     </div>
//   );
// };

// function graphColor(
//   context: any,
//   projection: Projection,
//   selectedStat: Stat,
//   colors: string[],
//   i: number
// ) {
//   //TODO: update color after hidden/unhidden
//   const index = context.dataIndex;
//   const value = context.dataset.data[index];
//   const proposition = projection.propositions.find(
//     (prop) => prop.statType === selectedStat
//   );
//   if (proposition === undefined) {
//     let color = `${colors[i % colors.length]}`;
//     return color;
//   }
//   return value < proposition.target
//     ? "rgb(200,50,50)"
//     : value > proposition.target
//     ? "rgb(0,200,100)"
//     : "rgb(150,150,150)";
// }

// export default PlayerStatsChart;
