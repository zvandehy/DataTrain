import React from 'react'
import { useQuery, gql } from '@apollo/client';

const PLAYERGAMES_QUERY = gql`
query games($playerID: Int!){
    games(playerID: $playerID) {
        season
        assists
        points
        total_rebounds
        field_goal_percentage
        opponent
    }
}`
function average(stat, data) {
    const sum = data.games.reduce((a,b) => a + b[stat], data.games[0][stat]);
    const avg = sum / data.games.length;
    return  Math.round((avg + Number.EPSILON) * 100) / 100
}
const Player = (props) => {
    const {player} = props;
    const {loading, error, data} = useQuery(PLAYERGAMES_QUERY, {variables: {playerID: player.playerID}});
    if (loading) return <tr><td>Loading...</td></tr>;
    if (error) return <tr><td>Error! {error.message}. {loading}. {data}</td></tr>;
    return (
    <tr><th>{player.first_name} {player.last_name}</th><td>{average("assists", data)}</td><td>{average("points", data)}</td><td>{average("total_rebounds", data)}</td></tr>
    )
}

export default Player