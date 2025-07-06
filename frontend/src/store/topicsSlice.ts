import { createSlice, createAsyncThunk, type PayloadAction } from '@reduxjs/toolkit';
import { fetchTopics as apiFetchTopics, createTopic as apiCreateTopic } from '../services/api';
import type { Topic } from '../types/Topic';

interface TopicsState {
  topics: Topic[];
  loading: boolean;
  error: string | null;
  createLoading: boolean;
  createError: string | null;
}

const initialState: TopicsState = {
  topics: [],
  loading: false,
  error: null,
  createLoading: false,
  createError: null,
};

export const fetchTopics = createAsyncThunk(
  'topics/fetchTopics',
  async (_, { rejectWithValue }) => {
    try {
      const topics = await apiFetchTopics();
      return topics;
    } catch (error) {
      return rejectWithValue(error instanceof Error ? error.message : 'Failed to fetch topics');
    }
  }
);

export const createTopic = createAsyncThunk(
  'topics/createTopic',
  async (name: string, { rejectWithValue }) => {
    try {
      const topic = await apiCreateTopic(name);
      return topic;
    } catch (error) {
      return rejectWithValue(error instanceof Error ? error.message : 'Failed to create topic');
    }
  }
);

const topicsSlice = createSlice({
  name: 'topics',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    clearCreateError: (state) => {
      state.createError = null;
    },
    updateTopicStatus: (state, action: PayloadAction<{ id: number; status: string }>) => {
      const topic = state.topics.find(t => t.id === action.payload.id);
      if (topic) {
        topic.status = action.payload.status;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchTopics.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchTopics.fulfilled, (state, action: PayloadAction<Topic[]>) => {
        state.loading = false;
        state.topics = action.payload;
        state.error = null;
      })
      .addCase(fetchTopics.rejected, (state, action) => {
        state.loading = false;
        state.error = action.payload as string;
      })
      .addCase(createTopic.pending, (state) => {
        state.createLoading = true;
        state.createError = null;
      })
      .addCase(createTopic.fulfilled, (state, action: PayloadAction<Topic>) => {
        state.createLoading = false;
        state.topics.push(action.payload);
        state.createError = null;
      })
      .addCase(createTopic.rejected, (state, action) => {
        state.createLoading = false;
        state.createError = action.payload as string;
      });
  },
});

export const { clearError, clearCreateError, updateTopicStatus } = topicsSlice.actions;
export default topicsSlice.reducer; 