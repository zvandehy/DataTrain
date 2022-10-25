import { ApolloProvider } from "@apollo/client";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import DashboardPage from "../../pages/dashboard/dashboard.page";
import apolloClient from "../../shared/apollo-client";
import { COLORS } from "../../shared/styles/constants";

import "../shared/styles";
import "./app.component.css";
declare module "@mui/material/styles" {
  interface Palette {
    negative: Palette["primary"];
    positive: Palette["primary"];
    push: Palette["primary"];
  }

  // allow configuration using `createTheme`
  interface PaletteOptions {
    negative?: PaletteOptions["primary"];
    positive?: PaletteOptions["primary"];
    push?: PaletteOptions["primary"];
  }
}

// Update the Button's color prop options
declare module "@mui/material/Button" {
  interface ButtonPropsColorOverrides {
    positive: true;
    negative: true;
    push: true;
  }
}

declare module "@mui/material/Typography" {
  interface ButtonPropsColorOverrides {
    positive: true;
    negative: true;
    push: true;
  }
}

export const theme = createTheme({
  palette: {
    mode: "dark",
    primary: {
      main: COLORS.PRIMARY,
    },
    secondary: {
      main: COLORS.SECONDARY,
    },
    background: {
      default: "#212121",
      paper: "#424242",
    },
    warning: {
      main: "#fdd835",
    },
    positive: {
      main: COLORS.HIGHER,
      contrastText: "#000",
    },
    negative: {
      main: COLORS.LOWER,
      contrastText: "#FFF",
    },
    push: {
      main: COLORS.PUSH,
      contrastText: "#000",
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
      <ApolloProvider client={apolloClient}>
        <Router>
          <Routes>
            {/* <Route path="/nba/players/:id" element={<PlayerPage />}></Route>
            <Route path="/wnba/players/:id" element={<PlayerPage />}></Route>
            <Route path="/" element={<Home />} />
            <Route path="/wnba" element={<Home />} />
            <Route path="/nba" element={<Home />} /> */}
            <Route path="/dashboard" element={<DashboardPage />} />
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
