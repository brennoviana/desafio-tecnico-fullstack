import { createSlice, createAsyncThunk, type PayloadAction } from '@reduxjs/toolkit';
import { openVotingSession as apiOpenVotingSession } from '../services/api';
import { updateTopicStatus, fetchTopics } from './topicsSlice';
import type { AppDispatch } from '.';

interface ActiveSession {
  topicId: number;
  endTime: number;
}

interface SessionState {
  activeSessions: Record<number, ActiveSession>;
  loading: boolean;
  error: string | null;
}

const STORAGE_KEY = 'voting_sessions';

// Função para salvar sessões no localStorage
const saveSessionsToStorage = (sessions: Record<number, ActiveSession>) => {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(sessions));
  } catch (error) {
    console.error('Erro ao salvar sessões no localStorage:', error);
  }
};

const loadSessionsFromStorage = (): Record<number, ActiveSession> => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    return stored ? JSON.parse(stored) : {};
  } catch (error) {
    console.error('Erro ao carregar sessões do localStorage:', error);
    return {};
  }
};

const checkExpiredSessions = () => {
  return (dispatch: AppDispatch) => {
    const now = Date.now();
    const sessions = loadSessionsFromStorage();
    const expiredTopicIds: number[] = [];
    const activeSessions: Record<number, ActiveSession> = {};
    
    Object.values(sessions).forEach(session => {
      if (now >= session.endTime) {
        expiredTopicIds.push(session.topicId);
      } else {
        activeSessions[session.topicId] = session;
      }
    });
    
    if (expiredTopicIds.length > 0) {
      
      saveSessionsToStorage(activeSessions);
      
      dispatch(sessionSlice.actions.closeMultipleSessions(expiredTopicIds));
      
      expiredTopicIds.forEach(topicId => {
        dispatch(updateTopicStatus({ id: topicId, status: 'Votação Encerrada' }));
      });
      
      dispatch(fetchTopics());
    }
  };
};

const initialState: SessionState = {
  activeSessions: loadSessionsFromStorage(),
  loading: false,
  error: null,
};

const sessionSlice = createSlice({
  name: 'session',
  initialState,
  reducers: {
    closeVotingSession: (state, action: PayloadAction<number>) => {
      const topicId = action.payload;
      delete state.activeSessions[topicId];
      saveSessionsToStorage(state.activeSessions);
    },
    closeMultipleSessions: (state, action: PayloadAction<number[]>) => {
      const topicIds = action.payload;
      topicIds.forEach(topicId => {
        delete state.activeSessions[topicId];
      });
    },
    clearError: (state) => {
      state.error = null;
    },
    initializeSessions: (state) => {
      state.activeSessions = loadSessionsFromStorage();
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(openVotingSession.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(openVotingSession.fulfilled, (state, action) => {
        state.loading = false;
        const { topicId, endTime } = action.payload;
        
        state.activeSessions[topicId] = {
          topicId,
          endTime,
        };
        
        saveSessionsToStorage(state.activeSessions);
      })
      .addCase(openVotingSession.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      });
  },
});

export const closeVotingSession = (topicId: number) => {
  return (dispatch: AppDispatch) => {
    dispatch(updateTopicStatus({ id: topicId, status: 'Votação Encerrada' }));
    dispatch(sessionSlice.actions.closeVotingSession(topicId));
  };
};

export const openVotingSession = createAsyncThunk(
  'session/openVotingSession',
  async (
    { topicId, duration }: { topicId: number; duration: number },
    { dispatch, rejectWithValue }
  ) => {
    try {
      await apiOpenVotingSession(topicId, duration);
      
      dispatch(updateTopicStatus({ id: topicId, status: 'Sessão Aberta' }));
      
      const endTime = Date.now() + duration * 60 * 1000;
      
      return {
        topicId,
        endTime,
      };
    } catch (error) {
      return rejectWithValue(error instanceof Error ? error.message : 'Failed to open voting session');
    }
  }
);

export const { clearError, initializeSessions } = sessionSlice.actions;
export { checkExpiredSessions };
export default sessionSlice.reducer; 