import { useState } from 'react';
import { createTopic } from '../services/api';
import { useAppSelector } from '../hooks/redux';

interface AddTopicButtonProps {
  onTopicAdded: () => void;
}

export const AddTopicButton: React.FC<AddTopicButtonProps> = ({ onTopicAdded }) => {
  const [isFormVisible, setIsFormVisible] = useState(false);
  const [topicName, setTopicName] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { token } = useAppSelector((state) => state.auth);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!topicName.trim()) {
      setError('Nome do tópico é obrigatório');
      return;
    }

    if (!token) {
      setError('Sessão expirada. Faça login novamente.');
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      await createTopic(topicName.trim(), token);
      setTopicName('');
      setIsFormVisible(false);
      onTopicAdded();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao criar tópico');
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    setIsFormVisible(false);
    setTopicName('');
    setError(null);
  };

  if (!isFormVisible) {
    return (
      <button
        onClick={() => setIsFormVisible(true)}
        className="btn btn-primary"
      >
        + Adicionar Tópico
      </button>
    );
  }

  return (
    <div className="card" style={{ marginBottom: '1rem' }}>
      <div className="card-body">
        <h3 className="mb-4">Novo Tópico</h3>
        
        {error && (
          <div className="alert alert-danger mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="form-group mb-4">
            <label htmlFor="topicName" className="form-label">
              Nome do Tópico:
            </label>
            <input
              type="text"
              id="topicName"
              value={topicName}
              onChange={(e) => setTopicName(e.target.value)}
              className="form-input"
              placeholder="Ex: Aprovação do novo estatuto"
              disabled={isLoading}
            />
          </div>

          <div className="flex gap-2">
            <button
              type="submit"
              disabled={isLoading}
              className="btn btn-primary"
            >
              {isLoading ? 'Criando...' : 'Criar Tópico'}
            </button>
            <button
              type="button"
              onClick={handleCancel}
              disabled={isLoading}
              className="btn btn-secondary"
            >
              Cancelar
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}; 