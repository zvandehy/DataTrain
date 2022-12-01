export interface PlayerFilter {
  name?: string;
  playerID?: number;
  season?: string;
  position?: string;
  seasons?: SeasonOption[];
  positionStrict?: Position;
  positionStrictMatch?: boolean;
  positionLoose?: Position;
  positionLooseMatch?: boolean;
  teamABR?: string;
  teamID?: number;
  startDate?: string;
  endDate?: string;
  statFilters?: StatFilter[];
  withPropositions?: ProjectionFilter;
  withGames?: GameFilter;
}

export enum Position {
  G = "G",
  F = "F",
  C = "C",
  G_F = "G_F",
  F_G = "F_G",
  F_C = "F_C",
  C_F = "C_F",
}

export interface ProjectionFilter {
  sportsbook?: string;
  playerName?: string;
  playerID?: number;
  startDate?: string;
  endDate?: string;
  teamID?: number;
  opponentID?: number;
}

export interface TeamFilter {
  name?: string;
  teamID?: number;
  abbreviation?: string;
}

export interface GameFilter {
  teamID?: number;
  opponentID?: number;
  opponentMatch?: boolean;
  playerID?: number;
  gameID?: string;
  seasons?: SeasonOption[];
  previousSeasonMatch?: boolean;
  seasonMatch?: boolean;
  startDate?: string;
  endDate?: string;
  lastX?: number;
  gameType?: GameTypeOption;
  gameTypeMatch?: boolean;
  homeOrAway?: HomeOrAwayOption;
  homeOrAwayMatch?: boolean;
  statFilters?: StatFilter[];
  outcome?: OutcomeOption;
}

export interface TeamFilter {
  name?: string;
  teamID?: number;
  abbreviation?: string;
}

export enum SeasonOption {
  SEASON_2020_21 = "SEASON_2020_21",
  SEASON_2021_22 = "SEASON_2021_22",
  SEASON_2022_23 = "SEASON_2022_23",
}

export enum GameTypeOption {
  REGULAR_SEASON = "REGULAR_SEASON",
  PLAYOFFS = "PLAYOFFS",
}

export enum HomeOrAwayOption {
  HOME = "HOME",
  AWAY = "AWAY",
}

export enum OutcomeOption {
  WIN = "WIN",
  LOSS = "LOSS",
  PENDING = "PENDING",
}
export interface StatFilter {
  period?: Period;
  stat: string; // Maybe this should be an enum?
  mode: StatMode;
  operator: Operator;
  value: number;
}

export interface Period {
  startDate?: string;
  endDate?: string;
  seasons?: SeasonOption[];
  limit?: number;
}

export enum Operator {
  GT = "GT",
  GTE = "GTE",
  LT = "LT",
  LTE = "LTE",
  EQ = "EQ",
  NEQ = "NEQ",
}

export enum StatMode {
  PER_GAME = "PER_GAME",
  PER_36 = "PER_36",
  PER_MINUTE = "PER_MINUTE",
  TOTAL = "TOTAL",
}
