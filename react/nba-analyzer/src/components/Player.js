import React, {useState, useEffect, useRef} from 'react'
import {useLocation} from 'react-router-dom';
import {FormatDate} from '../utils.js'
import PlayerStatsChart from './PlayerStatsChart'
import DataListInput from "react-datalist-input";
import {std, mean} from 'mathjs'
import NormalDistribution from 'normal-distribution'
import {gql, useQuery} from '@apollo/client';
import Playercard from './Playercard.js';

import "../styles/player.css"

const Player = () => {
    let location = useLocation()
    const date = FormatDate(new Date())
    console.log(date)
    const playerID = parseInt(location.pathname.split("/")[location.pathname.split("/").length-1])
    const query = gql` query Player($playerID: Int!, $date:String!) {       
        player(input:{playerID: $playerID}) {
            name
            projections(input:{sportsbook:"PrizePicks", startDate:$date}) {
                date
                opponent {
                    abbreviation
                }
                targets {
                    target
                    type
                }
            }
            games(input: {season:"2021-22"}) {
                points
                assists
                total_rebounds
                offensive_rebounds
                defensive_rebounds
                steals
                blocks
                turnovers
                # win_or_loss
                opponent {
                    abbreviation
                }
                minutes
                date
                field_goals_attempted
                field_goal_percentage
                field_goals_made
               
            }
        }
      }`;
    const { loading, error, data } = useQuery(query, {variables: {playerID:playerID, date:date}});
    // const [games, setGames] = useState([]);
    // const [player, setPlayer] = useState('');

    // useEffect(() => {
    //     if (data) { 
    //         setPlayer(data.player)
    //         setGames(player.games)
    //         setGames(games.sort(function(a, b) {
    //             var c = new Date(a.date);
    //             var d = new Date(b.date);
    //             return c-d;}));
    //     }
    // },
    // [data, games, player]
    // );
    
    if (loading) return 'Loading...';
    if (error) {
        return `Error! ${error.message}. ${loading}. ${data}`;
    }

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
