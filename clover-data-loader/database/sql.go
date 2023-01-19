package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/clover-data-loader/cloverdata"
)

type Database struct {
	*sqlx.DB
}

func NewDB(dataSourceName string) (*Database, error) {
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) UploadPlayerGames(playerGames []cloverdata.PlayerGame) error {
	batchSize := 50
	// upload playerGames in batches of batchSize
	for i := 0; i < len(playerGames); i += batchSize {
		end := i + batchSize
		if end > len(playerGames) {
			end = len(playerGames)
		}
		err := db.uploadPlayerGames(playerGames[i:end])
		if err != nil {
			logrus.Errorf("failed uploading playergames. Uploaded %d of %d", i, len(playerGames))
			return err
		}
	}
	return nil
}

func (db *Database) uploadPlayerGames(playerGames []cloverdata.PlayerGame) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, pg := range playerGames {
		query := `REPLACE INTO playergames (
			assists, blocks, date, defensiveRebounds, fieldGoalPct, fieldGoalsAttempted, fieldGoalsMade, 
			freeThrowsAttempted, freeThrowsMade, freeThrowPct, gameID, homeAway, minutes, offensiveRebounds, 
			opponentID, outcome, personalFouls, playerID, playoffs, points, rebounds, season, steals, 
			teamID, threePointPct, threePointersAttempted, threePointersMade, turnovers, assistPct, usage, 
			trueShootingPct, effectiveFieldGoalPct, defensiveReboundPct, offensiveReboundPct, potentialAssists, 
			passes, assistConversionRate, personalFoulsDrawn, margin, CreatedAt, UpdatedAt
		) VALUES (
			:assists, :blocks, :date, :defensiveRebounds, :fieldGoalPct, :fieldGoalsAttempted, :fieldGoalsMade, 
			:freeThrowsAttempted, :freeThrowsMade, :freeThrowPct, :gameID, :homeAway, :minutes, :offensiveRebounds, 
			:opponentID, :outcome, :personalFouls, :playerID, :playoffs, :points, :rebounds, :season, :steals, 
			:teamID, :threePointPct, :threePointersAttempted, :threePointersMade, :turnovers, :assistPct, :usage, 
			:trueShootingPct, :effectiveFieldGoalPct, :defensiveReboundPct, :offensiveReboundPct, :potentialAssists, 
			:passes, :assistConversionRate, :personalFoulsDrawn, :margin, :CreatedAt, :UpdatedAt
		)`
		_, err := tx.NamedExec(query, pg)
		if err != nil {
			logrus.Errorf("couldn't upload player game: %+v", pg)
			return err
		}
	}
	logrus.Infof("uploaded %d player games", len(playerGames))
	return tx.Commit()
}

func (db *Database) UploadStandardizedPlayerComparisons(date string) error {
	res := db.MustExec(cloverdata.SQLUploadPlayerComparisons(date))
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("Uploaded %2.2d league standardized stats for %s", rowsAffected, date)
	return nil
}

func (db *Database) GetMissingPlayerIDs(playerIDs []int) ([]int, error) {
	query := "SELECT playerID FROM players"
	var dbPlayers []int
	err := db.Select(&dbPlayers, query)
	if err != nil {
		return nil, err
	}
	dbPlayerMap := make(map[int]bool)
	for _, id := range dbPlayers {
		dbPlayerMap[id] = true
	}
	var missingIDs []int
	for _, id := range playerIDs {
		if _, ok := dbPlayerMap[id]; !ok {
			missingIDs = append(missingIDs, id)
		}
	}
	return missingIDs, nil
}

func (db *Database) UploadPlayers(players []cloverdata.Player) error {
	batchSize := 50
	// upload playerGames in batches of batchSize
	for i := 0; i < len(players); i += batchSize {
		end := i + batchSize
		if end > len(players) {
			end = len(players)
		}
		err := db.uploadPlayers(players[i:end])
		if err != nil {
			logrus.Errorf("failed uploading players. Uploaded %d of %d", i, len(players))
			return err
		}
	}
	return nil
}

func (db *Database) uploadPlayers(players []cloverdata.Player) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, p := range players {
		query := `REPLACE INTO players (
			firstName, lastName, name, playerID, position, teamABR, teamID, height, heightInches, weight, 
			CreatedAt, UpdatedAt, league
		) VALUES (
			:firstName, :lastName, :name, :playerID, :position, :teamABR, :teamID, :height, :heightInches, :weight, 
			:CreatedAt, :UpdatedAt, :league
		)`
		_, err := tx.NamedExec(query, p)
		if err != nil {
			logrus.Errorf("couldn't upload player: %+v", p)
			return err
		}
	}
	logrus.Infof("uploaded %d players", len(players))
	return tx.Commit()
}

func (db *Database) UploadTeamGames(teamGames []cloverdata.TeamGame) error {
	batchSize := 50
	// upload playerGames in batches of batchSize
	for i := 0; i < len(teamGames); i += batchSize {
		end := i + batchSize
		if end > len(teamGames) {
			end = len(teamGames)
		}
		err := db.uploadTeamGames(teamGames[i:end])
		if err != nil {
			logrus.Errorf("failed uploading teamgames. Uploaded %d of %d", i, len(teamGames))
			return err
		}
	}
	return nil
}

func (db *Database) uploadTeamGames(teamGames []cloverdata.TeamGame) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, tg := range teamGames {
		query := `REPLACE INTO teamgames (
			gameID, teamID, opponentID, homeAway, date, outcome, margin, playoffs, assists, blocks, defensiveRebounds, offensiveRebounds,
			fieldGoalPct, fieldGoalsAttempted, fieldGoalsMade, freeThrowPct, freeThrowsAttempted, freeThrowsMade, threePointPct, threePointersAttempted, threePointersMade, points, personalFouls, rebounds, steals, turnovers, season,
			defensiveReboundPct, offensiveReboundPct, plusMinusPerHundred, possessions, effectiveFieldGoalPct,
			CreatedAt, UpdatedAt
		) VALUES (
			:gameID, :teamID, :opponentID, :homeAway, :date, :outcome, :margin, :playoffs, :assists, :blocks, :defensiveRebounds, :offensiveRebounds,
			:fieldGoalPct, :fieldGoalsAttempted, :fieldGoalsMade, :freeThrowPct, :freeThrowsAttempted, :freeThrowsMade, :threePointPct, :threePointersAttempted, :threePointersMade, :points, :personalFouls, :rebounds, :steals, :turnovers, :season,
			:defensiveReboundPct, :offensiveReboundPct, :plusMinusPerHundred, :possessions, :effectiveFieldGoalPct,
			:CreatedAt, :UpdatedAt
		)`
		_, err := tx.NamedExec(query, tg)
		if err != nil {
			logrus.Errorf("couldn't upload teamgame: %+v", tg)
			return err
		}
	}
	logrus.Infof("uploaded %d teamgames", len(teamGames))
	return tx.Commit()
}
