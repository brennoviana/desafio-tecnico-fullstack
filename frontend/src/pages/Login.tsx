import React, { useState } from "react";
import { useAuth } from "../context/AuthContext";
import api from "../api/axios";

const Login: React.FC = () => {
  const { login } = useAuth();
  const [cpf, setCpf] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    try {
      const response = await api.post("/login", { cpf, password });
      login(response.data.token);
    } catch (err: any) {
      setError("Usuário ou senha inválidos");
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        placeholder="CPF"
        value={cpf}
        onChange={e => setCpf(e.target.value)}
        required
      />
      <input
        type="password"
        placeholder="Senha"
        value={password}
        onChange={e => setPassword(e.target.value)}
        required
      />
      <button type="submit">Entrar</button>
      {error && <div style={{ color: "red" }}>{error}</div>}
    </form>
  );
};

export default Login; 