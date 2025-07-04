import React, { useEffect, useState } from "react";
import api from "../api/axios";
import { useAuth } from "../context/AuthContext";

interface Topic {
  id: number;
  title: string;
  description: string;
}

const Topics: React.FC = () => {
  const [topics, setTopics] = useState<Topic[]>([]);
  const [loading, setLoading] = useState(true);
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    const fetchTopics = async () => {
      try {
        const response = await api.get("/topics");
        setTopics(response.data);
      } catch (err) {
        // Trate o erro conforme necessário
      } finally {
        setLoading(false);
      }
    };
    fetchTopics();
  }, []);

  if (loading) return <div>Carregando tópicos...</div>;

  return (
    <div>
      <h2>Tópicos</h2>
      {isAuthenticated && <button>Criar novo tópico</button>}
      <ul>
        {topics.map((topic) => (
          <li key={topic.id}>
            <strong>{topic.title}</strong> - {topic.description}
            {isAuthenticated && <button style={{ marginLeft: 8 }}>Votar</button>}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Topics; 