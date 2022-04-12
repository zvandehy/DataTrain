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
        <div>
            <div className="player-card">
                <h1>{data.player.name}</h1>
            </div>
            {/* <Playercard playerProp={data} key={data.player.playerID}/> */}
            {/* <DataListInput
                placeholder="Select a prop type"
                items={selectPropTypes}
                onSelect={(event) => setPropType(event.label)}
                clearInputOnClick={true}
                suppressReselect={false}
            />
            <p>Avergage {propType}: {m}</p>
            <p>Std Dev: {stddev}</p>
            <p>Target: {target}</p>
            <p>Over: {more}</p>
            <p>Under: {less}</p> */}
            <PlayerStatsChart games={data.player.games}/>
            {/* <div className="line-chart" style={{height:"500px"}}>
                <ResponsiveLine
                    data={chartData}
                    margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                    xScale={{ type: 'point' }}
                    axisBottom={{tickRotation:-90, tickSize:5, tickPadding: 0, format: function(value) {return `${new Date(value).getMonth()+1}/${new Date(value).getDate()+1}`}}}
                    yScale={{ type: 'linear', min: 0, max: 'auto', stacked: false, reverse: false }}
                    yFormat=" >-.2f"
                    markers={[{axis:"y", value: m, legend: "average"}, {axis:"y", value: target, legend: "target"}]}
                    legends={[
                        {
                            anchor: 'bottom-right',
                            direction: 'column',
                            justify: false,
                            translateX: 100,
                            translateY: 0,
                            itemsSpacing: 0,
                            itemDirection: 'left-to-right',
                            itemWidth: 80,
                            itemHeight: 20,
                            itemOpacity: 0.75,
                            symbolSize: 12,
                            symbolShape: 'circle',
                            symbolBorderColor: 'rgba(0, 0, 0, .5)',
                            effects: [
                                {
                                    on: 'hover',
                                    style: {
                                        itemBackground: 'rgba(0, 0, 0, .03)',
                                        itemOpacity: 1
                                    }
                                }
                            ]
                        }
                    ]}
                    useMesh={true}
                    isInteractive={true}
                    pointLabelYOffset={0}
                    enableSlices="x"
                    sliceTooltip={({ slice }) => {
                        return (
                            <div
                                style={{
                                    background: 'white',
                                    padding: '9px 12px',
                                    border: '1px solid #ccc',
                                }}
                            >
                                <div>Date: {slice.points.filter((point) => point.serieId === propType)[0].data.x}</div>
                                <div>GameID: {sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index].gameID}</div>
                                <div>Opponent: {sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index].opponent.abbreviation}</div>
                                <div>Season: {sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index].season}</div>
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === "minutes"})[0].serieColor}`}}>
                                    <strong>minutes</strong> {sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index].minutes}
                                </div>
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === propType})[0].serieColor}`}}>
                                    <strong>{propType}</strong> {GetPropScore(sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index],propType)}
                                </div>
                            </div>
                        )
                    }}
                />
            </div>
            <div className="line-chart" style={{height:"500px"}}>
                <ResponsiveLine data={countData}
                     margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                     xScale={{ type: 'point'}}
                     yScale={{ type: 'linear', min: 0, max: 'auto', stacked: false, reverse: false }}
                     yFormat=" >-.2f"
                     useMesh={true}
                    isInteractive={true}
                    pointLabelYOffset={0}
                    enableSlices="x"
                    areaOpacity={0.7}
                    enableArea={true}
                    axisBottom={{tickRotation:-90, tickSize:5, tickPadding: 0, format: " >-.2f"}}
                    sliceTooltip={({ slice }) => {
                        return (
                            <div
                                style={{
                                    background: 'white',
                                    padding: '9px 12px',
                                    border: '1px solid #ccc',
                                }}
                            >
                                <div style={{color:'grey'}}>
                                    <strong>zscore</strong> {normDist.zScore(slice.points.filter((point) => point.serieId === "normal")[0].data.x)}
                                </div>
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === "normal"})[0].serieColor}`}}>
                                    <strong>cdf</strong> {normDist.cdf(slice.points.filter((point) => point.serieId === "normal")[0].data.x)}
                                </div>
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === propType})[0].serieColor}`}}>
                                    <strong># games</strong> {countMap[slice.points.filter((point) => point.serieId === "normal")[0].data.x]?? 0}
                                </div>
                            </div>
                        )
                    }}
                />
            </div> */}
        </div>
    )
}

export default Player
