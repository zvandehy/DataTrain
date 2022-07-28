import { Game } from "./graphql/game.interface";
import { Proposition } from "./graphql/projection.interface";

interface IStat {
  display: string;
  abbreviation: string;
  label: string;
}

export class Stat {
  display: string;
  abbreviation: string;
  label: string;
  score: (game: Game) => number;
  average: (games: Game[]) => number;
  median: (games: Game[]) => number;
  constructor(
    stat: IStat,
    score?: (game: Game) => number,
    average?: (games: Game[]) => number,
    median?: (games: Game[]) => number
  ) {
    this.display = stat.display;
    this.abbreviation = stat.abbreviation;
    this.label = stat.label;
    this.score =
      score || ((game: Game) => defaultScore(game, this.label)(game));
    this.average =
      average || ((games: Game[]) => defaultAverage(games, this.score));
    this.median =
      median || ((games: Game[]) => defaultMedian(games, this.score));
  }
}

function defaultScore(game: Game, stat: string): (game: Game) => number {
  return (game: Game) => +game[stat as keyof Game];
}

function defaultAverage(
  games: Game[],
  scoreFn: (game: Game) => number
): number {
  return +(
    games.reduce((acc, game) => acc + scoreFn(game), 0) / games.length
  ).toFixed(2);
}

function defaultMedian(games: Game[], scoreFn: (game: Game) => number): number {
  if (games.length === 0) {
    return 0;
  }
  return +games
    .map((game) => scoreFn(game))
    .sort((a, b) => a - b)
    [Math.floor(games.length / 2)].toFixed(2);
}

export const Points: Stat = new Stat({
  display: "Points",
  abbreviation: "PTS",
  label: "points",
});

export const Rebounds: Stat = new Stat({
  display: "Rebounds",
  abbreviation: "REB",
  label: "rebounds",
});
export const Assists: Stat = new Stat({
  display: "Assists",
  abbreviation: "AST",
  label: "assists",
});
export const Steals: Stat = new Stat({
  display: "Steals",
  abbreviation: "STL",
  label: "steals",
});
export const Blocks: Stat = new Stat({
  display: "Blocks",
  abbreviation: "BLK",
  label: "blocks",
});
export const Turnovers: Stat = new Stat({
  display: "Turnovers",
  abbreviation: "TOV",
  label: "turnovers",
});
export const Fantasy: Stat = new Stat(
  {
    display: "Fantasy",
    abbreviation: "FAN",
    label: "fantasy_score",
  },
  (game: Game) =>
    +(
      game.points +
      game.rebounds * 1.2 +
      game.assists * 1.5 +
      (game.steals + game.blocks) * 3 -
      game.turnovers
    ).toFixed(2)
);
export const ThreeFGM: Stat = new Stat({
  display: "3-Pointers Made",
  abbreviation: "3PM",
  label: "three_pointers_made",
});
export const ThreeFGA: Stat = new Stat({
  display: "3-Pointers Attempted",
  abbreviation: "3PA",
  label: "three_pointers_attempted",
});
export const FreeThrowsMade: Stat = new Stat({
  display: "Free Throws Made",
  abbreviation: "FTM",
  label: "free_throws_made",
});
export const FreeThrowsAttempted: Stat = new Stat({
  display: "Free Throws Attempted",
  abbreviation: "FTA",
  label: "free_throws_attempted",
});
export const OffensiveRebounds: Stat = new Stat({
  display: "Offensive Rebounds",
  abbreviation: "OREB",
  label: "offensive_rebounds",
});
export const DefensiveRebounds: Stat = new Stat({
  display: "Defensive Rebounds",
  abbreviation: "DREB",
  label: "defensive_rebounds",
});
export const BlocksSteals: Stat = new Stat(
  {
    display: "Blocks+Steals",
    abbreviation: "BLK+STL",
    label: "blks+stls",
  },
  (game: Game) => game.blocks + game.steals
);
export const PointsRebounds: Stat = new Stat(
  {
    display: "Points+Rebounds",
    abbreviation: "PTS+REB",
    label: "pts+rebs",
  },
  (game: Game) => game.points + game.rebounds
);
export const PointsReboundsAssists: Stat = new Stat(
  {
    display: "Pts+Rebs+Asts",
    abbreviation: "PRA",
    label: "pts+rebs+asts",
  },
  (game: Game) => game.points + game.rebounds + game.assists
);
export const PointsAssists: Stat = new Stat(
  {
    display: "Points+Assists",
    abbreviation: "PTS+AST",
    label: "pts+asts",
  },
  (game: Game) => game.points + game.assists
);
export const ReboundsAssists: Stat = new Stat(
  {
    display: "Rebounds+Assists",
    abbreviation: "REB+AST",
    label: "rebs+asts",
  },
  (game: Game) => game.rebounds + game.assists
);
export const Unknown: Stat = new Stat(
  {
    display: "Unknown",
    abbreviation: "UNKNOWN",
    label: "unknown",
  },
  (game: Game) => 0
);

export const LookupStats: Record<string, Stat> = {
  Points: Points,
  Rebounds: Rebounds,
  Assists: Assists,
  "Blks+Stls": BlocksSteals,
  "Pts+Rebs+Asts": PointsReboundsAssists,
  "Pts+Rebs": PointsRebounds,
  "Pts+Asts": PointsAssists,
  "Rebs+Asts": ReboundsAssists,
  rebs_asts: ReboundsAssists,
  "Fantasy Score": Fantasy,
  "3-PT Made": ThreeFGM,
  "Free Throws Made": FreeThrowsMade,
  Turnovers: Turnovers,
};

export function GetStat(proposition: Proposition): Stat {
  let stat: Stat = LookupStats[proposition.type];
  if (!stat) {
    console.log("UNKNOWN type: ", proposition.type);
    return Unknown;
  }
  return stat;
}
