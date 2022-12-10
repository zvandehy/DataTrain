import {
  Grid,
  Typography,
  Stack,
  Button,
  Card,
  Box,
  Select,
  InputLabel,
  MenuItem,
} from "@mui/material";
import moment from "moment";
import { useRef, useState } from "react";
import { useGetPropositions } from "../../../hooks/useGetPropositions";
import { DEFAULT_MODEL } from "../../../shared/constants";
import { TotalPropsCard } from "../../cards/total-props-card.component.";
import { ModelAccuracyByPctDiff } from "../accuracy-by-percent-diff-chart.component";
import { ModelAccuracyByStatType } from "../accuracy-by-type-chart.component";
import ModelAccuracyChart from "../model-accuracy-chart-component";

interface HistoricalChartsProps {
  initialDate: string;
}

const HistoricalCharts: React.FC<HistoricalChartsProps> = ({
  initialDate,
}: HistoricalChartsProps) => {
  const dateRange = useRef({ startDate: initialDate, endDate: initialDate });
  const [loadingMore, setLoadingMore] = useState(false);
  const [stat, setStat] = useState<string>("any");
  const {
    loading,
    error,
    data: propositions,
    fetchMore,
  } = useGetPropositions({
    startDate: initialDate,
    endDate: initialDate,
    customModel: DEFAULT_MODEL,
  });

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    if (!propositions || !propositions.length) {
      return <div>Error loading data</div>;
    }
  }

  // get unique stat types
  const statTypes = propositions.reduce((acc: string[], prop) => {
    if (!acc.includes(prop.type)) {
      acc.push(prop.type);
    }
    return acc;
  }, []);

  const statOptions = [{ value: "any", label: "Any" }].concat(
    ...statTypes.map((type) => ({
      value: type,
      label: type.toUpperCase().replaceAll("_", "+"),
    }))
  );

  const filteredProps = propositions.filter((prop) => {
    return prop.type === stat || stat === "any";
  });

  return (
    <>
      <Grid item xs={12} md={7} lg={8} mt={2}>
        <Grid
          container
          alignItems="center"
          justifyContent="space-between"
          pb={1}
        >
          <Grid item>
            <Typography variant="h5">
              Model Accuracy Between{" "}
              {moment(dateRange.current.startDate).format("MMM DD")} -{" "}
              {moment(initialDate).format("MMM DD")}{" "}
            </Typography>
          </Grid>
          <Grid item>
            <Stack direction="row" alignItems="center" spacing={1}>
              <Select
                variant="outlined"
                size="small"
                labelId="select-stat"
                id="stat-selector"
                value={stat}
                onChange={(event) => setStat(event.target.value as string)}
              >
                {statOptions.map((option) => (
                  <MenuItem key={option.value} value={option.value}>
                    {option.label}
                  </MenuItem>
                ))}
              </Select>
              <Button
                size="small"
                variant={"contained"}
                onClick={() => {
                  setLoadingMore(true);
                  const newStartDate = moment(dateRange.current.startDate)
                    .subtract(1, "day")
                    .format("YYYY-MM-DD");
                  const newEndDate = moment(dateRange.current.endDate)
                    .subtract(1, "day")
                    .format("YYYY-MM-DD");
                  fetchMore({
                    variables: {
                      startDate: newStartDate,
                      endDate: newEndDate,
                    },
                    updateQuery(prev, { fetchMoreResult }) {
                      if (!fetchMoreResult || !fetchMoreResult?.propositions)
                        return prev;
                      return {
                        ...prev,
                        propositions: [
                          ...prev?.propositions,
                          ...fetchMoreResult.propositions,
                        ],
                      };
                    },
                  })
                    .catch((err) => {
                      console.log("Error Fetching More", err);
                    })
                    .finally(() => {
                      setLoadingMore(false);
                      dateRange.current.startDate = newStartDate;
                      dateRange.current.endDate = newEndDate;
                    });
                }}
              >
                Load +1 Day
              </Button>
            </Stack>
          </Grid>
        </Grid>
        <Card>
          <Box className={loadingMore ? "loading-data" : ""}>
            <ModelAccuracyChart propositions={filteredProps} />
          </Box>
        </Card>
      </Grid>
      <Grid item xs={12} md={5} lg={4} mt={"auto"}>
        <Card>
          <Box
            sx={{ p: 3, pb: 0, minHeight: "400px" }}
            className={loadingMore ? "loading-data" : ""}
          >
            <ModelAccuracyByStatType propositions={filteredProps} />
          </Box>
        </Card>
      </Grid>
      <Grid item container xs={12}>
        <Grid item xs={12} md={5} lg={4} mt={2} pl={1} pr={1}>
          <Card>
            <Box
              sx={{ p: 3, pb: 0, minHeight: "300px" }}
              className={loadingMore ? "loading-data" : ""}
            >
              <TotalPropsCard
                title={"All Props"}
                propositions={filteredProps}
              />
            </Box>
          </Card>
        </Grid>
        <Grid item xs={12} md={5} lg={4} mt={2} pl={1} pr={1}>
          <Card>
            <Box
              sx={{ p: 3, pb: 0, minHeight: "300px" }}
              className={loadingMore ? "loading-data" : ""}
            ></Box>
          </Card>
        </Grid>
        <Grid item xs={12} md={5} lg={4} mt={2} pl={1} pr={1}>
          <Card>
            <Box
              sx={{ p: 3, pb: 0, minHeight: "300px" }}
              className={loadingMore ? "loading-data" : ""}
            >
              <ModelAccuracyByPctDiff
                propositions={filteredProps}
                stepSize={10}
                steps={9}
              />
            </Box>
            <Box />
          </Card>
        </Grid>
      </Grid>
    </>
  );
};

export default HistoricalCharts;
