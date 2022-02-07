import React, {useState, useEffect, useCallback} from 'react'
import PlayerPreview  from './PlayerPreview'
import DataListInput from "react-datalist-input";
import {HOME_QUERY, showPlayerPreview} from '../utils.js'
import { useQuery } from '@apollo/client';

//TODO: Current performance of the page is really bad....
const Players = () => {
    const [lookup, setLookup] = useState('');
    const [showPlayers, setShowPlayers] = useState([]);
    const [propType, setPropType] = useState('Points');
    const { loading, error, data } = useQuery(HOME_QUERY);

    useEffect(() => {
        let team = localStorage.getItem('lookup')
        let type = localStorage.getItem('propType')
        if (data) { 
            setLookup(team)
            setPropType(type)
            console.group("Data updated")
            console.log(data)
            let filterCleaning = data.prizepicks.filter(item => item.player.playerID !== 0)
            console.log(filterCleaning.length)
            let filteredByTeam = lookup !== '' ? filterCleaning.filter(item => item.player.teamABR === lookup) : filterCleaning;
            console.log(lookup, filteredByTeam)
            let filteredByPropType = propType !== '' ? filteredByTeam.filter(item => item.type === propType) : filteredByTeam;
            // let p = filteredByPropType.map((prizepick) => {let prop = showPlayerPreview(prizepick.player, data.prizepicks, proptype); if (prop) return {player:player, prop:prop}}).filter((item) => item !== undefined);
            // p.sort((a, b) => a.prop.target > b.prop.target)
            filteredByPropType.sort((a, b) => a.target > b.target)
            console.log(propType, filteredByPropType)
            setShowPlayers(filteredByPropType)
            console.groupEnd()
        }
    },
    [data, propType, lookup]
    );

    const onSelectTeam = useCallback((selectedItem) => {
        let selected = lookup === selectedItem.label ? "" : selectedItem.label
        localStorage.setItem('lookup', selected)
        if (data) {
            setLookup(selected)
        }
      }, [data, lookup]);

      const onSelectProp = useCallback((selectedItem) => {
        localStorage.setItem('propType', selectedItem.label)
        if (data) { 
            setPropType(selectedItem.label)
        }
      }, [data]);
      
    if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}. ${loading}. ${data}`;
    console.log(data)
    const selectTeams =
          data.teams.map((team) => ({
            // required: what to show to the user
            label: team.abbreviation,
            // required: key to identify the item within the array
            key: team.teamID,
          }));

    
    return (
        <div className="players">
           <div className="teams-dropdown">
             <DataListInput
                placeholder="Select a team"
                items={selectTeams}
                onSelect={onSelectTeam}
                clearInputOnClick={true}
                suppressReselect={false}
            />
             <DataListInput
                placeholder="Select a propType"
                items={[{label:"Points", key:"Points"}, {label: "Assists", key: "Assists"}]}
                onSelect={onSelectProp}
                clearInputOnClick={true}
            />
        </div>
            <table>
                <thead>
                    <tr><th>Player</th><th>Target</th><th>Mean</th><th>Over</th><th>Under</th><th>Median</th><th>MAD</th></tr>
                </thead>
                <tbody>
                    {showPlayers.length > 0 ? showPlayers.map((item) => {return (<PlayerPreview playerProp={item} key={item.player.playerID}/>)}) : <tr><td>No Players to Show</td></tr>}
                </tbody>
            </table>
        </div>
    )
}

export default Players
