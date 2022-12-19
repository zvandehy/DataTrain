import {
  Tab,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Tabs,
} from "@mui/material";
import { useEffect, useState } from "react";
import { useGetPlayerPropositions } from "../../hooks/useGetPlayerPropositions";
import { ModelInput } from "../../shared/interfaces/custom-prediction.interface";
import {
  PropBreakdown,
  Proposition,
} from "../../shared/interfaces/graphql/proposition.interface";

interface PlayerModalProps {
  playerID: number;
  startDate: string;
  triggeredProp?: Proposition;
  customModel: ModelInput;
}

const PlayerModal: React.FC<PlayerModalProps> = ({
  playerID,
  startDate,
  triggeredProp,
  customModel,
}: PlayerModalProps) => {
  const [dateRange, setDateRange] = useState({ startDate, endDate: startDate });
  const [selectedProp, setSelectedProp] = useState<Proposition | undefined>(
    triggeredProp
  );

  //use effect to set the selected prop to the triggered prop if it exists
  useEffect(() => {
    if (triggeredProp) {
      setSelectedProp(triggeredProp);
    }
  }, [triggeredProp]);

  const {
    data: propositions,
    loading,
    error,
  } = useGetPlayerPropositions({
    playerID,
    startDate: dateRange.startDate,
    endDate: dateRange.endDate,
    customModel,
  });

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    return <div>Error: {error.message}</div>;
  }
  if (!propositions || propositions.length === 0) {
    return <div>No propositions found</div>;
  }

  let similarPlayerBreakdown: PropBreakdown = {
    name: "Similar Player Breakdown",
    over: 0,
    under: 0,
    push: 0,
    overPct: 0,
    underPct: 0,
    pushPct: 0,
    derivedAverage: 0,
    weight: 0,
    pctChange: 0,
    base: 0,
  };

  let similarPlayerBreakdowns = selectedProp?.prediction.breakdowns?.filter(
    (p) => p.name?.includes(" vs Opponent")
  );
  similarPlayerBreakdowns?.forEach((breakdown) => {
    similarPlayerBreakdown.over += breakdown.over;
    similarPlayerBreakdown.under += breakdown.under;
    similarPlayerBreakdown.push += breakdown.push;
    similarPlayerBreakdown.weight += breakdown.weight;
    similarPlayerBreakdown.base += breakdown.base;
    similarPlayerBreakdown.derivedAverage += breakdown.derivedAverage;
  });
  similarPlayerBreakdown.overPct =
    (similarPlayerBreakdown.over /
      (similarPlayerBreakdown.over +
        similarPlayerBreakdown.under +
        similarPlayerBreakdown.push)) *
    100;
  similarPlayerBreakdown.underPct =
    (similarPlayerBreakdown.under /
      (similarPlayerBreakdown.over +
        similarPlayerBreakdown.under +
        similarPlayerBreakdown.push)) *
    100;
  similarPlayerBreakdown.pushPct =
    (similarPlayerBreakdown.push /
      (similarPlayerBreakdown.over +
        similarPlayerBreakdown.under +
        similarPlayerBreakdown.push)) *
    100;
  similarPlayerBreakdown.base =
    similarPlayerBreakdown.base / (similarPlayerBreakdowns?.length ?? 0);
  console.log(
    similarPlayerBreakdown.derivedAverage,
    similarPlayerBreakdowns?.length
  );
  similarPlayerBreakdown.derivedAverage =
    similarPlayerBreakdown.derivedAverage /
    (similarPlayerBreakdowns?.length ?? 0);
  similarPlayerBreakdown.pctChange =
    (similarPlayerBreakdown.derivedAverage - similarPlayerBreakdown.base) /
    similarPlayerBreakdown.base;

  // create tabs for each stat type found in the propositions
  // select the tab based on the prop passed in, or default to the prop with the highest prediction significance
  // create a table for the selected prop stat type with the following columns:
  // date, opponent, target, outcome, actual result, prediction, significance, std dev, cumulative over, cumulative under, cumulative push
  // the rows should be each of the gameBreakdowns and playerBreakdown for the selected prop stat type
  // with an expand for the derived games of the breakdowns
  return (
    <>
      <Tabs value={selectedProp?.type ?? "points"} sx={{ overflowX: "scroll" }}>
        {propositions.map((prop) => (
          <Tab
            key={prop.type}
            value={prop.type}
            label={prop.type.toUpperCase()}
            onClick={() => {
              setSelectedProp(prop);
            }}
          />
        ))}
      </Tabs>
      <Table>
        <TableHead>
          <TableRow>
            <td>Target: {selectedProp?.target}</td>
            <td>Actual Result: {selectedProp?.actualResult}</td>
            <td>Prediction: {selectedProp?.prediction?.estimation}</td>
            <td>Significance: {selectedProp?.prediction?.significance}</td>
            <td>Std Dev: {selectedProp?.prediction?.stdDev}</td>
            <td>Wager Outcome: {selectedProp?.prediction.wagerOutcome}</td>
            <td>{selectedProp?.sportsbook}</td>
          </TableRow>
          <TableRow>
            <TableCell>Breakdown</TableCell>
            <TableCell>Over (%)</TableCell>
            <TableCell>Under (%)</TableCell>
            <TableCell>Push (%)</TableCell>
            <TableCell>DerivedAverage</TableCell>
            <TableCell>Base</TableCell>
            <TableCell>PctChange</TableCell>
            <TableCell>Weight</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {selectedProp?.prediction?.breakdowns?.map((breakdown) => (
            <TableRow>
              <TableCell>{breakdown.name}</TableCell>
              <TableCell>
                {breakdown.over} ({(breakdown.overPct * 100).toFixed()}%)
              </TableCell>
              <TableCell>
                {breakdown.under} ({(breakdown.underPct * 100).toFixed()}%)
              </TableCell>
              <TableCell>
                {breakdown.push} ({(breakdown.pushPct * 100).toFixed()}%)
              </TableCell>
              <TableCell>{breakdown.derivedAverage}</TableCell>
              <TableCell>{breakdown.base}</TableCell>
              <TableCell>{breakdown.pctChange}</TableCell>
              <TableCell>{breakdown.weight}</TableCell>
            </TableRow>
          ))}
          <TableRow>
            <TableCell>{similarPlayerBreakdown.name}</TableCell>
            <TableCell>
              {similarPlayerBreakdown.over} (
              {similarPlayerBreakdown.overPct.toFixed()}%)
            </TableCell>
            <TableCell>
              {similarPlayerBreakdown.under} (
              {similarPlayerBreakdown.underPct.toFixed()}%)
            </TableCell>
            <TableCell>
              {similarPlayerBreakdown.push} (
              {similarPlayerBreakdown.pushPct.toFixed()}%)
            </TableCell>
            <TableCell>
              {similarPlayerBreakdown.derivedAverage.toFixed(2)}
            </TableCell>
            <TableCell>{similarPlayerBreakdown.base.toFixed(2)}</TableCell>
            <TableCell>
              {(similarPlayerBreakdown.pctChange * 100).toFixed(2)}
            </TableCell>
            <TableCell>{similarPlayerBreakdown.weight.toFixed(2)}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </>
  );
};

export default PlayerModal;
