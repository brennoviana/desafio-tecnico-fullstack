import { useEffect } from 'react';
import { useAppDispatch } from './redux';
import { checkExpiredSessions } from '../store/sessionSlice';

export const useSessionChecker = () => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(checkExpiredSessions());

    const handleVisibilityChange = () => {
      if (!document.hidden) {
        dispatch(checkExpiredSessions());
      }
    };

    const handleFocus = () => {
      dispatch(checkExpiredSessions());
    };

    const intervalId = setInterval(() => {
      dispatch(checkExpiredSessions());
    }, 15000);

    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === 'voting_sessions') {
        dispatch(checkExpiredSessions());
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    window.addEventListener('focus', handleFocus);
    window.addEventListener('storage', handleStorageChange);

    return () => {
      clearInterval(intervalId);
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      window.removeEventListener('focus', handleFocus);
      window.removeEventListener('storage', handleStorageChange);
    };
  }, [dispatch]);

  return null;
}; 