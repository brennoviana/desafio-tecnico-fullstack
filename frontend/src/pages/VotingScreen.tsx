import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { vote } from '../services/api';
import { useAppSelector } from '../hooks/redux';
import type { Topic } from '../types/Topic';

export const VotingScreen: React.FC = () => {
  const { topicId } = useParams<{ topicId: string }>();
  const [topic, setTopic] = useState<Topic | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [voting, setVoting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);
  const [hasVoted, setHasVoted] = useState<boolean>(false);
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
        setError('Erro ao carregar dados da votação');
      } finally {
        setLoading(false);
      }
    };

    loadTopic();
  }, [topicId, isAuthenticated, navigate]);



  const handleVote = async (choice: 'Sim' | 'Não') => {
    if (!topicId || hasVoted) return;

    try {
      setError(null);
      setVoting(true);
      
      await vote(parseInt(topicId), choice);
      setSuccess(true);
      setHasVoted(true);
    } catch (err) {
      if (err instanceof Error && err.message.includes('já votou')) {
        setHasVoted(true);
        setError('Você já votou nesta pauta');
      } else {
        setError(err instanceof Error ? err.message : 'Erro ao votar');
      }
    } finally {
      setVoting(false);
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
            <p className="mb-6 text-muted">Você precisa estar logado para votar.</p>
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
          <h1>Votação</h1>
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
            Voto registrado com sucesso!
            <div className="mt-4">
              <Link to={`/topic/${topicId}/results`} className="btn btn-info">
                Ver Resultados
              </Link>
            </div>
          </div>
        )}

        {topic && (
          <div className="card">
            <div className="card-body">
              <h2 className="mb-4">{topic.name}</h2>
              
              {!hasVoted && !success && !error && (
                <div>
                  <h3 className="mb-6 text-center">Como você vota?</h3>
                  <div className="voting-buttons">
                    <button
                      onClick={() => handleVote('Sim')}
                      disabled={voting}
                      className="btn btn-success vote-btn"
                    >
                      {voting ? 'Votando...' : 'SIM'}
                    </button>
                    <button
                      onClick={() => handleVote('Não')}
                      disabled={voting}
                      className="btn btn-danger vote-btn"
                    >
                      {voting ? 'Votando...' : 'NÃO'}
                    </button>
                  </div>
                </div>
              )}

              {hasVoted && !success && (
                <div className="text-center card" style={{ background: 'var(--gray-50)' }}>
                  <div className="card-body">
                    <p className="text-lg mb-4">
                      Você já votou nesta pauta.
                    </p>
                    <Link to={`/topic/${topicId}/results`} className="btn btn-info">
                      Ver Resultados
                    </Link>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}; 