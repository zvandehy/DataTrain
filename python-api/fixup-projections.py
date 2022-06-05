import datetime
from pymongo import MongoClient
url = ""
with open("app.env") as file:
    db_source = file.read()
    url = db_source.split("DB_SOURCE=\"")[1][:-1]

client = MongoClient(url)
result = client['nba']['projections'].find({})
count = 0
# get time right now with format 2006-01-02T15:04:05.999Z07:00
timeNow = datetime.datetime.now().strftime("%Y-%m-%dT%H:%M:%S.000Z")
print(timeNow)
for doc in result:
    count += 1
#     opponent = ""
#     if "opponent" in doc:
#         opponent = doc["opponent"]
#     else:
#         opponent = doc["opponentabr"]
#     start = ""
#     if "startTime" in doc:
#         start = doc["startTime"]
#     else:
#         start = doc["starttime"]
#     date = ""
#     if "date" not in doc:
#         date = doc["starttime"]
#         date = date.split("T")[0]
#     else:
#         date = doc["date"]
    propositions = []
    for item in doc["propositions"]:
        proposition = {
            "target": item["target"],
            "type": item["type"],
            "sportsbook": item["sportsbook"],
            "lastModified": timeNow}
        propositions.append(proposition)
    new_doc = {
        "_id": doc["_id"],
        "playername": doc["playername"],
        "date": doc["date"],
        "startTime": doc["startTime"],
        "opponent": doc["opponent"],
        "propositions": propositions
    }
    client['nba']['projections'].find_one_and_replace(
        {"_id": doc["_id"]}, new_doc)
    print(count)
print("DONE")
