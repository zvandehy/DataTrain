import React, {useState} from 'react'
import { Line, Bar, Chart } from 'react-chartjs-2';
import zoomPlugin from 'chartjs-plugin-zoom';
import { GetPropScore } from '../utils';
// import {floor,random} from 'math'
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
  } from 'chart.js';

ChartJS.register(
  CategoryScale,
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
    let {games} = props;
    // games = games.slice(30,40)
    const options =  { responsive: true,
    interaction: {
      mode: 'index',
      intersect: false,
    },
    stacked: false,
    plugins: {
        zoom: {
          pan: {
            // pan options and/or events
            enabled:true,
            mode:"x"
          },
          limits: {
            // axis limits
            minRange: 5,
          },
          zoom: {
            // zoom options and/or events
            wheel: {enabled:true, mode:"x", speed: 0.01},
            drag: {enabled:true, mode:"x"},
            mode:"x",
          }
        },
      title: {
        display: true,
        text: 'Chart.js Line Chart - Multi Axis',
      },
    //   tooltip: {
    //     // Disable the on-canvas tooltip
    //     enabled: true,
    //     callbacks: {
    //         title: () => "Title",
    //     }
    //   },
    },
    scales: {
        y: {display:false, min:0, type:"linear"},
    },
  };
// let games = ['Jun', 'Jul', 'Aug'];
const dates = games.map((game) => game.date);
let type = "points"
const lines = ["minutes", "field_goal_percentage"]
let datasets = Object.getOwnPropertyNames(games[0]).filter(item => item !== "__typename" && item !== 'date').map(
    (prop, i) => {
        const r = Math.floor(Math.random() * 255);
        const g = Math.floor(Math.random() * 255);
        const b = Math.floor(Math.random() * 255);
        return {
        id: prop,
        label: prop,
        type: lines.find(item => item === prop) ? 'line' : 'line',
        fill: true,
        data: games.map(game=>GetPropScore(game,prop)),
        // yAxisID: prop.indexOf("percent") > -1 ? "percent" : "y",
        hidden:true,
        backgroundColor: `rgba(${r}, ${g}, ${b}, 0.5)`,
        color: `rgba(${r}, ${g}, ${b}, 0.5)`,
        fillColor: `rgba(${r}, ${g}, ${b}, 0.5)`,
        borderWidth:2,
        // borderColor: function(context) {
        //     const index = context.dataIndex;
        //     const value = context.dataset.data[index];
        //     if (type !== prop || true) {
        //         return `rgba(${r}, ${g}, ${b}, 0.5)`
        //     }
        //     return value < target ? 'rgba(255,0,0,0.5)' :  // draw negative values in red
        //         'rgba(0,255,0,0.5)';
        // },
    }}
);

let target = 20;
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
        datasetIdKey='id'
        options={options}
        data={{
            labels: dates.map((item, i) => i),
            datasets: datasets,
        }}
        />
  </div>
      
  )
}

export default PlayerStatsChart
