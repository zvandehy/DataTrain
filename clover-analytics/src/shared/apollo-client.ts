import { ApolloClient, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
  uri: "https://clover-backend.fly.dev/wnba/query",
  // uri: "http://localhost:8080/wnba/query",
  // uri: "https://www.clover-analytics.com/api/wnba",
  cache: new InMemoryCache(),
  headers: {
    "Access-Control-Allow-Origin": "*",
  },
});

export default client;
