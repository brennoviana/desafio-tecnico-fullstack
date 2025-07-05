import { useEffect, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { fetchTopics, getSessionStatus } from '../services/api';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { logoutUser } from '../store/authSlice';
import type { Topic } from '../types/Topic';

interface TopicWithStatus extends Topic {
  status: 'Awaiting Opening' | 'Open Session' | 'Voting Closed';
  sessionData?: {
    open_at: number;
    close_at: number;
  } | null;
}

export const Dashboard: React.FC = () => {
  const [topics, setTopics] = useState<TopicWithStatus[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const dispatch = useAppDispatch();
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);
  const navigate = useNavigate();

  useEffect(() => {
    const loadTopicsWithStatus = async () => {
      try {
        const topicsData = await fetchTopics();
        const topicsWithStatus = await Promise.all(
          topicsData.map(async (topic) => {
            try {
              const sessionData = await getSessionStatus(topic.id);
              const now = Date.now() / 1000;
              
              let status: 'Awaiting Opening' | 'Open Session' | 'Voting Closed';
              if (!sessionData) {
                status = 'Awaiting Opening';
              } else if (now >= sessionData.open_at && now <= sessionData.close_at) {
                status = 'Open Session';
              } else if (now > sessionData.close_at) {
                status = 'Voting Closed';
              } else {
                status = 'Awaiting Opening';
              }
              
              return {
                ...topic,
                status,
                sessionData
              };
            } catch {
              return {
                ...topic,
                status: 'Awaiting Opening' as const
              };
            }
          })
        );
        setTopics(topicsWithStatus);
      } catch (err) {
        console.error(err);
        setError('Erro ao carregar tópicos.');
      } finally {
        setLoading(false);
      }
    };

    loadTopicsWithStatus();
  }, []);

  const handleLogout = () => {
    dispatch(logoutUser());
    navigate('/dashboard');
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
            <div>
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
                  
                  <div className="mb-4">
                    <span className={`status-badge ${
                      topic.status === 'Awaiting Opening' ? 'status-awaiting' :
                      topic.status === 'Open Session' ? 'status-open' : 'status-closed'
                    }`}>
                      {topic.status}
                    </span>
                  </div>

                  {topic.sessionData && (
                    <div className="mb-4 text-sm text-muted">
                      <p className="mb-2">
                        <strong>Abertura:</strong> {new Date(topic.sessionData.open_at * 1000).toLocaleString()}
                      </p>
                      <p>
                        <strong>Fechamento:</strong> {new Date(topic.sessionData.close_at * 1000).toLocaleString()}
                      </p>
                    </div>
                  )}

                  <div className="flex flex-col gap-2">
                    {topic.status === 'Open Session' && isAuthenticated && (
                      <Link to={`/topic/${topic.id}/vote`} className="btn btn-primary text-center">
                        Votar
                      </Link>
                    )}
                    
                    {topic.status === 'Voting Closed' && (
                      <Link to={`/topic/${topic.id}/results`} className="btn btn-info text-center">
                        Ver Resultados
                      </Link>
                    )}

                    {topic.status === 'Awaiting Opening' && isAuthenticated && (
                      <button
                        onClick={() => navigate(`/topic/${topic.id}/manage`)}
                        className="btn btn-warning"
                      >
                        Abrir Sessão
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