import { PlayerGame } from "./graphql/types";

export interface Stat {
  display: string;
  label: string;
  altLabel?: string;
  abbreviation: string;
  relatedStats: Stat[];

  GetScore(game: PlayerGame): number;
  AverageScore(games: PlayerGame[]): number;
}

export class BasicStat implements Stat {
  display: string;
  label: string;
  abbreviation: string;
  relatedStats: Stat[];
  altLabel?: string;

  constructor(
    display: string,
    label: string,
    abbreviation: string,
    relatedStats: Stat[],
    altLabel?: string
  ) {
    this.display = display;
    this.label = label;
    this.abbreviation = abbreviation;
    this.relatedStats = relatedStats;
    this.altLabel = altLabel;
  }

  GetScore(game: PlayerGame): number {
    return game[this.label as keyof PlayerGame] as number;
  }
  AverageScore(games: PlayerGame[]): number {
    return (
      games.map((game) => this.GetScore(game)).reduce((a, b) => a + b, 0) /
      games.length
    );
  }
}

export class CalculateStat implements Stat {
  display: string;
  label: string;
  abbreviation: string;
  relatedStats: Stat[];
  altLabel?: string;
  scoreFn: (game: PlayerGame) => number;
  avgScoreFn: (games: PlayerGame[]) => number;

  constructor(
    display: string,
    label: string,
    abbreviation: string,
    relatedStats: Stat[],
    scoreFn: (game: PlayerGame) => number,
    avgScoreFn?: (games: PlayerGame[]) => number,
    altLabel?: string
  ) {
    this.display = display;
    this.label = label;
    this.abbreviation = abbreviation;
    this.relatedStats = relatedStats;
    this.scoreFn = scoreFn;
    this.avgScoreFn = avgScoreFn ?? this.defaultAvgScore;
    this.altLabel = altLabel;
  }

  GetScore(game: PlayerGame): number {
    return this.scoreFn(game);
  }
  AverageScore(games: PlayerGame[]): number {
    return this.avgScoreFn(games);
  }
  defaultAvgScore(games: PlayerGame[]): number {
    return (
      games.map((game) => this.GetScore(game)).reduce((a, b) => a + b, 0) /
      games.length
    );
  }
}

export const POINTS: Stat = new BasicStat("Points", "points", "PTS", []);

export const REBOUNDS: Stat = new BasicStat("Rebounds", "rebounds", "REB", []);

export const ASSISTS: Stat = new BasicStat("Assists", "assists", "AST", []);

export const STEALS: Stat = new BasicStat("Steals", "steals", "STL", []);

export const BLOCKS: Stat = new BasicStat("Blocks", "blocks", "BLK", []);

export const TURNOVERS: Stat = new BasicStat(
  "Turnovers",
  "turnovers",
  "TOV",
  []
);

export const BLOCKS_AND_STEALS: Stat = new CalculateStat(
  "BLKS+STLS",
  "blks+stls",
  "B+S",
  [BLOCKS, STEALS],
  (game: PlayerGame) => game.blocks + game.steals
);

export const REBOUNDS_AND_ASSISTS: Stat = new CalculateStat(
  "REB+AST",
  "rebs+asts",
  "R+A",
  [REBOUNDS, ASSISTS],
  (game: PlayerGame) => game.rebounds + game.assists
);

export const POINTS_AND_REBOUNDS: Stat = new CalculateStat(
  "PTS+REB",
  "pts+rebs",
  "P+R",
  [POINTS, REBOUNDS],
  (game: PlayerGame) => game.points + game.rebounds
);

export const POINTS_AND_ASSISTS: Stat = new CalculateStat(
  "PTS+AST",
  "pts+asts",
  "P+A",
  [REBOUNDS, ASSISTS],
  (game: PlayerGame) => game.rebounds + game.assists
);

export const THREES_MADE: Stat = new BasicStat(
  "3-Pointers Made",
  "three_pointers_made",
  "3PM",
  [],
  "3-pt made"
);

// assist_percentage: Scalars["Float"];
// defensive_rebound_percentage: Scalars["Float"];
// defensive_rebounds: Scalars["Int"];
// effective_field_goal_percentage: Scalars["Float"];
// field_goal_percentage: Scalars["Float"];
// field_goals_attempted: Scalars["Int"];
// field_goals_made: Scalars["Int"];
// free_throws_attempted: Scalars["Int"];
// free_throws_made: Scalars["Int"];
// free_throws_percentage: Scalars["Float"];
// minutes: Scalars["String"];
// offensive_rebound_percentage: Scalars["Float"];
// offensive_rebounds: Scalars["Int"];
// personal_fouls_drawn: Scalars["Int"];
// personal_fouls: Scalars["Int"];
// three_point_percentage: Scalars["Float"];
// three_pointers_attempted: Scalars["Int"];
// three_pointers_made: Scalars["Int"];
// true_shooting_percentage: Scalars["Float"];
// usage: Scalars["Float"];
