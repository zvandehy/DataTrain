export function FormatDate(date: Date): string {
  let ret = date.toISOString().split("T")[0];
  // const yyyy = date.getFullYear();
  // const mm = String(date.getMonth() + 1).padStart(2, "0");
  // const dd = String(date.getDate()).padStart(2, "0");
  // let ret = `${yyyy}-${mm}-${dd}`;
  return ret;
}

//returns negative if 1 is less than 2
export function CompareDates(date1: string, date2: string) {
  var a = new Date(date1);
  var b = new Date(date2);
  return a.getTime() - b.getTime();
}
