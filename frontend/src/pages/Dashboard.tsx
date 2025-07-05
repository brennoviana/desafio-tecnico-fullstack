import { useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { logoutUser } from '../store/authSlice';
import { fetchTopics } from '../store/topicsSlice';
import { AddTopicButton } from '../components/AddTopicButton';

export const Dashboard: React.FC = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  
  const { isAuthenticated } = useAppSelector((state) => state.auth);
  const { topics, loading, error } = useAppSelector((state) => state.topics);

  useEffect(() => {
    dispatch(fetchTopics());
  }, [dispatch]);

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
            <div className="flex gap-2">
              {isAuthenticated && (
                <AddTopicButton />
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
                  <div className="flex justify-between items-start mb-4">
                    <h3>{topic.name}</h3>
                    <span className={`badge ${topic.status === 'Aguardando Abertura' ? 'badge-warning' : 
                                           topic.status === 'Sessão Aberta' ? 'badge-success' : 
                                           topic.status === 'Votação Encerrada' ? 'badge-danger' : 'badge-secondary'}`}>
                      {topic.status}
                    </span>
                  </div>

                  <div className="flex flex-col gap-2">
                    {isAuthenticated && topic.status === 'Sessão Aberta' && (
                      <Link to={`/topic/${topic.id}/vote`} className="btn btn-primary text-center">
                        Votar
                      </Link>
                    )}
                    
                    {/* Results are always available for everyone */}
                    {topic.status === 'Votação Encerrada' && (
                      <Link to={`/topic/${topic.id}/results`} className="btn btn-info text-center">
                        Ver Resultados
                      </Link>
                    )}

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