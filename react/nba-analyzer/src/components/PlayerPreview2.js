import React, {useState} from 'react'
import {Link} from 'react-router-dom';
import {mean, round} from'mathjs'
import {GetPropScore} from "../utils"

const PlayerPreview2 = (props) => {
    const {playerProp} = props;
    // playerProp.targets.sort((a, b) => a.type > b.type);
    let player=playerProp.player
    let seasonData=playerProp.player.games.filter((game) => game.season === "2021-22").sort(function(a, b) {
        var c = new Date(a.date);
        var d = new Date(b.date);
        return c-d;
    });
    const [type, setType] = useState("points");

    const counts = [0, -30, -10, -5];
    
    if (!seasonData) {
        return <li></li>
    }
    
    return (
    <>
    <div className="playercard">
      <img className="player-photo" alt={player.name} src={"https://ak-static.cms.nba.com/wp-content/uploads/headshots/nba/latest/260x190/"+player.playerID+".png"}></img>
      <div className="player-info">
        <h2 className="player-name">
        <Link to={`/players/${player.playerID}`} state={{playerProp:playerProp}}>{player.name}</Link>
          <span className="player-position">({player.position ? player.position : "X"})</span>
        </h2>
        <p className="player-team">{player.currentTeam.abbreviation}</p>
        <p>vs <span className="player-opponent">{playerProp.opponent.abbreviation}</span></p>
        <ul className="player-since">
          <li><button onClick={(event) => setType("points")}>PTS</button></li>
          <li><button onClick={(event) => setType("pts+rebs+asts")}>PRA</button></li>
          <li><button onClick={(event) => setType("rebounds")}>REB</button></li>
          <li><button onClick={(event) => setType("assists")}>AST</button></li>
        </ul>
      </div>
      <table  className="player-stats">
        <thead>
          <tr>
            <th></th>
            <th>AVG</th>
            <th>TAR</th>
            <th>OVR</th>
            <th>PCT</th>
          </tr>
        </thead>
        <tbody>
        {counts.map((count) => {
            const data = seasonData.slice(count).map(game => GetPropScore(game, type));
            const target = playerProp.targets.filter(projection => {console.log(playerProp.targets, projection, type, projection.type.toLowerCase(), type.toLowerCase(), projection.type.toLowerCase() === type.toLowerCase(), ); return projection.type.toLowerCase() === type.toLowerCase()})[0]?.target ?? "-";
            const over = target ? data.filter(score => score > target).length : "-";
            const pct = target ? round((data.filter(score => score > target).length / data.length) *100,2) : "-"
            return (<>
            {data.length > 0 ? 
                <tr key={player.playerID + type + count}>
                    <th>{count ? count * -1 : data.length}</th>
                    <td>{round(mean(data), 1)}</td>
                    <td>{target}</td>
                    <td>{over}</td>
                    <td className={"player-over-pct " + (target === "-" ? "" : pct > 60 ? "high" : pct >= 50 ? "med" : "low")}>{pct + "%"}</td>
                </tr> : <tr></tr>}
            </>)
            })
        }
        </tbody>
      </table>
    </div>
    </>
    )
}

export default PlayerPreview2