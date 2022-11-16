export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Player = {
  __typename?: "Player";
  name: Scalars["String"];
  playerID: Scalars["Int"];
  seasons: Array<Scalars["String"]>;
  position: Scalars["String"];
  currentTeam: Team;
  games: Array<PlayerGame>;
  injuries?: Maybe<Array<Injury>>;
  projections: Array<Projection>;
  height: Scalars["String"];
  weight: Scalars["Int"];
  similarPlayers: Array<Player>;
};

export type PlayerGamesArgs = {
  input: GameFilter;
};

export type PlayerProjectionsArgs = {
  input: ProjectionFilter;
};

export type PlayerSimilarPlayersArgs = {
  input: GameFilter;
};

export type Team = {
  __typename?: "Team";
  name: Scalars["String"];
  teamID: Scalars["Int"];
  abbreviation: Scalars["String"];
  location: Scalars["String"];
  numWins: Scalars["Int"];
  numLoss: Scalars["Int"];
  games: Array<TeamGame>;
  players: Array<Player>;
  injuries?: Maybe<Array<Injury>>;
};

export type TeamGamesArgs = {
  input: GameFilter;
};

export type TeamGame = {
  __typename?: "TeamGame";
  date: Scalars["String"];
  margin: Scalars["Int"];
  points: Scalars["Int"];
  assists: Scalars["Int"];
  rebounds: Scalars["Int"];
  offensive_rebounds: Scalars["Int"];
  defensive_rebounds: Scalars["Int"];
  turnovers: Scalars["Int"];
  steals: Scalars["Int"];
  blocks: Scalars["Int"];
  three_pointers_attempted: Scalars["Int"];
  three_pointers_made: Scalars["Int"];
  defensive_rating: Scalars["Float"];
  defensive_rebound_percentage: Scalars["Float"];
  field_goal_percentage: Scalars["Float"];
  field_goals_attempted: Scalars["Int"];
  gameID: Scalars["String"];
  home_or_away: Scalars["String"];
  field_goals_made: Scalars["Int"];
  free_throws_attempted: Scalars["Int"];
  free_throws_made: Scalars["Int"];
  free_throws_percentage: Scalars["Float"];
  offensive_rebound_percentage: Scalars["Float"];
  opponent: Team;
  opponent_assists: Scalars["Int"];
  opponent_effective_field_goal_percentage: Scalars["Float"];
  opponent_field_goals_attempted: Scalars["Int"];
  opponent_free_throws_attempted: Scalars["Int"];
  opponent_points: Scalars["Int"];
  opponent_rebounds: Scalars["Int"];
  opponent_three_pointers_attempted: Scalars["Int"];
  opponent_three_pointers_made: Scalars["Int"];
  plus_minus_per_hundred: Scalars["Float"];
  possessions: Scalars["Int"];
  playoffs: Scalars["Boolean"];
  personal_fouls: Scalars["Int"];
  personal_fouls_drawn: Scalars["Int"];
  season: Scalars["String"];
  win_or_loss: Scalars["String"];
  playersInGame: PlayersInGame;
};

export type PlayerGame = {
  __typename?: "PlayerGame";
  assist_percentage: Scalars["Float"];
  assists: Scalars["Int"];
  date: Scalars["String"];
  defensive_rebound_percentage: Scalars["Float"];
  defensive_rebounds: Scalars["Int"];
  effective_field_goal_percentage: Scalars["Float"];
  field_goal_percentage: Scalars["Float"];
  field_goals_attempted: Scalars["Int"];
  field_goals_made: Scalars["Int"];
  free_throws_attempted: Scalars["Int"];
  free_throws_made: Scalars["Int"];
  free_throws_percentage: Scalars["Float"];
  gameID: Scalars["String"];
  home_or_away: Scalars["String"];
  win_or_loss: Scalars["String"];
  margin: Scalars["Int"];
  minutes: Scalars["String"];
  offensive_rebound_percentage: Scalars["Float"];
  offensive_rebounds: Scalars["Int"];
  opponent: Team;
  opponentStats: TeamGame;
  team: Team;
  teamStats: TeamGame;
  personal_fouls_drawn: Scalars["Int"];
  personal_fouls: Scalars["Int"];
  points: Scalars["Int"];
  player: Player;
  playoffs: Scalars["Boolean"];
  season: Scalars["String"];
  three_point_percentage: Scalars["Float"];
  three_pointers_attempted: Scalars["Int"];
  three_pointers_made: Scalars["Int"];
  rebounds: Scalars["Int"];
  true_shooting_percentage: Scalars["Float"];
  turnovers: Scalars["Int"];
  usage: Scalars["Float"];
  blocks: Scalars["Int"];
  steals: Scalars["Int"];
  playersInGame: PlayersInGame;
  projections: Array<Maybe<Projection>>;
};

export type Projection = {
  __typename?: "Projection";
  player: Player;
  opponent: Team;
  propositions: Array<Maybe<Proposition>>;
  result?: Maybe<PlayerGame>;
  startTime: Scalars["String"];
  date: Scalars["String"];
};

export type Proposition = {
  __typename?: "Proposition";
  target: Scalars["Float"];
  type: Scalars["String"];
  sportsbook: Scalars["String"];
  predictions: Array<Prediction>;
};

export type Prediction = {
  __typename?: "Prediction";
  model: Scalars["String"];
  overUnderPrediction: Scalars["String"];
  totalPrediction?: Maybe<Scalars["Float"]>;
};

export type Injury = {
  __typename?: "Injury";
  startDate: Scalars["String"];
  returnDate: Scalars["String"];
  status: Scalars["String"];
  player: Player;
};

export type Target = {
  __typename?: "Target";
  target: Scalars["Float"];
  type: Scalars["String"];
};

export type PlayersInGame = {
  __typename?: "PlayersInGame";
  team: Array<Player>;
  opponent: Array<Player>;
};

export type Query = {
  __typename?: "Query";
  players: Array<Player>;
  filterPlayers: Array<Player>;
  player: Player;
  teams: Array<Team>;
  filterTeams: Array<Team>;
  team: Team;
  teamGames: Array<TeamGame>;
  playerGames: Array<PlayerGame>;
  projections: Array<Projection>;
};

export type QueryFilterPlayersArgs = {
  input: PlayerFilter;
};

export type QueryPlayerArgs = {
  input: PlayerFilter;
};

export type QueryFilterTeamsArgs = {
  input: TeamFilter;
};

export type QueryTeamArgs = {
  input: TeamFilter;
};

export type QueryTeamGamesArgs = {
  input: GameFilter;
};

export type QueryPlayerGamesArgs = {
  input: GameFilter;
};

export type QueryProjectionsArgs = {
  input: ProjectionFilter;
};

export type ProjectionFilter = {
  sportsbook?: InputMaybe<Scalars["String"]>;
  playerName?: InputMaybe<Scalars["String"]>;
  playerID?: InputMaybe<Scalars["Int"]>;
  startDate?: InputMaybe<Scalars["String"]>;
  endDate?: InputMaybe<Scalars["String"]>;
  teamID?: InputMaybe<Scalars["Int"]>;
  opponentID?: InputMaybe<Scalars["Int"]>;
};

export type PlayerFilter = {
  name?: InputMaybe<Scalars["String"]>;
  playerID?: InputMaybe<Scalars["Int"]>;
  season?: InputMaybe<Scalars["String"]>;
  position?: InputMaybe<Scalars["String"]>;
  teamABR?: InputMaybe<Scalars["String"]>;
  teamID?: InputMaybe<Scalars["Int"]>;
};

export type TeamFilter = {
  name?: InputMaybe<Scalars["String"]>;
  teamID?: InputMaybe<Scalars["Int"]>;
  abbreviation?: InputMaybe<Scalars["String"]>;
};

export type GameFilter = {
  teamID?: InputMaybe<Scalars["Int"]>;
  playerID?: InputMaybe<Scalars["Int"]>;
  gameID?: InputMaybe<Scalars["String"]>;
  season?: InputMaybe<Scalars["String"]>;
  startDate?: InputMaybe<Scalars["String"]>;
  endDate?: InputMaybe<Scalars["String"]>;
};
