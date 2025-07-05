import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { registerUser, clearError } from '../store/authSlice';
import type { RegisterCredentials } from '../types/Auth';

export const RegisterPage: React.FC = () => {
  const [credentials, setCredentials] = useState<RegisterCredentials>({
    name: '',
    cpf: '',
    password: '',
  });
  const dispatch = useAppDispatch();
  const { loading, error, isAuthenticated } = useAppSelector((state) => state.auth);
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/topics');
    }
  }, [isAuthenticated, navigate]);

  useEffect(() => {
    return () => {
      dispatch(clearError());
    };
  }, [dispatch]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    dispatch(clearError());
    dispatch(registerUser(credentials));
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setCredentials(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  return (
    <div className="centered-page">
      <div className="card" style={{ width: '100%', maxWidth: '400px' }}>
        <div className="card-body">
          <h1 className="text-center mb-6">Registrar</h1>
          
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label htmlFor="name" className="form-label">
                Nome:
              </label>
              <input
                type="text"
                id="name"
                name="name"
                value={credentials.name}
                onChange={handleChange}
                required
                className="form-input"
              />
            </div>

            <div className="form-group">
              <label htmlFor="cpf" className="form-label">
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
                className="form-input"
              />
              <small className="text-muted text-sm">
                Digite apenas números (11 dígitos)
              </small>
            </div>

            <div className="form-group">
              <label htmlFor="password" className="form-label">
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
                className="form-input"
              />
              <small className="text-muted text-sm">
                Mínimo 6 caracteres
              </small>
            </div>

            {error && (
              <div className="alert alert-danger">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="btn btn-success btn-lg"
              style={{ width: '100%' }}
            >
              {loading ? 'Registrando...' : 'Registrar'}
            </button>
          </form>

          <div className="text-center mt-6">
            <p className="text-muted">
              Já tem uma conta?{' '}
              <Link to="/login" className="text-primary">
                Faça login
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}; 