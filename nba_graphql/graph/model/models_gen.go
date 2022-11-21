// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AverageStats struct {
	Assists                float64 `json:"assists"`
	Blocks                 float64 `json:"blocks"`
	DefensiveRebounds      float64 `json:"defensive_rebounds"`
	FieldGoalsAttempted    float64 `json:"field_goals_attempted"`
	FieldGoalsMade         float64 `json:"field_goals_made"`
	FreeThrowsAttempted    float64 `json:"free_throws_attempted"`
	FreeThrowsMade         float64 `json:"free_throws_made"`
	GamesPlayed            float64 `json:"games_played"`
	Height                 float64 `json:"height"`
	Minutes                float64 `json:"minutes"`
	OffensiveRebounds      float64 `json:"offensive_rebounds"`
	PersonalFoulsDrawn     float64 `json:"personal_fouls_drawn"`
	PersonalFouls          float64 `json:"personal_fouls"`
	Points                 float64 `json:"points"`
	Rebounds               float64 `json:"rebounds"`
	Steals                 float64 `json:"steals"`
	ThreePointersAttempted float64 `json:"three_pointers_attempted"`
	ThreePointersMade      float64 `json:"three_pointers_made"`
	Turnovers              float64 `json:"turnovers"`
	Weight                 float64 `json:"weight"`
	FantasyScore           float64 `json:"fantasy_score"`
	PointsAssists          float64 `json:"points_assists"`
	PointsRebounds         float64 `json:"points_rebounds"`
	PointsReboundsAssists  float64 `json:"points_rebounds_assists"`
	ReboundsAssists        float64 `json:"rebounds_assists"`
	BlocksSteals           float64 `json:"blocks_steals"`
	DoubleDouble           float64 `json:"double_double"`
	Passes                 float64 `json:"passes"`
	PotentialAssists       float64 `json:"potential_assists"`
}

type GameBreakdownInput struct {
	Name   string      `json:"name"`
	Filter *GameFilter `json:"filter"`
	Weight float64     `json:"weight"`
}

type GamePrediction struct {
	Estimation         *AverageStats             `json:"estimation"`
	EstimationAccuracy *AverageStats             `json:"estimationAccuracy"`
	Fragments          []*GamePredictionFragment `json:"fragments"`
}

type GamePredictionFragment struct {
	Name         string         `json:"name"`
	Derived      *AverageStats  `json:"derived"`
	DerivedGames []*PlayerGame  `json:"derivedGames"`
	Base         *AverageStats  `json:"base"`
	PctChange    *AverageStats  `json:"pctChange"`
	Weight       float64        `json:"weight"`
	Propositions []*Proposition `json:"propositions"`
}

type ModelInput struct {
	Model              *string               `json:"model"`
	GameBreakdowns     []*GameBreakdownInput `json:"gameBreakdowns"`
	SimilarPlayerInput *SimilarPlayerInput   `json:"similarPlayerInput"`
	SimilarTeamInput   *SimilarTeamInput     `json:"similarTeamInput"`
}

type PropBreakdown struct {
	Name           string        `json:"name"`
	Over           int           `json:"over"`
	Under          int           `json:"under"`
	Push           int           `json:"push"`
	OverPct        float64       `json:"overPct"`
	UnderPct       float64       `json:"underPct"`
	PushPct        float64       `json:"pushPct"`
	DerivedAverage float64       `json:"derivedAverage"`
	Weight         float64       `json:"weight"`
	PctChange      float64       `json:"pctChange"`
	Base           float64       `json:"base"`
	DerivedGames   []*PlayerGame `json:"derivedGames"`
}

type PropPrediction struct {
	Estimation         float64          `json:"estimation"`
	EstimationAccuracy *float64         `json:"estimationAccuracy"`
	Significance       float64          `json:"significance"`
	CumulativeOver     int              `json:"cumulativeOver"`
	CumulativeUnder    int              `json:"cumulativeUnder"`
	CumulativePush     int              `json:"cumulativePush"`
	CumulativeOverPct  float64          `json:"cumulativeOverPct"`
	CumulativeUnderPct float64          `json:"cumulativeUnderPct"`
	CumulativePushPct  float64          `json:"cumulativePushPct"`
	Wager              Wager            `json:"wager"`
	WagerOutcome       WagerOutcome     `json:"wagerOutcome"`
	Breakdowns         []*PropBreakdown `json:"breakdowns"`
}

type PropositionFilter struct {
	Sportsbook      *SportsbookOption `json:"sportsbook"`
	PropositionType *Stat             `json:"propositionType"`
	StartDate       *string           `json:"startDate"`
	EndDate         *string           `json:"endDate"`
	PlayerID        *int              `json:"PlayerID"`
	PlayerName      *string           `json:"PlayerName"`
	TeamID          *int              `json:"TeamID"`
	TeamName        *string           `json:"TeamName"`
}

type SimilarPlayerInput struct {
	Limit            int           `json:"limit"`
	StatsOfInterest  []Stat        `json:"statsOfInterest"`
	PlayerPoolFilter *PlayerFilter `json:"playerPoolFilter"`
	Weight           float64       `json:"weight"`
}

type SimilarTeamInput struct {
	Limit           int           `json:"limit"`
	StatsOfInterest []Stat        `json:"statsOfInterest"`
	TeamPoolFilter  []*TeamFilter `json:"teamPoolFilter"`
	Period          *Period       `json:"period"`
	Weight          float64       `json:"weight"`
}

type StatFilter struct {
	Period   *Period  `json:"period"`
	Stat     Stat     `json:"stat"`
	Mode     StatMode `json:"mode"`
	Operator Operator `json:"operator"`
	Value    float64  `json:"value"`
}

type TeamFilter struct {
	Name         *string `json:"name"`
	TeamID       *int    `json:"teamID"`
	Abbreviation *string `json:"abbreviation"`
}

type GameOutcome string

const (
	GameOutcomeWin     GameOutcome = "WIN"
	GameOutcomeLoss    GameOutcome = "LOSS"
	GameOutcomePending GameOutcome = "PENDING"
)

var AllGameOutcome = []GameOutcome{
	GameOutcomeWin,
	GameOutcomeLoss,
	GameOutcomePending,
}

func (e GameOutcome) IsValid() bool {
	switch e {
	case GameOutcomeWin, GameOutcomeLoss, GameOutcomePending:
		return true
	}
	return false
}

func (e GameOutcome) String() string {
	return string(e)
}

func (e *GameOutcome) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GameOutcome(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GameOutcome", str)
	}
	return nil
}

func (e GameOutcome) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type GameType string

const (
	GameTypeRegularSeason GameType = "REGULAR_SEASON"
	GameTypePlayoffs      GameType = "PLAYOFFS"
)

var AllGameType = []GameType{
	GameTypeRegularSeason,
	GameTypePlayoffs,
}

func (e GameType) IsValid() bool {
	switch e {
	case GameTypeRegularSeason, GameTypePlayoffs:
		return true
	}
	return false
}

func (e GameType) String() string {
	return string(e)
}

func (e *GameType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GameType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GameType", str)
	}
	return nil
}

func (e GameType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type HomeOrAway string

const (
	HomeOrAwayHome HomeOrAway = "HOME"
	HomeOrAwayAway HomeOrAway = "AWAY"
)

var AllHomeOrAway = []HomeOrAway{
	HomeOrAwayHome,
	HomeOrAwayAway,
}

func (e HomeOrAway) IsValid() bool {
	switch e {
	case HomeOrAwayHome, HomeOrAwayAway:
		return true
	}
	return false
}

func (e HomeOrAway) String() string {
	return string(e)
}

func (e *HomeOrAway) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HomeOrAway(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HomeOrAway", str)
	}
	return nil
}

func (e HomeOrAway) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Operator string

const (
	OperatorGt  Operator = "GT"
	OperatorGte Operator = "GTE"
	OperatorLt  Operator = "LT"
	OperatorLte Operator = "LTE"
	OperatorEq  Operator = "EQ"
	OperatorNeq Operator = "NEQ"
)

var AllOperator = []Operator{
	OperatorGt,
	OperatorGte,
	OperatorLt,
	OperatorLte,
	OperatorEq,
	OperatorNeq,
}

func (e Operator) IsValid() bool {
	switch e {
	case OperatorGt, OperatorGte, OperatorLt, OperatorLte, OperatorEq, OperatorNeq:
		return true
	}
	return false
}

func (e Operator) String() string {
	return string(e)
}

func (e *Operator) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Operator(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Operator", str)
	}
	return nil
}

func (e Operator) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Position string

const (
	PositionG  Position = "G"
	PositionF  Position = "F"
	PositionC  Position = "C"
	PositionGF Position = "G_F"
	PositionFG Position = "F_G"
	PositionFC Position = "F_C"
	PositionCF Position = "C_F"
)

var AllPosition = []Position{
	PositionG,
	PositionF,
	PositionC,
	PositionGF,
	PositionFG,
	PositionFC,
	PositionCF,
}

func (e Position) IsValid() bool {
	switch e {
	case PositionG, PositionF, PositionC, PositionGF, PositionFG, PositionFC, PositionCF:
		return true
	}
	return false
}

func (e Position) String() string {
	return string(e)
}

func (e *Position) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Position(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Position", str)
	}
	return nil
}

func (e Position) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PropOutcome string

const (
	PropOutcomeOver    PropOutcome = "OVER"
	PropOutcomeUnder   PropOutcome = "UNDER"
	PropOutcomePush    PropOutcome = "PUSH"
	PropOutcomePending PropOutcome = "PENDING"
)

var AllPropOutcome = []PropOutcome{
	PropOutcomeOver,
	PropOutcomeUnder,
	PropOutcomePush,
	PropOutcomePending,
}

func (e PropOutcome) IsValid() bool {
	switch e {
	case PropOutcomeOver, PropOutcomeUnder, PropOutcomePush, PropOutcomePending:
		return true
	}
	return false
}

func (e PropOutcome) String() string {
	return string(e)
}

func (e *PropOutcome) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PropOutcome(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PropOutcome", str)
	}
	return nil
}

func (e PropOutcome) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SportsbookOption string

const (
	SportsbookOptionPrizePicks      SportsbookOption = "PrizePicks"
	SportsbookOptionUnderdogFantasy SportsbookOption = "UnderdogFantasy"
)

var AllSportsbookOption = []SportsbookOption{
	SportsbookOptionPrizePicks,
	SportsbookOptionUnderdogFantasy,
}

func (e SportsbookOption) IsValid() bool {
	switch e {
	case SportsbookOptionPrizePicks, SportsbookOptionUnderdogFantasy:
		return true
	}
	return false
}

func (e SportsbookOption) String() string {
	return string(e)
}

func (e *SportsbookOption) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SportsbookOption(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SportsbookOption", str)
	}
	return nil
}

func (e SportsbookOption) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type StatMode string

const (
	StatModePerGame   StatMode = "PER_GAME"
	StatModePer36     StatMode = "PER_36"
	StatModePerMinute StatMode = "PER_MINUTE"
	StatModeTotal     StatMode = "TOTAL"
)

var AllStatMode = []StatMode{
	StatModePerGame,
	StatModePer36,
	StatModePerMinute,
	StatModeTotal,
}

func (e StatMode) IsValid() bool {
	switch e {
	case StatModePerGame, StatModePer36, StatModePerMinute, StatModeTotal:
		return true
	}
	return false
}

func (e StatMode) String() string {
	return string(e)
}

func (e *StatMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StatMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StatMode", str)
	}
	return nil
}

func (e StatMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Wager string

const (
	WagerOver  Wager = "OVER"
	WagerUnder Wager = "UNDER"
)

var AllWager = []Wager{
	WagerOver,
	WagerUnder,
}

func (e Wager) IsValid() bool {
	switch e {
	case WagerOver, WagerUnder:
		return true
	}
	return false
}

func (e Wager) String() string {
	return string(e)
}

func (e *Wager) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Wager(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Wager", str)
	}
	return nil
}

func (e Wager) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type WagerOutcome string

const (
	WagerOutcomeHit     WagerOutcome = "HIT"
	WagerOutcomeMiss    WagerOutcome = "MISS"
	WagerOutcomePush    WagerOutcome = "PUSH"
	WagerOutcomePending WagerOutcome = "PENDING"
)

var AllWagerOutcome = []WagerOutcome{
	WagerOutcomeHit,
	WagerOutcomeMiss,
	WagerOutcomePush,
	WagerOutcomePending,
}

func (e WagerOutcome) IsValid() bool {
	switch e {
	case WagerOutcomeHit, WagerOutcomeMiss, WagerOutcomePush, WagerOutcomePending:
		return true
	}
	return false
}

func (e WagerOutcome) String() string {
	return string(e)
}

func (e *WagerOutcome) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = WagerOutcome(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid WagerOutcome", str)
	}
	return nil
}

func (e WagerOutcome) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
