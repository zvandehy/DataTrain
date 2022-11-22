package database

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
	"github.com/zvandehy/DataTrain/nba_graphql/util"
)

type SQLClient struct {
	League  string
	Queries int
	*sqlx.DB
	// PlayerSimilarity model.PlayerSnapshots
	// TeamSimilarity   model.TeamSnapshots
	// PlayerCache      map[string][]*model.Player
	// TeamCache        map[string][]*model.Team
}

func NewSQLClient(league string) (*SQLClient, error) {
	db, err := GetDatabase(league)
	if err != nil {
		return nil, err
	}
	return &SQLClient{
		League:  league,
		DB:      db,
		Queries: 0,
	}, nil
}

func GetDatabase(league string) (*sqlx.DB, error) {
	dsn := os.Getenv("DSN")
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s/%s?tls=true&parseTime=true", dsn, strings.ToLower(league)))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, err
}

func (c *SQLClient) GetLeague() string {
	return c.League
}

func (c *SQLClient) CountQueries() int {
	return c.Queries
}

func (c *SQLClient) SetQueries(q int) {
	c.Queries = q
}

func (c *SQLClient) AddQuery() {
	c.Queries++
}

func (c *SQLClient) SaveDBPropositions(ctx context.Context, propositions []*model.DBProposition) (int, error) {
	tx := c.MustBegin()
	for _, proposition := range propositions {
		tx.MustExec("REPLACE INTO propositions (playerID, gameID, playername, opponentID, lastModified, sportsbook, statType, target) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", proposition.PlayerID, proposition.GameID, proposition.PlayerName, proposition.OpponentID, proposition.LastModified, proposition.Sportsbook, proposition.StatType, proposition.Target)
	}
	err := tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return len(propositions), nil
}

func (c *SQLClient) SavePropositions(ctx context.Context, propositions []*model.Proposition) (int, error) {
	tx := c.MustBegin()
	for _, proposition := range propositions {
		tx.MustExec("REPLACE INTO propositions (playerID, gameID, playername, opponentID, lastModified, sportsbook, statType, target) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", proposition.PlayerID, proposition.GameID, proposition.PlayerName, proposition.OpponentID, proposition.LastModified, proposition.Sportsbook, proposition.TypeRaw, proposition.Target)
	}
	err := tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return len(propositions), nil
}

var playerIDs = map[int]model.Player{}
var playerNames = map[string]model.Player{}

func (c *SQLClient) GetPlayers(ctx context.Context, withGames bool, playerFilters ...*model.PlayerFilter) ([]*model.Player, error) {
	c.AddQuery()
	start := time.Now()
	or := []string{}
	args := []interface{}{}
	for _, playerFilter := range playerFilters {
		where := []string{}
		if playerFilter != nil {
			if playerFilter.PlayerID != nil {
				where = append(where, "playerID = ?")
				args = append(args, *playerFilter.PlayerID)
			}
			if playerFilter.Name != nil {
				where = append(where, "name = ?")
				args = append(args, *playerFilter.Name)
			}
			if playerFilter.TeamID != nil {
				where = append(where, "teamID = ?")
				args = append(args, *playerFilter.TeamID)
			}
			if playerFilter.TeamAbr != nil {
				where = append(where, "teamABR = ?")
				args = append(args, *playerFilter.TeamAbr)
			}
			if playerFilter.PositionStrict != nil {
				switch *playerFilter.PositionStrict {
				case model.PositionG:
					where = append(where, "position = ?")
					args = append(args, model.PositionG)
				case model.PositionF:
					where = append(where, "position = ?")
					args = append(args, model.PositionF)
				case model.PositionC:
					where = append(where, "position = ?")
					args = append(args, model.PositionC)
				case model.PositionFG:
					fallthrough
				case model.PositionGF:
					where = append(where, "(position = ? OR position = ?)")
					args = append(args, model.PositionFG, model.PositionGF)
				case model.PositionCF:
					fallthrough
				case model.PositionFC:
					where = append(where, "(position = ? OR position = ?)")
					args = append(args, model.PositionCF, model.PositionFC)
				}

				where = append(where, "position = ?")
				args = append(args, *playerFilter.PositionStrict)
			}
			if playerFilter.PositionLoose != nil {
				switch *playerFilter.PositionLoose {
				case model.PositionG:
					where = append(where, "(position = ? OR position = ? OR position = ?)")
					args = append(args, model.PositionG, model.PositionGF, model.PositionFG)
				case model.PositionF:
					where = append(where, "(position = ? OR position = ? OR position = ? OR position = ? OR position = ?)")
					args = append(args, model.PositionF, model.PositionFG, model.PositionGF, model.PositionCF, model.PositionFC)
				case model.PositionC:
					where = append(where, "(position = ? OR position = ? OR position = ?)")
					args = append(args, model.PositionC, model.PositionCF, model.PositionFC)
				case model.PositionFG:
					fallthrough
				case model.PositionGF:
					where = append(where, "(position = ? OR position = ? OR position = ? OR position = ?)")
					args = append(args, model.PositionFG, model.PositionGF, model.PositionF, model.PositionG)
				case model.PositionCF:
					fallthrough
				case model.PositionFC:
					where = append(where, "(position = ? OR position = ? OR position = ? OR position = ?)")
					args = append(args, model.PositionFC, model.PositionCF, model.PositionF, model.PositionC)
				}
			}
			// TODO: PositionStrictMatch, positionloose, positionLooseMatch
			if len(where) > 0 {
				or = append(or, fmt.Sprintf("(%s)", strings.Join(where, " AND ")))
			}
		}
	}
	players := []*model.Player{}
	for _, filter := range playerFilters {
		if filter.PositionLoose != nil || filter.PositionStrict != nil {
			players = []*model.Player{}
			break
		}
		if filter.PlayerID != nil {
			if player, ok := playerIDs[*filter.PlayerID]; ok {
				players = append(players, &player)
			}
		} else if filter.Name != nil {
			if player, ok := playerNames[*filter.Name]; ok {
				players = append(players, &player)
			}
		}
	}
	cachedPlayers := len(players) == len(playerFilters) && len(players) > 0
	if !cachedPlayers {
		query := "SELECT * FROM players"
		if len(or) > 0 {
			query = fmt.Sprintf("%s WHERE %s", query, strings.Join(or, " OR "))
		}
		players = []*model.Player{}
		err := c.SelectContext(ctx, &players, query, args...)
		if err != nil {
			logrus.Warnf("query: %s", query)
			logrus.Warnf("args: %v", args)
			return nil, fmt.Errorf("failed to get players: %w", err)
		}
		for _, player := range players {
			playerIDs[player.PlayerID] = *player
			playerNames[player.Name] = *player
		}
	}
	games := []*model.PlayerGame{}
	if withGames {
		ids := []interface{}{}
		idqs := []string{}
		for _, player := range players {
			ids = append(ids, player.PlayerID)
			idqs = append(idqs, "?")
		}
		err := c.SelectContext(ctx, &games, fmt.Sprintf("SELECT * FROM playergames WHERE playerID IN (%s)", strings.Join(idqs, ", ")), ids...)
		if err != nil {
			return nil, fmt.Errorf("failed to get player games: %w", err)
		}
		for _, player := range players {
			for _, game := range games {
				if player.PlayerID == game.PlayerID {
					player.GamesCache = append(player.GamesCache, game)
				}
			}
		}
		if len(playerFilters) < 5 {
			logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Players with games: %v", len(players), playerFilters), start))
		} else {
			logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Players with games from %d inputs", len(games), len(playerFilters)), start))
		}
		return players, nil
	} else if cachedPlayers {
		return players, nil
	}
	if len(playerFilters) < 5 {
		logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Players: %v", len(players), playerFilters), start))
	} else {
		logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Players from %d inputs", len(players), len(playerFilters)), start))
	}
	return players, nil
}

func (c *SQLClient) SaveUpcomingGames(ctx context.Context, games []*model.PlayerGame) (int, error) {
	tx := c.MustBegin()
	for _, game := range games {
		g := &model.PlayerGame{
			Date:       game.Date,
			GameID:     game.GameID,
			HomeOrAway: game.HomeOrAway,
			OpponentID: game.OpponentID,
			TeamID:     game.TeamID,
			PlayerID:   game.PlayerID,
			Season:     game.Season,
			Playoffs:   game.Playoffs,
			Outcome:    "PENDING",
		}
		//TODO: team games
		// tg := &model.TeamGame{
		// 	Date:       game.Date,
		// 	GameID:     game.GameID,
		// 	HomeOrAway: game.HomeOrAway,
		// 	OpponentID: game.OpponentID,
		// 	TeamID:     game.TeamID,
		// 	Season:     game.Season,
		// 	Playoffs:   game.Playoffs,
		// 	Outcome:    "PENDING",
		// }
		// oppg := &model.TeamGame{}
		tx.MustExec("INSERT INTO playergames (date, gameID, homeAway, opponentID, teamID, playerID, season, playoffs, outcome) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", g.Date, g.GameID, g.HomeOrAway, g.OpponentID, g.TeamID, g.PlayerID, g.Season, g.Playoffs, g.Outcome)
		//tx.MustExce("REPLACE INTO teamgames") // team
		//tx.MustExce("REPLACE INTO teamgames") // opponent
	}
	err := tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return len(games), nil
}

func (c *SQLClient) GetPlayerGames(ctx context.Context, inputs ...model.GameFilter) (games []*model.PlayerGame, err error) {
	c.AddQuery()
	start := time.Now()
	or := []string{}
	args := []interface{}{}
	for _, input := range inputs {
		where := []string{}
		if input.PlayerID != nil {
			where = append(where, "playerID = ?")
			args = append(args, *input.PlayerID)
		}
		if input.GameID != nil {
			where = append(where, "gameID = ?")
			args = append(args, *input.GameID)
		}
		if input.TeamID != nil {
			where = append(where, "teamID = ?")
			args = append(args, *input.TeamID)
		}
		if input.OpponentID != nil {
			where = append(where, "opponentID = ?")
			args = append(args, *input.OpponentID)
		}
		// Seasons         *[]SeasonOption `json:"seasons"` /TODO: VERIFY THAT THIS WORKS
		if input.Seasons != nil {
			qs := strings.Repeat("?,", len(*input.Seasons))
			where = append(where, fmt.Sprintf("season IN (%s)", qs[:len(qs)-1]))
			for _, season := range *input.Seasons {
				args = append(args, season)
			}
		}
		if input.StartDate != nil && input.EndDate != nil {
			if *input.StartDate == *input.EndDate {
				where = append(where, "date >= ? AND date < ?")
				// start of day
				date, err := time.Parse(util.DATE_FORMAT, *input.StartDate)
				if err != nil {
					return nil, fmt.Errorf("failed to parse date: %w", err)
				}
				end := date.Add(24 * time.Hour)
				args = append(args, date, end)
			} else {
				where = append(where, "date BETWEEN ? AND ?")
				args = append(args, *input.StartDate, *input.EndDate)
			}
		}
		if input.StartDate != nil && input.EndDate == nil {
			where = append(where, "date >= ?")
			args = append(args, *input.StartDate)
		}
		if input.StartDate == nil && input.EndDate != nil {
			where = append(where, "date <= ?")
			args = append(args, *input.EndDate)
		}

		if input.HomeOrAway != nil {
			where = append(where, "homeAway = ?")
			args = append(args, *input.HomeOrAway)
		}
		if input.Outcome != nil {
			where = append(where, "outcome = ?")
			args = append(args, *input.Outcome)
		}

		if len(where) > 0 {
			or = append(or, fmt.Sprintf("(%s)", strings.Join(where, " AND ")))
		}
	}
	query := fmt.Sprintf("SELECT * FROM playergames WHERE %s ORDER BY date", strings.Join(or, " OR "))
	err = c.SelectContext(ctx, &games, query, args...)
	if err != nil {
		logrus.Warnf("failed to get playergames using query: %v | %+v", query, args)
		return nil, fmt.Errorf("failed to get player games: %w", err)
	}
	if len(games) == 0 {
		logrus.Warnf("received 0 playergames using query: %v | %+v", query, args)
	}
	if len(inputs) < 5 {
		logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Player Games: %v", len(games), inputs), start))
	} else {
		logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Player Games from %d inputs", len(games), len(inputs)), start))
	}
	return games, nil
}

func (c *SQLClient) GetPropositions(ctx context.Context, propositionFilter *model.PropositionFilter) ([]*model.Proposition, error) {
	c.AddQuery()
	start := time.Now()
	where := []string{}
	args := []interface{}{}
	if propositionFilter.PlayerID != nil {
		where = append(where, "playerID = ?")
		args = append(args, *propositionFilter.PlayerID)
	}
	if propositionFilter.TeamID != nil {
		where = append(where, "teamID = ?")
		args = append(args, *propositionFilter.TeamID)
	}
	// if propositionFilter.GameID != nil {
	// 	where = append(where, "gameID = ?")
	// 	args = append(args, *propositionFilter.GameID)
	// }
	if propositionFilter.StartDate != nil && propositionFilter.EndDate != nil {
		if *propositionFilter.StartDate == *propositionFilter.EndDate {
			where = append(where, "date = ?")
			args = append(args, *propositionFilter.StartDate)
		} else {
			where = append(where, "date BETWEEN ? AND ?")
			args = append(args, *propositionFilter.StartDate, *propositionFilter.EndDate)
		}
	}
	if propositionFilter.StartDate != nil && propositionFilter.EndDate == nil {
		where = append(where, "date >= ?")
		args = append(args, *propositionFilter.StartDate)
	}
	if propositionFilter.StartDate == nil && propositionFilter.EndDate != nil {
		where = append(where, "date <= ?")
		args = append(args, *propositionFilter.EndDate)
	}
	if propositionFilter.PropositionType != nil {
		// TODO: fix this to get the correct type
		where = append(where, "type = ?")
		args = append(args, *propositionFilter.PropositionType)
	}
	query := "SELECT statType, target, sportsbook, lastModified, pg.* FROM propositions pr JOIN playergames pg USING (playerID, gameID)"
	if len(where) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(where, " AND "))
	}
	// query += " LIMIT 10"

	var rawResults []*struct {
		Type         string                 `db:"statType"`
		Target       float64                `db:"target"`
		Sportsbook   model.SportsbookOption `db:"sportsbook"`
		LastModified *time.Time             `db:"lastModified"`
		// PropDate     *time.Time             `db:"propdate"`
		// playerName   string                 `db:"playerName"`
		*model.PlayerGame
	}
	// 	playerID
	// 	gameID
	// 	opponentID
	// 	playerName
	// 	date
	// 	target
	// 	sportsbook
	// 	lastModified
	// 	CreatedAt
	// 	UpdatedAt
	// 	gameID	varchar(255)	NO	PRI
	// playerID	int	NO	PRI
	// opponentID	int	NO	MUL
	// teamID	int	NO	MUL
	// outcome	enum('WIN','LOSS','PENDING')	YES
	// homeAway	enum('HOME','AWAY')	YES
	// date	datetime	NO	MUL
	// season	varchar(255)	YES
	// assists	int	YES
	// defensiveReboundPct	float	YES
	// defensiveRebounds	int	YES
	// offensiveReboundPct	float	YES
	// offensiveRebounds	int	YES
	// effectiveFieldGoalPct	float	YES
	// fieldGoalPct	float	YES
	// fieldGoalsAttempted	int	YES
	// fieldGoalsMade	int	YES
	// freeThrowsAttempted	int	YES
	// freeThrowsMade	int	YES
	// freeThrowPct	float	YES
	// margin	int	YES
	// personalFoulsDrawn	int	YES
	// personalFouls	int	YES
	// points	int	YES		0
	// playoffs	tinyint(1)	YES
	// threePointPct	float	YES
	// threePointersAttempted	int	YES
	// threePointersMade	int	YES
	// rebounds	int	YES
	// trueShootingPct	float	YES
	// turnovers	int	YES
	// blocks	int	YES
	// steals	int	YES
	// usage	float	YES
	// potentialAssists	int	YES
	// assistConversionRate	float	YES
	// assistPct	float	YES
	// minutes	float	YES		0
	// passes	int	YES		0
	// UpdatedAt	datetime	YES		CURRENT_TIMESTAMP	DEFAULT_GENERATED on update CURRENT_TIMESTAMP
	// CreatedAt	datetime	YES		CURRENT_TIMESTAMP	DEFAULT_GENERATED

	err := c.SelectContext(ctx, &rawResults, query, args...)
	if err != nil {
		logrus.Warnf("failed to get propositions using query: %v | %+v", query, args)
		return nil, fmt.Errorf("failed to get propositions: %w", err)
	}
	var propositions []*model.Proposition
	for _, rawResult := range rawResults {
		stat, err := model.NewStat(rawResult.Type)
		if err != nil {
			logrus.Warnf("failed to get stat type: %v", rawResult.Type)
			continue
		}

		proposition := &model.Proposition{
			Game:         rawResult.PlayerGame,
			TypeRaw:      rawResult.Type,
			Target:       rawResult.Target,
			Sportsbook:   rawResult.Sportsbook,
			LastModified: rawResult.LastModified,
			Type:         stat,
			Outcome:      model.PropOutcomePending,
		}
		if proposition.Game != nil && proposition.Game.Outcome != model.GameOutcomePending.String() {
			score := proposition.Game.Score(stat)
			proposition.ActualResult = &score
			if score > proposition.Target {
				proposition.Outcome = model.PropOutcomeOver
			} else if score < proposition.Target {
				proposition.Outcome = model.PropOutcomeUnder
			} else {
				proposition.Outcome = model.PropOutcomePush
			}
			proposition.Accuracy = (score - proposition.Target) / proposition.Target
		}
		propositions = append(propositions, proposition)
	}
	logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Propositions: %v", len(propositions), propositionFilter), start))
	return propositions, nil
}

var teamIDs = map[int]model.Team{}
var teamABRs = map[string]model.Team{}

func (c *SQLClient) GetTeams(ctx context.Context, withGames bool, teamFilters ...*model.TeamFilter) ([]*model.Team, error) {
	c.AddQuery()
	start := time.Now()
	teams := []*model.Team{}
	or := []string{}
	args := []interface{}{}
	for _, teamFilter := range teamFilters {
		where := []string{}
		if teamFilter.TeamID != nil {
			where = append(where, "teamID = ?")
			args = append(args, *teamFilter.TeamID)
		}
		if teamFilter.Name != nil {
			where = append(where, "name = ?")
			args = append(args, *teamFilter.Name)
		}
		if teamFilter.Abbreviation != nil {
			where = append(where, "abbreviation = ?")
			args = append(args, *teamFilter.Abbreviation)
		}
		if len(where) > 0 {
			or = append(or, fmt.Sprintf("(%s)", strings.Join(where, " AND ")))
		}
	}
	if withGames {
		// query := fmt.Sprintf("SELECT * FROM teams JOIN teamgames USING (teamID) WHERE %s", strings.Join(where, " AND "))
		// err := c.Select(&teamwithgames, query, args...)
		panic("get teams with games not implemented") // TODO: Implement
	}
	for _, team := range teamFilters {
		if team.TeamID != nil {
			if t, ok := teamIDs[*team.TeamID]; ok {
				teams = append(teams, &t)
			}
		} else if team.Abbreviation != nil {
			if t, ok := teamABRs[*team.Abbreviation]; ok {
				teams = append(teams, &t)
			}
		}
	}
	if len(teams) == len(teamFilters) {
		return teams, nil
	}
	teams = []*model.Team{}
	query := fmt.Sprintf("SELECT * FROM teams WHERE %s", strings.Join(or, " OR "))
	if len(or) == 0 {
		query = "SELECT * FROM teams"
	}
	// logrus.Warn(query, args)
	err := c.Select(&teams, query, args...)
	if err != nil {
		logrus.Warnf("failed to get teams using query: %v | %+v", query, args)
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}
	for _, team := range teams {
		teamIDs[team.TeamID] = *team
		teamABRs[team.Abbreviation] = *team
	}
	logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Teams: %v", len(teams), teamFilters), start))
	return teams, nil
}

func (c *SQLClient) GetSimilarPlayers(ctx context.Context, similarToPlayerID int, input *model.SimilarPlayerInput, endDate *time.Time) ([]*model.Player, error) {
	c.AddQuery()
	start := time.Now()
	end := endDate.Format(util.DATE_FORMAT)
	// stats := []string{"height", "points", "assists", "rebounds", "offensiveRebounds", "defensiveRebounds", "threePointersMade", "threePointersAttempted"}
	stats := []string{"points", "assists", "weight", "heightInches", "rebounds", "fieldGoalsAttempted", "threePointersMade", "threePointersAttempted", "offensiveRebounds", "defensiveRebounds"}
	summation := make([]string, len(stats))
	avg := make([]string, len(stats))
	std := make([]string, len(stats))
	selector := make([]string, len(stats))
	for i, stat := range stats {
		summation[i] = fmt.Sprintf(`((avg(%[1]s)-AVG_%[2]s) / STD_%[2]s*(avg(%[1]s)-AVG_%[2]s) / STD_%[2]s)`, stat, strings.ToUpper(stat))
		selector[i] = fmt.Sprintf(`avg(%[1]s) AS avg%[1]s`, stat)
		avg[i] = fmt.Sprintf(`avg(%[1]s) as AVG_%[2]s`, stat, strings.ToUpper(stat))
		std[i] = fmt.Sprintf(`stddev(%[1]s) as STD_%[2]s`, stat, strings.ToUpper(stat))
	}
	limit := 6
	if (*input).Limit != 0 {
		limit = input.Limit + 1
	}
	// TODO: Allow similarity to be based off of this seasons, this and last season, or all time
	gameFilter := fmt.Sprintf("date < %s AND season = %s", end, "2022-23")
	query := fmt.Sprintf(`
	SELECT p.name, playerID, count(*) AS games, %[6]s,
		SQRT(%[2]s) AS DISTANCE
	FROM playergames 
		JOIN (SELECT 
			%[3]s 
		FROM playergames JOIN players USING (playerID) WHERE playerID=%[1]d
		AND (%[8]s) ) AS from_player
		JOIN (SELECT 
			%[4]s 
		FROM playergames JOIN players USING (playerID)
		WHERE (%[8]s) ) AS from_league
	JOIN players p USING (playerID) WHERE (%[8]s)
	GROUP BY playerID, AVG_%[5]s, STD_%[5]s
	HAVING avg(points)>0 AND games>10
	ORDER BY DISTANCE ASC
	LIMIT %[7]d;`, similarToPlayerID, strings.Join(summation, "+"), strings.Join(avg, ","), strings.Join(std, ","), strings.ToUpper(stats[0]), strings.Join(selector, ","), limit, gameFilter)
	// logrus.Warn(query)
	playerDistances := []struct {
		Id                     int     `db:"playerID"`
		Name                   string  `db:"name"`
		NGames                 int     `db:"games"`
		Distance               float64 `db:"DISTANCE"`
		Points                 float64 `db:"avgpoints"`
		Height                 float64 `db:"avgheightInches"`
		Weight                 float64 `db:"avgweight"`
		Assists                float64 `db:"avgassists"`
		Rebounds               float64 `db:"avgrebounds"`
		FieldGoalsAttempted    float64 `db:"avgfieldGoalsAttempted"`
		ThreePointersMade      float64 `db:"avgthreePointersMade"`
		ThreePointersAttempted float64 `db:"avgthreePointersAttempted"`
		OffensiveRebounds      float64 `db:"avgoffensiveRebounds"`
		DefensiveRebounds      float64 `db:"avgdefensiveRebounds"`
	}{}
	err := c.Select(&playerDistances, query)
	if err != nil || len(playerDistances) == 0 {
		logrus.Warnf("failed to get similar players using query: %v", query)
		return nil, fmt.Errorf("failed to get similar players: %w", err)
	}
	logrus.Infof("Players most similar to %d based off of: %v", similarToPlayerID, stats)
	logrus.Infof("(%2.2d) %20.20s: %s %s %s %s %s %s %s %s\n", 0, "Player Name (#games)", "DISTANCE", "POINTS", "WEIGHT", "HEIGHT", "ASSISTS", "REBOUNDS", "  FGA", "  3PM  ")
	for i, pDistance := range playerDistances {
		logrus.Infof("(%2.2d) %15.15s (%2.2d): %8.3f %6.2f %6.0f %6.0f %7.1f %8.1f %5.1f %6.2f\n", i, pDistance.Name, pDistance.NGames, pDistance.Distance, pDistance.Points, pDistance.Weight, pDistance.Height, pDistance.Assists, pDistance.Rebounds, pDistance.FieldGoalsAttempted, pDistance.ThreePointersMade)
	}
	playerFilters := make([]*model.PlayerFilter, len(playerDistances))
	// TODO: Allow similarity to be based off of this seasons, this and last season, or all time
	seasons := []model.SeasonOption{}
	seasons = append(seasons, model.SEASON_2022_23)
	for i := range playerDistances {
		pFilter := *input.PlayerPoolFilter
		pFilter.PlayerID = &playerDistances[i].Id
		pFilter.EndDate = &end
		pFilter.Seasons = &seasons
		playerFilters[i] = &pFilter
	}
	players, err := c.GetPlayers(ctx, true, playerFilters...)
	if err != nil {
		return nil, fmt.Errorf("failed to get players from similar player ids: %w", err)
	}
	sort.Slice(players, func(i, j int) bool {
		var iDistance, jDistance float64
		for _, pDistance := range playerDistances {
			if pDistance.Id == players[i].PlayerID {
				iDistance = pDistance.Distance
			}
			if pDistance.Id == players[j].PlayerID {
				jDistance = pDistance.Distance
			}
		}
		return iDistance < jDistance
	})
	logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Similar Players: %s | %v", len(players), players[0].Name, input), start))
	return players[1:], nil
}

func (c *SQLClient) GetSimilarTeams(ctx context.Context, similarToTeamID int, input *model.SimilarTeamInput, endDate string) ([]*model.Team, error) {
	panic("not implemented") // TODO: Implement
}

func (c *SQLClient) GetPropositionsByPlayerGame(ctx context.Context, game model.PlayerGame) ([]*model.Proposition, error) {
	c.AddQuery()
	start := time.Now()
	where := []string{}
	args := []interface{}{}
	if game.PlayerID != 0 {
		where = append(where, "playerID=?")
		args = append(args, game.PlayerID)
	}
	if game.GameID != "" {
		where = append(where, "gameID=?")
		args = append(args, game.GameID)
	}
	query := fmt.Sprintf(`
	SELECT * FROM propositions
	WHERE %s
	`, strings.Join(where, " AND "))
	propositions := []*model.Proposition{}
	err := c.Select(&propositions, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get propositions: %w", err)
	}
	logrus.Info(util.TimeLog(fmt.Sprintf("Query (%d) Propositions from game: %v", len(propositions), game), start))
	return propositions, nil
}
