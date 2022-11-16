import React from "react";
import ReactDOM from "react-dom";
import "./styles/index.css";
import App from "./App";
import { BrowserRouter } from "react-router-dom";

// 1
import {
  ApolloProvider,
  ApolloClient,
  createHttpLink,
  InMemoryCache,
} from "@apollo/client";

const wnbaGQL = createHttpLink({
  uri: "https://datatrain-nba-yxh2z.ondigitalocean.app/wnba/query",
  // uri: "http://localhost:8080/wnba/query",
});

// 3
export const WNBAClient = new ApolloClient({
  link: wnbaGQL,
  cache: new InMemoryCache(),
});

// 2
const nbaGQL = createHttpLink({
  uri: "https://datatrain-nba-yxh2z.ondigitalocean.app/nba/query",
  // uri: "http://localhost:8080/nba/query",
});

// 3
export const NBAClient = new ApolloClient({
  link: nbaGQL,
  cache: new InMemoryCache(),
});

// 4
ReactDOM.render(
  <BrowserRouter>
    <ApolloProvider client={NBAClient}>
      <App />
    </ApolloProvider>
  </BrowserRouter>,
  document.getElementById("root")
);
