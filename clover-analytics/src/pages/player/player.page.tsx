import moment from "moment";
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import AutocompleteFilter from "../../ components/autocomplete-filter/autocomplete-filter.component";
import PlayerContext from "../../ components/player-detailed/player-context/player-context.component";
import PlayerProfileChart from "../../ components/player-detailed/player-profile/player-profile.component";
import StatSelectButtons from "../../ components/playercard-list/playercard/stat-select-buttons/stat-select-buttons.component";
import { useGetPlayerDetails } from "../../hooks/useGetPlayerDetail";
import { CompareDates } from "../../shared/functions/dates.fn";
import {
  GetMaxConfidence,
  UpdatePropositionWithPrediction,
} from "../../shared/functions/predictions.fn";
import { GetImpliedTarget } from "../../shared/functions/target.fn";
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
import "./player.page.css";

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

const PlayerPage: React.FC = () => {
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
  const [similarPlayersToggle, toggleSimilarPlayers] = useState(false);
  const [similarTeamsToggle, toggleSimilarTeams] = useState(false);

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
          projection
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

  useEffect(() => {
    if (player) {
      setProjection(getProjection(date, player.projections));
    }
  }, [date, player]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error</div>;
  if (!player) return <div>No player found</div>;
  let filteredGames = player.games.filter((game) =>
    moment(game.date).isBefore(date)
  );
  console.log(filteredGames);
  let gameOptions = player.games.map((game) => {
    return {
      id: game.date,
      label: game.date,
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
        label: projection.date,
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
            setDate(date || new Date());
          }}
          width={180}
        />
      </div>
      <PlayerContext
        player={player}
        selectedDate={date}
        setDate={setDate}
        projection={projection}
        game={player.games.find((game) =>
          moment(game.date).isSame(date, "day")
        )}
      />
      <PlayerProfileChart player={player} filteredGames={filteredGames} />

      {projection && proposition ? (
        <StatSelectButtons
          propositions={projection.propositions}
          selectedStat={proposition.statType}
          onStatSelect={onStatSelect}
        />
      ) : (
        <></>
      )}
      {/*<Prediction
        propositions={projection.propositions}
        predictions={predictions}
        selected={stat}
        game={game}
    /> */}
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
      {/*<PlayerStatsChart
        games={games}
        predictions={predictions}
        selected={stat}
      /> */}
    </div>
  );
};

export default PlayerPage;
