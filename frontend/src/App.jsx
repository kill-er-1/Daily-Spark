import React, { useState } from "react";
import LoginPage from "./pages/LoginPage.jsx";
import RegisterPage from "./pages/RegisterPage.jsx";

export default function App() {
  const [view, setView] = useState("login"); // login | register

  return (
    <>
      {view === "login" ? (
        <LoginPage onSwitch={(to) => setView(to)} />
      ) : (
        <RegisterPage onSwitch={(to) => setView(to)} />
      )}
    </>
  );
}
