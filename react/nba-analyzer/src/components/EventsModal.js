import React from "react";
import EventItem from "./Eventitem";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPersonCircleExclamation as injury } from "@fortawesome/free-solid-svg-icons";

const EventsModal = (props) => {
  const { player, opponent } = props;
  return (
    <button className="events-icon" href="#">
      <FontAwesomeIcon icon={injury} />
      <table className="events">
        <tbody>
          {player.currentTeam.injuries?.map((injury, i) => {
            let status = injury.status;
            let date = injury.startDate;
            if (injury.returnDate !== "") {
              status = "Returned";
              date = injury.returnDate;
            }
            return (
              <EventItem
                key={"team injury " + i}
                team={player.currentTeam.abbreviation}
                playername={injury.player.name}
                updateMsg={status}
                date={date}
              />
            );
          })}
          {opponent.injuries?.map((injury, i) => {
            let status = injury.status;
            let date = injury.startDate;
            if (injury.returnDate !== "") {
              status = "Returned";
              date = injury.returnDate;
            }
            return (
              <EventItem
                key={"opp injury " + i}
                team={opponent.abbreviation}
                playername={injury.player.name}
                updateMsg={status}
                date={date}
              />
            );
          })}
        </tbody>
      </table>
    </button>
  );
};

export default EventsModal;
