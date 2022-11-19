import { Game } from "./game.interface";
import { Team } from "./team.interface";

export interface Player {
  name: string;
  playerID: number;
  image: string;
  seasons: string[];
  position: string;
  currentTeam: Team;
  team: Team;
  games: Game[];
  height: string;
  weight: number;
  similarPlayers: Player[];
}
