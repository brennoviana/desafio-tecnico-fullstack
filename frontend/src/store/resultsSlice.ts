import { createSlice, createAsyncThunk, type PayloadAction } from '@reduxjs/toolkit';
import { vote as apiVote, getVoteResult as apiGetVoteResult } from '../services/api';

interface VoteResult {
  Sim: number;
  Não: number;
}

interface TopicResult {
  topicId: number;
  result: VoteResult;
  loading: boolean;
  error: string | null;
}

interface ResultsState {
  results: Record<number, TopicResult>;
  voteLoading: boolean;
  voteError: string | null;
}

const initialState: ResultsState = {
  results: {},
  voteLoading: false,
  voteError: null,
};

export const submitVote = createAsyncThunk(
  'results/submitVote',
  async ({ topicId, choice }: { topicId: number; choice: 'Sim' | 'Não' }, { rejectWithValue }) => {
    try {
      await apiVote(topicId, choice);
      return { topicId, choice };
    } catch (error) {
      return rejectWithValue(error instanceof Error ? error.message : 'Failed to submit vote');
    }
  }
);

export const fetchVoteResults = createAsyncThunk(
  'results/fetchVoteResults',
  async (topicId: number, { rejectWithValue }) => {
    try {
      const result = await apiGetVoteResult(topicId);
      return { topicId, result };
    } catch (error) {
      return rejectWithValue(error instanceof Error ? error.message : 'Failed to fetch results');
    }
  }
);

const resultsSlice = createSlice({
  name: 'results',
  initialState,
  reducers: {
    clearVoteError: (state) => {
      state.voteError = null;
    },
    clearResultError: (state, action: PayloadAction<number>) => {
      const topicId = action.payload;
      if (state.results[topicId]) {
        state.results[topicId].error = null;
      }
    },
    updateResults: (state, action: PayloadAction<{ topicId: number; result: VoteResult }>) => {
      const { topicId, result } = action.payload;
      if (state.results[topicId]) {
        state.results[topicId].result = result;
      } else {
        state.results[topicId] = {
          topicId,
          result,
          loading: false,
          error: null,
        };
      }
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(submitVote.pending, (state) => {
        state.voteLoading = true;
        state.voteError = null;
      })
      .addCase(submitVote.fulfilled, (state) => {
        state.voteLoading = false;
        state.voteError = null;
      })
      .addCase(submitVote.rejected, (state, action) => {
        state.voteLoading = false;
        state.voteError = action.payload as string;
      })
      // Fetch vote results
      .addCase(fetchVoteResults.pending, (state, action) => {
        const topicId = action.meta.arg;
        state.results[topicId] = {
          topicId,
          result: { Sim: 0, Não: 0 },
          loading: true,
          error: null,
        };
      })
      .addCase(fetchVoteResults.fulfilled, (state, action: PayloadAction<{ topicId: number; result: VoteResult }>) => {
        const { topicId, result } = action.payload;
        state.results[topicId] = {
          topicId,
          result,
          loading: false,
          error: null,
        };
      })
      .addCase(fetchVoteResults.rejected, (state, action) => {
        const topicId = action.meta.arg;
        if (state.results[topicId]) {
          state.results[topicId].loading = false;
          state.results[topicId].error = action.payload as string;
        }
      });
  },
});

export const { clearVoteError, clearResultError, updateResults } = resultsSlice.actions;
export default resultsSlice.reducer; 