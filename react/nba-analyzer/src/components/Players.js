import React from 'react'
import Player  from './Player'

import { useQuery, gql } from '@apollo/client';

const PLAYERS_QUERY = gql`
{
  players {
    playerID
    first_name
    last_name
    seasons
    teamABR
  }
}`

const Players = () => {
    const { loading, error, data } = useQuery(PLAYERS_QUERY);

    if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}. ${loading}. ${data}`;
    return (
        <div className="players">
            <table>
                <thead>
                    <tr><th>Player</th><th>Assists</th><th>Points</th><th>Rebounds</th></tr>
                </thead>
                <tbody>
                <Player player={data.players[0]}/>
                    {/* {data.players.map((player) => (<Player player={player}/>))} */}
                </tbody>
            </table>
        </div>
    )
}

export default Players
