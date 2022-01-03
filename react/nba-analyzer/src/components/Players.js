import React from 'react'
import PlayerPreview  from './PlayerPreview'
import {PLAYERS_QUERY} from '../utils.js'
import { useQuery } from '@apollo/client';

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
                    {data.players.slice(0,10).map((player) => (<PlayerPreview player={player} key={player.playerID}/>))}
                </tbody>
            </table>
        </div>
    )
}

export default Players
