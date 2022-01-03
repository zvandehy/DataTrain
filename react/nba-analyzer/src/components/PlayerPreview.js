import React from 'react'
import {Link} from 'react-router-dom';
import { useQuery } from '@apollo/client';
import {average, PLAYER_PREVIEW_QUERY} from '../utils.js'

const PlayerPreview = (props) => {
    const {player} = props;
    const {loading, error, data} = useQuery(PLAYER_PREVIEW_QUERY, {variables: {playerID: player.playerID}});
    if (loading) return <tr><td>Loading...</td></tr>;
    if (error) return <tr><td>Error! {error.message}. {loading}. {data}</td></tr>;
    return (
    <><tr className="player"><th><Link to={`/players/${player.playerID}`}>{player.first_name} {player.last_name} <span className="team">{player.teamABR}</span></Link></th><td>{average("assists", data.games.filter((game) => game.season === "2021-22"))}</td><td>{average("points", data.games.filter((game) => game.season === "2021-22"))}</td><td>{average("total_rebounds", data.games.filter((game) => game.season === "2021-22"))}</td></tr></>
    )
}

export default PlayerPreview