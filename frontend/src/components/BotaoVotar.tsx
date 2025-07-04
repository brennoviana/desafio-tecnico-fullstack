import React from "react";
import { useAuth } from "../context/AuthContext";

const BotaoVotar: React.FC = () => {
  const { isAuthenticated } = useAuth();
  if (!isAuthenticated) return null;
  return <button>Votar</button>;
};

export default BotaoVotar; 