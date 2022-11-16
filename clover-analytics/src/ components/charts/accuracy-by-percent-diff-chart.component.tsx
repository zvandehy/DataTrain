import React from "react";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Chart } from "react-chartjs-2";

import { COLORS } from "../../shared/styles/constants";
import { PropositionA } from "../../shared/interfaces/graphql/game.interface";

ChartJS.register(ArcElement, Tooltip, Legend);

export interface ModelAccuracyByPctDiffProps {
  propositions: PropositionA[];
  stepSize?: number;
  steps?: number;
}

export const ModelAccuracyByPctDiff: React.FC<ModelAccuracyByPctDiffProps> = ({
  propositions,
  stepSize = 3,
  steps = 5,
}: ModelAccuracyByPctDiffProps) => {
  const steppers = Array.from(
    { length: steps * 2 + 1 },
    (_, i) => (i - steps) * stepSize
  );

  // get all ranges between steppers, including +/- infinity
  const ranges: { min: number; max: number }[] = [];
  steppers.forEach((step, i) => {
    if (i === 0) {
      ranges.push({ min: -Infinity, max: step });
      ranges.push({ min: step, max: step + stepSize });
    } else if (i === steppers.length - 1) {
      ranges.push({ min: step, max: Infinity });
    } else {
      ranges.push({ min: step, max: steppers[i + 1] });
    }
  });

  let breakdowns: {
    label: string;
    hits: number;
    tick: string;
    misses: number;
    pushes: number;
    pending: number;
  }[] = [];

  ranges.forEach((x, i) => {
    const filtered = propositions.filter((prop) => {
      return (
        Math.abs(prop.predictionTargetDiffPCT) >= x.min &&
        Math.abs(prop.predictionTargetDiffPCT) < x.max
      );
    });
    const hits = filtered.filter((prop) => prop.predictionHit === "HIT").length;
    const misses = filtered.filter(
      (prop) => prop.predictionHit === "MISS"
    ).length;
    const pushes = filtered.filter(
      (prop) => prop.predictionHit === "PUSH"
    ).length;
    const pending = filtered.filter(
      (prop) => prop.predictionHit === "PENDING"
    ).length;
    let label = `${x.min}% to ${x.max}%`;
    if (x.min === -Infinity) {
      label = `<= ${x.max}%`;
    } else if (x.max === Infinity) {
      label = `>= ${x.min}%`;
    }
    if (filtered.length > 0) {
      breakdowns.push({
        label: label,
        tick: `${x.min}`,
        hits,
        misses,
        pushes,
        pending,
      });
    }
  });

  const accuracy = breakdowns.map((item) => {
    const { hits, misses, pushes } = item;
    return +(((hits + pushes) / (hits + misses + pushes)) * 100).toFixed();
  });

  const total = breakdowns.map((item) => {
    const { hits, misses, pushes, pending } = item;
    return hits + misses + pushes + pending;
  });

  const hits = breakdowns.map((item) => {
    const { hits } = item;
    return hits;
  });

  const misses = breakdowns.map((item) => {
    const { misses } = item;
    return misses;
  });

  const pushes = breakdowns.map((item) => {
    const { pushes } = item;
    return pushes;
  });

  const pending = breakdowns.map((item) => {
    const { pending } = item;
    return pending;
  });

  return (
    <Chart
      type="bar"
      title={"TODAY'S ACCURACY BY MODEL PCT DIFF"}
      data={{
        labels: breakdowns.map((item) => item.label),
        datasets: [
          {
            label: "Hits",
            data: hits,
            type: "bar",
            backgroundColor: COLORS.HIGHER,
            barPercentage: 1,
            categoryPercentage: 1,
            order: 1,
          },
          {
            label: "Misses",
            data: misses,
            type: "bar",
            backgroundColor: COLORS.LOWER,
            barPercentage: 1,
            categoryPercentage: 1,
            order: 1,
          },

          {
            label: "Push",
            data: pushes,
            type: "bar",
            backgroundColor: COLORS.PUSH,
            order: 1,
            barPercentage: 1,
            categoryPercentage: 1,
          },
          {
            label: "Pending",
            data: pending,
            type: "bar",
            backgroundColor: "gray",
            barPercentage: 1,
            categoryPercentage: 1,
            order: 1,
          },
          {
            label: "Total",
            data: total,
            type: "bar",
            backgroundColor: COLORS.SECONDARY,
            barPercentage: 1,
            categoryPercentage: 1,
            order: 1,
            hidden: true,
          },
          {
            label: "Accuracy",
            data: accuracy,
            type: "line",
            backgroundColor: COLORS.PRIMARY,
            borderColor: COLORS.PRIMARY,
            order: 0,
            yAxisID: "y2",
          },
        ],
      }}
      options={{
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          tooltip: {
            mode: "index",
          },
          title: {
            display: true,
            text: "MODEL ACCURACY BY PERCENT CHANGE",
            color: "white",
          },
          legend: {
            labels: { color: "white" },
            display: true,
          },
        },
        scales: {
          x: {
            stacked: true,
            ticks: {
              color: "white",
              align: "end",
            },
            display: false,
          },
          x1: {
            stacked: true,
            ticks: {
              color: "white",
              maxTicksLimit: breakdowns.length - 1,
              autoSkip: false,
            },
            title: {
              display: true,
              text: "Distance from Proposition Target (%)",
              color: "white",
            },
            labels: [...breakdowns.map((x) => `${x.tick}%`), "∞"],
          },
          y: {
            stacked: true,
            title: {
              display: true,
              text: "# Props",
              color: "white",
            },
            ticks: {
              color: "white",
              align: "start",
            },
            grid: {
              display: false,
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
    />
  );
};
