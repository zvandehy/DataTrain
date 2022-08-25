import { ApolloClient, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
  uri: "https://datatrain-nba-yxh2z.ondigitalocean.app/wnba/query",
  // uri: "http://localhost:8080/wnba/query",
  cache: new InMemoryCache(),
});

export default client;
