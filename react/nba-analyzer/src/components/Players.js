import React, {useState, useEffect, useCallback} from 'react'
import Playercard  from './Playercard'
import DataListInput from "react-datalist-input";
import {HOME_QUERY} from '../utils.js'
import { useQuery } from '@apollo/client';
import "../styles/players.css"

const Players = () => {
    const [lookup, setLookup] = useState('');
    const [showPlayers, setShowPlayers] = useState([]);
    const { loading, error, data } = useQuery(HOME_QUERY);
    useEffect(() => {
        let team = localStorage.getItem('lookup');
        if (data) { 
            setLookup(team)
            console.group("Data updated")
            console.log(data)
            let filterCleaning = data.projections.filter(item => item.player.playerID !== 0)
            console.log("filterCleaning", filterCleaning.length)
            let filteredByTeam = lookup ? filterCleaning.filter(item => item.player.currentTeam.abbreviation === lookup) : filterCleaning;
            console.log(lookup, filteredByTeam)
            // let p = filteredByPropType.map((prizepick) => {let prop = showPlayerPreview(prizepick.player, data.prizepicks, proptype); if (prop) return {player:player, prop:prop}}).filter((item) => item !== undefined);
            // p.sort((a, b) => a.prop.target > b.prop.target)
            // filteredByTeam.sort((a, b) => a.target > b.target)
            setShowPlayers(filteredByTeam)
            console.groupEnd()
        }
    },
    [data, lookup]
    );

    const onSelectTeam = useCallback((selectedItem) => {
        let selected = lookup === selectedItem.label ? "" : selectedItem.label
        localStorage.setItem('lookup', selected)
        if (data) {
            setLookup(selected)
        }
      }, [data, lookup]);
      
    if (loading) return 'Loading...';
    if (error) {
        return `Error! ${error.message}. ${loading}. ${data}`;
    }
    console.log(data)
    const selectTeams =
          data.teams.map((team) => ({
            // required: what to show to the user
            label: team.abbreviation,
            // required: key to identify the item within the array
            key: team.teamID,
          }));

    // console.log(showPlayers)
    
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
        </div>
            <ul className="players-list">
                {showPlayers.length > 0 ? showPlayers.map((item) => <Playercard playerProp={item} key={item.player.playerID}/>) : <li>No Players to Show</li>}
                {/* <Playercard playerProp={showPlayers[0]} key={showPlayers[0].player.playerID}/> */}
                {/* <Playercard confidence={55}/>
                <Playercard confidence={45}/>
                <Playercard/>
                <Playercard/> */}
            </ul>
        </div>
    )
}

export default Players
