export function ColorPct(value: number) {
  let ret = "med";
  if (value <= 0.4) {
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
