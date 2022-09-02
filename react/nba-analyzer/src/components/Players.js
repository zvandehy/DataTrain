import React, { useState, useEffect, useCallback } from "react";
import Playercard from "./Playercard";
import DataListInput from "react-datalist-input";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faSortUp as asc,
  faSortDown as desc,
} from "@fortawesome/free-solid-svg-icons";
import {
  HOME_QUERY,
  FormatDate,
  ParseDate,
  GetSelectableTeams,
  AveragePropScore,
} from "../utils.js";
import { StatObjects } from "../predictions.js";
import { useQuery } from "@apollo/client";
import "../styles/players.css";

function GetSortMethod(method) {
  switch (method) {
    case "player":
      return (a, b) => {
        if (a.player.name < b.player.name) return -1;
        if (a.player.name > b.player.name) return 1;
        return 0;
      };
    case "position":
      return (a, b) => {
        if (a.player.position < b.player.position) return -1;
        if (a.player.position > b.player.position) return 1;
        return 0;
      };
    case "team":
      return (a, b) => {
        if (
          a.player.currentTeam.abbreviation < b.player.currentTeam.abbreviation
        )
          return -1;
        if (
          a.player.currentTeam.abbreviation > b.player.currentTeam.abbreviation
        )
          return 1;
        return 0;
      };
    case "opponent":
      return (a, b) => {
        if (a.opponent.abbreviation < b.opponent.abbreviation) return -1;
        if (a.opponent.abbreviation > b.opponent.abbreviation) return 1;
        return 0;
      };
    default:
      return (a, b) => {
        return (
          AveragePropScore(a.player.games, method) -
          AveragePropScore(b.player.games, method)
        );
      };
  }
}

const Players = (prop) => {
  const { client, league } = prop;
  const season = league === "nba" ? "2021-22" : "2022-23";
  const [teamLookup, setTeamLookup] = useState("");
  const [showPlayers, setShowPlayers] = useState([]);
  const [statPreference, setStatPreference] = useState("");
  const [positionLookup, setPositionLookup] = useState("");
  const [searchLookup, setSearchLookup] = useState("");
  const [sortMethod, setSortMethod] = useState("team");
  const [sortOrder, setSortOrder] = useState("DESC");
  const [seasonType, setSeasonType] = useState("");
  const [date, setDate] = useState(FormatDate(new Date()));
  const { loading, error, data, refetch } = useQuery(HOME_QUERY, {
    variables: { date: date, season: season },
    client: client,
  });
  useEffect(() => {
    // let team = localStorage.getItem("lookup");
    if (data) {
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
                item.propositions.filter((target) => {
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
      let filterBySearch =
        searchLookup !== ""
          ? filterCleaning.filter((projection) => {
              return (
                projection.player.name
                  .toLowerCase()
                  .search(searchLookup.toLowerCase()) !== -1 ||
                projection.player.currentTeam.abbreviation
                  .toLowerCase()
                  .search(searchLookup.toLowerCase()) !== -1 ||
                projection.opponent.abbreviation
                  .toLowerCase()
                  .search(searchLookup.toLowerCase()) !== -1 ||
                projection.opponent.name
                  .toLowerCase()
                  .search(searchLookup.toLowerCase()) !== -1 ||
                projection.player.currentTeam.name
                  .toLowerCase()
                  .search(searchLookup.toLowerCase()) !== -1
              );
            })
          : filterByPosition;

      let sorted = filterBySearch.sort((a, b) => {
        return GetSortMethod(sortMethod)(a, b);
      });

      sorted = sortOrder === "DESC" ? sorted.reverse() : sorted;
      setShowPlayers(sorted);
    }
  }, [
    data,
    teamLookup,
    statPreference,
    positionLookup,
    searchLookup,
    sortMethod,
    sortOrder,
  ]);

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
          setTeamLookup("");
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

  function onSelectSort(input) {
    setSortMethod(input.value);
  }

  if (loading) return "Loading...";
  if (error) {
    return `Error! ${error.message}. ${loading}. ${data}`;
  }

  const selectTeams = GetSelectableTeams(data.teams);
  let selectableStats = StatObjects.map((item) => ({
    key: item.label,
    label: item.label,
    value: item,
  }));
  selectableStats.unshift({ key: "ANY", label: "ANY", value: "" });

  let selectSorts = [
    { key: "team", label: "Team", value: "team" },
    { key: "opponent", label: "Opponent", value: "opponent" },
    { key: "player", label: "Player Name", value: "player" },
    { key: "position", label: "Position", value: "position" },
  ];
  selectSorts = selectSorts.concat(
    StatObjects.map((item) => {
      return { key: item.recognize, label: item.label, value: item.recognize };
    })
  );
  return (
    <div className="players">
      <div className="teams-dropdown">
        <a className="league-toggle" href={league === "nba" ? "/wnba" : "/nba"}>
          {league === "nba" ? "WNBA" : "NBA"}
        </a>
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
        <div className="sorter">
          <DataListInput
            placeholder="Sort players by"
            items={selectSorts}
            onSelect={onSelectSort}
            clearInputOnClick={true}
            suppressReselect={false}
            value={sortMethod.toUpperCase()}
          />
          <button
            onClick={() => setSortOrder(sortOrder === "DESC" ? "ASC" : "DESC")}
          >
            {sortOrder === "DESC" ? (
              <FontAwesomeIcon icon={asc} />
            ) : (
              <FontAwesomeIcon icon={desc} />
            )}
          </button>
        </div>
        <button
          onClick={() =>
            setSeasonType(
              seasonType === "REG"
                ? "PLAYOFFS"
                : seasonType === "PLAYOFFS"
                ? ""
                : "REG"
            )
          }
        >
          {seasonType === "" ? "ANY" : seasonType}
        </button>
        <input
          type={"text"}
          className={"right"}
          placeholder={"Search by keyword"}
          onChange={(e) => setSearchLookup(e.target.value)}
        />
        {/* TODO: Add Sort by projection value (instead of season average) & sort by prediction (overUnder / confidence) */}
      </div>
      <ul className="players-list">
        {showPlayers.length > 0 ? (
          showPlayers.map((item) => (
            <Playercard
              league={league}
              projection={item}
              player={item.player}
              date={date}
              key={item.player.playerID}
              statPreference={statPreference}
              seasonType={seasonType}
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
