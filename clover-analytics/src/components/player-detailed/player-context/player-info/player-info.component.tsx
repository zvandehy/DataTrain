import { useEffect, useState } from "react";
import { Player } from "../../../../shared/interfaces/graphql/player.interface";
import AutocompleteFilter from "../../../autocomplete-filter/autocomplete-filter.component";

import "./player-info.component.css";
import moment from "moment";
import EventsModal from "../../../events-modal/events-modal.component";
import { Projection } from "../../../../shared/interfaces/graphql/projection.interface";
import { Game } from "../../../../shared/interfaces/graphql/game.interface";

interface PlayerInfoProps {
  player: Player;
  projection: Projection | undefined;
  game: Game | undefined;
}
let league = "wnba";

const PlayerInfo: React.FC<PlayerInfoProps> = ({
  player,
  projection,
  game,
}: PlayerInfoProps) => {
  const opponent = projection?.opponent ?? game?.opponent;
  return (
    <div className="player-info">
      <h2 className="player-name">
        {player.name}
        <span className="player-position tall titlecase">
          {player.position}
        </span>
      </h2>
      <div className="player-matchup">
        <span>
          <img
            loading="lazy"
            className="team-logo"
            alt={`team-logo for ${player.currentTeam.abbreviation}`}
            src={
              league === "nba"
                ? `https://cdn.nba.com/logos/nba/${player.currentTeam.teamID}/primary/D/logo.svg`
                : `https://${player.currentTeam.name.toLowerCase()}.wnba.com/wp-content/themes/wnba-parent/img/logos/${player.currentTeam.name.toLowerCase()}-primary-logo.svg`
            }
          ></img>
          <span className="team-abr">{player.currentTeam.abbreviation}</span>
        </span>
        <span className="versus">vs.</span>
        <span>
          <img
            loading="lazy"
            className="team-logo"
            alt={`team-logo for ${opponent?.abbreviation}`}
            src={
              league === "nba"
                ? `https://cdn.nba.com/logos/wnba/${opponent?.teamID}/primary/D/logo.svg`
                : `https://${opponent?.name.toLowerCase()}.wnba.com/wp-content/themes/wnba-parent/img/logos/${opponent?.name.toLowerCase()}-primary-logo.svg`
            }
          ></img>
          <span className="opp-abr">{opponent?.abbreviation}</span>
        </span>
        {projection ? <EventsModal projection={projection} /> : <></>}
      </div>
    </div>
  );
};

export default PlayerInfo;
