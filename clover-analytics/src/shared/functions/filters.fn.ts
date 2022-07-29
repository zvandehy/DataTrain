import { GameFilter } from "../interfaces/graphql/filters.interface";
import { Game } from "../interfaces/graphql/game.interface";
import { CompareDates } from "./dates.fn";
import {
  Projection,
  Proposition,
} from "../interfaces/graphql/projection.interface";
import {
  ListFilterOptions,
  ListSortOptions,
} from "../interfaces/listFilter.interface";

export function FilterGames(games: Game[], gameFilter: GameFilter): Game[] {
  const filteredGames = games
    .filter((game) => {
      //match the game if the season is the same
      let seasonMatch = true;
      if (gameFilter.season) {
        seasonMatch = gameFilter.season === game.season;
        if (!seasonMatch) {
          return seasonMatch;
        }
      }
      //match the game if the date is after the start date filter
      let gameIsAfterStartDate = true;
      if (gameFilter.startDate) {
        gameIsAfterStartDate = GameIsAfterDate(
          game,
          gameFilter.startDate,
          true
        );
      }
      if (!gameIsAfterStartDate) {
        return gameIsAfterStartDate;
      }
      //match the game if the date is before the end date filter
      let gameIsBeforeEndDate = true;
      if (gameFilter.endDate) {
        gameIsBeforeEndDate = GameIsBeforeDate(game, gameFilter.endDate, false);
      }
      if (!gameIsBeforeEndDate) {
        return gameIsBeforeEndDate;
      }
      return game.season === gameFilter.season;
    })
    .sort((a, b) => CompareDates(a.date, b.date));
  return filteredGames;
}

function GameIsAfterDate(
  game: Game,
  date: string,
  inclusive: boolean
): boolean {
  const cmp = CompareDates(game.date, date);
  if (cmp < 0) {
    //game.date is before date
    return false;
  } else if (cmp === 0) {
    //game.date is equal to date
    return inclusive;
  }
  return true;
}

function GameIsBeforeDate(
  game: Game,
  date: string,
  inclusive: boolean
): boolean {
  const cmp = CompareDates(game.date, date);
  if (cmp < 0) {
    //game.date is before date
    return true;
  } else if (cmp === 0) {
    //game.date is equal to date
    return inclusive;
  }
  return false;
}

export function Match(
  projection: Projection,
  filterOptions: ListFilterOptions
): boolean {
  const { lookup, statType } = filterOptions;
  if (statType?.label) {
    return (
      MatchLookup(projection, lookup) &&
      projection.propositions.some(
        (proposition) => proposition.statType === statType
      )
    );
  }
  return MatchLookup(projection, lookup);
}

function MatchLookup(projection: Projection, lookup: string): boolean {
  return (
    projection.player.name.toLowerCase().search(lookup.toLowerCase()) !== -1 ||
    projection.player.currentTeam.abbreviation
      .toLowerCase()
      .search(lookup.toLowerCase()) !== -1 ||
    projection.opponent.abbreviation
      .toLowerCase()
      .search(lookup.toLowerCase()) !== -1 ||
    projection.opponent.name.toLowerCase().search(lookup.toLowerCase()) !==
      -1 ||
    projection.player.currentTeam.name
      .toLowerCase()
      .search(lookup.toLowerCase()) !== -1
  );
}

export function SortProjections(
  projections: Projection[],
  sortOptions: ListSortOptions
): Projection[] {
  const { sortOrder, statType, sortBy } = sortOptions;
  return projections.sort((a, b) => {
    if (sortBy === "confidence") {
      let aProps = a.propositions;
      let bProps = b.propositions;
      // if statType is set, then sort the projections by the statType
      if (statType !== undefined) {
        aProps = a.propositions.filter((prop) => prop.statType === statType);
        bProps = b.propositions.filter((prop) => prop.statType === statType);
      }
      const aMax = Math.max(
        ...aProps.map((p: Proposition) => p.customPrediction.confidence)
      );
      const bMax = Math.max(
        ...bProps.map((p: Proposition) => p.customPrediction.confidence)
      );
      return bMax - aMax;
    }
    return a.player.currentTeam.name.localeCompare(b.player.currentTeam.name);
  });
}
