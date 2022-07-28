import React from "react";
import CssBaseline from "@mui/material/CssBaseline";

import App from "./ components/app/app.component";
import { createRoot } from "react-dom/client";
const container = document.getElementById("root");
const root = createRoot(container!);
root.render(
  <React.StrictMode>
    <CssBaseline />
    <App />
  </React.StrictMode>
);
