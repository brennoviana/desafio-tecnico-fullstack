import { useEffect, useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { fetchTopics } from '../services/api';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { logoutUser } from '../store/authSlice';
import type { Topic } from '../types/Topic';
import { TopicList } from '../components/TopicList';

export const TopicsPage: React.FC = () => {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const dispatch = useAppDispatch();
  const { user, isAuthenticated } = useAppSelector((state) => state.auth);
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
    dispatch(logoutUser());
    navigate('/topics');
  };

  return (
    <div style={{ padding: '2rem' }}>
      {/* Header with conditional content */}
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
          {isAuthenticated && user && (
            <p style={{ margin: '0.5rem 0 0 0', color: '#666' }}>
              Bem-vindo, {user.name}!
            </p>
          )}
          {!isAuthenticated && (
            <p style={{ margin: '0.5rem 0 0 0', color: '#666' }}>
              Faça login para votar nos tópicos
            </p>
          )}
        </div>
        <div>
          {isAuthenticated ? (
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
          ) : (
            <div style={{ display: 'flex', gap: '0.5rem' }}>
              <Link 
                to="/login"
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor: '#007bff',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  textDecoration: 'none',
                  fontSize: '0.9rem'
                }}
              >
                Login
              </Link>
              <Link 
                to="/register"
                style={{
                  padding: '0.5rem 1rem',
                  backgroundColor: '#28a745',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  textDecoration: 'none',
                  fontSize: '0.9rem'
                }}
              >
                Registrar
              </Link>
            </div>
          )}
        </div>
      </div>

      {/* Topics content */}
      {loading && <p>Carregando...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {!loading && !error && <TopicList topics={topics} isAuthenticated={isAuthenticated} />}
    </div>
  );
};
