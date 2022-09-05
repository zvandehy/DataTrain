import moment from "moment";
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import AutocompleteFilter from "../../ components/autocomplete-filter/autocomplete-filter.component";
import PlayerContext from "../../ components/player-detailed/player-context/player-context.component";
import PlayerProfileChart from "../../ components/player-detailed/player-profile/player-profile.component";
import PlayerPage from "../../ components/player-detailed/player/player.page";
import PlayerStatsPreview from "../../ components/player-stats-table/player-stats-preview/player-stats-preview.component";
import StatSelectButtons from "../../ components/playercard-list/playercard/stat-select-buttons/stat-select-buttons.component";
import { useGetPlayerDetails } from "../../hooks/useGetPlayerDetail";
import { CompareDates } from "../../shared/functions/dates.fn";
import { FindProjectionByDate } from "../../shared/functions/findProjection.fn";
import {
  GetMaxConfidence,
  UpdatePropositionWithPrediction,
} from "../../shared/functions/predictions.fn";
import { GetImpliedTarget } from "../../shared/functions/target.fn";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
import {
  GameFilter,
  ProjectionFilter,
} from "../../shared/interfaces/graphql/filters.interface";
import { Game } from "../../shared/interfaces/graphql/game.interface";
import { Player } from "../../shared/interfaces/graphql/player.interface";
import {
  Projection,
  Proposition,
} from "../../shared/interfaces/graphql/projection.interface";
import { Team } from "../../shared/interfaces/graphql/team.interface";
import { Option } from "../../shared/interfaces/option.interface";
import { Minutes, Stat } from "../../shared/interfaces/stat.interface";
import "./player-wrapper.page.css";

function getProposition(
  projection: Projection | undefined,
  stat: Stat | undefined
): Proposition | undefined {
  if (projection) {
    if (stat) {
      let statProp = projection.propositions.find((p) => p.statType === stat);
      if (statProp) {
        return statProp;
      }
    }
    return GetMaxConfidence(projection.propositions);
  }
  return undefined;
}

const PlayerPageWrapper: React.FC = () => {
  const { id } = useParams();
  let playerID = id ? parseInt(id) : 0;
  // TODO: Ensure that date without projection/game works
  const [date, setDate] = useState(new Date());
  const [season, setSeason] = useState("2022-23");
  const [sportsbook, setSportsbook] = useState("");
  const [similarPlayersToggle, toggleSimilarPlayers] = useState(true);
  const [similarTeamsToggle, toggleSimilarTeams] = useState(true);
  const [customPredictionModel, setCustomPredictionModel] =
    useState<CustomCalculation>({
      includePush: true,
      includeOnDifferentTeam: true,
      recency: [
        { count: 0, weight: 0.2 },
        { count: -20, weight: 0.1 },
        { count: -10, weight: 0.1 },
        { count: -5, weight: 0.12 },
      ],
      similarPlayers: { count: 10, weight: 0.12 },
      similarTeams: { count: 3, weight: 0.15 },
      opponentWeight: 0.21,
      // homeAwayWeight:0.1,
    });

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
  //SEASONS
  const seasons: Option<string>[] = [
    { label: "2022-23 (Current)", id: "2022-23" },
    { label: "2021-22", id: "2021-22" },
  ];
  const onSeasonChange = (value: string) => {
    setSeason(value);
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
