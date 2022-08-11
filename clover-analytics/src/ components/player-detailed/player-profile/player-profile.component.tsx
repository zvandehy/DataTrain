import "./player-profile.component.css";
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
import { Player } from "../../../shared/interfaces/graphql/player.interface";
import {
  Assists,
  Blocks,
  Points,
  Rebounds,
  Steals,
  ThreeFGA,
} from "../../../shared/interfaces/stat.interface";

import { Game } from "../../../shared/interfaces/graphql/game.interface";

ChartJS.register(
  RadialLinearScale,
  PointElement,
  LineElement,
  Filler,
  Tooltip,
  Legend
);

interface PlayerProfileChartProps {
  player: Player;
  filteredGames: Game[];
}

const PlayerProfileChart: React.FC<PlayerProfileChartProps> = ({
  player,
  filteredGames: games,
}: PlayerProfileChartProps) => {
  const profileChartStats = [
    Points,
    Assists,
    ThreeFGA,
    Rebounds,
    Blocks,
    Steals,
  ];

  interface PercentOfTeam {
    average: number;
    percent: number;
    label: string;
  }

  const percentOfTeam: PercentOfTeam[] = profileChartStats.map((stat) => {
    const avg = stat.average(games);
    const pct = +(
      avg / +stat.teamAverage(games.map((game) => game.teamStats)).toFixed(2)
    );
    return {
      label: stat.abbreviation,
      average: avg,
      percent: pct,
    };
  });

  const data = {
    labels: percentOfTeam.map((item) => item.label),
    datasets: [
      {
        label: "Percent of Team",
        data: percentOfTeam.map((item) => +(item.percent * 100).toFixed(2)),
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
          color: "rgba(255,255,255,1)",
        },
        max: 50,
        min: 0,
        ticks: {
          display: false,
        },
        pointLabels: {
          color: "rgba(255,255,255,1)",
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
