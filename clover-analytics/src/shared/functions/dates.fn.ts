//returns negative if 1 is before 2
export function CompareDates(date1: string, date2: string) {
  var a = new Date(date1);
  var b = new Date(date2);
  return a.getTime() - b.getTime();
}
