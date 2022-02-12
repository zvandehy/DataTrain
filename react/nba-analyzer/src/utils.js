import {gql} from '@apollo/client';

export function GetPropScore(game, propType) {
    switch (propType.toLowerCase()) {
        case "pts+rebs+asts":
            return game["points"] + game["total_rebounds"] + game["assists"]
        case "rebounds":
            return game["total_rebounds"]
        default:
            return game[propType]
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
  prizepicks {
    player {
        playerID
        first_name
        last_name
        seasons
        teamABR
    }
    playerGames {
        points
        assists
        assist_percentage
        total_rebounds
        defensive_rebounds
        offensive_rebounds
        three_pointers_attempted
        three_pointers_made
        free_throws_attempted
        free_throws_made
        season
        minutes
        date
    }
    type
    target
    opponent
  }
  teams {
      abbreviation
      teamID
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
