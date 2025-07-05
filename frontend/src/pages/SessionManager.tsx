import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { openVotingSession } from '../services/api';
import { useAppSelector } from '../hooks/redux';
import type { Topic } from '../types/Topic';

export const SessionManager: React.FC = () => {
  const { topicId } = useParams<{ topicId: string }>();
  const [topic, setTopic] = useState<Topic | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [opening, setOpening] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);
  const [duration, setDuration] = useState<number>(1);
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);
  const navigate = useNavigate();

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    const loadTopic = async () => {
      if (!topicId) {
        setError('ID do tópico não encontrado');
        setLoading(false);
        return;
      }

      try {
        // Load topic data
        const response = await fetch(`http://localhost:8080/api/topics`);
        const responseData = await response.json();
        const topicsData = responseData.data || [];
        const currentTopic = topicsData.find((t: Topic) => t.id === parseInt(topicId));
        
        if (!currentTopic) {
          setError('Tópico não encontrado');
          setLoading(false);
          return;
        }

        setTopic(currentTopic);
      } catch (err) {
        console.error(err);
        setError('Erro ao carregar dados do tópico');
      } finally {
        setLoading(false);
      }
    };

    loadTopic();
  }, [topicId, isAuthenticated, navigate]);

  const handleOpenSession = async () => {
    if (!topicId) return;

    try {
      setError(null);
      setOpening(true);
      
      await openVotingSession(parseInt(topicId), duration);
      setSuccess(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao abrir sessão');
    } finally {
      setOpening(false);
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
      <div className="container" style={{ maxWidth: '600px' }}>
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

        {error && (
          <div className="alert alert-danger">
            {error}
          </div>
        )}

        {success && (
          <div className="alert alert-success">
            Sessão de votação aberta com sucesso!
          </div>
        )}

        {topic && (
          <div className="card">
            <div className="card-body">
              <h2 className="mb-6">{topic.name}</h2>
              
              {/* Session Management */}
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
                      className="form-input"
                      style={{ width: '120px' }}
                    />
                    <span className="text-muted text-sm">
                      (máximo 60 minutos)
                    </span>
                  </div>
                </div>

                <button
                  onClick={handleOpenSession}
                  disabled={opening}
                  className="btn btn-warning btn-lg"
                >
                  {opening ? 'Abrindo sessão...' : 'Abrir Sessão de Votação'}
                </button>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-4 flex-mobile-col gap-mobile-4">
                <Link to={`/topic/${topicId}/vote`} className="btn btn-primary btn-lg">
                  Ir para Votação
                </Link>
                
                <Link to={`/topic/${topicId}/results`} className="btn btn-info btn-lg">
                  Ver Resultados
                </Link>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}; 