export function ShortenName(name: string, limit: number): string {
  if (name.length <= limit) return name;
  const first_last = name.split(" ", 2);
  if (first_last[1].length + 2 <= limit)
    return first_last[0][0] + ". " + first_last[1];
  if (first_last[0].length + 1 <= limit)
    return first_last[0] + " " + first_last[1][0] + ".";
  return name.substring(0, limit) + "...";
}
