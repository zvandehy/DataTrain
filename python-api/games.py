# %%
from pymongo import MongoClient
url = 'mongodb://zeke:helloworld@104.248.53.50:27017/nba'
client = MongoClient(url)
db = client.nba

# %%
from nba_api.stats.static import players
from nba_api.stats.endpoints import commonplayerinfo
import time


# nba_players = players.get_active_players()
# print(nba_players[0]['id'])
# time.sleep(1)
# data = commonplayerinfo.CommonPlayerInfo(nba_players[0]['id'])
# print(data.common_player_info.get_dict())
# time.sleep(1)

# %%
from nba_api.stats.static import teams

def lookupTeam(abbreviation):
    original = abbreviation
    if abbreviation == "SEA":
        abbreviation = "OKC"
    if abbreviation == "NOH" or abbreviation == "NOK":
        abbreviation = "NOP"
    if abbreviation == "NJN":
        abbreviation = "BKN"
    ret = [team for team in teams.get_teams() if team['abbreviation'] == abbreviation]
    if len(ret) == 0:
        print(original)
    return ret[0]

# for team in teams.get_teams():
#     team_data = {
#         "id": team['id'],
#         'name': team['nickname'],
#         'abbreviation': team['abbreviation'],
#         "city": team["city"],
#     }
#     if team["abbreviation"] == "OKC":
#         team_data["alternate_abbreviations"] = ["SEA"]
#     if team["abbreviation"] == "NOP":
#         team_data["alternate_abbreviations"] = ["NOK", "NOH"]
#     if team["abbreviation"] == "BKN":
#         team_data["alternate_abbreviations"] = ["NJN"]
#     db.teams.insert_one(team_data)


# %%
# from nba_api.stats.endpoints import commonallplayers
# from nba_api.stats.endpoints import playercareerstats
# from nba_api.stats.library.parameters import PerModeSimple

# all_players = commonallplayers.CommonAllPlayers()
# data = all_players.get_normalized_dict()
# count = 0
# for player in data['CommonAllPlayers']:
#     # if player['TO_YEAR'] == '2021':
#         # print(player)
#     if player["ROSTERSTATUS"] == 1:
#         names = player["DISPLAY_FIRST_LAST"].split(" ", 1)
#         player_data = {
#             "id":player["PERSON_ID"],
#             "first_name":names[0],
#             "last_name":names[1],
#             "seasons":[],
#             "games": []
#         }
#         if db.players.count_documents({"id":player["PERSON_ID"]}) > 0:
#             continue
#         player_career = playercareerstats.PlayerCareerStats(player["PERSON_ID"], per_mode36=PerModeSimple.per_game)
#         regular_seasons = player_career.get_normalized_dict()["SeasonTotalsRegularSeason"]
#         time.sleep(.8)
#         for season in regular_seasons:
#             if season["MIN"] > 10.0:
#                 player_data["seasons"].append(season["SEASON_ID"])
#         if len(player_data["seasons"]) > 0:
#             db.players.insert_one(player_data)
#             count+=1
#             print(str(count) + ": " + player["DISPLAY_FIRST_LAST"])
    

# %%
# from nba_api.stats.endpoints import leaguegamefinder
# from nba_api.stats.endpoints import boxscoreadvancedv2
# from nba_api.stats.endpoints import boxscoremiscv2

# sleeptime = .6

# players = db.players.find({})
# for player in players:
#     time.sleep(sleeptime)
#     for season in player["seasons"]:
#         if int(season.split("-")[0]) < 2020:
#             print("skip season: ", season)
#             continue
#         gamefinder = leaguegamefinder.LeagueGameFinder(player_id_nullable=player["id"], league_id_nullable="00", season_type_nullable="Regular Season", season_nullable=season)
#         games = gamefinder.get_normalized_dict()["LeagueGameFinderResults"]
#         for game in games:
#             if db.games.count_documents({"player": player["id"], "game":game["GAME_ID"]}) > 0:
#                 continue
#             game_data = {
#                 "game":game["GAME_ID"],
#                 "player":player["id"],
#                 "team": lookupTeam(game["MATCHUP"].split(" ")[0])['id'],
#                 "opponent":lookupTeam(game["MATCHUP"].split(" ")[2])['id'],
#                 "date": game["GAME_DATE"],
#                 "season": game["SEASON_ID"],
#                 "minutes": game["MIN"],
#                 "points": game["PTS"],
#                 "field_goals_made":game["FGM"], 
#                 "field_goals_attempted":game["FGA"], 
#                 "field_goal_percentage":game["FG_PCT"], 
#                 "assists":game["AST"], 
#                 "turnovers":game["TOV"], 
#                 "defensive_rebounds":game["DREB"], 
#                 "offensive_rebounds":game["OREB"], 
#                 "total_rebounds":game["REB"], 
#                 "free_throws_made":game["FTM"], 
#                 "free_throws_attempted":game["FTA"],
#                 "free_throws_percentage":game["FT_PCT"],
#                 "personal_fouls":game["PF"], 
#                 "three_pointers_made":game["FG3M"], 
#                 "three_pointers_attempted":game["FG3A"], 
#                 "three_point_percentage":game["FG3_PCT"], 
#             }
#             time.sleep(sleeptime)
#             advanced = boxscoreadvancedv2.BoxScoreAdvancedV2(game_data["game"]).get_normalized_dict()
#             advanced = advanced["PlayerStats"]
#             for game in advanced:
#                 if game["PLAYER_ID"] == player['id']:
#                     game_data["assist_percentage"] = game["AST_PCT"]
#                     game_data["defensive_rebound_percentage"] = game["DREB_PCT"]
#                     game_data["offensive_rebound_percentage"] = game["OREB_PCT"]
#                     game_data["usage"] = game["USG_PCT"]
#                     game_data["effective_field_goal_percentage"] = game["EFG_PCT"]
#                     game_data["true_shooting_percentage"] = game["TS_PCT"]
#             time.sleep(sleeptime)
#             misc = boxscoremiscv2.BoxScoreMiscV2(game_data["game"]).get_normalized_dict()
#             misc = misc["sqlPlayersMisc"]
#             for game in misc:
#                 if game["PLAYER_ID"] == player['id']:
#                     game_data["personal_fouls_drawn"] = game["PFD"]
#             db.games.insert_one(game_data)
#             print(player["first_name"], game_data["season"], game_data["date"])

# %%
# get other players from same game
from nba_api.stats.endpoints import leaguegamefinder
from nba_api.stats.endpoints import boxscoretraditionalv2
from nba_api.stats.endpoints import boxscoreadvancedv2
from nba_api.stats.endpoints import boxscoremiscv2
import os

print(os.getpid())

sleeptime = .6

players = db.players.find()
for player in players:
    print("From player: ", player["first_name"])
    time.sleep(sleeptime)
    gamefinder = leaguegamefinder.LeagueGameFinder(player_id_nullable=player["id"], league_id_nullable="00", season_type_nullable="Regular Season")
    games = gamefinder.get_normalized_dict()["LeagueGameFinderResults"]
    print(len(games), " total games")
    games = [game for game in games if int(game["SEASON_ID"][1:]) >= 2015 and int(game["SEASON_ID"][1:]) <= 2019]
    print(len(games), " recent games")

    for g in games:
        # print(g)
        transform_season_id = g["SEASON_ID"][1:]+"-"+str(int(g["SEASON_ID"][1:])+1)[2:]
        #skip any season not in player.seasons (they didn't play > 10 min)
        if not transform_season_id in player["seasons"]:
            continue
        if db.games.count_documents({"player": player["id"], "game":g["GAME_ID"]}) > 0:
            continue
        time.sleep(sleeptime)
        boxscorefinder = boxscoretraditionalv2.BoxScoreTraditionalV2(g["GAME_ID"]).get_normalized_dict()
        time.sleep(sleeptime)
        advanced = boxscoreadvancedv2.BoxScoreAdvancedV2(g["GAME_ID"]).get_normalized_dict()["PlayerStats"]
        time.sleep(sleeptime)
        misc = boxscoremiscv2.BoxScoreMiscV2(g["GAME_ID"]).get_normalized_dict()["sqlPlayersMisc"]
        for game in boxscorefinder["PlayerStats"]:
            #if the player played in te game, and the player exists with the provided season (averaged > 10 min in the season), and this player's game stats haven't been added already
            if game["MIN"] != None and db.players.count_documents({"id": game["PLAYER_ID"], "seasons":transform_season_id}) > 0 and db.games.count_documents({"player": game["PLAYER_ID"], "game":game["GAME_ID"]}) == 0:
                team1 = lookupTeam(g["MATCHUP"].split(" ")[0])['id']
                team2 = lookupTeam(g["MATCHUP"].split(" ")[2])['id']
                home_or_away = ""
                team = ""
                opponent = ""
                if team1 == game["TEAM_ID"]:
                    team = team1
                    opponent = team2
                    if g["MATCHUP"].split(" ")[1] == "@":
                        #team1 is away
                        home_or_away = "away"
                    else:
                        #team1 is home
                        home_or_away = "home"
                elif team2 == game["TEAM_ID"]:
                    team = team2
                    opponent = team1
                    if g["MATCHUP"].split(" ")[1] == "@":
                        #team1 is away
                        home_or_away = "home"
                    else:
                        #team1 is home
                        home_or_away = "away"
                else:
                    print("ERROR: Team not found")
                game_data = {
                    "game":game["GAME_ID"],
                    "player":game["PLAYER_ID"],
                    "team": team,
                    "home_or_away":home_or_away,
                    "opponent":opponent,
                    "date": g["GAME_DATE"],
                    "season": transform_season_id,
                    "minutes": game["MIN"],
                    "points": game["PTS"],
                    "field_goals_made":game["FGM"], 
                    "field_goals_attempted":game["FGA"], 
                    "field_goal_percentage":game["FG_PCT"], 
                    "assists":game["AST"], 
                    "turnovers":game["TO"], 
                    "defensive_rebounds":game["DREB"], 
                    "offensive_rebounds":game["OREB"], 
                    "total_rebounds":game["REB"], 
                    "free_throws_made":game["FTM"], 
                    "free_throws_attempted":game["FTA"],
                    "free_throws_percentage":game["FT_PCT"],
                    "personal_fouls":game["PF"], 
                    "three_pointers_made":game["FG3M"], 
                    "three_pointers_attempted":game["FG3A"], 
                    "three_point_percentage":game["FG3_PCT"], 
                }
                for advanced_game in advanced:
                    if game_data['player'] == advanced_game["PLAYER_ID"]:
                        game_data["assist_percentage"] = advanced_game["AST_PCT"]
                        game_data["defensive_rebound_percentage"] = advanced_game["DREB_PCT"]
                        game_data["offensive_rebound_percentage"] = advanced_game["OREB_PCT"]
                        game_data["usage"] = advanced_game["USG_PCT"]
                        game_data["effective_field_goal_percentage"] = advanced_game["EFG_PCT"]
                        game_data["true_shooting_percentage"] = advanced_game["TS_PCT"]
                for misc_game in misc:
                    if game_data['player'] == misc_game["PLAYER_ID"]:
                        game_data["personal_fouls_drawn"] = misc_game["PFD"]
                db.games.insert_one(game_data)
                print(game["PLAYER_NAME"], game_data["season"], game_data["date"], game_data["minutes"])

# %%
# # get other players from same game
# from nba_api.stats.endpoints import leaguegamefinder
# from nba_api.stats.endpoints import boxscoretraditionalv2
# from nba_api.stats.endpoints import boxscoreadvancedv2
# from nba_api.stats.endpoints import boxscoremiscv2
# import sys

# sleeptime = .6

# teams = db.teams.find()
# for team in teams:
#     if team["abbreviation"] != "DEN":
#         continue
#     print("From team: ", team["name"])
#     time.sleep(sleeptime)
#     gamefinder = leaguegamefinder.LeagueGameFinder(team_id_nullable=team["id"], league_id_nullable="00", season_type_nullable="Regular Season")
#     games = gamefinder.get_normalized_dict()["LeagueGameFinderResults"]
#     print(len(games), " total games")
#     games = [game for game in games if int(game["SEASON_ID"][1:]) >= 2020]
#     print(len(games), " recent games")

#     for g in games:
#         transform_season_id = g["SEASON_ID"][1:]+"-"+str(int(g["SEASON_ID"][1:])+1)[2:]

#         if db.games.count_documents({"player": player["id"], "game":g["GAME_ID"]}) > 0:
#             continue

#         time.sleep(sleeptime)
#         boxscorefinder = boxscoretraditionalv2.BoxScoreTraditionalV2(g["GAME_ID"]).get_normalized_dict()
#         time.sleep(sleeptime)
#         advanced = boxscoreadvancedv2.BoxScoreAdvancedV2(g["GAME_ID"]).get_normalized_dict()
#         time.sleep(sleeptime)
#         misc = boxscoremiscv2.BoxScoreMiscV2(g["GAME_ID"]).get_normalized_dict()

#         for game in boxscorefinder["TeamStats"]:
#             #TODO: check this
#             if team["id"] != game["TEAM_ID"]:
#                 print("ERROR - team ID's do not match. have: ", team["id"], " but got: ", game["TEAM_ID"])
#                 sys.exit(0)
#             #if the team's game stats haven't been added already
#             if db.teamgames.count_documents({"team": game["TEAM_ID"], "game":game["GAME_ID"]}) == 0:
#                 # opponent = lookupTeam(g["MATCHUP"].split(" ")[2])['id']
#                 home_or_away = ""
#                 opponent = ""
#                 # if g["MATCHUP"].split(" ")[1] == "@":
#                 #     #team1 is away
#                 #     home_or_away = "away"
#                 # else:
#                 #     #team1 is home
#                 #     home_or_away = "home"
#                 game_data = {
#                     # "game":game["GAME_ID"],
#                     # "home_or_away":home_or_away,
#                     # "win_or_loss":game["WIN"],
#                     # "opponent":opponent,
#                     # "opponent_record":"W-L",
#                     # "opponent_record_ratio":.500,
#                     # "date": g["GAME_DATE"],
#                     # "season": transform_season_id,


#                     # field_goals_attempted: Int!
#                     # field_goals_made: Int!
#                     # field_goal_percentage: Float!
#                     # pace: Float!
#                     # defensive_rating: Float!
#                     # defensive_rebound_percentage: Float!
#                     # offensive_rebound_percentage: Float!
#                     # personal_fouls: Int!
#                     # personal_fouls_drawn: Int!
#                     # opponent_points: Int!
#                     # opponent_effective_field_goal_percentage: Float!
#                     # opponent_assists: Int!
#                     # opponent_rebounds: Int!
#                     # opponent_three_pointers_made: Int!
#                     # opponent_three_pointers_attempted: Int!
#                     # opponent_field_goals_attempted: Int!
#                     # opponent_three_point_frequency: Float!
#                     # opponent_free_throws_attempted: Int!
#                 }
#                 print(game)
#                 for advanced_game in advanced["TeamStats"]:
#                     print(advanced_game)
#                     # if game_data['player'] == advanced_game["PLAYER_ID"]:
#                     #     game_data["assist_percentage"] = advanced_game["AST_PCT"]
#                     #     game_data["defensive_rebound_percentage"] = advanced_game["DREB_PCT"]
#                     #     game_data["offensive_rebound_percentage"] = advanced_game["OREB_PCT"]
#                     #     game_data["usage"] = advanced_game["USG_PCT"]
#                     #     game_data["effective_field_goal_percentage"] = advanced_game["EFG_PCT"]
#                     #     game_data["true_shooting_percentage"] = advanced_game["TS_PCT"]
#                 print(misc)
#                 for misc_game in misc["sqlTeamMisc"]:
#                     print(misc_game)
#                     # if game_data['player'] == misc_game["PLAYER_ID"]:
#                     #     game_data["personal_fouls_drawn"] = misc_game["PFD"]
#                 # db.teamgames.insert_one(game_data)
#                 # print(team["name"], game_data["season"], game_data["date"], game_data["minutes"])
#                 print("----------------------------------------------------------------")

# %% [markdown]
# 


