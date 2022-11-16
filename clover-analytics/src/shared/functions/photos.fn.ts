export default function getPlayerPhotoUrl(
  playerId: number,
  league: string
): string {
  return `https://ak-static.cms.nba.com/wp-content/uploads/headshots/${league.toLowerCase()}/latest/260x190/${playerId}.png`;
}
