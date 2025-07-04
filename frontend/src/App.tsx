import React from "react";
import { useAuth } from "./context/AuthContext";
import Login from "./pages/Login";
import Topics from "./pages/Topics";
import { BrowserRouter as Router, Routes, Route, Navigate, useNavigate } from "react-router-dom";
import './App.css'

const TopBar: React.FC = () => {
  const { isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  return (
    <div style={{ marginBottom: 16 }}>
      {isAuthenticated ? (
        <button onClick={logout}>Sair</button>
      ) : (
        <button onClick={() => navigate("/login")}>Entrar</button>
      )}
    </div>
  );
};

const App: React.FC = () => {
  return (
    <Router>
      <TopBar />
      <Routes>
        <Route path="/topics" element={<Topics />} />
        <Route path="/login" element={<Login />} />
        <Route path="*" element={<Navigate to="/topics" />} />
      </Routes>
    </Router>
  );
};

export default App;
