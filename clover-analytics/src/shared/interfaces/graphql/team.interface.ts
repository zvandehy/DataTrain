import { Injury } from "./injury.interface";
import { Player } from "./player.interface";
import { TeamGame } from "./teamgame.interface";

export interface Team {
  name: string;
  teamID: number;
  abbreviation: string;
  location: string;
  games: TeamGame[];
  players: Player[];
  injuries: Injury[];
}
