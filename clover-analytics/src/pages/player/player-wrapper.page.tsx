import { Button } from "@mui/material";
import moment from "moment";
import React, { useState } from "react";
import { useParams } from "react-router-dom";
import CustomModelDialog from "../../ components/custom-model-dialog/custom-model-dialog.component";
import { INITIAL_CUSTOM_MODEL_STATE } from "../../ components/custom-model-dialog/custom-model-dialog.reducer";
import PlayerPage from "../../ components/player-detailed/player/player.page";
import { useGetPlayerDetails } from "../../hooks/useGetPlayerDetail";
import { FindProjectionByDate } from "../../shared/functions/findProjection.fn";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
import {
  GameFilter,
  ProjectionFilter,
} from "../../shared/interfaces/graphql/filters.interface";
import { Minutes } from "../../shared/interfaces/stat.interface";
import "./player-wrapper.page.css";

const PlayerPageWrapper: React.FC = () => {
  const { id } = useParams();
  let playerID = id ? parseInt(id) : 0;
  // TODO: Ensure that date without projection/game works
  const [date, setDate] = useState(new Date());
  const [season, setSeason] = useState("2022-23");
  const [sportsbook, setSportsbook] = useState("");
  const [openCustomModel, setOpenCustomModel] = useState(false);
  const [customPredictionModel, setCustomPredictionModel] =
    useState<CustomCalculation>(INITIAL_CUSTOM_MODEL_STATE);
  const close = () => {
    setOpenCustomModel(false);
  };
  let projectionFilter: ProjectionFilter = {};
  if (sportsbook) {
    projectionFilter.sportsbook = sportsbook;
  }
  const predictionFilter: GameFilter = {
    season: season,
    endDate: moment(date).format("YYYY-MM-DD"),
  };
  const gameFilter: GameFilter = {
    endDate: moment(date).format("YYYY-MM-DD"),
    statFilters: [{ stat: Minutes, min: 10 }],
  };

  const {
    loading,
    error,
    data: player,
  } = useGetPlayerDetails({
    playerID: playerID,
    predictionFilter,
    gameFilter,
    customModel: customPredictionModel,
  });

  if (loading) return <div>{loading}</div>;
  if (error) return <div>{error}</div>;
  if (!player) return <div>No player found</div>;
  const projection = FindProjectionByDate(date, player.projections, player);
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
        setCustomModel={setCustomPredictionModel}
      />
      <PlayerPage
        player={player}
        selectedProjection={projection}
        setSelectedDate={setDate}
        gameFilter={gameFilter}
        customModel={customPredictionModel}
      />
    </div>
  );
};

export default PlayerPageWrapper;
