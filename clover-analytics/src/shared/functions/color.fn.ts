export function ColorPct(value: number) {
  if (value > 1) {
    value = value / 100;
  }
  let ret = "med";
  if (value <= 0.45) {
    ret = "low";
  } else if (value >= 0.6) {
    ret = "high";
  }
  return ret;
}

export function ColorResult(result: string) {
  let ret = "med";
  if (result.toUpperCase() === "UNDER") {
    ret = "low";
  } else if (result.toUpperCase() === "OVER") {
    ret = "high";
  }
  return ret;
}

export function ColorCompare(value: number, target: number) {
  let ret = "med";
  if (value < target) {
    ret = "low";
  } else if (value > target) {
    ret = "high";
  }
  return ret;
}
