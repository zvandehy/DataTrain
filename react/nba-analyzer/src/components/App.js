import React from 'react';
import {Routes, Route } from 'react-router-dom';
import Players from './Players';
import Player from './Player';
import "../fonts/Oswald-Regular.ttf"

const App = () => {
  return (
    <div className="App">
      <Routes>
        <Route exact path="/" element={<Players/>}/>
        <Route path="/players/:id" element={<Player/>}>
          
        </Route>
      </Routes>
    </div>
  )
}

export default App;
