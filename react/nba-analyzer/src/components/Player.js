import React, {useState, useEffect, useRef} from 'react'
import {useLocation} from 'react-router-dom';
import {GetPropScore} from '../utils.js'
import PlayerStatsChart from './PlayerStatsChart'
import DataListInput from "react-datalist-input";
import {std, mean} from 'mathjs'
import NormalDistribution from 'normal-distribution'
import {gql, useQuery} from '@apollo/client';
import Playercard from './Playercard.js';

import "../styles/player.css"

const Player = () => {
    let location = useLocation()
    const playerID = parseInt(location.pathname.split("/")[location.pathname.split("/").length-1])
    const query = gql` query Player($playerID: Int!) {
        player(input:{playerID: $playerID}) {
            name
            games(input: {season:"2021-22"}) {
                points
                assists
                minutes
                date
                field_goals_attempted
                field_goal_percentage
                field_goals_made
            }
        }
      }`;
    const { loading, error, data } = useQuery(query, {variables: {playerID}});
    // const [games, setGames] = useState([]);
    // const [player, setPlayer] = useState('');

    // useEffect(() => {
    //     console.log(data.player.games)
    //     if (data) { 
    //         setPlayer(data.player)
    //         setGames(player.games)
    //         setGames(games.sort(function(a, b) {
    //             var c = new Date(a.date);
    //             var d = new Date(b.date);
    //             return c-d;}));
    //         console.log(games)
    //     }
    // },
    // [data, games, player]
    // );
    
    if (loading) return 'Loading...';
    if (error) {
        return `Error! ${error.message}. ${loading}. ${data}`;
    }
    // console.log(games)

    return (
        <div className="player-page">
            <div className="player-card">
                <h1>{data.player.name}</h1>
            </div>
            <PlayerStatsChart games={data.player.games}/>
        </div>
    )
}

export default Player
