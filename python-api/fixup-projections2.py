import datetime
from os import dup
from pymongo import MongoClient
url = ""
with open("app.env") as file:
    db_source = file.read()
    url = db_source.split("DB_SOURCE=\"")[1][:-1]

client = MongoClient(url)
result = client['nba']['projections'].find()
count = 0
timeNow = datetime.datetime.now().strftime("%Y-%m-%dT%H:%M:%S")
for doc in result:
    count += 1
    duplicates = client['nba']['projections'].find(
        {"date": doc["date"], "playername": doc["playername"]})
    maxDoc = doc
    maxProps = 0
    for duplicate in duplicates:
        if len(duplicate["propositions"]) > maxProps:
            maxDoc = duplicate
            maxProps = len(duplicate["propositions"])
    client['nba']['projections'].delete_many(
        {"date": doc["date"], "playername": doc["playername"]})
    client['nba']['projections'].insert_one(maxDoc)
    print(count)
result = client['nba']['projections'].find()
count = 0
for doc in result:
    count += 1
print(count, " remaining")
print("DONE")
