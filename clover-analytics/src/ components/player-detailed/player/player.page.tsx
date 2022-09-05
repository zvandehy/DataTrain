import moment from "moment";
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  GetMaxConfidence,
  UpdatePropositionWithPrediction,
} from "../../../shared/functions/predictions.fn";
import { GetImpliedTarget } from "../../../shared/functions/target.fn";
import {
  GameFilter,
  ProjectionFilter,
} from "../../../shared/interfaces/graphql/filters.interface";
import { Player } from "../../../shared/interfaces/graphql/player.interface";
import {
  Projection,
  Proposition,
} from "../../../shared/interfaces/graphql/projection.interface";
import { Points, Stat } from "../../../shared/interfaces/stat.interface";
import { Option } from "../../../shared/interfaces/option.interface";
import "./player.page.css";
import AutocompleteFilter from "../../autocomplete-filter/autocomplete-filter.component";
import PlayerContext from "../player-context/player-context.component";
import StatSelectButtons from "../../playercard-list/playercard/stat-select-buttons/stat-select-buttons.component";
import PlayerStatsPreview from "../../player-stats-table/player-stats-preview/player-stats-preview.component";
import Prediction from "../../prediction/prediction.component";
import PlayerStatsChart from "../player-stats-chart/player-stats-chart.component";
import { CustomCalculation } from "../../../shared/interfaces/custom-prediction.interface";
import { FindProjectionByDate } from "../../../shared/functions/findProjection.fn";

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

interface PlayerPageProps {
  player: Player;
  selectedProjection: Projection;
  setSelectedDate: (date: Date) => void;
  gameFilter: GameFilter;
  customModel: CustomCalculation;
}
let league = "wnba";

const PlayerPage: React.FC<PlayerPageProps> = ({
  player,
  selectedProjection,
  setSelectedDate,
  gameFilter,
  customModel,
}: PlayerPageProps) => {
  const [statType, setStatType] = useState(undefined as Stat | undefined);
  const [proposition, setProposition] = useState(
    getProposition(selectedProjection, statType)
  );
  // const [date, setDate] = useState(moment(new Date()).add(1, "days").toDate()); //Ensure that date without projection/game works
  const [season, setSeason] = useState("2022-23");
  const [sportsbook, setSportsbook] = useState("");
  const [similarPlayersToggle, toggleSimilarPlayers] = useState(false);
  const [similarTeamsToggle, toggleSimilarTeams] = useState(false);

  const onStatSelect = (stat: Stat) => {
    setStatType(stat);
  };
  //SEASONS
  const seasons: Option<string>[] = [
    { label: "2022-23 (Current)", id: "2022-23" },
    { label: "2021-22", id: "2021-22" },
  ];
  const onSeasonChange = (value: string) => {
    setSeason(value);
  };

  useEffect(() => {
    if (player) {
      if (statType) {
        let customTarget = GetImpliedTarget(selectedProjection, statType);
        let customProp: Proposition = {
          target: customTarget || 0,
          statType: statType,
          type: statType.label,
          sportsbook: customTarget ? "Implied" : "None",
          lastModified: new Date(),
          predictions: [],
          customPrediction: {
            model: "Custom",
            overUnderPrediction: "",
            confidence: 0,
            totalPrediction: 0,
            recencyFragments: [],
          },
        };
        const foundProp =
          selectedProjection.propositions.find(
            (p) => p.statType === statType
          ) ||
          UpdatePropositionWithPrediction(
            customProp,
            gameFilter,
            selectedProjection,
            customModel
          );
        setProposition(foundProp);
      }
    }
  }, [customModel, gameFilter, player, selectedProjection, statType]);

  let filteredGames = player.games.filter((game) => {
    return moment(game.date).isBefore(selectedProjection.date);
  });
  let gameOptions = player.games.map((game) => {
    return {
      id: game.date,
      label: `${game.date} vs ${game.opponent.name}`,
    };
  });
  let missingProjections = player.projections.filter(
    (projection) =>
      !gameOptions.some((gameOption) =>
        moment(gameOption.id).isSame(projection.date)
      )
  );
  gameOptions = [
    ...gameOptions,
    ...missingProjections.map((projection) => {
      return {
        id: projection.date,
        label: `${projection.date} vs ${projection.opponent.name}`,
      };
    }),
  ];
  gameOptions = gameOptions.sort((a, b) => {
    return moment(b.id).unix() - moment(a.id).unix();
  });
  return (
    <div className="player-page">
      <div id={"games-dropdown"}>
        <AutocompleteFilter
          label="Game"
          options={gameOptions}
          onChange={(date: Date) => {
            setSelectedDate(date || new Date());
          }}
          width={180}
          value={selectedProjection.date}
        />
      </div>
      <PlayerContext
        player={player}
        selectedDate={moment(selectedProjection.date).toDate()}
        setDate={setSelectedDate}
        projection={selectedProjection}
        game={player.games.find((game) =>
          moment(game.date).isSame(selectedProjection.date, "day")
        )}
      />
      {/* <PlayerProfileChart player={player} filteredGames={player.games} /> */}

      {proposition ? (
        <>
          <StatSelectButtons
            propositions={selectedProjection.propositions} // TODO: select active proposition for each statType (most recent 'last modified')
            selectedStat={proposition.statType}
            selectedProp={proposition}
            onStatSelect={onStatSelect}
          />
          <Prediction
            projection={selectedProjection}
            selectedProp={proposition}
            selectedStat={proposition.statType}
            onPredictionSelect={setProposition}
          />
          <PlayerStatsPreview // TODO: variable similarity metrics
            // TODO: more in depth table
            selectedProp={proposition}
            projection={{ ...selectedProjection, player: player }}
            customModel={customModel}
          />
          <PlayerStatsChart
            games={filteredGames}
            selectedProjection={selectedProjection}
            selectedStat={proposition?.statType}
            gameFilter={gameFilter}
          />
        </>
      ) : (
        <></>
      )}

      {/*<PlayerStatsTable
        predictions={predictions}
        selected={stat}
        matchups={matchups}
        games={statData}
        player={data.player}
        similar={data.player.similarPlayers}
        similarTeams={projection.opponent.similarTeams}
        opponent={
          projection?.opponent ??
          game?.opponent ??
          games[games.length - 1].opponent
        }
      /> */}
    </div>
  );
};

export default PlayerPage;
