import { useState } from 'react';
import type { Topic } from '../types/Topic';
import { vote, getVoteResult, openVotingSession } from '../services/api';

interface TopicListProps {
  topics: Topic[];
  isAuthenticated: boolean;
}

export const TopicList: React.FC<TopicListProps> = ({ topics, isAuthenticated }) => {
  const [votingStates, setVotingStates] = useState<{ [key: number]: boolean }>({});
  const [voteResults, setVoteResults] = useState<{ [key: number]: { Sim: number; Não: number } }>({});
  const [sessionStates, setSessionStates] = useState<{ [key: number]: boolean }>({});
  const [error, setError] = useState<string | null>(null);

  const handleVote = async (topicId: number, choice: 'Sim' | 'Não') => {
    try {
      setError(null);
      setVotingStates(prev => ({ ...prev, [topicId]: true }));
      
      await vote(topicId, choice);
      
      // Load results after voting
      const results = await getVoteResult(topicId);
      setVoteResults(prev => ({ ...prev, [topicId]: results }));
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao votar');
    } finally {
      setVotingStates(prev => ({ ...prev, [topicId]: false }));
    }
  };

  const handleShowResults = async (topicId: number) => {
    try {
      setError(null);
      const results = await getVoteResult(topicId);
      setVoteResults(prev => ({ ...prev, [topicId]: results }));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao carregar resultados');
    }
  };

  const handleOpenSession = async (topicId: number) => {
    try {
      setError(null);
      setSessionStates(prev => ({ ...prev, [topicId]: true }));
      
      await openVotingSession(topicId, 1);
    
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao abrir sessão');
    } finally {
      setSessionStates(prev => ({ ...prev, [topicId]: false }));
    }
  };

  if (topics.length === 0) {
    return <p>Nenhum tópico encontrado.</p>;
  }

  return (
    <div>
      {error && (
        <div style={{ 
          color: 'red', 
          marginBottom: '1rem', 
          padding: '0.5rem',
          backgroundColor: '#ffebee',
          borderRadius: '4px',
          border: '1px solid #ffcdd2'
        }}>
          {error}
        </div>
      )}
      
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {topics.map((topic) => (
          <li key={topic.id} style={{ 
            marginBottom: '2rem',
            padding: '1.5rem',
            border: '1px solid #ddd',
            borderRadius: '8px',
            backgroundColor: '#fff'
          }}>
            <h3 style={{ margin: '0 0 1rem 0' }}>{topic.name}</h3>
            
            {isAuthenticated ? (
              <div>
                <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '1rem' }}>
                  <button
                    onClick={() => handleVote(topic.id, 'Sim')}
                    disabled={votingStates[topic.id]}
                    style={{
                      padding: '0.5rem 1rem',
                      backgroundColor: '#28a745',
                      color: 'white',
                      border: 'none',
                      borderRadius: '4px',
                      cursor: votingStates[topic.id] ? 'not-allowed' : 'pointer',
                      opacity: votingStates[topic.id] ? 0.6 : 1,
                      fontSize: '0.9rem'
                    }}
                  >
                    {votingStates[topic.id] ? 'Votando...' : 'Votar Sim'}
                  </button>
                  <button
                    onClick={() => handleVote(topic.id, 'Não')}
                    disabled={votingStates[topic.id]}
                    style={{
                      padding: '0.5rem 1rem',
                      backgroundColor: '#dc3545',
                      color: 'white',
                      border: 'none',
                      borderRadius: '4px',
                      cursor: votingStates[topic.id] ? 'not-allowed' : 'pointer',
                      opacity: votingStates[topic.id] ? 0.6 : 1,
                      fontSize: '0.9rem'
                    }}
                  >
                    {votingStates[topic.id] ? 'Votando...' : 'Votar Não'}
                  </button>
                </div>
                <button
                  onClick={() => handleOpenSession(topic.id)}
                  disabled={sessionStates[topic.id]}
                  style={{
                    padding: '0.5rem 1rem',
                    backgroundColor: '#ffc107',
                    color: '#212529',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: sessionStates[topic.id] ? 'not-allowed' : 'pointer',
                    opacity: sessionStates[topic.id] ? 0.6 : 1,
                    fontSize: '0.9rem',
                    marginBottom: '1rem'
                  }}
                >
                  {sessionStates[topic.id] ? 'Abrindo sessão...' : 'Abrir Sessão de Votação'}
                </button>
              </div>
            ) : (
              <p style={{ color: '#666', fontStyle: 'italic' }}>
                Faça login para votar neste tópico
              </p>
            )}
            
            <button
              onClick={() => handleShowResults(topic.id)}
              style={{
                padding: '0.5rem 1rem',
                backgroundColor: '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
                fontSize: '0.9rem',
                marginRight: '0.5rem'
              }}
            >
              Ver Resultados
            </button>
            
            {voteResults[topic.id] && (
              <div style={{ 
                marginTop: '1rem',
                padding: '1rem',
                backgroundColor: '#f8f9fa',
                borderRadius: '4px',
                border: '1px solid #e9ecef'
              }}>
                <h4 style={{ margin: '0 0 0.5rem 0' }}>Resultados:</h4>
                <p style={{ margin: '0.25rem 0' }}>
                  <strong>Sim:</strong> {voteResults[topic.id].Sim} votos
                </p>
                <p style={{ margin: '0.25rem 0' }}>
                  <strong>Não:</strong> {voteResults[topic.id].Não} votos
                </p>
                <p style={{ margin: '0.5rem 0 0 0', fontSize: '0.9rem', color: '#666' }}>
                  Total: {voteResults[topic.id].Sim + voteResults[topic.id].Não} votos
                </p>
              </div>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
};
