import { Player } from "./player.interface";
import { PlayersInGame } from "./playersInGame.interface";
import { Projection } from "./projection.interface";
import { Team } from "./team.interface";
import { TeamGame } from "./teamgame.interface";

export interface Game {
  assist_percentage: number;
  assists: number;
  date: string;
  defensive_rebound_percentage: number;
  defensive_rebounds: number;
  effective_field_goal_percentage: number;
  field_goal_percentage: number;
  field_goals_attempted: number;
  field_goals_made: number;
  free_throws_attempted: number;
  free_throws_made: number;
  free_throws_percentage: number;
  gameID: string;
  home_or_away: string;
  win_or_loss: string;
  margin: number;
  minutes: string;
  offensive_rebound_percentage: number;
  offensive_rebounds: number;
  opponent: Team;
  opponentStats: TeamGame;
  team: Team;
  teamStats: TeamGame;
  personal_fouls_drawn: number;
  personal_fouls: number;
  points: number;
  player: Player;
  playoffs: boolean;
  season: string;
  three_point_percentage: number;
  three_pointers_attempted: number;
  three_pointers_made: number;
  rebounds: number;
  true_shooting_percentage: number;
  turnovers: number;
  usage: number;
  blocks: number;
  steals: number;
  playersInGame: PlayersInGame;
  projections: Projection[];
}
//   constructor(game: Game) {
//     this.assist_percentage = game.assist_percentage;
//     this.assists = game.assists;
//     this.date = game.date;
//     this.defensive_rebound_percentage = game.defensive_rebound_percentage;
//     this.defensive_rebounds = game.defensive_rebounds;
//     this.effective_field_goal_percentage = game.effective_field_goal_percentage;
//     this.field_goal_percentage = game.field_goal_percentage;
//     this.field_goals_attempted = game.field_goals_attempted;
//     this.field_goals_made = game.field_goals_made;
//     this.free_throws_attempted = game.free_throws_attempted;
//     this.free_throws_made = game.free_throws_made;
//     this.free_throws_percentage = game.free_throws_percentage;
//     this.gameID = game.gameID;
//     this.home_or_away = game.home_or_away;
//     this.win_or_loss = game.win_or_loss;
//     this.margin = game.margin;
//     this.minutes = game.minutes;
//     this.offensive_rebound_percentage = game.offensive_rebound_percentage;
//     this.offensive_rebounds = game.offensive_rebounds;
//     this.opponent = game.opponent;
//     this.opponentStats = game.opponentStats;
//     this.team = game.team;
//     this.teamStats = game.teamStats;
//     this.personal_fouls_drawn = game.personal_fouls_drawn;
//     this.personal_fouls = game.personal_fouls;
//     this.points = game.points;
//     this.player = game.player;
//     this.playoffs = game.playoffs;
//     this.season = game.season;
//     this.three_point_percentage = game.three_point_percentage;
//     this.three_pointers_attempted = game.three_pointers_attempted;
//     this.three_pointers_made = game.three_pointers_made;
//     this.rebounds = game.rebounds;
//     this.true_shooting_percentage = game.true_shooting_percentage;
//     this.turnovers = game.turnovers;
//     this.usage = game.usage;
//     this.blocks = game.blocks;
//     this.steals = game.steals;
//     this.playersInGame = game.playersInGame;
//     this.projections = game.projections;
//   }

//   fantasy(): number {
//     return (
//       this.points +
//       this.rebounds * 1.2 +
//       this.assists * 1.5 +
//       this.blocks * 3 +
//       this.steals * 3 -
//       this.turnovers
//     );
//   }

//   pra(): number {
//     return this.points + this.assists + this.rebounds;
//   }

//   ptsAst
// }
