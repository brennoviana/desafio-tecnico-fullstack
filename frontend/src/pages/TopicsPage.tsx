import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchTopics } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import type { Topic } from '../types/Topic';
import { TopicList } from '../components/TopicList';

export const TopicsPage: React.FC = () => {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    const loadTopics = async () => {
      try {
        const data = await fetchTopics();
        setTopics(data);
      } catch (err) {
        console.error(err);
        setError('Erro ao carregar tópicos.');
      } finally {
        setLoading(false);
      }
    };

    loadTopics();
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div style={{ padding: '2rem' }}>
      {/* Header with user info and logout */}
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        marginBottom: '2rem',
        padding: '1rem',
        backgroundColor: '#f8f9fa',
        borderRadius: '8px'
      }}>
        <div>
          <h1 style={{ margin: 0 }}>Lista de Tópicos</h1>
          {user && (
            <p style={{ margin: '0.5rem 0 0 0', color: '#666' }}>
              Bem-vindo, {user.name}!
            </p>
          )}
        </div>
        <button
          onClick={handleLogout}
          style={{
            padding: '0.5rem 1rem',
            backgroundColor: '#dc3545',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
            fontSize: '0.9rem'
          }}
        >
          Logout
        </button>
      </div>

      {/* Topics content */}
      {loading && <p>Carregando...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {!loading && !error && <TopicList topics={topics} />}
    </div>
  );
};
