import { useState } from "react";
import { ProjectionQueryResult } from "../../hooks/useGetProjections";
import { Match, SortProjections } from "../../shared/functions/filters.fn";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
import { ScoreType } from "../../shared/interfaces/score-type.enum";
import { Stat } from "../../shared/interfaces/stat.interface";
import ProjectionsSummary from "../projections-summary/projections-summary.component";
import PlayerListFilters from "./list-filters/list-filters.component";
import "./playercard-list.component.css";
import PlayerCard from "./playercard/playercard.component";

interface PlayerCardListProps {
  projectionQueryResult: ProjectionQueryResult;
  customModel: CustomCalculation;
}

const PlayerCardList: React.FC<PlayerCardListProps> = ({
  projectionQueryResult,
  customModel,
}: PlayerCardListProps) => {
  const [lookup, setLookup] = useState("");
  const [sortType, setSortType] = useState("");
  const [statType, setStatType] = useState(undefined as Stat | undefined);
  const { loading, error, data: projections } = projectionQueryResult;
  if (loading) {
    return (
      <>
        <PlayerListFilters
          onSearchChange={setLookup}
          onSortSelect={setSortType}
          onStatSelect={setStatType}
        />
        {loading}
      </>
    );
  }
  if (error) {
    return <>{error}</>;
  }
  console.log(projections.filter((p) => p.date === undefined));
  let filteredProjections = projections.filter((projection) => {
    if (statType !== undefined) {
      return Match(projection, { lookup: lookup, statType: statType as Stat });
    }
    return Match(projection, { lookup: lookup });
  });
  filteredProjections = SortProjections(filteredProjections, {
    sortBy: sortType,
    statType: statType,
  });
  console.log(filteredProjections);

  return (
    <>
      <PlayerListFilters
        onSearchChange={setLookup}
        onSortSelect={setSortType}
        onStatSelect={setStatType}
      />
      <ProjectionsSummary
        projections={filteredProjections}
        filteredStat={statType}
      />
      <div id="player-list">
        {filteredProjections.length > 0 ? (
          filteredProjections.map((projection) => {
            return (
              <PlayerCard
                key={`${projection.player.playerID} ${projection.date}`}
                projection={projection}
                filteredStatType={statType}
                customModel={customModel}
              />
            );
          })
        ) : (
          <div className={"no-results"}>
            <h1>No Projections Found</h1>
            <p>
              Try another date! If you think this is an error then please
              contact support.
            </p>
          </div>
        )}
      </div>
    </>
  );
};

export default PlayerCardList;
