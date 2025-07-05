import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { fetchTopics } from '../store/topicsSlice';
import { submitVote } from '../store/resultsSlice';

export const VotingScreen: React.FC = () => {
  const { topicId } = useParams<{ topicId: string }>();
  const [hasVoted, setHasVoted] = useState<boolean>(false);
  const [success, setSuccess] = useState<boolean>(false);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);
  const { topics, loading, error } = useAppSelector((state) => state.topics);
  const { voteLoading, voteError } = useAppSelector((state) => state.results);
  
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

  const handleVote = async (choice: 'Sim' | 'Não') => {
    if (!topicId || hasVoted) return;

    try {
      await dispatch(submitVote({
        topicId: parseInt(topicId),
        choice
      })).unwrap();
      
      setSuccess(true);
      setHasVoted(true);
    } catch (err) {
      if (err instanceof Error && err.message.includes('já votou')) {
        setHasVoted(true);
      }
      // Error is handled by Redux state
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

  if (!topic) {
    return (
      <div className="centered-page">
        <div className="card">
          <div className="card-body text-center">
            <h2 className="mb-4">Tópico não encontrado</h2>
            <p className="mb-6 text-muted">O tópico que você está tentando acessar não existe.</p>
            <Link to="/dashboard" className="btn btn-primary">
              Voltar ao Dashboard
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

        {(error || voteError) && (
          <div className="alert alert-danger">
            {error || voteError}
          </div>
        )}

        {success && (
          <div className="alert alert-success">
            Voto registrado com sucesso!
          </div>
        )}

        <div className="card">
          <div className="card-body">
            <div className="flex justify-between items-start mb-4">
              <h2>{topic.name}</h2>
              <span className={`badge ${topic.status === 'Aguardando Abertura' ? 'badge-warning' : 
                                     topic.status === 'Sessão Aberta' ? 'badge-success' : 
                                     topic.status === 'Votação Encerrada' ? 'badge-danger' : 'badge-secondary'}`}>
                {topic.status}
              </span>
            </div>
            
            {topic.status !== 'Sessão Aberta' && (
              <div className="alert alert-warning text-center">
                <h3>Votação não disponível</h3>
                <p>Esta pauta não está aberta para votação no momento.</p>
                <p><strong>Status:</strong> {topic.status}</p>
              </div>
            )}

            {topic.status === 'Sessão Aberta' && !hasVoted && !success && (
              <div>
                <h3 className="mb-6 text-center">Como você vota?</h3>
                <div className="voting-buttons">
                  <button
                    onClick={() => handleVote('Sim')}
                    disabled={voteLoading}
                    className="btn btn-success vote-btn"
                  >
                    {voteLoading ? 'Votando...' : 'SIM'}
                  </button>
                  <button
                    onClick={() => handleVote('Não')}
                    disabled={voteLoading}
                    className="btn btn-danger vote-btn"
                  >
                    {voteLoading ? 'Votando...' : 'NÃO'}
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
      </div>
    </div>
  );
}; 