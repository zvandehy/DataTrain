package database

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/zvandehy/DataTrain/nba_graphql/graph/model"
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

func (c *SQLClient) SavePropositions(ctx context.Context, propositions []*model.DBProposition) (int, error) {
	tx := c.MustBegin()
	for _, proposition := range propositions {
		tx.MustExec("REPLACE INTO propositions (playerID, gameID, playername, date, opponentID, lastModified, sportsbook, statType, target) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", proposition.PlayerID, proposition.GameID, proposition.PlayerName, proposition.Date, proposition.OpponentID, proposition.LastModified, proposition.Sportsbook, proposition.StatType, proposition.Target)
	}
	err := tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return len(propositions), nil
}

func (c *SQLClient) GetPlayers(ctx context.Context, withGames bool, playerFilters ...*model.PlayerFilter) ([]*model.Player, error) {
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
	query := "SELECT * FROM players"
	if len(or) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(or, " OR "))
	}
	players := []*model.Player{}
	err := c.SelectContext(ctx, &players, query, args...)
	if err != nil {
		logrus.Warnf("query: %s", query)
		logrus.Warnf("args: %v", args)
		return nil, fmt.Errorf("failed to get players: %w", err)
	}

	games := []*model.PlayerGame{}
	ids := []interface{}{}
	idqs := []string{}
	for _, player := range players {
		ids = append(ids, player.PlayerID)
		idqs = append(idqs, "?")
	}
	if withGames {
		err = c.SelectContext(ctx, &games, fmt.Sprintf("SELECT * FROM playergames WHERE playerID IN (%s)", strings.Join(idqs, ", ")), ids...)
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
		fmt.Println(len(games))
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
		tx.MustExec("INSERT IGNORE INTO playergames (date, gameID, homeAway, opponentID, teamID, playerID, season, playoffs, outcome) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", g.Date, g.GameID, g.HomeOrAway, g.OpponentID, g.TeamID, g.PlayerID, g.Season, g.Playoffs, g.Outcome)
		//tx.MustExce("REPLACE INTO teamgames") // team
		//tx.MustExce("REPLACE INTO teamgames") // opponent
	}
	err := tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return len(games), nil
}

func (c *SQLClient) GetPlayerGames(ctx context.Context, input *model.GameFilter) (games []*model.PlayerGame, err error) {
	where := []string{}
	args := []interface{}{}
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
		seasons := []string{}
		for _, season := range *input.Seasons {
			seasons = append(seasons, string(season))
		}
		s := fmt.Sprintf("season IN (%s)", strings.Join(seasons, ","))
		where = append(where, s)
		args = append(args, seasons)
	}
	if input.StartDate != nil && input.EndDate != nil {
		if *input.StartDate == *input.EndDate {
			where = append(where, "date = ?")
			args = append(args, *input.StartDate)
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

	// OpponentMatch   *bool           `json:"opponentMatch"`
	// GameType        *GameType       `json:"gameType"`
	// GameTypeMatch   *bool           `json:"gameTypeMatch"`
	// HomeOrAwayMatch *bool           `json:"homeOrAwayMatch"`
	// StatFilters     *[]*StatFilter  `json:"statFilters"`
	// LastX           *int            `json:"lastX"`

	if input.LastX != nil || input.StatFilters != nil || input.HomeOrAwayMatch != nil || input.GameTypeMatch != nil || input.OpponentMatch != nil || input.GameType != nil {
		panic("game filter not implemented") // TODO: IMPLEMENT
	}

	query := fmt.Sprintf("SELECT * FROM playergames WHERE %s", strings.Join(where, " AND "))
	err = c.SelectContext(ctx, &games, query, args...)
	if err != nil {
		logrus.Warnf("failed to get playergames using query: %v | %+v", query, args)
		return nil, fmt.Errorf("failed to get player games: %w", err)
	}
	return games, nil
}

func (c *SQLClient) GetPropositions(ctx context.Context, propositionFilter *model.PropositionFilter) ([]*model.Proposition, error) {
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
	query := fmt.Sprintf("SELECT * FROM propositions WHERE %s", strings.Join(where, " AND "))
	var propositions []*model.Proposition
	err := c.SelectContext(ctx, &propositions, query, args...)
	if err != nil {
		logrus.Warnf("failed to get propositions using query: %v | %+v", query, args)
		return nil, fmt.Errorf("failed to get propositions: %w", err)
	}
	return propositions, nil
}

func (c *SQLClient) GetTeams(ctx context.Context, withGames bool, teamFilters ...*model.TeamFilter) ([]*model.Team, error) {
	if len(teamFilters) > 1 {
		panic("multiple team filter not implemented") // TODO: Implement
	}
	teams := []model.Team{}
	where := []string{}
	args := []interface{}{}
	if len(teamFilters) > 0 {
		if (teamFilters)[0].TeamID != nil {
			where = append(where, "teamID = ?")
			args = append(args, *(teamFilters)[0].TeamID)
		}
		if (teamFilters)[0].Name != nil {
			where = append(where, "name = ?")
			args = append(args, *(teamFilters)[0].Name)
		}
		if (teamFilters)[0].Abbreviation != nil {
			where = append(where, "abbreviation = ?")
			args = append(args, *(teamFilters)[0].Abbreviation)
		}
	}
	if withGames {
		// query := fmt.Sprintf("SELECT * FROM teams JOIN teamgames USING (teamID) WHERE %s", strings.Join(where, " AND "))
		// err := c.Select(&teamwithgames, query, args...)
		panic("get teams with games not implemented") // TODO: Implement
	}
	query := fmt.Sprintf("SELECT * FROM teams WHERE %s", strings.Join(where, " AND "))
	if len(where) == 0 {
		query = "SELECT * FROM teams"
	}
	logrus.Debug(query, args)
	err := c.Select(&teams, query, args...)
	if err != nil {
		logrus.Warnf("failed to get teams using query: %v | %+v", query, args)
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}
	ret := []*model.Team{}
	for _, team := range teams {
		ret = append(ret, &team)
	}
	return ret, nil
	// return , nil
}

func (c *SQLClient) GetSimilarPlayers(ctx context.Context, similarToPlayerID int, input *model.SimilarPlayerInput, endDate string) ([]*model.Player, error) {
	logrus.Info("start getting similar players")
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
		limit = input.Limit
	}
	query := fmt.Sprintf(`
	SELECT p.name, playerID, count(*) AS games, %[6]s,
		SQRT(%[2]s) AS DISTANCE
	FROM playergames 
		JOIN (SELECT 
			%[3]s 
		FROM playergames JOIN players USING (playerID) WHERE playerID=%[1]d
		AND (playoffs=False OR playoffs=TRUE) ) AS from_player
		JOIN (SELECT 
			%[4]s 
		FROM playergames JOIN players USING (playerID)
		WHERE (playoffs=False OR playoffs=TRUE) ) AS from_league
	JOIN players p USING (playerID) WHERE (playoffs=False OR playoffs=TRUE) AND playerID <> %[1]d
	GROUP BY playerID, AVG_%[5]s, STD_%[5]s
	HAVING avg(points)>0 AND games>10
	ORDER BY DISTANCE ASC
	LIMIT %[7]d;`, similarToPlayerID, strings.Join(summation, "+"), strings.Join(avg, ","), strings.Join(std, ","), strings.ToUpper(stats[0]), strings.Join(selector, ","), limit)

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
	for i := range playerDistances {
		playerFilters[i] = &model.PlayerFilter{
			PlayerID: &playerDistances[i].Id,
		}
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
	return players, nil
}

func (c *SQLClient) GetSimilarTeams(ctx context.Context, similarToTeamID int, input *model.SimilarTeamInput, endDate string) ([]*model.Team, error) {
	panic("not implemented") // TODO: Implement
}

func (c *SQLClient) GetPropositionsByPlayerGame(ctx context.Context, game model.PlayerGame) ([]*model.Proposition, error) {
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
	return propositions, nil
}
