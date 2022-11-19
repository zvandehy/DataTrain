import { Games } from "@material-ui/icons";
import {
  Box,
  Button,
  Card,
  Grid,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableHead,
  Typography,
  useTheme,
} from "@mui/material";
import moment from "moment";
import { useState } from "react";
import PrimarySearchAppBar from "../../ components/appbar/appbar.component";
import { FeaturedPropCard } from "../../ components/cards/featured-prop-card.component";
// import { FeaturedPropCard } from "../../ components/cards/featured-prop-card.component";
import { TotalPropsCard } from "../../ components/cards/total-props-card.component.";
import { ModelAccuracyByPctDiff } from "../../ components/charts/accuracy-by-percent-diff-chart.component";
import { ModelAccuracyByStatType } from "../../ components/charts/accuracy-by-type-chart.component";
import ModelAccuracyChart from "../../ components/charts/model-accuracy-chart-component";
// import PlayerRow from "../../ components/player-row/player-row.component";
import { useGetPropositions } from "../../hooks/useGetPropositions";
import { DEFAULT_MODEL } from "../../shared/constants";
import { Player } from "../../shared/interfaces/graphql/player.interface";
import {
  ComparePropByPredictionDeviation,
  Proposition,
} from "../../shared/interfaces/graphql/proposition.interface";

// ==============================|| DASHBOARD ||============================== //

const DashboardPage = () => {
  const [slot, setSlot] = useState("week");
  const [startDate, setStartDate] = useState("2022-11-18");

  const theme = useTheme();

  const {
    loading,
    error,
    data: propositions,
  } = useGetPropositions({
    startDate: startDate,
    endDate: startDate,
    customModel: DEFAULT_MODEL,
  });

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    return <div>Error: {error.message}</div>;
  }
  if (!propositions || propositions.length === 0) {
    return <div>No propositions found</div>;
  }

  let gameIDs: string[] = [];

  // each player's highest discrepancy proposition
  let mapTopProps: { [playerID: number]: Proposition } = {};
  propositions.forEach((prop) => {
    if (!gameIDs.includes(prop.game.gameID)) {
      gameIDs.push(prop.game.gameID);
    }
    if (mapTopProps[prop.game?.player?.playerID] === undefined) {
      mapTopProps[prop.game?.player?.playerID] = prop;
    } else {
      if (
        ComparePropByPredictionDeviation(
          prop,
          mapTopProps[prop.game?.player?.playerID]
        ) > 0
      ) {
        mapTopProps[prop.game?.player?.playerID] = prop;
      }
    }
  });

  const topProps: Proposition[] = Object.values(mapTopProps);

  return (
    <Box>
      <PrimarySearchAppBar />
      <Grid container p={1} spacing={1}>
        {/* row 1 */}
        <Grid item xs={12} sx={{ mb: -2.25 }}>
          <Typography
            variant="h5"
            sx={{ fontWeight: "bold" }}
            textTransform={"uppercase"}
            pl={2}
          >
            Dashboard
          </Typography>
        </Grid>
        {/* row 2 */}
        <Grid container xs={12} md={4} sx={{ m: "auto" }}>
          <Grid
            item
            xs={6}
            md={12}
            sx={{ p: 1, margin: "auto", "& > *": { mb: 1, mt: 1 } }}
          >
            <TotalPropsCard
              title={"All Props"}
              propositions={propositions}
              nGames={gameIDs.length}
              nPlayers={Object.keys(mapTopProps).length}
            />
          </Grid>
          <Grid
            item
            xs={6}
            md={12}
            sx={{ p: 1, margin: "auto", "& > *": { mb: 1, mt: 1 } }}
          >
            <TotalPropsCard
              title={"Top Props"}
              propositions={topProps}
              nGames={gameIDs.length}
              nPlayers={Object.keys(mapTopProps).length}
            />
          </Grid>
          {/* <TotalPropsCard total={50} /> */}
        </Grid>
        <Grid container item xs={12} md={8}>
          {topProps
            .sort((a, b) => {
              return ComparePropByPredictionDeviation(b, a);
            })
            .slice(0, 8)
            .map((prop, i) => {
              return (
                <Grid
                  item
                  xs={12}
                  sm={6}
                  sx={{ p: 1, margin: "auto" }}
                  key={prop.game.player.playerID + "-" + i}
                >
                  <FeaturedPropCard prop={prop} rank={i + 1} />
                </Grid>
              );
            })}
        </Grid>

        {/* row 2 */}
        <Grid item xs={12} md={7} lg={8} mt={2}>
          <Grid container alignItems="center" justifyContent="space-between">
            <Grid item>
              <Typography variant="h5">Model Accuracy</Typography>
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
              </Stack>
            </Grid>
          </Grid>
          <Card sx={{ mt: 1.5 }}>
            <Box sx={{ pt: 1, pr: 2 }}>
              <Box slot={slot}>
                <ModelAccuracyChart
                  endDate={moment(startDate)
                    .add(1, "days")
                    .format("YYYY-MM-DD")}
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
      </Grid>

      {/* row 3 */}
      <Grid item xs={12} md={7} lg={8} paddingX={2}>
        <Table sx={{ backgroundColor: theme.palette.background.paper }}>
          <TableHead>
            <TableCell>Expand</TableCell>
            <TableCell>Player</TableCell>
            <TableCell>Team</TableCell>
            <TableCell>Matchup</TableCell>
            <TableCell sx={{ borderLeft: "1px solid" }}>Sportsbook</TableCell>
            <TableCell>Stat</TableCell>
            <TableCell>Target</TableCell>
            <TableCell>Prediction</TableCell>
            <TableCell>% DIFF</TableCell>
            <TableCell>Predicted Min</TableCell>
            <TableCell>Predicted/Min</TableCell>
            <TableCell sx={{ borderLeft: "1px solid" }}>Outcome</TableCell>
            <TableCell>Min</TableCell>
            <TableCell>Actual/Min</TableCell>
          </TableHead>
          <TableBody>
            {/* {players.map((player) => {
              if (player.games === undefined || player.games.length === 0) {
                console.warn("no games for player", player.name);
                return null;
              }
              const game = player.games[0];
              if (
                game.prediction.propositions === undefined ||
                game.prediction.propositions.length === 0
              ) {
                console.warn(
                  "no propositions for player",
                  player,
                  player.games
                );
                return null;
              }
              return <PlayerRow player={player} game={game} />;
            })} */}
          </TableBody>
        </Table>
      </Grid>
      <Grid item xs={12} bgcolor={"gray"} mt={1}>
        <Typography textAlign={"center"} p={1}>
          FOOTER
        </Typography>
      </Grid>
    </Box>
  );
};

export default DashboardPage;
