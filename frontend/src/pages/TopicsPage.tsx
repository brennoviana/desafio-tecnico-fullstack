import { useEffect, useState } from 'react';
import { fetchTopics } from '../services/api';
import type { Topic } from '../types/Topic';
import { TopicList } from '../components/TopicList';

export const TopicsPage: React.FC = () => {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

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

  return (
    <div style={{ padding: '2rem' }}>
      <h1>Lista de Tópicos</h1>
      {loading && <p>Carregando...</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {!loading && !error && <TopicList topics={topics} />}
    </div>
  );
};
