SELECT  (AVG(points)-AVG_points)/STDDEV_points
FROM playergames
JOIN
(
	SELECT  *
	FROM standardized
	WHERE date = Cast("2022-11-27" AS Date)
	AND duration = "2022-23" 
) USING (date)
JOIN players USING
(playerID
)
WHERE date < Cast("2022-11-27" AS Date)
AND season = "2022-23"
AND name = "Nikola Jokic";