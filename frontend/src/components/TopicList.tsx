import type { Topic } from '../types/Topic';

interface TopicListProps {
  topics: Topic[];
}

export const TopicList: React.FC<TopicListProps> = ({ topics }) => {
  if (topics.length === 0) {
    return <p>Nenhum tópico encontrado.</p>;
  }

  return (
    <ul>
      {topics.map((topic) => (
        <li key={topic.id} style={{ marginBottom: '1rem' }}>
          <h3>{topic.name}</h3>
        </li>
      ))}
    </ul>
  );
};
