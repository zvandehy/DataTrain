import { Grid, useTheme } from "@mui/material";
import React from "react";
import { Card, CardContent, CardHeader, Typography } from "@mui/material";
import { OverUnderPieChart } from "../charts/over-under-pie-chart.component";
import { Proposition } from "../../shared/interfaces/graphql/proposition.interface";

export interface TotalPropsCardProps {
  propositions: Proposition[];
  title: string;
  nGames: number;
  nPlayers: number;
}

export const TotalPropsCard: React.FC<TotalPropsCardProps> = ({
  propositions,
  title,
  nGames,
  nPlayers,
}: TotalPropsCardProps) => {
  const theme = useTheme();

  return (
    <Card
      elevation={0}
      sx={{
        border: "1px solid",
        borderRadius: 2,
        borderColor: theme.palette.divider,
        boxShadow: "inherit",
        ":hover": {
          boxShadow: "inherit",
        },
        "& pre": {
          m: 0,
          p: "16px !important",
          fontFamily: theme.typography.fontFamily,
          fontSize: "0.75rem",
        },
      }}
    >
      <CardHeader
        sx={{
          "&": {
            pb: 0,
          },
        }}
        title={<Typography variant="h6">{title}</Typography>}
        // subheader={}
        // action={<OverUnderPieChart />}
      />
      <CardContent sx={{}}>
        <Grid container>
          <Grid item xs={5} sx={{ margin: "auto" }}>
            <Typography
              sx={{
                typography: { xs: "subtitle2", sm: "subtitle1" },
              }}
              noWrap={true}
              variant="subtitle2"
            >
              {nGames} Games
            </Typography>
            <Typography
              sx={{
                typography: { xs: "subtitle2", sm: "subtitle1" },
              }}
              noWrap={true}
              variant="subtitle2"
            >
              {nPlayers} Players
            </Typography>
            <Typography
              sx={{
                typography: { xs: "subtitle2", sm: "subtitle1" },
              }}
              noWrap={true}
              variant="subtitle2"
            >
              {propositions.length} Props
            </Typography>
          </Grid>
          <Grid xs={4} sx={{ margin: "auto" }}>
            <OverUnderPieChart propositions={propositions} />
          </Grid>
        </Grid>
      </CardContent>

      {/* card footer - clipboard & highlighter  */}
      {/* {codeHighlight && (
          <>
            <Divider sx={{ borderStyle: "dashed" }} />
            <Highlighter codeHighlight={codeHighlight} main>
              {children}
            </Highlighter>
          </>
        )} */}
    </Card>
  );
};
