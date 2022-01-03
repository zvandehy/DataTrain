import React from 'react'
import {useParams} from 'react-router-dom';
import { useQuery } from '@apollo/client';
import {average, PLAYERGAMES_QUERY} from '../utils.js'
import { ResponsiveLine } from '@nivo/line'

const Player = () => {
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
    //points data
    sortedGames.forEach((game) => {
        chartData.filter(series => series.id === "points")[0].data.push({y:game.points,x:game.date})
        chartData.filter(series => series.id === "minutes")[0].data.push({y:parseInt(game.minutes.split(":")[0])/10,x:game.date})
    })
    return (
        <div>
            <div className="player-card">
                <h1>{player.first_name} {player.last_name}</h1>
            </div>
            <p>Points: {avgPoints}</p>
            <div className="line-chart" style={{height:"500px"}}>
                <ResponsiveLine
                    data={chartData}
                    margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
                    xScale={{ type: 'point' }}
                    axisBottom={{tickRotation:-90, tickSize:5, tickPadding: 0, format: function(value) {return `${new Date(value).getMonth()+1}/${new Date(value).getDate()+1}`}}}
                    yScale={{ type: 'linear', min: 0, max: 'auto', stacked: false, reverse: false }}
                    yFormat=" >-.2f"
                    markers={[{axis:"y", value: avgPoints, legend: "average"}, {axis:"y", value: 10.0, legend: "target"}]}
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
                                <div style={{color: `${slice.points.filter((point) => {console.log(point); return point.serieId === "minutes"})[0].serieColor}`}}>
                                    <strong>minutes</strong> {sortedGames[slice.points.filter((point) => point.serieId === "points")[0].index].minutes}
                                </div>
                                {slice.points.map(point => (
                                    <div
                                        key={point.id}
                                        style={{
                                            color: point.serieColor,
                                            padding: '3px 0',
                                        }}
                                    >
                                        <strong>{point.serieId}</strong> {point.data.yFormatted}
                                    </div>
                                ))}
                            </div>
                        )
                    }}
                />
            </div>
        </div>
    )
}

export default Player
