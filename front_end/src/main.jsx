import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.jsx";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import Login from "./components/Login.jsx";
import UserDetail from "./components/UserDetail.jsx";
import UsersList from "./components/UsersList.jsx";
import UsersLog from "./components/UsersLog.jsx";
import Welcome from "./components/Welcome.jsx";
import CreateUserForm from "./components/CreateUserForm.jsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,

    children: [
      { path: "/", element: <Welcome /> },

      {
        path: "/users",
        element: <UsersList />,
      },
      { path: "user/:id", element: <UserDetail /> },
      { path: "/userlogs", element: <UsersLog /> },
      // { path: "/login", element: <Login /> },
      { path: "/create", element: <CreateUserForm /> },

    ],
  },
]);
ReactDOM.createRoot(document.getElementById("root")).render(
  <RouterProvider router={router} />
);
