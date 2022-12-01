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
import { PrimarySearchAppBar } from "../../ components/appbar/appbar.component";
import { FeaturedPropCard } from "../../ components/cards/featured-prop-card.component";
// import { FeaturedPropCard } from "../../ components/cards/featured-prop-card.component";
import { TotalPropsCard } from "../../ components/cards/total-props-card.component.";
import { ModelAccuracyByPctDiff } from "../../ components/charts/accuracy-by-percent-diff-chart.component";
import { ModelAccuracyByStatType } from "../../ components/charts/accuracy-by-type-chart.component";
import HistoricalCharts from "../../ components/charts/historical-charts/historical-charts-container.component";
import ModelAccuracyChart from "../../ components/charts/model-accuracy-chart-component";
import PlayerRow from "../../ components/player-row/player-row.component";
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
  const [slot, setSlot] = useState<"month" | "week" | "day">("day");
  const [startDate, setStartDate] = useState(moment().format("YYYY-MM-DD"));

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
  // if (!propositions || propositions.length === 0) {
  //   return <div>No propositions found</div>;
  // }

  let gameIDs: string[] = [];

  // each player's highest discrepancy proposition
  let playerPropositions: { [playerGameID: string]: Proposition[] } = {};
  propositions.forEach((prop) => {
    if (!gameIDs.includes(prop.game.gameID)) {
      gameIDs.push(prop.game.gameID);
    }
    if (
      !playerPropositions[prop.game.gameID + "" + prop.game.player.playerID]
    ) {
      playerPropositions[prop.game.gameID + "" + prop.game.player.playerID] = [
        prop,
      ];
    } else {
      playerPropositions[
        prop.game.gameID + "" + prop.game.player.playerID
      ].push(prop);
      playerPropositions[
        prop.game.gameID + "" + prop.game.player.playerID
      ].sort((a, b) => ComparePropByPredictionDeviation(b, a));
    }
  });

  const topProps: Proposition[] = Object.entries(playerPropositions).map(
    (entry) => entry[1][0]
  );

  return (
    <Box>
      <PrimarySearchAppBar date={startDate} onDateSelect={setStartDate} />
      <Grid container p={1} spacing={1}>
        {/* row 1 */}
        <Grid item xs={12} sx={{ mb: -2.25 }}>
          <Typography
            variant="h5"
            sx={{ fontWeight: "bold" }}
            textTransform={"uppercase"}
            pl={2}
          >
            Dashboard {moment(startDate).format("MM/DD/YYYY")}
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
              nPlayers={Object.keys(playerPropositions).length}
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
              nPlayers={Object.keys(playerPropositions).length}
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

        <HistoricalCharts endDate={startDate} />
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
            <TableCell>Significance</TableCell>
            <TableCell>Std Dev</TableCell>
            <TableCell>N</TableCell>
            <TableCell sx={{ borderLeft: "1px solid" }}>Outcome</TableCell>
            {/* <TableCell>Min</TableCell>
            <TableCell>Actual/Min</TableCell> */}
          </TableHead>
          <TableBody>
            {Object.entries(playerPropositions)
              .sort((entryA, entryB) =>
                ComparePropByPredictionDeviation(entryB[1][0], entryA[1][0])
              )
              .map((entry) => {
                const playerProps = entry[1];
                return <PlayerRow propositions={playerProps} />;
              })}
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
