from gql import gql, Client
from gql.transport.aiohttp import AIOHTTPTransport

# Select your transport with a defined url endpoint
transport = AIOHTTPTransport(url="http://localhost:8080/query")

# Create a GraphQL client using the defined transport
client = Client(transport=transport, fetch_schema_from_transport=True)

# Provide a GraphQL query
query = gql(
    """
    query getProjections {
      player(input:{name:"Jayson Tatum"}) {
        name
        games(input: {season:"2021-22"}) {
        projections {
            sportsbook
            targets {
                target
                type
            }
            date
        }
        date
        }
      }
    }
"""
)

# Execute the query on the transport
result = client.execute(query)

print(result.values())
