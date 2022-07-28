import React, { useState } from "react";
import PlayerCardList from "../../ components/playercard-list/playercard-list.component";
import { FormatDate } from "../../shared/functions/dates.fn";
import "./home.page.css";

import {
  GameFilter,
  ProjectionFilter,
} from "../../shared/interfaces/graphql/filters.interface";
import { Option } from "../../shared/interfaces/option.interface";
import AutocompleteFilter from "../../ components/autocomplete-filter/autocomplete-filter.component";
import { useGetProjections } from "../../hooks/useGetProjections";
import { DesktopDatePicker } from "@mui/x-date-pickers/DesktopDatePicker";
import { LocalizationProvider } from "@mui/x-date-pickers/LocalizationProvider";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import { TextField } from "@mui/material";

const Home: React.FC = () => {
  const [date, setDate] = useState(new Date(new Date().toISOString()));
  const [season, setSeason] = useState("2022-23");
  const [sportsbook, setSportsbook] = useState("");
  const [similarPlayersToggle, toggleSimilarPlayers] = useState(false);
  const [similarTeamsToggle, toggleSimilarTeams] = useState(false);

  // if (FormatDate(date) !== "2022-07-24") {
  //   console.log("date", date);
  //   setDate("2022-07-24");
  // }

  let projectionFilter: ProjectionFilter = {
    startDate: FormatDate(date),
    endDate: FormatDate(date),
  };
  if (sportsbook) {
    projectionFilter.sportsbook = sportsbook;
  }
  const gameFilter: GameFilter = { season: season };
  console.log(projectionFilter, gameFilter);
  //SEASONS
  const seasons: Option<string>[] = [
    { label: "2022-23 (Current)", id: "2022-23" },
    { label: "2021-22", id: "2021-22" },
  ];
  const onSeasonChange = (value: string) => {
    setSeason(value);
  };
  const onDateChange = (newValue: Date | null) => {
    console.log(FormatDate(newValue!));
    setDate(newValue || new Date());
  };

  let result = useGetProjections({
    projectionFilter,
    gameFilter,
    predictionFilter: gameFilter,
  });

  //COMPONENT

  return (
    <div id="home-page">
      {/* Move to own component */}
      {/* Query Filters trigger a new query */}
      <div id="query-filters" className={"filters-wrapper"}>
        <AutocompleteFilter
          options={[
            { label: "PrizePicks", id: "PrizePicks" },
            { label: "Underdog", id: "UnderdogFantasy" },
          ]}
          onChange={setSportsbook}
          label="Sportsbook"
          width={180}
        />
        <AutocompleteFilter
          options={seasons}
          onChange={onSeasonChange}
          label="Season"
          width={160}
        />
        <LocalizationProvider dateAdapter={AdapterDayjs}>
          <DesktopDatePicker
            label="Date desktop"
            inputFormat="MM-DD-YYYY"
            value={date}
            onChange={onDateChange}
            PaperProps={{
              style: {
                color: "black",
              },
            }}
            renderInput={(params) => <TextField {...params} />}
          />
        </LocalizationProvider>
      </div>
      {result.data ? <PlayerCardList projectionQueryResult={result} /> : <></>}
    </div>
  );
};

export default Home;
