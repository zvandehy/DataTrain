import { Link } from "react-router-dom";
import { Projection } from "../../../../shared/interfaces/graphql/projection.interface";
import EventsModal from "../../../events-modal/events-modal.component";
import "./playercard-context.component.css";

interface PlayercardContextProps {
  projection: Projection;
}

const PlayercardContext: React.FC<PlayercardContextProps> = ({
  projection,
}: PlayercardContextProps) => {
  let league = "wnba";
  return (
    <div className="player-context">
      <img
        loading="lazy"
        className="player-photo"
        alt={projection.player.name}
        src={`https://ak-static.cms.nba.com/wp-content/uploads/headshots/${league}/latest/260x190/${projection.player.playerID}.png`}
      ></img>
      <div className="player-info">
        <h2 className="player-name">
          <Link to={`/${league}/players/${projection.player.playerID}`}>
            {projection.player.name}
          </Link>
          <span className="player-position tall titlecase">
            {projection.player.position}
          </span>
        </h2>
        <div className="player-matchup">
          <span>
            <img
              loading="lazy"
              className="team-logo"
              alt={`team-logo for ${projection.player.currentTeam.abbreviation}`}
              src={
                league === "nba"
                  ? `https://cdn.nba.com/logos/nba/${projection.player.currentTeam.teamID}/primary/D/logo.svg`
                  : `https://${projection.player.currentTeam.name.toLowerCase()}.wnba.com/wp-content/themes/wnba-parent/img/logos/${projection.player.currentTeam.name.toLowerCase()}-primary-logo.svg`
              }
            ></img>
            <span className="team-abr">
              {projection.player.currentTeam.abbreviation}
            </span>
          </span>
          <span className="versus">vs.</span>
          <span>
            <img
              loading="lazy"
              className="team-logo"
              alt={`team-logo for ${projection.opponent?.abbreviation}`}
              src={
                league === "nba"
                  ? `https://cdn.nba.com/logos/wnba/${projection.opponent?.teamID}/primary/D/logo.svg`
                  : `https://${projection.opponent?.name.toLowerCase()}.wnba.com/wp-content/themes/wnba-parent/img/logos/${projection.opponent?.name.toLowerCase()}-primary-logo.svg`
              }
            ></img>
            <span className="opp-abr">{projection.opponent?.abbreviation}</span>
          </span>
          <EventsModal projection={projection} />
        </div>
      </div>
    </div>
  );
};

export default PlayercardContext;
