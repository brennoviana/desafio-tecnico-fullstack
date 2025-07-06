import { useState } from 'react';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { createTopic } from '../store/topicsSlice';

export const AddTopicButton: React.FC = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [topicName, setTopicName] = useState('');
  const dispatch = useAppDispatch();
  
  const { token } = useAppSelector((state) => state.auth);
  const { createLoading, createError } = useAppSelector((state) => state.topics);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!topicName.trim() || !token) {
      return;
    }

    try {
      await dispatch(createTopic(topicName.trim())).unwrap();
      
      window.location.reload();
    } catch (err) {
      console.error('Failed to create topic:', err);
    }
  };

  const handleCancel = () => {
    setIsModalOpen(false);
    setTopicName('');
  };

  return (
    <>
      <button
        onClick={() => setIsModalOpen(true)}
        className="btn btn-primary"
      >
        + Adicionar Tópico
      </button>

      {isModalOpen && (
        <div 
          className="modal-overlay"
          onClick={handleCancel}
        >
          <div 
            className="modal-content"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="modal-title">Novo Tópico</h3>
            
            {createError && (
              <div className="alert alert-danger">
                {createError}
              </div>
            )}

            <form onSubmit={handleSubmit}>
              <div className="modal-form-group">
                <label htmlFor="topicName" className="modal-form-label">
                  Nome do Tópico:
                </label>
                <input
                  type="text"
                  id="topicName"
                  value={topicName}
                  onChange={(e) => setTopicName(e.target.value)}
                  className="modal-form-input"
                  placeholder="Ex: Aprovação do novo estatuto"
                  disabled={createLoading}
                  required
                  autoFocus
                />
              </div>

              <div className="modal-buttons">
                <button
                  type="button"
                  onClick={handleCancel}
                  disabled={createLoading}
                  className="modal-btn modal-btn-secondary"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={createLoading || !topicName.trim()}
                  className="modal-btn modal-btn-primary"
                >
                  {createLoading ? 'Criando...' : 'Criar Tópico'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
}; 