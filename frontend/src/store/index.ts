import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';
import topicsReducer from './topicsSlice';
import resultsReducer from './resultsSlice';
import sessionReducer from './sessionSlice';

export const store = configureStore({
  reducer: {
    auth: authReducer,
    topics: topicsReducer,
    results: resultsReducer,
    session: sessionReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch; 