import { Grid, Typography, Stack, Button, Card, Box } from "@mui/material";
import moment from "moment";
import { useState } from "react";
import { useGetPropositions } from "../../../hooks/useGetPropositions";
import { DEFAULT_MODEL } from "../../../shared/constants";
import { ModelAccuracyByPctDiff } from "../accuracy-by-percent-diff-chart.component";
import { ModelAccuracyByStatType } from "../accuracy-by-type-chart.component";
import ModelAccuracyChart from "../model-accuracy-chart-component";

interface HistoricalChartsProps {
  endDate: string;
}

const HistoricalCharts: React.FC<HistoricalChartsProps> = ({
  endDate,
}: HistoricalChartsProps) => {
  const [slot, setSlot] = useState<"month" | "week" | "day">("week");
  const startDate =
    slot === "day"
      ? endDate
      : slot === "week"
      ? moment(endDate).subtract(7, "days").format("YYYY-MM-DD")
      : moment(endDate).subtract(1, "month").format("YYYY-MM-DD");

  const {
    loading,
    error,
    data: propositions,
  } = useGetPropositions({
    startDate: startDate,
    endDate: endDate,
    customModel: DEFAULT_MODEL,
  });

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    return <div>Error: {error.message}</div>;
  }
  return (
    <>
      <Grid item xs={12} md={7} lg={8} mt={2}>
        <Grid container alignItems="center" justifyContent="space-between">
          <Grid item>
            <Typography variant="h5">Overall Historical Accuracy</Typography>
          </Grid>
          <Grid item>
            <Stack direction="row" alignItems="center" spacing={0}>
              <Button
                size="small"
                onClick={() => setSlot("month")}
                color={slot === "month" ? "primary" : "secondary"}
                variant={slot === "month" ? "outlined" : "text"}
              >
                Month
              </Button>
              <Button
                size="small"
                onClick={() => setSlot("week")}
                color={slot === "week" ? "primary" : "secondary"}
                variant={slot === "week" ? "outlined" : "text"}
              >
                Week
              </Button>
              <Button
                size="small"
                onClick={() => setSlot("day")}
                color={slot === "day" ? "primary" : "secondary"}
                variant={slot === "day" ? "outlined" : "text"}
              >
                Day
              </Button>
            </Stack>
          </Grid>
        </Grid>
        <Card sx={{ mt: 1.5 }}>
          <Box sx={{ pt: 1, pr: 2 }}>
            <Box slot={slot}>
              <ModelAccuracyChart
                propositions={propositions}
                startDate={startDate}
                endDate={endDate}
              />
            </Box>
          </Box>
        </Card>
      </Grid>
      <Grid item xs={12} md={5} lg={4} mt={2}>
        <Grid container alignItems="center" justifyContent="space-between">
          <Grid item>
            <Typography variant="h5">Stat Accuracy</Typography>
          </Grid>
        </Grid>
        <Card sx={{ mt: 2 }}>
          <Box sx={{ p: 3, pb: 0, minHeight: "300px" }}>
            <ModelAccuracyByStatType propositions={propositions} />
          </Box>
          <Box />
        </Card>
        <Typography variant="h5" mt={1}>
          Significance Accuracy
        </Typography>
        <Card sx={{ mt: 2 }}>
          <Box sx={{ p: 3, pb: 0, minHeight: "300px" }}>
            <ModelAccuracyByPctDiff
              propositions={propositions}
              stepSize={10}
              steps={5}
            />
          </Box>
          <Box />
        </Card>
      </Grid>
    </>
  );
};

export default HistoricalCharts;
