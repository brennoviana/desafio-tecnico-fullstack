import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';
import topicsReducer from './topicsSlice';
import resultsReducer from './resultsSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    topics: topicsReducer,
    results: resultsReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch; 