import React, {useState} from 'react'
import {useParams} from 'react-router-dom';
import { useQuery } from '@apollo/client';
import {average, PLAYERGAMES_QUERY} from '../utils.js'
import { ResponsiveLine } from '@nivo/line'
import {ScatterPlot} from '@nivo/scatterplot'
import {std, mean, round, variance, mad} from 'mathjs'
import NormalDistribution from 'normal-distribution'

//TODO: see Apexcharts 
//https://apexcharts.com/react-chart-demos/mixed-charts/multiple-yaxis/
const Player = () => {
    const [target, setTarget] = useState(10);
    const {id} = useParams();
    console.log(id)
    const {loading, error, data} = useQuery(PLAYERGAMES_QUERY, {variables: {playerID: id}});
    if (loading) return <p>Loading..</p>;
    if (error || data.players.length === 0) return <p>Error! {error.message}. {loading}. {data}</p>;
    const player = data.players.filter((player) => player.playerID === parseInt(id))[0]
    console.log(player)
    let sortedGames = data.games.slice().filter((game) => game.season==="2021-22").sort(function(a, b) {
        var c = new Date(a.date);
        var d = new Date(b.date);
        return c-d;
    });
    const avgPoints = average("points", data.games.filter((game) => game.season === "2021-22"))
    let chartData = [{id:"points", data: []}, {id:"minutes", data: []}];
    let pointsArr = []
    let countMap = {}
    //points data
    let max = 0
    let min = -1;
    sortedGames.forEach((game) => {
        chartData.filter(series => series.id === "points")[0].data.push({y:game.points,x:game.date})
        chartData.filter(series => series.id === "minutes")[0].data.push({y:parseInt(game.minutes.split(":")[0])/10,x:game.date})
        if (countMap[game.points]) {
            countMap[game.points]++
        } else {
            countMap[game.points]=1
        }
        pointsArr.push(game.points)
        if (game.points > max) {
            max = game.points
        }
        if (game.points < min || min === -1) {
            min = game.points
        }
    });
    let less = 0;
    let more = 0;
    for (let i = 0; i < max; i++) {
        if (!countMap[i]) {
            countMap[i] = 0
        }
        if (i<target) {
            less+=countMap[i]
        } else {
            more+=countMap[i]
        }
    }
    let countData = [{id:"points", data: []}, {id:"normal", data: []}];
    let stddev = std(pointsArr)
    let m = mean(pointsArr)
    let lambda = variance(pointsArr)
    const normDist = new NormalDistribution(m, stddev);
    // console.log(m)
    // console.log(stddev)
    // console.log(normDist.pdf(m))
    // console.log(normDist.zScore(5))
    // console.log(normDist.zScore(8))
    // console.log(normDist.zScore(10))
    // console.log(normDist.zScore(15))
    // console.log(normDist.zScore(20))
    // console.log(normDist.zScore(25))

    let scale = 10
   for(let i = 0; i< 3*stddev+m;i++) {
       let xval = i
       let yval = normDist.pdf(xval);
        if (xval >= 0 ) {
            countData[1].data.push({x:xval, y:yval})
        }
        countData[0].data.push({x:xval, y:countMap[xval]? countMap[xval] /100: 0})
   }
   
    return (
        <div>
             <input type="text" id="search" className="search" onChange={e => {
                setTarget(e.target.value)
            }}></input>
            <div className="player-card">
                <h1>{player.first_name} {player.last_name}</h1>
            </div>
            <p>Points: {avgPoints}</p>
            <p>Variance: {lambda}</p>
            <p>Over: {more}</p>
            <p>Under: {less}</p>
            <div className="line-chart" style={{height:"500px"}}>
                <ResponsiveLine
                    data={chartData}
                    margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                    xScale={{ type: 'point' }}
                    axisBottom={{tickRotation:-90, tickSize:5, tickPadding: 0, format: function(value) {return `${new Date(value).getMonth()+1}/${new Date(value).getDate()+1}`}}}
                    yScale={{ type: 'linear', min: 0, max: 'auto', stacked: false, reverse: false }}
                    yFormat=" >-.2f"
                    markers={[{axis:"y", value: avgPoints, legend: "average"}, {axis:"y", value: target, legend: "target"}]}
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
                                <div>Date: {slice.points.filter((point) => point.serieId === "points")[0].data.x}</div>
                                <div>GameID: {sortedGames[slice.points.filter((point) => point.serieId === "points")[0].index].gameID}</div>
                                <div>Opponent: {sortedGames[slice.points.filter((point) => point.serieId === "points")[0].index].opponent}</div>
                                <div>Season: {sortedGames[slice.points.filter((point) => point.serieId === "points")[0].index].season}</div>
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === "minutes"})[0].serieColor}`}}>
                                    <strong>minutes</strong> {sortedGames[slice.points.filter((point) => point.serieId === "points")[0].index].minutes}
                                </div>
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === "points"})[0].serieColor}`}}>
                                    <strong>points</strong> {sortedGames[slice.points.filter((point) => point.serieId === "points")[0].index].points}
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
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === "normal"})[0].serieColor}`}}>
                                    <strong>zscore</strong> {normDist.zScore(slice.points.filter((point) => point.serieId === "normal")[0].data.x)}
                                </div>
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === "normal"})[0].serieColor}`}}>
                                    <strong># games</strong> {countMap[slice.points.filter((point) => point.serieId === "normal")[0].data.x]?? 0}
                                </div>
                            </div>
                        )
                    }}
                />
            </div>
        </div>
    )
}

export default Player
