import { useEffect } from 'react';
import { useAppDispatch } from '../hooks/redux';
import { initializeAuth } from '../store/authSlice';

export const AuthInitializer: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(initializeAuth());
  }, [dispatch]);

  return <>{children}</>;
}; 