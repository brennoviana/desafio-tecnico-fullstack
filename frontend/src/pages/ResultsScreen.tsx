import { useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { fetchTopics } from '../store/topicsSlice';
import { fetchVoteResults } from '../store/resultsSlice';

export const ResultsScreen: React.FC = () => {
  const { topicId } = useParams<{ topicId: string }>();
  const dispatch = useAppDispatch();
  const numericTopicId = parseInt(topicId || '0');
  
  // Get state from Redux
  const { topics, loading: topicsLoading, error: topicsError } = useAppSelector((state) => state.topics);
  const { results } = useAppSelector((state) => state.results);
  
  const topic = topics.find(t => t.id === numericTopicId);
  const topicResult = results[numericTopicId];
  const voteResults = topicResult?.result;
  const resultsLoading = topicResult?.loading || false;
  const resultsError = topicResult?.error || null;

  useEffect(() => {
    // Load topics if not already loaded
    if (topics.length === 0) {
      dispatch(fetchTopics());
    }
    
    // Load vote results
    if (topicId) {
      dispatch(fetchVoteResults(numericTopicId));
    }
  }, [dispatch, topicId, numericTopicId, topics.length]);

  const getTotalVotes = () => {
    if (!voteResults) return 0;
    return voteResults.Sim + voteResults.Não;
  };

  const getPercentage = (votes: number, total: number) => {
    if (total === 0) return 0;
    return Math.round((votes / total) * 100);
  };

  if (topicsLoading || resultsLoading) {
    return (
      <div className="loading">
        <div className="spinner"></div>
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

  const total = getTotalVotes();
  const simPercentage = voteResults ? getPercentage(voteResults.Sim, total) : 0;
  const naoPercentage = voteResults ? getPercentage(voteResults.Não, total) : 0;

  return (
    <div className="page">
      <div className="container" style={{ maxWidth: '800px' }}>
        <Link to="/dashboard" className="back-link">
          ← Voltar ao Dashboard
        </Link>
        <div className="mb-6">
          <h1>Resultados da Votação</h1>
        </div>

        {(topicsError || resultsError) && (
          <div className="alert alert-danger">
            {topicsError || resultsError}
          </div>
        )}

        <div className="card">
          <div className="card-body">
            <div className="flex justify-between items-start mb-6">
              <h2>{topic.name}</h2>
              <span className={`badge ${topic.status === 'Aguardando Abertura' ? 'badge-warning' : 
                                     topic.status === 'Sessão Aberta' ? 'badge-success' : 
                                     topic.status === 'Votação Encerrada' ? 'badge-danger' : 'badge-secondary'}`}>
                {topic.status}
              </span>
            </div>
            
            {/* Vote Results */}
            {voteResults ? (
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
                          {voteResults.Sim}
                        </p>
                        <p className="vote-percentage" style={{ color: '#065f46' }}>
                          {simPercentage}%
                        </p>
                      </div>

                      <div className="vote-result-card vote-result-no">
                        <h4 className="mb-2">NÃO</h4>
                        <p className="vote-count" style={{ color: '#991b1b' }}>
                          {voteResults.Não}
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
                      <div className={`final-result ${voteResults.Sim > voteResults.Não ? 'approved' : 'rejected'}`}>
                        <h3 className="final-result-title">
                          🏆 Resultado Final
                        </h3>
                        <p className="final-result-text">
                          {voteResults.Sim > voteResults.Não ? 'APROVADO' : 
                           voteResults.Não > voteResults.Sim ? 'REJEITADO' : 'EMPATE'}
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
          </div>
        </div>
      </div>
    </div>
  );
}; 