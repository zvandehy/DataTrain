import { Button } from "@mui/material";
import moment from "moment";
import React, { useState } from "react";
import { useParams } from "react-router-dom";
import CustomModelDialog from "../../components/custom-model-dialog/custom-model-dialog.component";
import { INITIAL_CUSTOM_MODEL_STATE } from "../../components/custom-model-dialog/custom-model-dialog.reducer";
import PlayerPage from "../../components/player-detailed/player/player.page";
import { useGetPlayerDetails } from "../../hooks/useGetPlayerDetail";
import { DEFAULT_MODEL } from "../../shared/constants";
import { FindProjectionByDate } from "../../shared/functions/findProjection.fn";
import {
  CustomCalculation,
  ModelInput,
} from "../../shared/interfaces/custom-prediction.interface";
import {
  GameFilter,
  ProjectionFilter,
  SeasonOption,
} from "../../shared/interfaces/graphql/filters.interface";
import "./player-wrapper.page.css";

const PlayerPageWrapper: React.FC = () => {
  const { id } = useParams();
  let playerID = id ? parseInt(id) : 0;
  // TODO: Ensure that date without projection/game works
  const [date, setDate] = useState(
    localStorage.getObject("date") ?? new Date()
  );
  const onDateChange = (newValue: Date | null) => {
    const newDate: Date = newValue || new Date();
    setDate(newDate);
    localStorage.setObject("date", newDate);
  };
  const [season, setSeason] = useState(SeasonOption.SEASON_2022_23);
  const [sportsbook, setSportsbook] = useState("");
  const [openCustomModel, setOpenCustomModel] = useState(false);
  const [customPredictionModel, setCustomPredictionModel] =
    useState<ModelInput>(
      localStorage.getObject("customModel") ?? DEFAULT_MODEL
    );
  const close = () => {
    setOpenCustomModel(false);
  };
  const save = (value: ModelInput) => {
    localStorage.setObject("customModel", value);
    setCustomPredictionModel(value);
    close();
  };
  let projectionFilter: ProjectionFilter = {};
  if (sportsbook) {
    projectionFilter.sportsbook = sportsbook;
  }
  const predictionFilter: GameFilter = {
    seasons: [season],
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  const gameFilter: GameFilter = {
    endDate: moment(date).format("YYYY-MM-DD"),
    // statFilters: [{ stat: Minutes, min: 10 }],
  };

  // const {
  //   loading,
  //   error,
  //   data: player,
  // } = useGetPlayerDetails({
  //   playerID: playerID,
  //   predictionFilter,
  //   gameFilter,
  //   customModel: customPredictionModel,
  // });

  // if (loading) return <div>{loading}</div>;
  // if (error) return <div>{error}</div>;
  // if (!player) return <div>No player found</div>;

  return (
    <div className="player-wrapper">
      <Button
        variant={"outlined"}
        sx={{
          color: "white",
          borderColor: "white",
          "&:hover": { borderColor: "white" },
        }}
        onClick={() => setOpenCustomModel(true)}
      >
        Custom Model
      </Button>
      <CustomModelDialog
        open={openCustomModel}
        closeDialog={close}
        setCustomModel={save}
      />
      {/* <PlayerPage
        player={player}
        projection={FindProjectionByDate(date, player.projections, player)}
        setSelectedDate={onDateChange}
        gameFilter={gameFilter}
        customModel={customPredictionModel}
      /> */}
    </div>
  );
};

export default PlayerPageWrapper;
