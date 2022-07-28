import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import { ApolloProvider } from "@apollo/client";

import Header from "../header/header.component";
import About from "../../pages/about/about.page";
import Home from "../../pages/home/home.page";

import "../../shared/styles";
import "./app.component.css";
import "../../shared/fonts/Oswald-Regular.ttf";
import client from "../../shared/apollo-client";
import { createTheme, ThemeProvider } from "@mui/material/styles";

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
  return (
    <ThemeProvider theme={theme}>
      <ApolloProvider client={client}>
        <Router>
          <Routes>
            <Route path="/about" element={<About />} />
            <Route path="/" element={<Home />} />
          </Routes>
        </Router>
      </ApolloProvider>
    </ThemeProvider>
  );
};

export default App;
