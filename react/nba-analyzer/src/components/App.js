import React from "react";
import { Routes, Route } from "react-router-dom";
import Players from "./Players";
import Player from "./Player";
import "../fonts/Oswald-Regular.ttf";
import { NBAClient, WNBAClient } from "../index";

const App = () => {
  return (
    <div className="App">
      <Routes>
        <Route
          exact
          path="/"
          element={<Players client={NBAClient} league="nba" />}
        />
        {/* <Route
          path="/players/:id"
          element={<Player client={NBAClient} />}
        ></Route> */}
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
          exact
          path="/wnba"
          element={<Players client={WNBAClient} league="wnba" />}
        />
        <Route
          path="/wnba/players/:id"
          element={<Player client={WNBAClient} league="wnba" />}
        ></Route>
      </Routes>
    </div>
  );
};

export default App;
