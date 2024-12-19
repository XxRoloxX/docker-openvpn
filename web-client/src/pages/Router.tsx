// import * as React from "react";
// import { createRoot } from "react-dom/client";
import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  RouterProvider,
} from "react-router-dom";
import Clients from "./clients/Clients";
import SidebarLayout from "@/layouts/SidebarLayout";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<SidebarLayout />}>
      <Route path="" element={<Clients />} />
      <Route path="clients" element={<Clients />} />
    </Route>,
  ),
);

const Router = () => {
  return <RouterProvider router={router} />;
};

export default Router;
