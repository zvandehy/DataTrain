import React, {useState} from 'react'
import {useParams, useLocation} from 'react-router-dom';
import { useQuery } from '@apollo/client';
import {average, PLAYERGAMES_QUERY} from '../utils.js'
import { ResponsiveLine } from '@nivo/line'
import {ScatterPlot} from '@nivo/scatterplot'
import {std, mean, round, variance, mad} from 'mathjs'
import NormalDistribution from 'normal-distribution'

//TODO: see Apexcharts 
//https://apexcharts.com/react-chart-demos/mixed-charts/multiple-yaxis/
const Player = (props) => {
    const {id} = useParams();
    let location = useLocation()
    const {playerProp} = location.state
    let target = playerProp.target
    let propType = playerProp.type.toLowerCase();
    let player=playerProp.player
    let seasonData=playerProp.playerGames.filter((game) => game.season === "2021-22")
    let sortedGames = seasonData.sort(function(a, b) {
        var c = new Date(a.date);
        var d = new Date(b.date);
        return c-d;
    });
    const avgPoints = average(propType, sortedGames.filter((game) => game.season === "2021-22"));
    let chartData = [{id:propType, data: []}, {id:"minutes", data: []}];
    let pointsArr = []
    let countMap = {}
    let max = 0
    let min = -1;
    sortedGames.forEach((game) => {
        chartData.filter(series => series.id === propType)[0].data.push({y:game[propType],x:game.date})
        chartData.filter(series => series.id === "minutes")[0].data.push({y:parseInt(game.minutes.split(":")[0])/10,x:game.date})
        if (countMap[game[propType]]) {
            countMap[game[propType]]++
        } else {
            countMap[game[propType]]=1
        }
        pointsArr.push(game[propType])
        if (game[propType] > max) {
            max = game[propType]
        }
        if (game[propType] < min || min === -1) {
            min = game[propType]
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
    let countData = [{id:propType, data: []}, {id:"normal", data: []}];
    let stddev = std(pointsArr)
    let m = mean(pointsArr)
    let lambda = variance(pointsArr)
    const normDist = new NormalDistribution(m, stddev);

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
            <div className="player-card">
                <h1>{player.first_name} {player.last_name}</h1>
            </div>
            <p>Avergage {propType}: {avgPoints}</p>
            <p>Variance: {lambda}</p>
            <p>Target: {target}</p>
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
                                <div>Date: {slice.points.filter((point) => point.serieId === propType)[0].data.x}</div>
                                <div>GameID: {sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index].gameID}</div>
                                <div>Opponent: {sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index].opponent}</div>
                                <div>Season: {sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index].season}</div>
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === "minutes"})[0].serieColor}`}}>
                                    <strong>minutes</strong> {sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index].minutes}
                                </div>
                                <div style={{color: `${slice.points.filter((point) => {return point.serieId === propType})[0].serieColor}`}}>
                                    <strong>{propType}</strong> {sortedGames[slice.points.filter((point) => point.serieId === propType)[0].index][propType]}
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
            </div>
        </div>
    )
}

export default Player
