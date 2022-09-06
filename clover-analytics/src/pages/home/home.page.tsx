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
import { Button, TextField } from "@mui/material";
import { INITIAL_CUSTOM_MODEL_STATE } from "../../ components/custom-model-dialog/custom-model-dialog.reducer";
import moment from "moment";
import CustomModelDialog from "../../ components/custom-model-dialog/custom-model-dialog.component";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";

const Home: React.FC = () => {
  const [date, setDate] = useState(new Date());
  const [season, setSeason] = useState("2022-23");
  const [sportsbook, setSportsbook] = useState("");
  const [openCustomModel, setOpenCustomModel] = useState(false);
  const [customPredictionModel, setCustomPredictionModel] =
    useState<CustomCalculation>(INITIAL_CUSTOM_MODEL_STATE);

  let projectionFilter: ProjectionFilter = {
    startDate: moment(date).format("YYYY-MM-DD"),
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  if (sportsbook) {
    projectionFilter.sportsbook = sportsbook;
  }
  const gameFilter: GameFilter = {
    endDate: moment(date).format("YYYY-MM-DD"),
    // statFilters: [{ stat: Minutes, min: 10 }],
  };
  const predictionFilter: GameFilter = {
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
    predictionFilter,
    customModel: customPredictionModel,
  });

  const close = () => {
    setOpenCustomModel(false);
  };
  //COMPONENT

  return (
    <div id="home-page">
      {/* Move to own component */}
      {/* Query Filters trigger a new query */}
      <div id="query-filters" className={"filters-wrapper"}>
        <Button
          variant={"outlined"}
          onClick={() => setOpenCustomModel(true)}
          sx={{ marginRight: "1rem" }}
        >
          Custom Model
        </Button>
        <CustomModelDialog
          open={openCustomModel}
          closeDialog={close}
          setCustomModel={setCustomPredictionModel}
        />

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
          gameFilter={gameFilter}
        />
      ) : (
        <></>
      )}
    </div>
  );
};

export default Home;
