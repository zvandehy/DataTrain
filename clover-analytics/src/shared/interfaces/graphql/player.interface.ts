import { Game } from "./game.interface";
import { Injury } from "./injury.interface";
import { Projection } from "./projection.interface";
import { Team } from "./team.interface";

export interface Player {
  name: string;
  playerID: number;
  seasons: string[];
  position: string;
  currentTeam: Team;
  team: Team;
  games: Game[];
  injuries: Injury[];
  projections: Projection[];
  height: string;
  weight: number;
  similarPlayers: Player[];
}
