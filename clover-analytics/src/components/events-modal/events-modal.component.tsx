import EventItem from "./events-item/event-item.component";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPersonCircleExclamation as injury } from "@fortawesome/free-solid-svg-icons";
import "./events-modal.component.css";

interface EventsModalProps {
  // projection: Projection;
}

const EventsModal: React.FC<EventsModalProps> = ({}: // projection,
EventsModalProps) => {
  return (
    <button className="events-icon">
      <FontAwesomeIcon icon={injury} />
      <table className="events">
        <tbody>
          {/* {projection.player.currentTeam.injuries?.map((injury, i) => {
            return (
              <EventItem
                key={"team injury " + i}
                team={projection.player.currentTeam}
                injury={injury}
              />
            );
          })}
          {projection.opponent?.injuries?.map((injury, i) => {
            return (
              <EventItem
                key={"opp injury " + i}
                team={projection.opponent}
                injury={injury}
              />
            );
          })} */}
        </tbody>
      </table>
    </button>
  );
};

export default EventsModal;
