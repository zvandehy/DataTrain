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
import { Stat } from "../../shared/interfaces/stat.interface";
import "./player-wrapper.page.css";

// const [date, setDate] = useState(FormatDate(new Date()));
// //TODO: Handle error when pick game that isn't in seasonType
// const [seasonType, setSeasonType] = useState("");

function FormatDate(date: Date): string {
  return moment(date).format("YYYY-MM-DD");
}

function getProjection(
  date: Date,
  projections: Projection[]
): Projection | undefined {
  return projections.find((projection) => projection.date === FormatDate(date));
}

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
  const [date, setDate] = useState(new Date());
  const [projection, setProjection] = useState(getProjection(date, []));
  const [statType, setStatType] = useState(undefined as Stat | undefined);
  const [proposition, setProposition] = useState(
    getProposition(projection, statType)
  );
  // const [date, setDate] = useState(moment(new Date()).add(1, "days").toDate()); //Ensure that date without projection/game works
  const [season, setSeason] = useState("2022-23");
  const [sportsbook, setSportsbook] = useState("");
  const [similarPlayersToggle, toggleSimilarPlayers] = useState(true);
  const [similarTeamsToggle, toggleSimilarTeams] = useState(true);
  const [customPredictionModel, setCustomPredictionModel] =
    useState<CustomCalculation>({
      recency: [
        { count: 5, weight: 0.2 },
        { count: 15, weight: 0.2 },
        { count: 0, weight: 0.2 },
      ],
      similarPlayers: { count: 10, weight: 0.2 },
      similarTeams: { count: 3, weight: 0.2 },
      includePush: true,
    });

  const onStatSelect = (stat: Stat) => {
    if (projection) {
      setStatType(stat);
      let customTarget = GetImpliedTarget(projection, stat);
      let customProp: Proposition = {
        target: customTarget || 0,
        statType: stat,
        type: stat.label,
        sportsbook: customTarget ? "Implied" : "None",
        predictions: [],
        lastModified: new Date(),
        customPrediction: {
          model: "Custom",
          overUnderPrediction: "",
          confidence: 0,
          totalPrediction: 0,
          predictionFragments: [],
        },
      };
      const foundProp =
        projection.propositions.find((p) => p.statType === stat) ||
        UpdatePropositionWithPrediction(
          customProp,
          projection.player.games,
          projection,
          customPredictionModel
        );
      setProposition(foundProp);
    }
  };
  let projectionFilter: ProjectionFilter = {};
  if (sportsbook) {
    projectionFilter.sportsbook = sportsbook;
  }
  const gameFilter: GameFilter = {
    season: season,
    // endDate: moment(date).format("YYYY-MM-DD"),
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
    predictionFilter: gameFilter,
  });

  if (loading) return <div>{loading}</div>;
  if (error) return <div>{error}</div>;
  if (!player) return <div>No player found</div>;
  return (
    <div className="player-wrapper">
      <PlayerPage
        player={player}
        selectedDate={date}
        setSelectedDate={setDate}
      />
    </div>
  );
};

export default PlayerPageWrapper;
