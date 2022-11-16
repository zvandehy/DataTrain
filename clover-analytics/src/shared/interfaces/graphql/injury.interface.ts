import { Player } from "./player.interface";

export interface Injury {
  startDate: string;
  returnDate: string;
  status: string;
  player: Player;
}
