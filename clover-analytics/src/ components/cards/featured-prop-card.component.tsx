import { Avatar, Badge, CardActionArea, Grid, useTheme } from "@mui/material";
import React from "react";
import { Card, CardHeader, Typography } from "@mui/material";
import {
  OverUnderIcon,
  OverUnderTypography,
} from "../icons/overUnderIcon.component";
import { HitMissIcon, HitMissTypography } from "../icons/hit-miss.icon";
import { ShortenName } from "../../shared/functions/name.fn";
import { COLORS } from "../../shared/styles/constants";
import { Player } from "../../shared/interfaces/graphql/player.interface";
import { GetStatAbbreviation } from "../../shared/interfaces/stat.interface";
import {
  GetPropPredictionDeviation,
  Proposition,
} from "../../shared/interfaces/graphql/proposition.interface";

export interface FeaturedPropCardProps {
  prop: Proposition;
  rank: number;
}

export const FeaturedPropCard: React.FC<FeaturedPropCardProps> = ({
  prop,
  rank,
}: FeaturedPropCardProps) => {
  const theme = useTheme();
  const player = prop.game.player;
  const playername = ShortenName(player.name, 15);
  const game = prop.game;
  return (
    <Card
      elevation={0}
      sx={{
        border: "3px solid",
        borderRadius: 2,
        borderColor:
          prop?.prediction?.wagerOutcome === "PENDING"
            ? theme.palette.divider
            : prop?.prediction?.wagerOutcome === "MISS"
            ? COLORS.LOWER
            : prop?.prediction?.wagerOutcome === "HIT"
            ? COLORS.HIGHER
            : COLORS.PUSH,
        boxShadow: "inherit",
        ":hover": {
          boxShadow: "inherit",
        },
        "& pre": {
          m: 0,
          // p: "16px !important",
          fontFamily: theme.typography.fontFamily,
          fontSize: "0.75rem",
        },
      }}
    >
      <CardActionArea>
        <CardHeader
          sx={{
            p: 2.5,
            "& .MuiCardHeader-action": { m: "0px auto", alignSelf: "center" },
          }}
          avatar={
            <Badge
              badgeContent={rank}
              color="secondary"
              invisible={false}
              anchorOrigin={{
                vertical: "top",
                horizontal: "left",
              }}
            >
              <Avatar
                sx={{
                  borderRadius: 1,
                  width: 58,
                  height: 58,
                  bgcolor: COLORS.AVATAR,
                }}
                src={player.image}
                alt={"player" + { player }}
              >
                {rank}
              </Avatar>
            </Badge>
          }
          title={
            <Grid
              container
              direction="column"
              columnSpacing={1}
              gap={"inherit"}
              wrap="nowrap"
              alignContent={"start"}
            >
              <Grid item>
                <Typography
                  variant={"subtitle1"}
                  alignSelf={"start"}
                  sx={{ fontWeight: 500, lineHeight: 0, mb: 1 }}
                >
                  {playername}
                </Typography>
              </Grid>
              <Grid item>
                <Typography
                  variant={"overline"}
                  alignSelf={"start"}
                  sx={{ fontWeight: 400, lineHeight: 0, ml: 0.5 }}
                >
                  {game.home_or_away === "home" ? "vs" : "@"}{" "}
                  {game.opponent.abbreviation}
                </Typography>
              </Grid>
              <Grid
                item
                // textAlign={"start"}
                alignSelf={"start"}
                justifyContent={"center"}
              >
                <Grid
                  container
                  // alignContent={"center"}
                  alignItems={"center"}
                  columnGap={0.5}
                >
                  <OverUnderIcon
                    size={16}
                    overUnder={prop?.prediction?.wager}
                  />
                  <OverUnderTypography
                    size={16}
                    overUnder={prop?.prediction?.wager}
                  />
                  <Typography
                    fontSize={16}
                    fontWeight={600}
                    textTransform={"uppercase"}
                  >
                    {prop?.target} {GetStatAbbreviation(prop?.type ?? "")}
                  </Typography>
                </Grid>
              </Grid>
            </Grid>
          }
          action={
            <Grid
              container
              direction={"column"}
              textAlign={"center"}
              width={"100%"}
              margin={0}
            >
              <Typography
                textAlign="center"
                justifyContent={"center"}
                fontSize={15}
                fontWeight={500}
                variant={"body2"}
              >
                Model Prediction
              </Typography>
              <Typography
                textAlign="center"
                justifyContent={"center"}
                fontSize={15}
                fontWeight={500}
                textTransform={"uppercase"}
              >
                {prop?.prediction?.estimation.toFixed(2)}{" "}
                {GetStatAbbreviation(prop?.type ?? "")} (
                {prop?.prediction?.estimation > prop?.target ? "+" : "-"}
                {GetPropPredictionDeviation(prop)}%)
              </Typography>
              <Grid item alignSelf={"center"}>
                {prop?.prediction?.wagerOutcome !== "PENDING" ? (
                  <Grid container alignContent={"center"}>
                    <HitMissTypography
                      outcome={prop?.prediction?.wagerOutcome}
                      result={prop?.actualResult}
                    />
                    <HitMissIcon outcome={prop?.prediction?.wagerOutcome} />
                  </Grid>
                ) : (
                  <Grid container alignSelf={"center"}>
                    {/* TODO: use game startTime */}
                    {game.date}
                  </Grid>
                )}
              </Grid>
            </Grid>
          }
        />
      </CardActionArea>
      {/* <Divider sx={{ borderStyle: "dashed" }} /> */}
      {/* <CardContent sx={{ height: "50px", m: 0 }}>{"children"}</CardContent> */}
    </Card>
  );
};
