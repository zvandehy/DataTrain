import { ApolloClient, InMemoryCache } from "@apollo/client";

const apolloClient = new ApolloClient({
  uri: "https://clover-backend.fly.dev/nba/query",
  // uri: "http://localhost:8080/nba/query",
  // uri: "https://www.clover-analytics.com/api/wnba",
  cache: new InMemoryCache(),
  headers: {
    "Access-Control-Allow-Origin": "*",
  },
});

export default apolloClient;
