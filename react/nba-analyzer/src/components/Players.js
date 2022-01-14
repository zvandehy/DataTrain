import React, {useState} from 'react'
import PlayerPreview  from './PlayerPreview'
import {PLAYERS_QUERY} from '../utils.js'
import { useQuery } from '@apollo/client';

const Players = () => {
    const [lookup, setLookup] = useState('DEN');
    const { loading, error, data } = useQuery(PLAYERS_QUERY);

    if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}. ${loading}. ${data}`;
    
    
    
    console.log(data.players.filter((player) => player.teamABR === lookup))
    return (
        <div className="players">
            <input type="text" id="search" className="search" onChange={e => {
                setLookup(e.target.value)
            }}></input>
            <table>
                <thead>
                    <tr><th>Player</th><th>Mean</th><th>Variance</th><th>Median</th><th>MAD</th></tr>
                </thead>
                <tbody>
                    {data.players.filter((player) => player.teamABR === lookup).map((player) => (<PlayerPreview player={player} key={player.playerID}/>))}
                </tbody>
            </table>
        </div>
    )
}

export default Players
