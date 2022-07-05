from pymongo import MongoClient

import time
from datetime import datetime
import requests
from nba_api.stats.endpoints import playercareerstats
from nba_api.stats.library.parameters import PerModeSimple


lastModifiedDate = datetime.now()
url = "mongodb+srv://datatrain:6JUgI5GJlmIY6ro0@datatrain.i5rgk.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
client = MongoClient(url)
db = client.wnba

correctNames = {
    "Asia (AD) Durr": {"fn": "Asia", "ln": "Durr"},
    "Azurá Stevens": {"fn": "Azura", "ln": "Stevens"},
}

# add players to db
get_season = "2022"
players_url = "https://data.wnba.com/data/5s/v2015/json/mobile_teams/wnba/" + \
    get_season+"/players/10_player_info.json"
all_players = requests.get(players_url).json()["pls"]["pl"]
attempt = 0
sleeptime = 5
total_players = len(all_players)
count = 0
for attempt in range(5):
    print("attempt: " + str(attempt))
    try:
        for count in range(total_players):
            player = all_players[count]
            fn = player["fn"]
            ln = player["ln"]
            name = player["fn"] + " " + player["ln"]
            if name in correctNames:
                save_name = correctNames[name]
                name = save_name["fn"] + " " + save_name["ln"]
                fn = save_name["fn"]
                ln = save_name["ln"]
            # else:
            #     continue
            player_data = {
                "playerID": player["pid"],
                "name": name,
                "first_name": fn,
                "last_name": ln,
                "teamABR": player["ta"],
                "height": player["ht"],
                "weight": player["wt"],
                "position": player["pos"],
                "lastModifiedDate": lastModifiedDate,
                "seasons": [],
            }
            if player_data["teamABR"] == "TOT":
                print("skipping TOT player " + player_data["name"])
                count += 1
                continue
            player_career = playercareerstats.PlayerCareerStats(
                player["pid"], per_mode36=PerModeSimple.per_game, league_id_nullable=10)
            regular_seasons = player_career.get_normalized_dict()[
                "SeasonTotalsRegularSeason"]
            time.sleep(sleeptime)
            nonTOTseasons = [
                s for s in regular_seasons if s["TEAM_ABBREVIATION"] != "TOT"]
            for season in nonTOTseasons:
                player_data["seasons"].append(season["SEASON_ID"])
                if season["SEASON_ID"] == "2022-23":
                    player_data["teamABR"] = season["TEAM_ABBREVIATION"]
            # if player_data["weight"] == "" or player_data["height"] == "" or player_data["position"] == "":
            #     info = commonplayerinfo.CommonPlayerInfo(player["pid"], league_id_nullable=10).get_normalized_dict()["CommonPlayerInfo"][0]
            #     player_data["position"]=info["POSITION"]
            #     player_data["height"]=info["HEIGHT"]
            #     player_data["weight"]=int(info["WEIGHT"])
            db.players.update_one({"playerID": player_data["playerID"]}, {
                                  "$set": player_data}, True)
            count += 1
            print(count, player_data["first_name"],
                  player_data["last_name"], player_data["teamABR"])
        break
    except Exception as e:
        print("failed on attempt: " + str(attempt))
        print(e)
        print(player_data)
        time.sleep(sleeptime*10)
        sleeptime *= 2
    attempt += 1
