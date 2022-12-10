SELECT  players.name
       ,players.playerID
       ,AVG(points) AS vsOpp_points
       ,COUNT(*)    AS vsOpp_games
       ,averages.avg_points
       ,total_games
FROM playergames
JOIN players USING
(playerID
)
JOIN teams
ON teams.teamID = playergames.opponentID
JOIN
(
	SELECT  players.name
	       ,playerID
	       ,AVG(points) AS avg_points
	       ,COUNT(*)    AS total_games
	FROM playergames
	JOIN players USING
	(playerID
	)
	JOIN teams
	ON teams.teamID = playergames.opponentID
	WHERE season = "2022-23"
	AND position = "G"
	GROUP BY  playerID
) AS averages
ON averages.playerID = players.playerID
WHERE abbreviation = "LAL"
AND season = "2022-23"
AND position = "G"
GROUP BY  playerID;