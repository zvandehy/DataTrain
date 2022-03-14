import {gql} from '@apollo/client';

export function GetPropScore(game, propType) {
    switch (propType.toLowerCase()) {
        case "pts+rebs+asts":
            return game["points"] + game["total_rebounds"] + game["assists"]
        case "rebounds":
            return game["total_rebounds"]
        case "free throws made":
            return game["free_throws_made"]
        case "3-pt made":
            return game["three_pointers_made"]
        case "fantasy score":
            return game["points"] + game["total_rebounds"]*1.2 + game["assists"]*1.5 + game["blocks"]*3 + game["steals"]*3 - game["turnovers"]
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
    // console.log(player)
    // console.log(type)
    let ret = props.find((prop) => player.first_name + " " + player.last_name === prop.playerName && prop.type === type)
    // console.log(ret)
    return ret
}



// GRAPHQL Queries

export const HOME_QUERY = gql`
{
    teams {
        teamID
        abbreviation
    }
    projections(sportsbook: "PrizePicks") {
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
