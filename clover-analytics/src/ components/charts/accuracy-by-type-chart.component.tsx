import React from "react";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";
import { Bar } from "react-chartjs-2";

import { COLORS } from "../../shared/styles/constants";
import { GetStatAbbreviation } from "../../shared/interfaces/stat.interface";
import { Proposition } from "../../shared/interfaces/graphql/proposition.interface";

ChartJS.register(ArcElement, Tooltip, Legend);

export interface ModelAccuracyByStatTypeProps {
  propositions: Proposition[];
}

export const ModelAccuracyByStatType: React.FC<
  ModelAccuracyByStatTypeProps
> = ({ propositions }: ModelAccuracyByStatTypeProps) => {
  const statBreakdowns = new Map<
    string,
    {
      hits: number;
      misses: number;
      pushes: number;
      pending: number;
      total: number;
    }
  >();
  propositions.forEach((prop) => {
    const breakdown = statBreakdowns.get(prop.type) ?? {
      hits: 0,
      misses: 0,
      pushes: 0,
      pending: 0,
      total: 0,
    };
    if (prop.prediction?.wagerOutcome === "HIT") {
      breakdown.hits++;
    }
    if (prop.prediction?.wagerOutcome === "MISS") {
      breakdown.misses++;
    }
    if (prop.prediction?.wagerOutcome === "PUSH") {
      breakdown.pushes++;
    }
    if (prop.prediction?.wagerOutcome === "PENDING") {
      breakdown.pending++;
    }
    breakdown.total++;
    statBreakdowns.set(prop.type, breakdown);
  });

  const labels: string[] = [];
  const hits: number[] = [];
  const misses: number[] = [];
  const pushes: number[] = [];
  const pending: number[] = [];
  let keys = Array.from(statBreakdowns.keys()).sort((a, b) => {
    return (
      statBreakdowns.get(b)!.hits / statBreakdowns.get(b)!.total -
      statBreakdowns.get(a)!.hits / statBreakdowns.get(a)!.total
    );
  });
  keys.forEach((key) => {
    const breakdown = statBreakdowns.get(key)!;
    labels.push(GetStatAbbreviation(key));
    hits.push(breakdown.hits);
    misses.push(breakdown.misses);
    pushes.push(breakdown.pushes);
    pending.push(breakdown.pending);
  });

  return (
    <Bar
      title={"TODAY'S ACCURACY BY STAT TYPE"}
      data={{
        labels: labels,
        datasets: [
          {
            label: "Hit",
            data: hits,
            backgroundColor: COLORS.HIGHER,
          },
          {
            label: "Miss",
            data: misses,
            backgroundColor: COLORS.LOWER,
          },
          {
            label: "Push",
            data: pushes,
            backgroundColor: COLORS.PUSH,
          },
          {
            label: "Pending",
            data: pending,
            backgroundColor: "gray",
          },
        ],
      }}
      options={{
        indexAxis: "y",
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          tooltip: {
            mode: "index",
            callbacks: {
              title: function (context) {
                const index = labels.indexOf(context[0].label);
                let total = hits[index] + misses[index] + pushes[index];

                return (
                  context[0].label +
                  " (" +
                  ((hits[index] / total) * 100).toFixed(0) +
                  "%)"
                );
              },
            },
          },
          title: {
            display: true,
            text: "MODEL ACCURACY BY STAT TYPE",
            color: "white",
          },
          legend: {
            display: true,
            position: "right",
            title: {
              color: "white",
            },
            labels: {
              color: "white",
            },
          },
        },
        scales: {
          x: {
            stacked: true,
            title: {
              display: true,
              text: "# Props",
              color: "white",
            },
            ticks: {
              color: "white",
            },
            grid: {
              display: true,
              drawBorder: true,
              drawOnChartArea: true,
              lineWidth: 1,
              color: "white",
            },
          },
          y: {
            stacked: true,
            title: {
              color: "white",
            },
            ticks: {
              color: "white",
              autoSkip: false,
            },
            grid: {
              display: true,
              drawOnChartArea: true,
              drawBorder: true,
              //   color: "white",
              color(ctx, options) {
                console.log(ctx);
                if (ctx.index === 0 || ctx.tick === undefined) {
                  return "white";
                }
                return "";
              },
            },
          },
        },
      }}
    />
  );
};
