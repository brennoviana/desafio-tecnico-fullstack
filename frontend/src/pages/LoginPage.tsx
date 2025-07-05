import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Link, useNavigate } from 'react-router-dom';
import type { LoginCredentials } from '../types/Auth';

export const LoginPage: React.FC = () => {
  const [credentials, setCredentials] = useState<LoginCredentials>({
    cpf: '',
    password: '',
  });
  const [error, setError] = useState<string | null>(null);
  const { login, loading } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      await login(credentials);
      navigate('/topics');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao fazer login');
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCredentials(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  return (
    <div style={{ 
      maxWidth: '400px', 
      margin: '2rem auto', 
      padding: '2rem',
      border: '1px solid #ddd',
      borderRadius: '8px',
      boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
    }}>
      <h1 style={{ textAlign: 'center', marginBottom: '2rem' }}>Login</h1>
      
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: '1rem' }}>
          <label htmlFor="cpf" style={{ display: 'block', marginBottom: '0.5rem' }}>
            CPF:
          </label>
          <input
            type="text"
            id="cpf"
            name="cpf"
            value={credentials.cpf}
            onChange={handleChange}
            required
            pattern="[0-9]{11}"
            placeholder="12345678900"
            style={{
              width: '100%',
              padding: '0.5rem',
              border: '1px solid #ddd',
              borderRadius: '4px',
              fontSize: '1rem'
            }}
          />
          <small style={{ color: '#666', fontSize: '0.875rem' }}>
            Digite apenas números (11 dígitos)
          </small>
        </div>

        <div style={{ marginBottom: '1rem' }}>
          <label htmlFor="password" style={{ display: 'block', marginBottom: '0.5rem' }}>
            Senha:
          </label>
          <input
            type="password"
            id="password"
            name="password"
            value={credentials.password}
            onChange={handleChange}
            required
            style={{
              width: '100%',
              padding: '0.5rem',
              border: '1px solid #ddd',
              borderRadius: '4px',
              fontSize: '1rem'
            }}
          />
        </div>

        {error && (
          <div style={{ 
            color: 'red', 
            marginBottom: '1rem', 
            padding: '0.5rem',
            backgroundColor: '#ffebee',
            borderRadius: '4px',
            border: '1px solid #ffcdd2'
          }}>
            {error}
          </div>
        )}

        <button
          type="submit"
          disabled={loading}
          style={{
            width: '100%',
            padding: '0.75rem',
            backgroundColor: '#007bff',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            fontSize: '1rem',
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading ? 0.6 : 1
          }}
        >
          {loading ? 'Fazendo login...' : 'Entrar'}
        </button>
      </form>

      <div style={{ textAlign: 'center', marginTop: '1rem' }}>
        <p>
          Não tem uma conta?{' '}
          <Link to="/register" style={{ color: '#007bff', textDecoration: 'none' }}>
            Registre-se
          </Link>
        </p>
      </div>
    </div>
  );
}; 