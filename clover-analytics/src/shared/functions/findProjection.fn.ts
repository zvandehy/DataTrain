import moment from "moment";
import { Game } from "../interfaces/graphql/game.interface";
import { Player } from "../interfaces/graphql/player.interface";
import { Projection } from "../interfaces/graphql/projection.interface";

export function FindProjectionByDate(
  date: Date,
  projections: Projection[],
  player: Player
): Projection {
  let p =
    projections.find((projection) =>
      moment(projection.date).isSame(moment(date))
    ) ?? CreateMissingProjection(date, player);
  let dates = projections.map((projection) => projection.date);
  dates.push(...player.games.map((game) => game.date));
  return p;
}

export function CreateMissingProjection(
  date: Date,
  player: Player
): Projection {
  let game =
    player.games.find((game) => moment(game.date) === moment(date)) ??
    player.games[player.games.length - 1];
  let projection: Projection = {
    date: game.date,
    opponent: game.opponent,
    player: player,
    propositions: [],
    result: game,
    startTime: moment(game.date).format("YYYY-MM-DD hh:ss"),
  };
  return projection;
}
