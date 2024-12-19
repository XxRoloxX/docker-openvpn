import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.scss";
import Router from "./pages/Router.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Router />
  </StrictMode>,
);
