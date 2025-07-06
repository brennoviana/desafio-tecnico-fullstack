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
    }, 30000);

    document.addEventListener('visibilitychange', handleVisibilityChange);
    window.addEventListener('focus', handleFocus);

    return () => {
      clearInterval(intervalId);
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      window.removeEventListener('focus', handleFocus);
    };
  }, [dispatch]);

  return null;
}; 