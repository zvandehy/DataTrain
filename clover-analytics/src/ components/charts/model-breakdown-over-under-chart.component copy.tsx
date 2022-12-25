import { Box } from "@mui/system";
import { Chart } from "react-chartjs-2";
import { Proposition } from "../../shared/interfaces/graphql/proposition.interface";
import { COLORS } from "../../shared/styles/constants";

interface ModelBreakdownOverUnderChartProps {
  proposition?: Proposition;
}

const ModelBreakdownOverUnderChart: React.FC<
  ModelBreakdownOverUnderChartProps
> = ({ proposition }: ModelBreakdownOverUnderChartProps) => {
  const labels = proposition?.prediction?.breakdowns?.map(
    (breakdown) => breakdown.name
  );
  const overs = proposition?.prediction?.breakdowns?.map(
    (breakdown) => +breakdown.over.toFixed(2)
  );
  const unders = proposition?.prediction?.breakdowns?.map(
    (breakdown) => +breakdown.under.toFixed(2)
  );

  return (
    <Box>
      <Chart
        type="bar"
        title={"PLAYER PROP MODEL BREAKDOWN"}
        data={{
          labels,
          datasets: [
            {
              label: "# OVER",
              data: overs,
              type: "bar",
              backgroundColor: (context) => {
                const color = COLORS.HIGHER;
                const intensity = 200;
                return color + intensity.toString(16);
              },
            },
            {
              label: "# UNDER",
              data: unders,
              type: "bar",
              backgroundColor: (context) => {
                const color = COLORS.LOWER;
                const intensity = 200;
                return color + intensity.toString(16);
              },
            },
          ],
        }}
        options={{
          indexAxis: "y",
          responsive: true,
          plugins: {
            tooltip: {
              mode: "index",
            },
            title: {
              display: true,
              text: "Prediction Breakdown Over/Under",
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
            },
            y: {
              labels: labels?.map(
                (label) =>
                  label +
                  " (" +
                  (
                    (proposition?.prediction?.breakdowns?.find(
                      (x) => x.name === label
                    )?.overPct ?? 0) * 100
                  ).toFixed() +
                  "%)"
              ),
              ticks: {
                color: "white",
              },
            },
          },
        }}
        //   options={{
        //     responsive: true,
        //     maintainAspectRatio: false,
        //     interaction: { mode: "index", intersect: false },
        //     plugins: {
        //       tooltip: {
        //         mode: "index",
        //       },
        //       title: {
        //         display: true,
        //         text: "MODEL ACCURACY BY DIFFERENCE FROM TARGET",
        //         color: "white",
        //       },
        //       legend: {
        //         labels: { color: "white" },
        //         display: true,
        //       },
        //     },
        //     scales: {
        //       x: {
        //         stacked: true,
        //         ticks: {
        //           color: "white",
        //           align: "end",
        //         },
        //         display: false,
        //       },
        //       x1: {
        //         stacked: true,
        //         ticks: {
        //           color: "white",
        //           maxTicksLimit: breakdowns.length - 1,
        //           autoSkip: false,
        //         },
        //         title: {
        //           display: true,
        //           text: `Distance from Proposition Target${
        //             !absolute ? " (%)" : ""
        //           }`,
        //           color: "white",
        //         },
        //         labels: [
        //           ...breakdowns.map((x) => `${x.tick}${!absolute ? "%" : ""}`),
        //           "∞",
        //         ],
        //       },
        //       y: {
        //         stacked: true,
        //         title: {
        //           display: true,
        //           text: "# Props",
        //           color: "white",
        //         },
        //         ticks: {
        //           color: "white",
        //           align: "start",
        //         },
        //         grid: {
        //           display: false,
        //         },
        //       },
        //       y2: {
        //         position: "right",
        //         min: 0,
        //         max: 100,
        //         ticks: {
        //           color: "white",
        //         },
        //         title: {
        //           display: true,
        //           text: `ACCURACY ${!absolute ? "(%)" : ""}`,
        //           color: "white",
        //         },
        //       },
        //     },
        //   }}
      />
    </Box>
  );
};

export default ModelBreakdownOverUnderChart;
