import { Box } from "@mui/system";
import { ChartDataset } from "chart.js";
import { Chart } from "react-chartjs-2";
import { Proposition } from "../../shared/interfaces/graphql/proposition.interface";
import { COLORS } from "../../shared/styles/constants";

interface PredictionContributionBarProps {
  proposition?: Proposition;
}

const PredictionContributionBar: React.FC<PredictionContributionBarProps> = ({
  proposition,
}: PredictionContributionBarProps) => {
  const contributions: ChartDataset<"bar", number[]>[] =
    proposition?.prediction?.breakdowns?.map((breakdown) => {
      return {
        label: breakdown.name,
        type: "bar",
        data: [breakdown.contribution, 0],
        backgroundColor: (context) => {
          let color =
            breakdown.derivedAverage > proposition?.target
              ? COLORS.HIGHER
              : COLORS.LOWER;
          if (breakdown.name.includes(" vs Opponent")) {
            color = breakdown.pctChange > 0 ? COLORS.HIGHER : COLORS.LOWER;
          }
          // color opacity based on the weight
          const intensity = +(
            50 +
            (+breakdown.weight.toFixed() * (255 - 50)) / 100
          ).toFixed();
          return color + intensity.toString(16);
        },
        borderColor: "white",
        borderWidth: 1,
      };
    }) ?? [];

  const targetDataset: ChartDataset<"bar", number[]> = {
    label: "Target",
    type: "bar",
    data: [0, proposition?.target ?? 0],
    backgroundColor: () => {
      return (proposition?.prediction?.estimation ?? 0) >
        (proposition?.target ?? 0)
        ? COLORS.HIGHER
        : COLORS.LOWER;
    },
    borderColor: "white",
    borderWidth: 1,
  };

  return (
    <Box maxHeight={"150px"}>
      <Chart
        type="bar"
        data={{
          labels: ["Prediction Contributions", "Target"],
          datasets: contributions.concat(targetDataset),
        }}
        options={{
          indexAxis: "y",
          responsive: true,
          maintainAspectRatio: false,
          plugins: {
            tooltip: {
              mode: "point",
              callbacks: {
                // title: (context) => {
                //   const dataset = context[0].dataset;
                //   const index = context[0].dataIndex;
                //   const value = dataset.data[index];
                //   const breakdown = proposition?.prediction?.breakdowns?.find(
                //     (x) => x.name === dataset.label
                //   );
                //   return `${dataset.label} (${
                //     breakdown?.derivedAverage
                //   } * ${breakdown?.weight.toFixed(0)}% = ${value.toFixed(2)})`;
                // },
                // label: (context) => {
                //   const dataset = context.dataset;
                //   const index = context.dataIndex;
                //   const value = dataset.data[index];
                //   const percentOfTotal =
                //     (value / (proposition?.prediction?.estimation ?? 1)) * 100;
                //   return `${percentOfTotal.toFixed(
                //     2
                //   )}% of total prediction (${value.toFixed(2)})`;
                // },
              },
            },
            title: {
              display: true,
              text: "Prediction Model Contributions",
              color: "white",
            },
            legend: {
              labels: { color: "white" },
              display: false,
            },
          },
          color: "white",
          borderColor: "white",
          backgroundColor: "white",
          scales: {
            x: {
              grid: {
                tickColor: "white",
              },
              ticks: {
                color: "white",
              },
              stacked: true,
            },
            y: {
              stacked: true,
              ticks: {
                color: "white",
              },
            },
          },
        }}
      />
    </Box>
  );
};

export default PredictionContributionBar;
