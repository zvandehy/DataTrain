import { ApolloClient, InMemoryCache } from "@apollo/client";

const client = new ApolloClient({
  uri: "http://clover-backend.fly.dev/wnba/query",
  // uri: "http://localhost:8080/wnba/query",
  cache: new InMemoryCache(),
});

export default client;
