import { useEffect, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { fetchTopics } from '../services/api';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { logoutUser } from '../store/authSlice';
import { AddTopicButton } from '../components/AddTopicButton';
import type { Topic } from '../types/Topic';

export const Dashboard: React.FC = () => {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const dispatch = useAppDispatch();
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);
  const navigate = useNavigate();

  const loadTopics = async () => {
    try {
      const topicsData = await fetchTopics();
      setTopics(topicsData);
    } catch (err) {
      console.error(err);
      setError('Erro ao carregar tópicos.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTopics();
  }, []);

  const handleLogout = () => {
    dispatch(logoutUser());
    navigate('/dashboard');
  };

  const handleTopicAdded = () => {
    loadTopics();
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="container">
        {/* Header */}
        <div className="header">
          <div className="header-content">
            <div>
              <h1 className="header-title">Dashboard - Pautas de Votação</h1>
              {!isAuthenticated && (
                <p className="header-subtitle">
                  Faça login para votar nas pautas
                </p>
              )}
            </div>
            <div className="flex gap-2">
              {isAuthenticated && (
                <AddTopicButton onTopicAdded={handleTopicAdded} />
              )}
              {isAuthenticated ? (
                <button onClick={handleLogout} className="btn btn-danger">
                  Logout
                </button>
              ) : (
                <div className="flex gap-2">
                  <Link to="/login" className="btn btn-primary">
                    Login
                  </Link>
                  <Link to="/register" className="btn btn-success">
                    Registrar
                  </Link>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Topics Dashboard */}
        {error && (
          <div className="alert alert-danger">
            {error}
          </div>
        )}

        {topics.length === 0 ? (
          <p className="text-center text-muted">Nenhuma pauta encontrada.</p>
        ) : (
          <div className="grid grid-auto-fit">
            {topics.map((topic) => (
              <div key={topic.id} className="card">
                <div className="card-body">
                  <h3 className="mb-4">{topic.name}</h3>

                  <div className="flex flex-col gap-2">
                    {isAuthenticated && (
                      <Link to={`/topic/${topic.id}/vote`} className="btn btn-primary text-center">
                        Votar
                      </Link>
                    )}
                    
                    {/* Results are always available for everyone */}
                    <Link to={`/topic/${topic.id}/results`} className="btn btn-info text-center">
                      Ver Resultados
                    </Link>

                    {isAuthenticated && (
                      <button
                        onClick={() => navigate(`/topic/${topic.id}/manage`)}
                        className="btn btn-warning"
                      >
                        Gerenciar Sessão
                      </button>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}; 