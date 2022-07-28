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
  return games
    .filter((game) => {
      return game.season === gameFilter.season;
    })
    .sort((a, b) => CompareDates(a.date, b.date));
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
