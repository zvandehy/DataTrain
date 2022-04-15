import {gql} from '@apollo/client';
import {round} from 'mathjs';

export function FormatDate(date) {
    const yyyy = date.getFullYear();
    const mm = String(date.getMonth() + 1).padStart(2, '0');
    const dd = String(date.getDate()).padStart(2, '0');
    let ret = `${yyyy}-${mm}-${dd}`
    return ret
}

export function ParseDate(date) {
    return Date.parse(date.replace(/-/g, '\/').replace(/T.+/, ''));
}

//returns negative if 1 is less than 2
export function CompareDates(date1, date2) {
        var a = new Date(date1);
        var b = new Date(date2);
        return a-b;
}

export function GetSelectableTeams(teams) {
    let ret  = teams.map((team) => ({
        // required: what to show to the user
        label: team.abbreviation,
        // required: key to identify the item within the array
        key: team.teamID,
      }));
    ret.push({label:"ANY", key:"ANY"});
    ret.sort((a,b) => {return a.label > b.label ? 1 : -1});
    return ret
}

export function GetPropScore(game, propType) {
    switch (propType.toLowerCase()) {
        case "pts+rebs+asts":
            return game["points"] + game["total_rebounds"] + game["assists"]
        case "rebounds":
            return game["total_rebounds"]
        case "free throws made":
            return game["free_throws_made"]
        case "3-pt made":
        case "3 pointers":
            return game["three_pointers_made"]
        case "fantasy score":
            return round(game["points"] + game["total_rebounds"]*1.2 + game["assists"]*1.5 + game["blocks"]*3 + game["steals"]*3 - game["turnovers"],2)
        case "blks+stls":
            return game["blocks"] + game["steals"]
        case "double-double":
            let doubles = 0
            if (game["points"] >= 10) doubles++
            if (game["total_rebounds"] >= 10) doubles++
            if (game["assists"] >= 10) doubles++
            if (game["blocks"] >= 10) doubles++
            if (game["steals"] >= 10) doubles++
            if (doubles >= 2) return 1
            return 0
        case "minutes":
            return parseInt(game.minutes.split(":")[0])
        case "field_goal_percentage":
            return game["field_goal_percentage"]
        default:
            return game[propType.toLowerCase()] ?? 0
    }
}

export function GetShortType(type) {
    switch (type.toLowerCase()) {
        case "points":
            return "PTS"
        case "rebounds":
            return "REB"
        case "assists":
            return "AST"
        case "fantasy score":
            return "FAN"
        case "pts+rebs+asts":
            return "PRA"
        case "free throws made":
            return "FTM"
        default:
            return type
    }
}

export function showPlayerPreview(player, props, type) {
    let ret = props.find((prop) => player.first_name + " " + player.last_name === prop.playerName && prop.type === type)
    return ret
}

export function GetColor(type, num) {
    if (type.toLowerCase() === "pct") {
        if(num >= 60.0) {
            return 'high';
        } else if(num >= 50.0) {
            return 'med';
        }
        return 'low';
    }
    if (type.toLowerCase() === "over" && num.toLowerCase() === "over") {
        return "high"
    }
    return 'low';
}

// GRAPHQL Queries

export const HOME_QUERY = gql`
query HOME($date: String!){
    teams {
        teamID
        abbreviation
    }
    projections(input: {sportsbook: "PrizePicks", startDate: $date, endDate: $date}) {
        player {
            name
            position
            playerID
            currentTeam {
                abbreviation
                teamID
            }
            games(input: {season:"2021-22"}) {
                season
                date
                gameID
                opponent {
                    name
                    teamID
                    abbreviation
                }
                points
                assists
                total_rebounds
                defensive_rebounds
                offensive_rebounds
                three_pointers_attempted
                three_pointers_made
                free_throws_attempted
                free_throws_made
                minutes
                blocks
                turnovers
                steals
            }
        }
        opponent {
            abbreviation
            teamID
        }
        targets {
            target
            type
        }
    }
}`

//TODO: Add server support to query for single player
export const PLAYER_PREVIEW_QUERY = gql`
query games($playerID: Int!){
    games(playerID: $playerID) {
        assists
        points
        total_rebounds
        season
    }
}`

export const PLAYERGAMES_QUERY = gql`
query games($playerID: Int!){
    games(playerID: $playerID) {
        assist_percentage
        assists
        date
        defensive_rebound_percentage
        defensive_rebounds
        effective_field_goal_percentage
        field_goal_percentage
        field_goals_attempted
        field_goals_made
        free_throws_attempted
        free_throws_made
        free_throws_percentage
        gameID
        home_or_away
        minutes
        offensive_rebound_percentage
        offensive_rebounds
        opponent
        personal_fouls
        personal_fouls_drawn
        points
        season
        three_point_percentage
        three_pointers_attempted
        three_pointers_made
        total_rebounds
        true_shooting_percentage
        turnovers
        usage
    }
    players {
        playerID
        first_name
        last_name
        seasons
        teamABR
    }
}`
