import React, { useState, useEffect, useCallback } from "react";
import Playercard from "./Playercard";
import DataListInput from "react-datalist-input";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import {
  HOME_QUERY,
  FormatDate,
  ParseDate,
  GetSelectableTeams,
} from "../utils.js";
import { useQuery } from "@apollo/client";
import "../styles/players.css";

const Players = () => {
  const [teamLookup, setTeamLookup] = useState("ANY");
  const [showPlayers, setShowPlayers] = useState([]);
  const [date, setDate] = useState(FormatDate(new Date()));
  const { loading, error, data, refetch } = useQuery(HOME_QUERY, {
    variables: { date: date },
  });
  useEffect(() => {
    let team = localStorage.getItem("lookup");
    if (data) {
      setTeamLookup(team);
      let filterCleaning = data.projections.filter(
        (item) => item.player.playerID !== 0
      );
      let filteredByTeam = teamLookup
        ? filterCleaning.filter(
            (item) => item.player.currentTeam.abbreviation === teamLookup
          )
        : filterCleaning;
      setShowPlayers(filteredByTeam);
    }
  }, [data, teamLookup]);

  function changeDate(date) {
    date = FormatDate(date);
    setDate(date);
    refetch({ date: date });
  }

  const onSelectTeam = useCallback(
    (selectedItem) => {
      if (data) {
        let selected = selectedItem.label;
        if (selected === "ANY" || selected === teamLookup) {
          selected = "";
          setTeamLookup("ANY");
        } else {
          setTeamLookup(selected);
        }
        localStorage.setItem("lookup", selected);
      }
    },
    [data, teamLookup]
  );

  if (loading) return "Loading...";
  if (error) {
    return `Error! ${error.message}. ${loading}. ${data}`;
  }
  console.log(data);

  const selectTeams = GetSelectableTeams(data.teams);

  return (
    <div className="players">
      <div className="teams-dropdown">
        <DataListInput
          placeholder="Select a team"
          items={selectTeams}
          onSelect={onSelectTeam}
          clearInputOnClick={true}
          suppressReselect={false}
          value={teamLookup ?? "ANY"}
        />
        <DatePicker
          selected={ParseDate(date)}
          onChange={(date) => changeDate(date)}
        />
      </div>
      <ul className="players-list">
        {showPlayers.length > 0 ? (
          showPlayers.map((item) => (
            <Playercard
              projection={item}
              player={item.player}
              date={date}
              key={item.player.playerID}
            />
          ))
        ) : (
          <li>No Players to Show</li>
        )}
      </ul>
    </div>
  );
};

export default Players;
