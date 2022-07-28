import { PlayersInGame } from "./playersInGame.interface";
import { Team } from "./team.interface";

export interface TeamGame {
  date: string;
  margin: number;
  points: number;
  assists: number;
  rebounds: number;
  offensive_rebounds: number;
  defensive_rebounds: number;
  turnovers: number;
  steals: number;
  blocks: number;
  three_pointers_attempted: number;
  three_pointers_made: number;
  defensive_rating: number;
  defensive_rebound_percentage: number;
  field_goal_percentage: number;
  field_goals_attempted: number;
  gameID: string;
  home_or_away: string;
  field_goals_made: number;
  free_throws_attempted: number;
  free_throws_made: number;
  free_throws_percentage: number;
  offensive_rebound_percentage: number;
  opponent: Team;
  opponent_assists: number;
  opponent_effective_field_goal_percentage: number;
  opponent_field_goals_attempted: number;
  opponent_free_throws_attempted: number;
  opponent_points: number;
  opponent_rebounds: number;
  opponent_three_pointers_attempted: number;
  opponent_three_pointers_made: number;
  plus_minus_per_hundred: number;
  possessions: number;
  playoffs: boolean;
  personal_fouls: number;
  personal_fouls_drawn: number;
  season: string;
  win_or_loss: string;
  playersInGame: PlayersInGame;
}
