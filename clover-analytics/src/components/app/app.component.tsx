import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { ApolloProvider } from "@apollo/client";
import { createTheme, ThemeProvider } from "@mui/material/styles";

import Home from "../../pages/home/home.page";
import PlayerPage from "../../pages/player/player-wrapper.page";
import "../../shared/styles";
import "./app.component.css";
import "../../shared/fonts/Oswald-Regular.ttf";
import client from "../../shared/apollo-client";

const theme = createTheme({
  palette: {
    primary: {
      main: "#395c6b",
    },
    secondary: {
      main: "#f59a4b",
    },
    text: {
      primary: "#f8ffff",
    },
    info: {
      main: "#f8ffff",
    },
  },
  components: {
    MuiSelect: {
      defaultProps: {
        variant: "standard",
        style: {
          // fontSize: "inherit",
          color: "inherit",
        },
      },
      styleOverrides: {
        icon: {
          fill: "#f59a4b",
          fontSize: "1.75rem",
        },
      },
    },
    MuiAutocomplete: {
      defaultProps: {
        // variant: "",

        style: {
          // fontSize: "inherit",
          color: "white",
        },
      },
      styleOverrides: {
        // icon: {
        //   fill: "#f59a4b",
        //   fontSize: "1.75rem",
        // },
      },
    },
  },
});

const App: React.FC = () => {
  Storage.prototype.setObject = function (key: string, value: Object) {
    this.setItem(key, JSON.stringify(value));
  };
  Storage.prototype.getObject = function (key: string) {
    let value = this.getItem(key);
    return value && JSON.parse(value);
  };
  return (
    <ThemeProvider theme={theme}>
      <ApolloProvider client={client}>
        <Router>
          <Routes>
            <Route path="/nba/players/:id" element={<PlayerPage />}></Route>
            <Route path="/wnba/players/:id" element={<PlayerPage />}></Route>
            <Route path="/" element={<Home />} />
            <Route path="/wnba" element={<Home />} />
            <Route path="/nba" element={<Home />} />
            {/* <Route
          exact
          path="/"
          element={<Players client={NBAClient} league="nba" />}
        />
        <Route
          exact
          path="/nba"
          element={<Players client={NBAClient} league="nba" />}
        />
        <Route
          path="/nba/players/:id"
          element={<Player client={NBAClient} />}
        ></Route>
        <Route
          path="/wnba/players/:id"
          element={<Player client={WNBAClient} league="wnba" />}
        ></Route>
        <Route
          exact
          path="/wnba"
          element={<Players client={WNBAClient} league="wnba" />}
        />
        <Route
          path="/wnba/players/:id"
          element={<Player client={WNBAClient} league="wnba" />}
        ></Route> */}
          </Routes>
        </Router>
      </ApolloProvider>
    </ThemeProvider>
  );
};

export default App;
