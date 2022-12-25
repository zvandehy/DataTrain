import React from "react";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Doughnut } from "react-chartjs-2";
import { COLORS } from "../../shared/styles/constants";
import { Proposition } from "../../shared/interfaces/graphql/proposition.interface";

ChartJS.register(ArcElement, Tooltip, Legend);

export interface OverUnderPieChartProps {
  propositions: Proposition[];
}

export const OverUnderPieChart: React.FC<OverUnderPieChartProps> = ({
  propositions,
}: OverUnderPieChartProps) => {
  let hits = 0;
  let pushes = 0;
  let misses = 0;
  let pending = 0;
  propositions.forEach((proposition) => {
    if (proposition.prediction?.wagerOutcome === "HIT") {
      hits++;
    } else if (proposition.prediction?.wagerOutcome === "PUSH") {
      pushes++;
    } else if (proposition.prediction?.wagerOutcome === "MISS") {
      misses++;
    } else {
      pending++;
    }
  });

  const data = {
    labels: ["Hit", "Push", "Miss", "Pending"],
    datasets: [
      {
        label: "none",
        data: [hits, pushes, misses, pending],
        backgroundColor: [
          COLORS.HIGHER,
          COLORS.PUSH,
          COLORS.LOWER,
          "rgba(201, 203, 207, 0.5)",
        ],
        borderColor: [
          COLORS.HIGHER_DARK,
          COLORS.PUSH_DARK,
          COLORS.LOWER_DARK,
          "rgba(201, 203, 207, 1)",
        ],
        borderWidth: 1,
      },
    ],
  };

  return (
    <Doughnut
      data={data}
      options={{
        responsive: true,
        maintainAspectRatio: true,
        plugins: {
          legend: {
            display: false,
            position: "right",
          },
        },
      }}
    />
  );
};
