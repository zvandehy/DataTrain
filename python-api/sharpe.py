# THIS FILE SHOULD EVENTUALLY CALCULATE THE SHARPE VALUE OF OUR PREDICTIONS
import datetime
from gql import gql, Client
from gql.transport.aiohttp import AIOHTTPTransport

# Select your transport with a defined url endpoint
transport = AIOHTTPTransport(url="http://localhost:8080/query")

# Create a GraphQL client using the defined transport
client = Client(transport=transport, fetch_schema_from_transport=True)

# Provide a GraphQL query
query = gql(
    """
    query getProjections($date: String!) {
      projections(input:{startDate:$date endDate:$date}) {
        player {
          name
          games(input: {season:"2021-22"}) {
            points
            rebounds
            assists
            blocks
            steals
            turnovers
            three_pointers_made
            free_throws_made
            date
          }
        }
        propositions {
          sportsbook
          target
          type
        }
        date
        result {
          points
          rebounds
          assists
          blocks
          steals
          turnovers
          three_pointers_made
          free_throws_made
          date
        }
      }
    }
"""
)


def GetScores(games, stat):
    scores = []
    for game in games:
        score = GetScore(game, stat)
        if score == -99:
            return []
        scores.append(score)
    return scores


def GetScore(game, stat):
    if stat == 'Pts+Asts':
        return game['points'] + game['assists']
    elif stat == 'Pts+Rebs':
        return game['points'] + game['rebounds']
    elif stat == 'Pts+Rebs+Asts':
        return game['points'] + game['rebounds'] + game['assists']
    elif stat == 'Rebs+Asts':
        return game['assists'] + game['rebounds']
    elif stat == "Blks+Stls":
        return game['blocks'] + game['steals']
    elif stat == "3-PT Made":
        return game['three_pointers_made']
    elif stat == "Blocked Shots":
        return game['blocks']
    elif stat == "Turnovers":
        return game["turnovers"]
    elif stat == "Points":
        return game["points"]
    elif stat == "Rebounds":
        return game["rebounds"]
    elif stat == "Assists":
        return game["assists"]
    elif stat == "Steals":
        return game["steals"]
    elif stat == "Blocks":
        return game["blocks"]
    elif stat == "Free Throws Made":
        return game["free_throws_made"]
    elif stat == "Fantasy Score":
        return game["points"] + game["rebounds"]*1.2 + game["assists"] * 1.5 + game["blocks"]*3 + game["steals"]*3 + game["turnovers"]*-1
    else:
        return -99


# TODO: Get projection from backend instead of generating after receiving data
def WeightedAverage(scores, target):
    pct_o = 0
    pct_u = 0
    (over_5, under_5) = CalculatePercentOverUnder(scores[:5], target)
    (over_10, under_10) = CalculatePercentOverUnder(scores[:10], target)
    (over_30, under_30) = CalculatePercentOverUnder(scores[:30], target)
    (over, under) = CalculatePercentOverUnder(scores, target)
    weights = [0.3, 0.27, 0.25, 0.18]
    pct_o = (over_5*weights[0] + over_10*weights[1] +
             over_30*weights[2] + over*weights[3])
    pct_u = (under_5*weights[0] + under_10*weights[1] +
             under_30*weights[2] + under*weights[3])
    return (round(pct_o, 2), round(pct_u, 2))


def CalculatePercentOverUnder(scores, target):
    overs = 0
    unders = 0
    for score in scores:
        if score > target:
            overs += 1
        elif score < target:
            unders += 1
        else:
            overs += 0.5
            unders += 0.5
    return (overs/len(scores), unders/len(scores))


date = datetime.datetime.strptime("2022-03-09", "%Y-%m-%d")
overall_correct = 0
overall_incorrect = 0
overall_max_correct = 0
overall_max_incorrect = 0
while date < datetime.datetime.strptime("2022-06-06", "%Y-%m-%d"):
    print(datetime.datetime.strftime(date, "%Y-%m-%d"))

    # Execute the query on the transport
    result = client.execute(query, variable_values={
                            "date": datetime.datetime.strftime(date, "%Y-%m-%d")})

    total_correct = 0
    total_incorrect = 0
    total_max_correct = 0
    total_max_incorrect = 0
    for projection in result["projections"]:
        if "Starters" in projection["player"]["name"] or projection["result"] == None:
            continue
        correct = 0
        incorrect = 0
        maxType = ""
        max = 0
        maxResult = ""
        for proposition in projection["propositions"]:
            scores = GetScores(
                projection["player"]["games"], proposition["type"])
            if scores == []:
                print("Can't get scores for ",
                      proposition["type"], " continuing...")
                continue
            # 0 is recent, len-1 is oldest
            average = round(sum(scores) / len(scores), 2)
            (pct_o, pct_u) = WeightedAverage(scores, proposition["target"])
            if pct_o > max:
                max = pct_o
                maxType = proposition["type"]
                if GetScore(projection["result"], maxType) > proposition["target"]:
                    maxResult = "Correct"
                else:
                    maxResult = "Incorrect"
            if pct_u > max:
                max = pct_u
                maxType = proposition["type"]
                if GetScore(projection["result"], maxType) < proposition["target"]:
                    maxResult = "Correct"
                else:
                    maxResult = "Incorrect"

            # print(
            #     proposition["target"],
            #     proposition["type"],
            #     average,
            #     GetScore(projection["result"], proposition["type"])
            # )

            result_score = GetScore(projection["result"], proposition["type"])
            if result_score == -99:
                print("Can't get score for RESULT continuing...")
                continue
            if pct_o > 0.6:
                if result_score > proposition["target"]:
                    correct += 1
                else:
                    incorrect += 1
            if pct_u > 0.6:
                if result_score < proposition["target"]:
                    correct += 1
                else:
                    incorrect += 1
        if correct+incorrect > 0:
            total_correct += correct
            total_incorrect += incorrect
            print(projection["player"]["name"], "Correct: ", correct, " Incorrect: ", incorrect,
                  " Percent Correct: ", correct/(correct+incorrect), " Max Type: ", maxType, " (", max, ") Max Result: ", maxResult)
        if maxResult == "Correct":
            total_max_correct += 1
        else:
            total_max_incorrect += 1
    if total_correct+total_incorrect > 0:
        print("==========================================================")
        print("Total Correct: ", total_correct, " Total Incorrect: ", total_incorrect,
              " Total Percent Correct: ", total_correct/(total_correct+total_incorrect))
        print("Total Max Correct: ", total_max_correct,
              " Total Max Incorrect: ", total_max_incorrect, " Total Max Percent Correct: ", total_max_correct/(total_max_correct+total_max_incorrect))
    date += datetime.timedelta(days=1)
    overall_max_correct += total_max_correct
    overall_max_incorrect += total_max_incorrect
    overall_correct += total_correct
    overall_incorrect += total_incorrect
print("==========================================================")
print("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
print("==========================================================")
print("OVERALL Correct: ", overall_correct, " OVERALL Incorrect: ", overall_incorrect,
      " Total OVERALL Percent Correct: ", overall_correct/(overall_correct+overall_incorrect))
print("OVERALL Max Correct: ", overall_max_correct,
      " OVERALL Max Incorrect: ", overall_max_incorrect, " OVERALL Max Percent Correct: ", overall_max_correct/(overall_max_correct+overall_max_incorrect))

# OVERALL Correct:  4336  OVERALL Incorrect:  2627  Total OVERALL Percent Correct:  0.6227200919144047
# OVERALL Max Correct:  1963  OVERALL Max Incorrect:  1032  OVERALL Max Percent Correct:  0.6554257095158598
