package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/zvandehy/DataTrain/clover-data-loader/cloverdata"
)

type Database struct {
	*sqlx.DB
}

func New(dataSourceName string) (*Database, error) {
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (db *Database) UploadPlayerGames(playerGames []cloverdata.PlayerGame) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, pg := range playerGames {
		query := `INSERT INTO playergames (
			assists, blocks, date, defensiveRebounds, fieldGoalPct, fieldGoalsAttempted, fieldGoalsMade, 
			freeThrowsAttempted, freeThrowsMade, freeThrowPct, gameID, homeAway, minutes, offensiveRebounds, 
			opponentID, outcome, personalFouls, playerID, playoffs, points, rebounds, season, steals, 
			teamID, threePointPct, threePointersAttempted, threePointersMade, turnovers, assistPct, usage, 
			trueShootingPct, effectiveFieldGoalPct, defensiveReboundPct, offensiveReboundPct, potentialAssists, 
			passes, assistConversionRate, personalFoulsDrawn, margin, createdAt, updatedAt
		) VALUES (
			:assists, :blocks, :date, :defensiveRebounds, :fieldGoalPct, :fieldGoalsAttempted, :fieldGoalsMade, 
			:freeThrowsAttempted, :freeThrowsMade, :freeThrowPct, :gameID, :homeAway, :minutes, :offensiveRebounds, 
			:opponentID, :outcome, :personalFouls, :playerID, :playoffs, :points, :rebounds, :season, :steals, 
			:teamID, :threePointPct, :threePointersAttempted, :threePointersMade, :turnovers, :assistPct, :usage, 
			:trueShootingPct, :effectiveFieldGoalPct, :defensiveReboundPct, :offensiveReboundPct, :potentialAssists, 
			:passes, :assistConversionRate, :personalFoulsDrawn, :margin, :createdAt, :updatedAt
		)`
		_, err := tx.NamedExec(query, pg)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
