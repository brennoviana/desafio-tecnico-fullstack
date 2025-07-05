import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Link, useNavigate } from 'react-router-dom';
import type { RegisterCredentials } from '../types/Auth';

export const RegisterPage: React.FC = () => {
  const [credentials, setCredentials] = useState<RegisterCredentials>({
    name: '',
    cpf: '',
    password: '',
  });
  const [error, setError] = useState<string | null>(null);
  const { register, loading } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      await register(credentials);
      navigate('/topics');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao registrar');
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
      <h1 style={{ textAlign: 'center', marginBottom: '2rem' }}>Registrar</h1>
      
      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: '1rem' }}>
          <label htmlFor="name" style={{ display: 'block', marginBottom: '0.5rem' }}>
            Nome:
          </label>
          <input
            type="text"
            id="name"
            name="name"
            value={credentials.name}
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
            minLength={6}
            style={{
              width: '100%',
              padding: '0.5rem',
              border: '1px solid #ddd',
              borderRadius: '4px',
              fontSize: '1rem'
            }}
          />
          <small style={{ color: '#666', fontSize: '0.875rem' }}>
            Mínimo 6 caracteres
          </small>
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
            backgroundColor: '#28a745',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            fontSize: '1rem',
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading ? 0.6 : 1
          }}
        >
          {loading ? 'Registrando...' : 'Registrar'}
        </button>
      </form>

      <div style={{ textAlign: 'center', marginTop: '1rem' }}>
        <p>
          Já tem uma conta?{' '}
          <Link to="/login" style={{ color: '#007bff', textDecoration: 'none' }}>
            Faça login
          </Link>
        </p>
      </div>
    </div>
  );
}; 