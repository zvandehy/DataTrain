REPLACE INTO standardized ( date, duration, AVG_assists, STDDEV_assists, AVG_defensiveRebounds, STDDEV_defensiveRebounds, AVG_offensiveRebounds, STDDEV_offensiveRebounds, AVG_fieldGoalsAttempted, STDDEV_fieldGoalsAttempted, AVG_fieldGoalsMade, STDDEV_fieldGoalsMade, AVG_freeThrowsAttempted, STDDEV_freeThrowsAttempted, AVG_freeThrowsMade, STDDEV_freeThrowsMade, AVG_personalFoulsDrawn, STDDEV_personalFoulsDrawn, AVG_personalFouls, STDDEV_personalFouls, AVG_points, STDDEV_points, AVG_threePointersAttempted, STDDEV_threePointersAttempted, AVG_threePointersMade, STDDEV_threePointersMade, AVG_rebounds, STDDEV_rebounds, AVG_turnovers, STDDEV_turnovers, AVG_blocks, STDDEV_blocks, AVG_steals, STDDEV_steals, AVG_potentialAssists, STDDEV_potentialAssists, AVG_passes, STDDEV_passes, AVG_minutes, STDDEV_minutes, AVG_heightInches, STDDEV_heightInches, AVG_weight, STDDEV_weight)
SELECT  Cast("2022-11-27" AS Date)     AS date
       ,"2022-23"                      AS duration
       ,AVG(assists)                   AS AVG_assists
       ,stddev(assists)                AS STDDEV_assists
       ,AVG(defensiveRebounds)         AS AVG_defensiveRebounds
       ,stddev(defensiveRebounds)      AS STDDEV_defensiveRebounds
       ,AVG(offensiveRebounds)         AS AVG_offensiveRebounds
       ,stddev(offensiveRebounds)      AS STDDEV_offensiveRebounds
       ,AVG(fieldGoalsAttempted)       AS AVG_fieldGoalsAttempted
       ,stddev(fieldGoalsAttempted)    AS STDDEV_fieldGoalsAttempted
       ,AVG(fieldGoalsMade)            AS AVG_fieldGoalsMade
       ,stddev(fieldGoalsMade)         AS STDDEV_fieldGoalsMade
       ,AVG(freeThrowsAttempted)       AS AVG_freeThrowsAttempted
       ,stddev(freeThrowsAttempted)    AS STDDEV_freeThrowsAttempted
       ,AVG(freeThrowsMade)            AS AVG_freeThrowsMade
       ,stddev(freeThrowsMade)         AS STDDEV_freeThrowsMade
       ,AVG(personalFoulsDrawn)        AS AVG_personalFoulsDrawn
       ,stddev(personalFoulsDrawn)     AS STDDEV_personalFoulsDrawn
       ,AVG(personalFouls)             AS AVG_personalFouls
       ,stddev(personalFouls)          AS STDDEV_personalFouls
       ,AVG(points)                    AS AVG_points
       ,stddev(points)                 AS STDDEV_points
       ,AVG(threePointersAttempted)    AS AVG_threePointersAttempted
       ,stddev(threePointersAttempted) AS STDDEV_threePointersAttempted
       ,AVG(threePointersMade)         AS AVG_threePointersMade
       ,stddev(threePointersMade)      AS STDDEV_threePointersMade
       ,AVG(rebounds)                  AS AVG_rebounds
       ,stddev(rebounds)               AS STDDEV_rebounds
       ,AVG(turnovers)                 AS AVG_turnovers
       ,stddev(turnovers)              AS STDDEV_turnovers
       ,AVG(blocks)                    AS AVG_blocks
       ,stddev(blocks)                 AS STDDEV_blocks
       ,AVG(steals)                    AS AVG_steals
       ,stddev(steals)                 AS STDDEV_steals
       ,AVG(potentialAssists)          AS AVG_potentialAssists
       ,stddev(potentialAssists)       AS STDDEV_potentialAssists
       ,AVG(passes)                    AS AVG_passes
       ,stddev(passes)                 AS STDDEV_passes
       ,AVG(minutes)                   AS AVG_minutes
       ,stddev(minutes)                AS STDDEV_minutes
       ,AVG(heightInches)              AS AVG_heightInches
       ,stddev(heightInches)           AS STDDEV_heightInches
       ,AVG(weight)                    AS AVG_weight
       ,stddev(weight)                 AS STDDEV_weight
FROM playergames
JOIN players USING
(playerID
)
WHERE date < Cast("2022-11-27" AS Date)
AND season = "2022-23"
AND minutes > 10;