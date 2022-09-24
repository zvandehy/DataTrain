import moment from "moment";
import { useEffect, useMemo, useState } from "react";
import { Match, SortProjections } from "../../shared/functions/filters.fn";
import { CalculatePredictions } from "../../shared/functions/predictions.fn";
import { CustomCalculation } from "../../shared/interfaces/custom-prediction.interface";
import { DateRange } from "../../shared/interfaces/dateRange.interface";
import { Accuracy } from "../../shared/interfaces/accuracy.interface";
import { GameFilter } from "../../shared/interfaces/graphql/filters.interface";
import { Projection } from "../../shared/interfaces/graphql/projection.interface";
import { Stat } from "../../shared/interfaces/stat.interface";
import CalendarSummary from "../projections-summary/calendar-summary/calendar-summary.component";
import PlayerListFilters from "./list-filters/list-filters.component";
import "./playercard-list.component.css";
import PlayerCard from "./playercard/playercard.component";
import OverallAccuracyBreakdownTable from "../projections-summary/breakdown-table/overall-accuracy-breakdown.component";

interface PlayerCardListProps {
  projections: Projection[];
  customModel: CustomCalculation;
  gameFilter: GameFilter;
}

const PlayerCardList: React.FC<PlayerCardListProps> = ({
  projections,
  customModel,
  gameFilter,
}: PlayerCardListProps) => {
  const calculatedProjections = useMemo(() => {
    return CalculatePredictions(projections, gameFilter, customModel);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [JSON.stringify(customModel), JSON.stringify(gameFilter)]);
  const [lookup, setLookup] = useState("");
  const [sortType, setSortType] = useState("");
  const [statType, setStatType] = useState(undefined as Stat | undefined);

  let filteredProjections = calculatedProjections.filter((projection) => {
    if (statType !== undefined) {
      return Match(projection, { lookup: lookup, statType: statType as Stat });
    }
    return Match(projection, { lookup: lookup });
  });
  filteredProjections = SortProjections(filteredProjections, {
    sortBy: sortType,
    statType: statType,
  });

  const [dateRange, setDateRange] = useState(
    filteredProjections.length > 0
      ? ({
          start: moment(filteredProjections[0].startTime)
            .subtract(7, "days")
            .format("YYYY-MM-DD"),
          end: moment(filteredProjections[0].startTime).format("YYYY-MM-DD"),
        } as DateRange)
      : {
          start: moment(new Date()).subtract(7, "days").format("YYYY-MM-DD"),
          end: moment(new Date()).format("YYYY-MM-DD"),
        }
  );

  const onSetDates = (date: string) => {
    if (moment(date).isBefore(dateRange.start)) {
      setDateRange((prev) => ({
        ...prev,
        start: moment(date).format("YYYY-MM-DD"),
      }));
    }
    if (moment(date).isAfter(dateRange.end)) {
      setDateRange((prev) => ({
        ...prev,
        end: moment(date).format("YYYY-MM-DD"),
      }));
    }
  };

  // const [accuracies, setAccuracies] = useState([] as Accuracy[]);
  // const [totalAccuracy, setTotalAccuracy] = useState(new Accuracy());

  // useEffect(() => {
  //   setTotalAccuracy(new Accuracy());
  //   setDateRange(
  //     filteredProjections.length > 0
  //       ? ({
  //           start: moment(filteredProjections[0].startTime)
  //             .add(-1, "days")
  //             .format("YYYY-MM-DD"),
  //           end: moment(filteredProjections[0].startTime).format("YYYY-MM-DD"),
  //         } as DateRange)
  //       : {
  //           start: moment(new Date()).add(-1, "days").format("YYYY-MM-DD"),
  //           end: moment(new Date()).format("YYYY-MM-DD"),
  //         }
  //   );
  // }, [customModel, statType]);

  // const calendarSummary = useMemo(() => {
  //   console.log("reload calendar summary");
  //   return (
  //     // <CalendarSummary
  //     //   accuracy={totalAccuracy}
  //     //   dateRange={dateRange}
  //     //   setDates={onSetDates}
  //     //   setAccuracy={setTotalAccuracy}
  //     //   customModel={customModel}
  //     //   statType={statType}
  //     //   lookup={lookup}
  //     // ></CalendarSummary>
  //   );
  // }, [totalAccuracy.allProps, dateRange.start]);

  const totalTable = useMemo(() => {
    console.log("Reload table");
    return (
      <OverallAccuracyBreakdownTable
        customModel={customModel}
        lookup={lookup}
        dateRange={dateRange}
        setDates={setDateRange}
        hitCriteria={customModel.hitCriteria}
      />
    );
  }, [customModel, dateRange]);

  return (
    <>
      <PlayerListFilters
        onSearchChange={setLookup}
        onSortSelect={setSortType}
        onStatSelect={setStatType}
      />
      {/* {calendarSummary} */}
      {totalTable}
      <div id="player-list">
        {filteredProjections.length > 0 ? (
          filteredProjections.map((projection) => {
            return (
              <PlayerCard
                key={`${projection.player.playerID} ${projection.startTime}`}
                projection={projection}
                filteredStatType={statType}
                gameFilter={gameFilter}
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
