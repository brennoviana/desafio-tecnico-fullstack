import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getVoteResult, getSessionStatus } from '../services/api';
import { useAppSelector } from '../hooks/redux';
import type { Topic } from '../types/Topic';

interface VoteResults {
  Sim: number;
  Não: number;
}

export const ResultsScreen: React.FC = () => {
  const { topicId } = useParams<{ topicId: string }>();
  const [topic, setTopic] = useState<Topic | null>(null);
  const [results, setResults] = useState<VoteResults | null>(null);
  const [sessionData, setSessionData] = useState<{ open_at: number; close_at: number } | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const { isAuthenticated } = useAppSelector((state) => state.auth);

  useEffect(() => {
    const loadResultsData = async () => {
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

        // Load session data
        try {
          const session = await getSessionStatus(parseInt(topicId));
          setSessionData(session);
        } catch (err) {
          console.error('Session data not available:', err);
          setSessionData(null);
        }

        // Load vote results
        const voteResults = await getVoteResult(parseInt(topicId));
        setResults(voteResults);

      } catch (err) {
        console.error(err);
        setError('Erro ao carregar resultados da votação');
      } finally {
        setLoading(false);
      }
    };

    loadResultsData();
  }, [topicId]);

  const getSessionStatusText = () => {
    if (!sessionData) {
      return 'Sessão não iniciada';
    }

    const now = Date.now() / 1000;
    if (now < sessionData.open_at) {
      return 'Aguardando abertura';
    } else if (now >= sessionData.open_at && now <= sessionData.close_at) {
      return 'Sessão aberta';
    } else {
      return 'Votação encerrada';
    }
  };

  const getTotalVotes = () => {
    if (!results) return 0;
    return results.Sim + results.Não;
  };

  const getPercentage = (votes: number, total: number) => {
    if (total === 0) return 0;
    return Math.round((votes / total) * 100);
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
      </div>
    );
  }

  const total = getTotalVotes();
  const simPercentage = results ? getPercentage(results.Sim, total) : 0;
  const naoPercentage = results ? getPercentage(results.Não, total) : 0;

  return (
    <div className="page">
      <div className="container" style={{ maxWidth: '800px' }}>
        <Link to="/dashboard" className="back-link">
          ← Voltar ao Dashboard
        </Link>
        <div className="mb-6">
          <h1>Resultados da Votação</h1>
        </div>

        {error && (
          <div className="alert alert-danger">
            {error}
          </div>
        )}

        {topic && (
          <div className="card">
            <div className="card-body">
              <h2 className="mb-6">{topic.name}</h2>
              
              {/* Session Status */}
              <div className="card mb-6" style={{ background: 'var(--gray-50)' }}>
                <div className="card-body">
                  <h3 className="mb-2">Status da Sessão</h3>
                  <p className="font-semibold mb-2">
                    {getSessionStatusText()}
                  </p>
                  {sessionData && (
                    <div className="text-sm text-muted">
                      <p className="mb-1">
                        <strong>Abertura:</strong> {new Date(sessionData.open_at * 1000).toLocaleString()}
                      </p>
                      <p>
                        <strong>Fechamento:</strong> {new Date(sessionData.close_at * 1000).toLocaleString()}
                      </p>
                    </div>
                  )}
                </div>
              </div>

              {/* Vote Results */}
              {results ? (
                <div>
                  <h3 className="mb-6">Resultados</h3>
                  
                  {total === 0 ? (
                    <div className="text-center card" style={{ background: 'var(--gray-50)' }}>
                      <div className="card-body">
                        <p className="text-lg text-muted">
                          Nenhum voto registrado ainda
                        </p>
                      </div>
                    </div>
                  ) : (
                    <div>
                      {/* Vote Summary */}
                      <div className="vote-results">
                        <div className="vote-result-card vote-result-yes">
                          <h4 className="mb-2">SIM</h4>
                          <p className="vote-count" style={{ color: '#065f46' }}>
                            {results.Sim}
                          </p>
                          <p className="vote-percentage" style={{ color: '#065f46' }}>
                            {simPercentage}%
                          </p>
                        </div>

                        <div className="vote-result-card vote-result-no">
                          <h4 className="mb-2">NÃO</h4>
                          <p className="vote-count" style={{ color: '#991b1b' }}>
                            {results.Não}
                          </p>
                          <p className="vote-percentage" style={{ color: '#991b1b' }}>
                            {naoPercentage}%
                          </p>
                        </div>
                      </div>

                      {/* Visual Progress Bar */}
                      <div className="mb-6">
                        <h4 className="mb-4">Distribuição dos Votos</h4>
                        <div className="progress-bar">
                          <div
                            className="progress-segment progress-yes"
                            style={{ width: `${simPercentage}%` }}
                          >
                            {simPercentage > 15 ? `${simPercentage}%` : ''}
                          </div>
                          <div
                            className="progress-segment progress-no"
                            style={{ width: `${naoPercentage}%` }}
                          >
                            {naoPercentage > 15 ? `${naoPercentage}%` : ''}
                          </div>
                        </div>
                      </div>

                      {/* Total Votes */}
                      <div className="text-center card mb-6" style={{ background: 'var(--gray-100)' }}>
                        <div className="card-body">
                          <p className="text-xl font-bold">
                            Total de Votos: {total}
                          </p>
                        </div>
                      </div>

                      {/* Result Declaration */}
                      {total > 0 && (
                        <div className={`final-result ${results.Sim > results.Não ? 'approved' : 'rejected'}`}>
                          <h3 className="final-result-title">
                            🏆 Resultado Final
                          </h3>
                          <p className="final-result-text">
                            {results.Sim > results.Não ? 'APROVADO' : 
                             results.Não > results.Sim ? 'REJEITADO' : 'EMPATE'}
                          </p>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              ) : (
                <div className="text-center card" style={{ background: 'var(--gray-50)' }}>
                  <div className="card-body">
                    <p className="text-lg text-muted">
                      Carregando resultados...
                    </p>
                  </div>
                </div>
              )}

              {/* Action Buttons */}
              <div className="flex justify-center gap-4 mt-8 flex-mobile-col gap-mobile-4">
                {sessionData && (
                  Date.now() / 1000 >= sessionData.open_at && Date.now() / 1000 <= sessionData.close_at && isAuthenticated && (
                    <Link to={`/topic/${topicId}/vote`} className="btn btn-primary btn-lg">
                      Votar Agora
                    </Link>
                  )
                )}
                
                <button
                  onClick={() => window.location.reload()}
                  className="btn btn-secondary btn-lg"
                >
                  Atualizar Resultados
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}; 