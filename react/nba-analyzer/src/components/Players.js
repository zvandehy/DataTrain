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
import { StatObjects } from "../predictions.js";
import { useQuery } from "@apollo/client";
import "../styles/players.css";

const Players = () => {
  const [teamLookup, setTeamLookup] = useState("ANY");
  const [showPlayers, setShowPlayers] = useState([]);
  const [statPreference, setStatPreference] = useState("");
  const [positionLookup, setPositionLookup] = useState("");
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
      let filterByStat =
        statPreference !== ""
          ? filteredByTeam.filter(
              (item) =>
                item.targets.filter((target) => {
                  return (
                    target?.type.toLowerCase() === statPreference.recognize &&
                    target?.target
                  );
                }).length > 0
            )
          : filteredByTeam;
      let filterByPosition =
        positionLookup !== "ANY" && positionLookup !== ""
          ? filterByStat.filter(
              (item) => item.player.position === positionLookup
            )
          : filterByStat;
      setShowPlayers(filterByPosition);
    }
  }, [data, teamLookup, statPreference, positionLookup]);

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

  function onSelectStatPreference(input) {
    StatObjects.forEach((item) => {
      if (item.recognize === input.recognize) {
        setStatPreference(item);
      }
    });
    setStatPreference(input.value);
  }

  function onSelectPosition(input) {
    if (input.label === "ANY") {
      setPositionLookup("");
    } else {
      setPositionLookup(input.value);
    }
  }

  if (loading) return "Loading...";
  if (error) {
    return `Error! ${error.message}. ${loading}. ${data}`;
  }
  console.log(data);

  const selectTeams = GetSelectableTeams(data.teams);
  let selectableStats = StatObjects.map((item) => ({
    key: item.label,
    label: item.label,
    value: item,
  }));
  selectableStats.unshift({ key: "ANY", label: "ANY", value: "" });
  return (
    <div className="players">
      <div className="teams-dropdown">
        <DataListInput
          placeholder="Filter by team"
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
        <DataListInput
          key="stat-preference"
          placeholder="Filter by stat"
          items={selectableStats}
          onSelect={onSelectStatPreference}
          clearInputOnClick={true}
          suppressReselect={false}
          value={
            statPreference !== ""
              ? StatObjects.find(
                  (item) => item.recognize === statPreference.recognize
                ).label
              : ""
          }
        />
        <DataListInput
          placeholder="Filter by position"
          items={[
            { key: "ANY", label: "ANY", value: "" },
            { key: "G", label: "Guard", value: "G" },
            { key: "F", label: "Forward", value: "F" },
            { key: "C", label: "Center", value: "C" },
          ]}
          onSelect={onSelectPosition}
          clearInputOnClick={true}
          suppressReselect={false}
          value={positionLookup}
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
              statPreference={statPreference}
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
