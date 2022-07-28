export interface PlayerFilter {
  name?: string;
  playerID?: number;
  season?: string;
  position?: string;
  teamABR?: string;
  teamID?: number;
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
  playerID?: number;
  gameID?: string;
  season?: string;
  startDate?: string;
  endDate?: string;
}
