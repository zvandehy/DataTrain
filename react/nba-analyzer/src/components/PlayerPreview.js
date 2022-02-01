import React from 'react'
import {Link} from 'react-router-dom';
import { useQuery } from '@apollo/client';
import {average, PLAYER_PREVIEW_QUERY} from '../utils.js'
import {mean, variance, mad, median, round, std} from'mathjs'

const PlayerPreview = (props) => {
    const {playerProp} = props;
    let propType = playerProp.type.toLowerCase();
    let player=playerProp.player
    let seasonData=playerProp.playerGames.filter((game) => game.season === "2021-22")
    // let seasonData = data.games.filter((game) => game.season === "2021-22");
    let statData = [0];
    if (seasonData && seasonData.length > 0) {
        statData = seasonData.map((game) => game[propType]);
    }
    // console.log(statData)
    // console.groupEnd()
    return (
    <><tr className="player"><th><Link to={`/players/${player.playerID}`} state={{playerProp:playerProp}}>{player.first_name} {player.last_name} <span className="team">{player.teamABR}</span> vs. <span className="team opponent">{playerProp.opponent}</span></Link></th>
        {/* <td>{average("assists", data.games.filter((game) => game.season === "2021-22"))}</td>
        <td>{average("points", data.games.filter((game) => game.season === "2021-22"))}</td>
        <td>{average("total_rebounds", data.games.filter((game) => game.season === "2021-22"))}</td> */}
         <td>{playerProp.target}</td>
        <td>{round(mean(statData),2)}</td>
        <td>{round(std(statData),2)}</td>
        <td>{round(median(statData),2)}</td>
        <td>{round(mad(statData),2)}</td>
    </tr></>
    )
}

export default PlayerPreview