import React, {useState, useCallback} from 'react'
import {useLocation} from 'react-router-dom';
import {GetPropScore} from '../utils.js'
import { ResponsiveLine } from '@nivo/line'
import DataListInput from "react-datalist-input";
import {std, mean} from 'mathjs'
import NormalDistribution from 'normal-distribution'

//TODO: see Apexcharts 
//https://apexcharts.com/react-chart-demos/mixed-charts/multiple-yaxis/
const Player = (props) => {
    // const {id} = useParams();
    let location = useLocation()
    const [propType, setPropType] = useState("points");
    const {playerProp} = location.state
    let selectPropTypes = playerProp.targets.map(item => ({
        // required: what to show to the user
        label: item.type,
        // required: key to identify the item within the array
        key: item.type.toLowerCase(),
      }));
    let target = playerProp.targets.find(item => item.type === propType)?.target
    let player=playerProp.player
    let seasonData=playerProp.player.games.filter((game) => game.season === "2021-22")
    let sortedGames = seasonData.sort(function(a, b) {
        var c = new Date(a.date);
        var d = new Date(b.date);
        return c-d;
    });
    let chartData = [{id:propType, data: []}, {id:"minutes", data: []}];
    let propScoreArr = []
    let countMap = {}
    let max = 0
    let min = -1;
    sortedGames.forEach((game) => {
        const propScore = GetPropScore(game,propType)
        chartData.filter(series => series.id === propType)[0].data.push({y:propScore,x:game.date})
        chartData.filter(series => series.id === "minutes")[0].data.push({y:parseInt(game.minutes.split(":")[0])/10,x:game.date})
        if (countMap[propScore]) {
            countMap[propScore]++
        } else {
            countMap[propScore]=1
        }
        propScoreArr.push(propScore)
        if (propScore > max) {
            max = propScore
        }
        if (propScore < min || min === -1) {
            min = propScore
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
    let stddev = std(propScoreArr)
    let m = mean(propScoreArr)
    const normDist = new NormalDistribution(m, stddev);

   for(let i = 0; i< 3*stddev+m;i++) {
       let xval = i
       let yval = normDist.pdf(xval);
        if (xval >= 0 ) {
            countData[1].data.push({x:xval, y:yval})
        }
        countData[0].data.push({x:xval, y:countMap[xval]? countMap[xval] /100: 0})
   }

//    const onSelectPropType = useCallback((selectedItem) => {
    //     setPropType(selectedItem.label)
    // });
   
    return (
        <div>
            <div className="player-card">
                <h1>{player.name}</h1>
            </div>
            <DataListInput
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
            <p>Under: {less}</p>
            <div className="line-chart" style={{height:"500px"}}>
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
            </div>
        </div>
    )
}

export default Player
