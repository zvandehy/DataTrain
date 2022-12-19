import {
  Box,
  Grid,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  useTheme,
} from "@mui/material";
import moment from "moment";
import { useState } from "react";
import { PrimarySearchAppBar } from "../../ components/appbar/appbar.component";
import { FeaturedPropCard } from "../../ components/cards/featured-prop-card.component";
// import { FeaturedPropCard } from "../../ components/cards/featured-prop-card.component";
import { TotalPropsCard } from "../../ components/cards/total-props-card.component.";
import HistoricalCharts from "../../ components/charts/historical-charts/historical-charts-container.component";
import FullScreenDialog from "../../ components/layouts/full-screen-dialog.component";
import PlayerModal from "../../ components/player-modal/player-modal.component";
import PlayerRow from "../../ components/player-row/player-row.component";
// import PlayerRow from "../../ components/player-row/player-row.component";
import { useGetPropositions } from "../../hooks/useGetPropositions";
import { DEFAULT_MODEL } from "../../shared/constants";
import {
  ComparePropByPredictionDeviation,
  Proposition,
} from "../../shared/interfaces/graphql/proposition.interface";

// ==============================|| DASHBOARD ||============================== //

const DashboardPage = () => {
  const [startDate, setStartDate] = useState(moment().format("YYYY-MM-DD"));
  const [open, setOpen] = useState(false);
  const [selectedProp, setSelectedProp] = useState<Proposition>();

  const handleClickOpen = (prop: Proposition) => {
    setSelectedProp(prop);
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

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
      <FullScreenDialog
        open={open}
        handleClose={handleClose}
        title={`${selectedProp?.game.player.name} ${
          selectedProp?.game.home_or_away?.toUpperCase() === "HOME" ? "vs" : "@"
        } ${selectedProp?.game.opponent.name}`}
      >
        <PlayerModal
          playerID={selectedProp?.game.player.playerID ?? 0}
          triggeredProp={selectedProp}
          customModel={DEFAULT_MODEL}
          startDate={startDate}
        />
      </FullScreenDialog>
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
        <Grid container item xs={12} md={4} sx={{ m: "auto" }}>
          <Grid
            item
            xs={6}
            md={12}
            sx={{ p: 1, margin: "auto", "& > *": { mb: 1, mt: 1 } }}
          >
            <TotalPropsCard
              title={"All Props " + moment(startDate).format("MM/DD/YYYY")}
              propositions={propositions.filter((p) =>
                moment(p.game.date).isSame(startDate)
              )}
            />
          </Grid>
          <Grid
            item
            xs={6}
            md={12}
            sx={{ p: 1, margin: "auto", "& > *": { mb: 1, mt: 1 } }}
          >
            <TotalPropsCard
              title={"Top Props  " + moment(startDate).format("MM/DD/YYYY")}
              propositions={topProps.filter((p) =>
                moment(p.game.date).isSame(startDate)
              )}
            />
          </Grid>
          {/* <TotalPropsCard total={50} /> */}
        </Grid>
        <Grid container item xs={12} md={8}>
          {topProps
            .filter((p) => moment(p.game.date).isSame(startDate))
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
                  <FeaturedPropCard
                    onClick={handleClickOpen}
                    prop={prop}
                    rank={i + 1}
                  />
                </Grid>
              );
            })}
        </Grid>

        <HistoricalCharts initialDate={startDate} />
      </Grid>
      {/* row 3 */}
      <Grid
        item
        xs={12}
        md={7}
        lg={8}
        paddingX={2}
        style={{
          overflowX: "scroll",
        }}
      >
        <Table sx={{ backgroundColor: theme.palette.background.paper }}>
          <TableHead>
            <TableRow>
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
            </TableRow>
          </TableHead>
          <TableBody>
            {Object.entries(playerPropositions)
              .filter(
                (entry) =>
                  entry[1].filter((p) => moment(p.game.date).isSame(startDate))
                    .length > 0
              )
              .sort((entryA, entryB) =>
                ComparePropByPredictionDeviation(entryB[1][0], entryA[1][0])
              )
              .map((entry) => {
                const playerProps = entry[1];
                return (
                  <PlayerRow
                    key={entry[0]}
                    propositions={playerProps}
                    onClick={handleClickOpen}
                  />
                );
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
