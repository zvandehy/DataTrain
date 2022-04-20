from pymongo import MongoClient

from pymongo import MongoClient

# Requires the PyMongo package.
# https://api.mongodb.com/python/current
url = "mongodb+srv://datatrain:Zd3O9S1zdmyakh1E@datatrain.i5rgk.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
client = MongoClient(url)
running = True
while running:
    result = client['nba']['games'].aggregate([
        {
            '$match': {
                'season': '2021-22'
            }
        }, {
            '$group': {
                '_id': {
                    'player': '$player', 
                    'gameID': '$gameID'
                }, 
                'count': {
                    '$count': {}
                }, 
            }
        }, {
            '$match': {
                'count': {
                    '$gt': 1
                }
            }
        }
    ])
    count=0
    for doc in result:
        print(doc)
        count+=1
        client['nba']['games'].delete_one(doc["_id"])
    print(count)
    if(count == 0):
        running = False