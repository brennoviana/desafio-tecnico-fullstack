import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { openVotingSession } from '../store/sessionSlice';
import { fetchTopics } from '../store/topicsSlice';

export const SessionManager: React.FC = () => {
  const { topicId } = useParams<{ topicId: string }>();
  const [duration, setDuration] = useState<number>(1);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);
  const { topics, loading, error } = useAppSelector((state) => state.topics);
  const { loading: sessionLoading, error: sessionError } = useAppSelector((state) => state.session);
  
  const topic = topics.find(t => t.id === parseInt(topicId || '0'));

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    if (topics.length === 0) {
      dispatch(fetchTopics());
    }
  }, [isAuthenticated, navigate, dispatch, topics.length]);

  const handleOpenSession = async () => {
    if (!topicId) return;

    try {
      await dispatch(openVotingSession({ 
        topicId: parseInt(topicId), 
        duration 
      })).unwrap();
      
      navigate('/dashboard');
    } catch (err) {
      console.error('Error opening session:', err);
    }
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="centered-page">
        <div className="card">
          <div className="card-body text-center">
            <h2 className="mb-4">Acesso Negado</h2>
            <p className="mb-6 text-muted">Você precisa estar logado para gerenciar sessões.</p>
            <Link to="/login" className="btn btn-primary">
              Fazer Login
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="container container-md">
        <Link to="/dashboard" className="back-link">
          ← Voltar ao Dashboard
        </Link>
        <div className="mb-6">
          <h1>Gerenciar Sessão de Votação</h1>
          {user && (
            <p className="text-muted">
              Logado como: {user.name}
            </p>
          )}
        </div>

        {(error || sessionError) && (
          <div className="alert alert-danger">
            {error || sessionError}
          </div>
        )}

        {sessionLoading && (
          <div className="alert alert-info">
            Abrindo sessão de votação...
          </div>
        )}

        {topic && (
          <div className="card">
            <div className="card-body">
              <h2 className="mb-6">{topic.name}</h2>
              
              <div className="mb-6">
                <h3 className="mb-4">Abrir Sessão de Votação</h3>
                
                <div className="form-group">
                  <label htmlFor="duration" className="form-label">
                    Duração da sessão (minutos):
                  </label>
                  <div className="flex items-center gap-2">
                    <input
                      type="number"
                      id="duration"
                      min="1"
                      max="60"
                      value={duration}
                      onChange={(e) => setDuration(parseInt(e.target.value) || 1)}
                      className="form-input w-input-sm"
                    />
                  </div>
                </div>

                <button
                  onClick={handleOpenSession}
                  disabled={sessionLoading}
                  className="btn btn-warning btn-lg"
                >
                  Abrir Sessão
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}; 