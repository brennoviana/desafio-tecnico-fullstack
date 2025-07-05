import { useState } from 'react';
import { useAppDispatch, useAppSelector } from '../hooks/redux';
import { createTopic } from '../store/topicsSlice';

interface AddTopicButtonProps {
  onTopicAdded?: () => void;
}

export const AddTopicButton: React.FC<AddTopicButtonProps> = () => {
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
      {/* The button that stays in the header */}
      <button
        onClick={() => setIsModalOpen(true)}
        className="btn btn-primary"
      >
        + Adicionar Tópico
      </button>

      {/* Modal overlay */}
      {isModalOpen && (
        <div 
          className="modal-overlay"
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0, 0, 0, 0.5)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 1000
          }}
          onClick={handleCancel}
        >
          <div 
            className="modal-content"
            style={{
              backgroundColor: 'white',
              padding: '2rem',
              borderRadius: '8px',
              minWidth: '400px',
              maxWidth: '500px',
              width: '90%',
              boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)'
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <h3 style={{ marginTop: 0, marginBottom: '1.5rem' }}>Novo Tópico</h3>
            
            {createError && (
              <div className="alert alert-danger" style={{ marginBottom: '1.5rem' }}>
                {createError}
              </div>
            )}

            <form onSubmit={handleSubmit}>
              <div style={{ marginBottom: '1.5rem' }}>
                <label htmlFor="topicName" style={{ display: 'block', marginBottom: '0.5rem', fontWeight: '500' }}>
                  Nome do Tópico:
                </label>
                <input
                  type="text"
                  id="topicName"
                  value={topicName}
                  onChange={(e) => setTopicName(e.target.value)}
                  className="form-input"
                  style={{
                    width: '100%',
                    padding: '0.75rem',
                    border: '1px solid #ddd',
                    borderRadius: '4px',
                    fontSize: '1rem'
                  }}
                  placeholder="Ex: Aprovação do novo estatuto"
                  disabled={createLoading}
                  required
                  autoFocus
                />
              </div>

              <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'flex-end' }}>
                <button
                  type="button"
                  onClick={handleCancel}
                  disabled={createLoading}
                  className="btn btn-secondary"
                  style={{
                    padding: '0.75rem 1.5rem',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: 'pointer',
                    backgroundColor: '#6c757d',
                    color: 'white'
                  }}
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  disabled={createLoading || !topicName.trim()}
                  className="btn btn-primary"
                  style={{
                    padding: '0.75rem 1.5rem',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: createLoading || !topicName.trim() ? 'not-allowed' : 'pointer',
                    backgroundColor: createLoading || !topicName.trim() ? '#ccc' : '#007bff',
                    color: 'white'
                  }}
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