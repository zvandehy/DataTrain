import moment from "moment";
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
  GetMaxConfidence,
  UpdatePropositionWithPrediction,
} from "../../../shared/functions/predictions.fn";
import { GetImpliedTarget } from "../../../shared/functions/target.fn";
import { ProjectionFilter } from "../../../shared/interfaces/graphql/filters.interface";
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
import PlayerProfileChart from "../player-profile/player-profile.component";
import StatSelectButtons from "../../playercard-list/playercard/stat-select-buttons/stat-select-buttons.component";
import PlayerStatsPreview from "../../player-stats-table/player-stats-preview/player-stats-preview.component";
import Prediction from "../../prediction/prediction.component";
import PlayerStatsChart from "../player-stats-chart/player-stats-chart.component";

// const [date, setDate] = useState(FormatDate(new Date()));
// //TODO: Handle error when pick game that isn't in seasonType
// const [seasonType, setSeasonType] = useState("");

function FormatDate(date: Date): string {
  return moment(date).format("YYYY-MM-DD");
}

function getProjection(date: Date, projections: Projection[]): Projection {
  return (
    projections.find((projection) => projection.date === FormatDate(date)) ??
    projections.sort((a, b) => (a.date > b.date ? 1 : -1))[0]
  );
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

interface PlayerPageProps {
  player: Player;
  selectedDate: Date;
  setSelectedDate: (date: Date) => void;
}
let league = "wnba";

const PlayerPage: React.FC<PlayerPageProps> = ({
  player,
  selectedDate,
  setSelectedDate,
}: PlayerPageProps) => {
  const [projection, setProjection] = useState(getProjection(selectedDate, []));
  const [statType, setStatType] = useState(Points as Stat | undefined);
  const [proposition, setProposition] = useState(
    getProposition(projection, statType)
  );
  // const [date, setDate] = useState(moment(new Date()).add(1, "days").toDate()); //Ensure that date without projection/game works
  const [season, setSeason] = useState("2022-23");
  const [sportsbook, setSportsbook] = useState("");
  const [similarPlayersToggle, toggleSimilarPlayers] = useState(false);
  const [similarTeamsToggle, toggleSimilarTeams] = useState(false);

  const onStatSelect = (stat: Stat) => {
    if (projection) {
      setStatType(stat);
    }
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
      setProjection(getProjection(selectedDate, player.projections));
      if (projection && statType) {
        let customTarget = GetImpliedTarget(projection, statType);
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
            predictionFragments: [],
          },
        };
        const foundProp =
          projection.propositions.find((p) => p.statType === statType) ||
          UpdatePropositionWithPrediction(
            customProp,
            projection.player.games,
            projection
          );
        setProposition(foundProp);
      }
    }
  }, [selectedDate, player, projection, statType]);

  let filteredGames = player.games.filter((game) =>
    moment(game.date).isBefore(selectedDate)
  );
  let gameOptions = player.games.map((game) => {
    return {
      id: game.date,
      label: `${game.date} vs ${game.opponent.name}`,
    };
  });
  let missingProjections = player.projections.filter(
    (projection) => !gameOptions.some((game) => game.id === projection.date)
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
        />
      </div>
      <PlayerContext
        player={player}
        selectedDate={selectedDate}
        setDate={setSelectedDate}
        projection={projection}
        game={player.games.find((game) =>
          moment(game.date).isSame(selectedDate, "day")
        )}
      />
      <PlayerProfileChart player={player} filteredGames={player.games} />

      {projection && proposition ? (
        <>
          <StatSelectButtons
            propositions={projection.propositions} // TODO: select active proposition for each statType (most recent 'last modified')
            selectedStat={proposition.statType}
            onStatSelect={onStatSelect}
          />
          <Prediction
            projection={projection}
            selectedProp={proposition}
            selectedStat={proposition.statType}
            onPredictionSelect={setProposition}
          />
          <PlayerStatsPreview // TODO: variable similarity metrics
            // TODO: more in depth table
            selectedProp={proposition}
            projection={{ ...projection, player: player }}
          />
          <PlayerStatsChart
            player={player}
            selectedProjection={projection}
            selectedStat={proposition?.statType}
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
