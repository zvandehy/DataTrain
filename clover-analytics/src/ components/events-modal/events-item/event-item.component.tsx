import { Injury } from "../../../shared/interfaces/graphql/injury.interface";
import { Team } from "../../../shared/interfaces/graphql/team.interface";
import "./event-item.component.css";

interface EventItemProps {
  team: Team;
  injury: Injury;
}

const EventItem: React.FC<EventItemProps> = ({
  team,
  injury,
}: EventItemProps) => {
  let date = injury.startDate;
  if (injury.returnDate !== "") {
    injury.status = "Returned";
    date = injury.returnDate;
  }
  return (
    <tr>
      <th>
        <span className="team-abr bold">{team.abbreviation}</span>
      </th>
      <td>{injury.player.name}</td>
      <td>{injury.status}</td>
      <td>{date}</td>
    </tr>
  );
};

export default EventItem;
