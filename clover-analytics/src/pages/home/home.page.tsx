import React, { useState } from "react";
import PlayerCardList from "../../ components/playercard-list/playercard-list.component";
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
import { AdapterMoment } from "@mui/x-date-pickers/AdapterMoment";
import { Checkbox, FormControlLabel, TextField } from "@mui/material";
import { Label } from "@mui/icons-material";
import moment from "moment";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";

const Home: React.FC = () => {
  const [date, setDate] = useState(new Date());
  const [season, setSeason] = useState("2022-23");
  const [sportsbook, setSportsbook] = useState("");
  const [similarPlayersToggle, toggleSimilarPlayers] = useState(true);
  const [similarTeamsToggle, toggleSimilarTeams] = useState(true);
  const [customPredictionModel, setCustomPredictionModel] =
    useState<CustomCalculation>({
      recency: [
        { count: 0, weight: 0.25 },
        { count: -15, weight: 0.1 },
        { count: -5, weight: 0.1 },
      ],
      similarPlayers: { count: 10, weight: 0.2 },
      similarTeams: { count: 3, weight: 0.35 },
      includePush: true,
    });

  let projectionFilter: ProjectionFilter = {
    startDate: moment(date).format("YYYY-MM-DD"),
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  if (sportsbook) {
    projectionFilter.sportsbook = sportsbook;
  }
  const gameFilter: GameFilter = {
    season: season,
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  //SEASONS
  const seasons: Option<string>[] = [
    { label: "2022-23 (Current)", id: "2022-23" },
    { label: "2021-22", id: "2021-22" },
  ];
  const onSeasonChange = (value: string) => {
    setSeason(value);
  };
  const onDateChange = (newValue: Date | null) => {
    setDate(newValue || new Date());
  };

  let result = useGetProjections({
    projectionFilter,
    gameFilter,
    predictionFilter: gameFilter,
    similarPlayers: similarPlayersToggle,
    similarTeams: similarTeamsToggle,
    customModel: customPredictionModel,
  });

  //COMPONENT

  return (
    <div id="home-page">
      {/* Move to own component */}
      {/* Query Filters trigger a new query */}
      <div id="query-filters" className={"filters-wrapper"}>
        <div id="checkboxes">
          <FormControlLabel
            control={
              <Checkbox
                size="small"
                name="SimilarPlayers"
                value={similarPlayersToggle}
                sx={{
                  "&.Mui-checked": {
                    color: "var(--color-accent)",
                  },
                }}
                onChange={() => toggleSimilarPlayers(!similarPlayersToggle)}
              />
            }
            label="Similar Players"
          />
          <FormControlLabel
            control={
              <Checkbox
                size="small"
                name="SimilarTeams"
                value={similarTeamsToggle}
                sx={{
                  "&.Mui-checked": {
                    color: "var(--color-accent)",
                  },
                }}
                onChange={() => toggleSimilarTeams(!similarTeamsToggle)}
              />
            }
            label="Similar Teams"
          />
        </div>

        <AutocompleteFilter
          options={[
            { label: "PrizePicks", id: "PrizePicks" },
            { label: "Underdog", id: "UnderdogFantasy" },
          ]}
          onChange={setSportsbook}
          label="Sportsbook"
        />
        <AutocompleteFilter
          options={seasons}
          onChange={onSeasonChange}
          label="Season"
        />
        <LocalizationProvider dateAdapter={AdapterMoment}>
          <DesktopDatePicker
            label="Date"
            inputFormat="M-D-YY"
            value={date}
            onChange={onDateChange}
            PaperProps={{
              style: {
                color: "black",
              },
              sx: {
                "&input": {
                  "max-width": "180px",
                },
              },
            }}
            renderInput={(params) => <TextField {...params} />}
          />
        </LocalizationProvider>
      </div>
      {result.data ? (
        <PlayerCardList
          projectionQueryResult={result}
          customModel={customPredictionModel}
        />
      ) : (
        <></>
      )}
    </div>
  );
};

export default Home;
