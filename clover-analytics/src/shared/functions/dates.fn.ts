export function FormatDate(date: Date): string {
  // console.log("FormatDate: ", date, date.toISOString(), FormatDate(date));
  //subtract 1 day from the date to get the day before the date
  let dateBefore = new Date(new Date(date).getTime() - 24 * 60 * 60 * 1000);
  // const yyyy = date.getFullYear();
  // const mm = String(date.getMonth() + 1).padStart(2, "0");
  // const dd = String(date.getDate()).padStart(2, "0");
  // let ret = `${yyyy}-${mm}-${dd}`;
  return dateBefore.toISOString().split("T")[0];
}

//returns negative if 1 is before 2
export function CompareDates(date1: string, date2: string) {
  var a = new Date(date1);
  var b = new Date(date2);
  return a.getTime() - b.getTime();
}
