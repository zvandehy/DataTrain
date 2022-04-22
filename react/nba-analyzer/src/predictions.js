// Utility functions for calculating predictions
import { round, mean } from "mathjs";
import { GetPropScore } from "./utils";

//TODO: add counts option to custom filters
const counts = [0, -30, -10, -5];
const weights = [0.3, 0.27, 0.25, 0.18]; //TODO: determine best weights to use
const stats = [
  {
    label: "Points",
    abbreviation: "PTS",
    recognize: "points",
  },
  {
    label: "Assists",
    abbreviation: "AST",
    recognize: "assists",
  },
  {
    label: "3 Pointers",
    abbreviation: "3PM",
    recognize: "3-pt made",
  },
  {
    label: "PTS + REB + AST",
    abbreviation: "PRA",
    recognize: "pts+rebs+asts",
  },
  {
    label: "Rebounds",
    abbreviation: "REB",
    recognize: "rebounds",
  },
  {
    label: "Free Throws",
    abbreviation: "FTM",
    recognize: "free throws made",
  },
  {
    label: "Fantasy",
    abbreviation: "FAN",
    recognize: "fantasy score",
  },
  {
    label: "Blocks + Steals",
    abbreviation: "B+S",
    recognize: "blks+stls",
  },
  {
    label: "Double Double",
    abbreviation: "DD",
    recognize: "double-double",
  },
];

export function CalculatePredictions(projection, statData) {
  return stats.map((item) => {
    const target = projection
      ? getTarget(projection.targets, item.recognize)
      : [];
    const playerStats = getStats(statData, counts, item.recognize, target);
    const predictionAndConfidence = getPredictionAndConfidence(
      playerStats,
      weights
    );
    return {
      stat: item,
      target: target,
      overUnder: predictionAndConfidence[0],
      confidence: predictionAndConfidence[1],
      counts: playerStats,
    };
  });
}

//TODO: Make this more sophisticated
function getPredictionAndConfidence(stats, weights) {
  let conf_o = 0;
  let conf_u = 0;

  stats.forEach((item, i) => {
    conf_o += item.pct_o * weights[i];
    conf_u += item.pct_u * weights[i];
  });
  conf_o = round(conf_o, 2);
  conf_u = round(conf_u, 2);
  if (conf_o > conf_u) {
    return ["OVER", conf_o];
  }
  return ["UNDER", conf_u];
}

function getTarget(targets, stat) {
  //TODO: Update this using the constant list of types and score mappings to get all the types
  const exists = targets.filter(
    (item) => item.type.toLowerCase() === stat.toLowerCase()
  );
  if (exists.length === 1) {
    return exists[0].target;
  }
  return 0;
}

function getStats(games, counts, stat, target) {
  let stats = [];
  const scores = games.map((game) => GetPropScore(game, stat));
  counts.forEach((count) => {
    //TODO: Apply game stat filters
    const data = scores.slice(count);
    const avg = data.length ? round(mean(data), 2) : 0;
    const over = data.filter((score) => score > target).length;
    const under = data.filter((score) => score < target).length;
    const pct_o = round((over / data.length) * 100, 2);
    const pct_u = round((under / data.length) * 100, 2);
    stats.push({
      n: data.length,
      avg: avg,
      over: over,
      under: under,
      pct_o: pct_o,
      pct_u: pct_u,
    });
  });
  return stats;
}
