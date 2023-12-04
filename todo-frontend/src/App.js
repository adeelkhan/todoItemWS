import logo from "./logo.svg";
import "./App.css";
import ListItem from "./components/ListItem";
import Login from "./components/Login";
import LogOut from "./components/Logout";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthContext } from "./AuthContext";
import { useState } from "react";

function App() {
  const [username, setAuthUser] = useState("");
  return (
    <div>
      <AuthContext.Provider value={{ username, setAuthUser }}>
        <BrowserRouter>
          <Routes>
            <Route path="/list" element={<ListItem />} />
            <Route path="/login" element={<Login />} />
            <Route path="/logout" element={<LogOut />} />
          </Routes>
        </BrowserRouter>
      </AuthContext.Provider>
    </div>
  );
}

export default App;
