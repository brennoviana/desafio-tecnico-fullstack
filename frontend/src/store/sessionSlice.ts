import { createSlice, createAsyncThunk, type PayloadAction } from '@reduxjs/toolkit';
import { openVotingSession as apiOpenVotingSession } from '../services/api';
import { updateTopicStatus } from './topicsSlice';
import type { AppDispatch } from '.';

interface SessionState {
  activeSessions: Record<number, {
    topicId: number;
    timeoutId: number | null;
    endTime: Date;
  }>;
  loading: boolean;
  error: string | null;
}

const initialState: SessionState = {
  activeSessions: {},
  loading: false,
  error: null,
};

const sessionSlice = createSlice({
  name: 'session',
  initialState,
  reducers: {
    closeVotingSession: (state, action: PayloadAction<number>) => {
      const topicId = action.payload;
      const session = state.activeSessions[topicId];
      
      if (session) {
        if (session.timeoutId) {
          window.clearTimeout(session.timeoutId);
        }
        
        delete state.activeSessions[topicId];
      }
    },
    clearError: (state) => {
      state.error = null;
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
        const { topicId, timeoutId, endTime } = action.payload;
        
        state.activeSessions[topicId] = {
          topicId,
          timeoutId,
          endTime,
        };
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
      
      const endTime = new Date(Date.now() + duration * 60 * 1000);
      
      const timeoutId = window.setTimeout(() => {
        dispatch(updateTopicStatus({ id: topicId, status: 'Votação Encerrada' }));
        dispatch(sessionSlice.actions.closeVotingSession(topicId));
      }, duration * 60 * 1000);
      
      return {
        topicId,
        timeoutId,
        endTime,
      };
    } catch (error) {
      return rejectWithValue(error instanceof Error ? error.message : 'Failed to open voting session');
    }
  }
);

export const { clearError } = sessionSlice.actions;
export default sessionSlice.reducer; 