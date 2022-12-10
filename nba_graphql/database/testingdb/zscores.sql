SELECT  name                                                                                                                            AS fromPlayerID
    --    ,"2022-23"                                                                                                                           AS duration
    --    ,Cast("2022-11-26" AS Date)                                                                                                          AS endDate
       ,((AVG(assists)-LEAGUE_AVG_AST)/LEAGUE_STD_AST)                                    AS assists
    --    ,((AVG(defensiveRebounds)-LEAGUE_AVG_DREB)/LEAGUE_STD_DREB)                      AS defensiveRebounds
    --    ,((AVG(offensiveRebounds)-LEAGUE_AVG_OREB)/LEAGUE_STD_OREB)                      AS offensiveRebounds
       ,((AVG(fieldGoalsAttempted)-LEAGUE_AVG_FGA)/LEAGUE_STD_FGA)                        AS fieldGoalsAttempted
    --    ,((AVG(fieldGoalsMade)-LEAGUE_AVG_FGM)/LEAGUE_STD_FGM)                             AS fieldGoalsMade
    --    ,((AVG(freeThrowsAttempted)-LEAGUE_AVG_FTA)/LEAGUE_STD_FTA)                        AS freeThrowsAttempted
    --    ,((AVG(freeThrowsMade)-LEAGUE_AVG_FTM)/LEAGUE_STD_FTM)                             AS freeThrowsMade
    --    ,((AVG(personalFoulsDrawn)-LEAGUE_AVG_PFD)/LEAGUE_STD_PFD)                         AS personalFoulsDrawn
    --    ,((AVG(personalFouls)-LEAGUE_AVG_PF)/LEAGUE_STD_PF)                                  AS personalFouls
       ,((AVG(points)-LEAGUE_AVG_POINTS)/LEAGUE_STD_POINTS)                         AS points
    --    ,((AVG(threePointersAttempted)-LEAGUE_AVG_THREE_PA)/LEAGUE_STD_THREE_PA) AS threePointersAttempted
       ,((AVG(threePointersMade)-LEAGUE_AVG_THREE_PM)/LEAGUE_STD_THREE_PM)      AS threePointersMade
       ,((AVG(rebounds)-LEAGUE_AVG_REB)/LEAGUE_STD_REB)                                   AS rebounds
    --    ,((AVG(turnovers)-LEAGUE_AVG_TOV)/LEAGUE_STD_TOV)                                  AS turnovers
    --    ,((AVG(blocks)-LEAGUE_AVG_BLK)/LEAGUE_STD_BLK)                                     AS blocks
    --    ,((AVG(steals)-LEAGUE_AVG_STL)/LEAGUE_STD_STL)                                     AS steals
    --    ,((AVG(potentialAssists)-LEAGUE_AVG_POT_AST)/LEAGUE_STD_POT_AST)           AS potentialAssists
    --    ,((AVG(passes)-LEAGUE_AVG_PASS)/LEAGUE_STD_PASS)                                 AS passes
       ,((AVG(minutes)-LEAGUE_AVG_MIN)/LEAGUE_STD_MIN)                                    AS minutes
       ,((heightInches-LEAGUE_AVG_HEIGHT)/LEAGUE_STD_HEIGHT) AS height
       ,((weight-LEAGUE_AVG_WEIGHT)/LEAGUE_STD_WEIGHT) AS weight
FROM players
JOIN playergames USING
(playerID
)
JOIN
(
	SELECT  AVG(assists)                   AS LEAGUE_AVG_AST
	       ,stddev(assists)                AS LEAGUE_STD_AST
	       ,AVG(defensiveRebounds)         AS LEAGUE_AVG_DREB
	       ,stddev(defensiveRebounds)      AS LEAGUE_STD_DREB
	       ,AVG(offensiveRebounds)         AS LEAGUE_AVG_OREB
	       ,stddev(offensiveRebounds)      AS LEAGUE_STD_OREB
	       ,AVG(fieldGoalsAttempted)       AS LEAGUE_AVG_FGA
	       ,stddev(fieldGoalsAttempted)    AS LEAGUE_STD_FGA
	       ,AVG(fieldGoalsMade)            AS LEAGUE_AVG_FGM
	       ,stddev(fieldGoalsMade)         AS LEAGUE_STD_FGM
	       ,AVG(freeThrowsAttempted)       AS LEAGUE_AVG_FTA
	       ,stddev(freeThrowsAttempted)    AS LEAGUE_STD_FTA
	       ,AVG(freeThrowsMade)            AS LEAGUE_AVG_FTM
	       ,stddev(freeThrowsMade)         AS LEAGUE_STD_FTM
	       ,AVG(personalFoulsDrawn)        AS LEAGUE_AVG_PFD
	       ,stddev(personalFoulsDrawn)     AS LEAGUE_STD_PFD
	       ,AVG(personalFouls)             AS LEAGUE_AVG_PF
	       ,stddev(personalFouls)          AS LEAGUE_STD_PF
	       ,AVG(points)                    AS LEAGUE_AVG_POINTS
	       ,stddev(points)                 AS LEAGUE_STD_POINTS
	       ,AVG(threePointersAttempted)    AS LEAGUE_AVG_THREE_PA
	       ,stddev(threePointersAttempted) AS LEAGUE_STD_THREE_PA
	       ,AVG(threePointersMade)         AS LEAGUE_AVG_THREE_PM
	       ,stddev(threePointersMade)      AS LEAGUE_STD_THREE_PM
	       ,AVG(rebounds)                  AS LEAGUE_AVG_REB
	       ,stddev(rebounds)               AS LEAGUE_STD_REB
	       ,AVG(turnovers)                 AS LEAGUE_AVG_TOV
	       ,stddev(turnovers)              AS LEAGUE_STD_TOV
	       ,AVG(blocks)                    AS LEAGUE_AVG_BLK
	       ,stddev(blocks)                 AS LEAGUE_STD_BLK
	       ,AVG(steals)                    AS LEAGUE_AVG_STL
	       ,stddev(steals)                 AS LEAGUE_STD_STL
	       ,AVG(potentialAssists)          AS LEAGUE_AVG_POT_AST
	       ,stddev(potentialAssists)       AS LEAGUE_STD_POT_AST
	       ,AVG(passes)                    AS LEAGUE_AVG_PASS
	       ,stddev(passes)                 AS LEAGUE_STD_PASS
	       ,AVG(minutes)                   AS LEAGUE_AVG_MIN
	       ,stddev(minutes)                AS LEAGUE_STD_MIN
	       ,AVG(heightInches)              AS LEAGUE_AVG_HEIGHT
	       ,stddev(heightInches)           AS LEAGUE_STD_HEIGHT
	       ,AVG(weight)                    AS LEAGUE_AVG_WEIGHT
	       ,stddev(weight)                 AS LEAGUE_STD_WEIGHT
	FROM playergames
	JOIN players USING
	(playerID
	)
	WHERE date < Cast("2022-11-26" AS Date)
	AND season = "2022-23" 
) AS LEAGUE
WHERE date < Cast("2022-11-26" AS Date)
AND season = "2022-23"
AND name IN ("Nikola Jokic", "Joel Embiid", "Luka Doncic", "Ja Morant", "Trae Young", "Zion Williamson", "Devin Booker", "Jayson Tatum", "Zeke Nnaji")
GROUP BY  playerID
         ,LEAGUE_AVG_POINTS;